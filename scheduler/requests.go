package scheduler

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	models "github.com/nre-learning/antidote-core/db/models"
	"github.com/nre-learning/antidote-core/services"
)

func (s *AntidoteScheduler) handleRequestCREATE(newRequest services.LessonScheduleRequest) {

	nsName := generateNamespaceName(s.Config.InstanceID, newRequest.LiveLessonID)

	ll, err := s.Db.GetLiveLesson(newRequest.LiveLessonID)
	if err != nil {
		log.Errorf("Error getting livelesson: %v", err)
		return
	}

	err = s.createK8sStuff(newRequest)
	if err != nil {
		log.Errorf("Error creating lesson: %v", err)
		return
	}

	log.Debugf("Bootstrap complete for livelesson %s. Moving into BOOTING status", newRequest.LiveLessonID)
	err = s.Db.UpdateLiveLessonStatus(ll.ID, models.Status_BOOTING)
	if err != nil {
		log.Errorf("Error updating livelesson: %v", err)
		return
	}

	var success = false
	for i := 0; i < 600; i++ {
		time.Sleep(1 * time.Second)

		log.Debugf("About to test endpoint reachability for livelesson %s with endpoints %v", ll.ID, ll.LiveEndpoints)

		reachability := s.testEndpointReachability(ll)

		log.Debugf("Livelesson %s health check results: %v", ll.ID, reachability)

		// Update reachability status
		failed := false
		healthy := int32(0)
		total := int32(len(reachability))
		for _, reachable := range reachability {
			if reachable {
				healthy++
			} else {
				failed = true
			}
		}
		ll.HealthyTests = healthy
		ll.TotalTests = total

		// Begin again if one of the endpoints isn't reachable
		if failed {
			continue
		}

		success = true
		break

	}

	if !success {
		log.Errorf("Timeout waiting for livelesson %s to become reachable", ll.ID)
		err = s.Db.UpdateLiveLessonError(ll.ID, true)
		if err != nil {
			log.Errorf("Error updating livelesson: %v", err)
		}
		return
	}

	log.Debugf("Setting status for livelesson %s to CONFIGURATION", newRequest.LiveLessonID)
	err = s.Db.UpdateLiveLessonStatus(ll.ID, models.Status_CONFIGURATION)
	if err != nil {
		log.Errorf("Error updating livelesson %s: %v", ll.ID, err)
	}

	log.Infof("Performing configuration for livelesson %s", ll.ID)
	err = s.configureStuff(nsName, ll, newRequest)
	if err != nil {
		log.Errorf("Error configuring livelesson %s: %v", ll.ID, err)
		err = s.Db.UpdateLiveLessonError(ll.ID, true)
		if err != nil {
			log.Errorf("Error updating livelesson %s: %v", ll.ID, err)
		}
		return
	}

	// Set network policy ONLY after configuration has had a chance to take place. Once this is in place,
	// only config pods spawned by Jobs will have internet access, so if this takes place earlier, lessons
	// won't initially come up at all.
	if s.Config.AllowEgress {
		s.createNetworkPolicy(nsName)
	}

	log.Debugf("Setting livelesson %s to READY", newRequest.LiveLessonID)
	err = s.Db.UpdateLiveLessonStatus(ll.ID, models.Status_READY)
	if err != nil {
		log.Errorf("Error updating livelesson %s: %v", ll.ID, err)
	}

	s.NEC.Publish("antidote.lsr.completed", newRequest)
}

func (s *AntidoteScheduler) handleRequestMODIFY(newRequest services.LessonScheduleRequest) {

	nsName := generateNamespaceName(s.Config.InstanceID, newRequest.LiveLessonID)

	ll, err := s.Db.GetLiveLesson(newRequest.LiveLessonID)
	if err != nil {
		log.Errorf("Error getting livelesson: %v", err)
		return
	}

	log.Infof("Performing configuration for stage %d of livelesson %s", newRequest.Stage, newRequest.LiveLessonID)
	err = s.configureStuff(nsName, ll, newRequest)
	if err != nil {
		log.Errorf("Error configuring livelesson %s: %v", ll.ID, err)
		err = s.Db.UpdateLiveLessonError(ll.ID, true)
		if err != nil {
			log.Errorf("Error updating livelesson %s: %v", ll.ID, err)
		}
		return
	}

	log.Debugf("Setting livelesson %s to READY", newRequest.LiveLessonID)
	err = s.Db.UpdateLiveLessonStatus(ll.ID, models.Status_READY)
	if err != nil {
		log.Errorf("Error updating livelesson %s: %v", ll.ID, err)
	}

	err = s.boopNamespace(nsName)
	if err != nil {
		log.Errorf("Problem modify-booping %s: %v", nsName, err)
	}
}

func (s *AntidoteScheduler) handleRequestBOOP(newRequest services.LessonScheduleRequest) {
	nsName := generateNamespaceName(s.Config.InstanceID, newRequest.LiveLessonID)

	err := s.boopNamespace(nsName)
	if err != nil {
		log.Errorf("Problem booping %s: %v", nsName, err)
	}
}

// handleRequestDELETE handles a livelesson deletion request by first sending a delete request
// for the corresponding namespace, and then cleaning up local state.
func (s *AntidoteScheduler) handleRequestDELETE(newRequest services.LessonScheduleRequest) {
	nsName := generateNamespaceName(s.Config.InstanceID, newRequest.LiveLessonID)
	err := s.deleteNamespace(nsName)
	if err != nil {
		log.Errorf("Unable to delete namespace %s: %v", nsName, err)
		return
	}
	err = s.Db.DeleteLiveLesson(newRequest.LiveLessonID)
	if err != nil {
		log.Errorf("Error getting livelesson: %v", err)
	}
}

// createK8sStuff is a high-level workflow for creating all of the things necessary for a new instance
// of a livelesson. Pods, services, networks, networkpolicies, ingresses, etc to support a new running
// lesson are all created as part of this workflow.
func (s *AntidoteScheduler) createK8sStuff(req services.LessonScheduleRequest) error {

	ns, err := s.createNamespace(req)
	if err != nil {
		log.Error(err)
	}

	err = s.syncSecret(ns.ObjectMeta.Name)
	if err != nil {
		log.Errorf("Unable to sync secret into this namespace. Ingress-based resources (like http presentations or jupyter notebooks) may not work: %v", err)
	}

	lesson, err := s.Db.GetLesson(req.LessonSlug)
	if err != nil {
		return err
	}

	ll, err := s.Db.GetLiveLesson(req.LiveLessonID)
	if err != nil {
		return err
	}

	// Append endpoint and create ingress for jupyter lab guide if necessary
	if usesJupyterLabGuide(lesson) {

		jupyterEp := &models.LiveEndpoint{
			Name:  "jupyterlabguide",
			Image: fmt.Sprintf("antidotelabs/jupyter:%s", s.BuildInfo["imageVersion"]),
			Ports: []int32{8888},
		}
		ll.LiveEndpoints[jupyterEp.Name] = jupyterEp

		_, err := s.createIngress(
			ns.ObjectMeta.Name,
			jupyterEp,
			&models.LivePresentation{
				Name: "web",
				Port: 8888,
			},
		)
		if err != nil {
			return fmt.Errorf("Unable to create ingress resource - %v", err)
		}
	}

	// Create networks from connections property
	for c := range lesson.Connections {
		connection := lesson.Connections[c]
		_, err := s.createNetwork(c, fmt.Sprintf("%s-%s-net", connection.A, connection.B), req)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	createdPods := make(map[string]*corev1.Pod)

	// Create pods and services
	for d := range ll.LiveEndpoints {
		ep := ll.LiveEndpoints[d]

		// createPod doesn't try to ensure a certain pod status. That's done later
		newPod, err := s.createPod(
			ep,
			getMemberNetworks(ep.Name, lesson.Connections),
			req,
		)
		if err != nil {
			log.Error(err)
			return err
		}

		// TODO(mierdin): get data on created pods like which host they were deployed on. Make sure this is placed in the trace/span

		createdPods[newPod.ObjectMeta.Name] = newPod

		// Expose via service if needed
		if len(newPod.Spec.Containers[0].Ports) > 0 {
			svc, err := s.createService(
				newPod,
				req,
			)
			if err != nil {
				log.Error(err)
				return err
			}

			// Update livelesson liveendpoint with cluster IP
			s.Db.UpdateLiveLessonEndpointIP(ll.ID, ep.Name, svc.Spec.ClusterIP)
			if err != nil {
				log.Error("Unable to update livelesson endpoint with clusterIP")
			}

		}

		// Create ingresses for http presentations
		for pr := range ep.Presentations {
			p := ep.Presentations[pr]

			if p.Type == "http" {
				_, err := s.createIngress(
					ns.ObjectMeta.Name,
					ep,
					p,
				)
				if err != nil {
					log.Error(err)
					return err
				}
			}
		}
	}

	// Before moving forward with network-based health checks, let's look back at the pods
	// we've deployed, and wait until they're in a "Running" status. This allows us to keep a hold
	// of maximum amounts of context for troubleshooting while we have it
	wg := new(sync.WaitGroup)
	wg.Add(len(createdPods))

	failLesson := false

	for name, pod := range createdPods {
		go func() {
			defer wg.Done()

			for i := 0; i < 300; i++ {

				rdy, err := s.getPodStatus(pod)
				if err != nil {
					failLesson = true
					return
				}

				if rdy {
					delete(createdPods, name)
					return
				}

				time.Sleep(1 * time.Second)
			}

			log.Infof("Timed out waiting for %s to start", name)
			failLesson = true
			return
		}()
	}

	wg.Wait()

	// At this point, the only pods left in createdPods should be ones that failed to ready
	if failLesson || len(createdPods) > 0 {
		log.Infof("Failed pods: %v", createdPods)
		return errors.New("Some pods failed to start")
	}

	return nil
}
