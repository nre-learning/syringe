package scheduler

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// nukeFromOrbit seeks out all syringe-managed namespaces, and deletes them.
// This will effectively reset the cluster to a state with all of the remaining infrastructure
// in place, but no running labs. Syringe doesn't manage itself, or any other Antidote services.
func (ls *LabScheduler) nukeFromOrbit() error {

	coreclient, err := corev1client.NewForConfig(ls.Config)
	if err != nil {
		panic(err)
	}
	nameSpaces, err := coreclient.Namespaces().List(metav1.ListOptions{
		// VERY Important. Only delete those with this label, otherwise you'll nuke the cluster.
		LabelSelector: "syringeManaged",
	})
	if err != nil {
		return err
	}

	// No need to nuke if no syringe namespaces exist
	if len(nameSpaces.Items) == 0 {
		log.Info("No syringe-managed namespaces found. Starting normally.")
		return nil
	}

	log.Warn("Nuking all syringe-managed namespaces")
	var wg sync.WaitGroup
	wg.Add(len(nameSpaces.Items))
	for n := range nameSpaces.Items {

		nsName := nameSpaces.Items[n].ObjectMeta.Name
		go func() {
			defer wg.Done()
			ls.deleteNamespace(nsName)
		}()
	}
	wg.Wait()
	log.Info("Nuke complete. It was the only way to be sure...")
	return nil
}

func (ls *LabScheduler) deleteNamespace(name string) error {

	coreclient, err := corev1client.NewForConfig(ls.Config)
	if err != nil {
		panic(err)
	}

	err = coreclient.Namespaces().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	// Wait for the namespace to be deleted
	deleteTimeoutSecs := 120
	for i := 0; i < deleteTimeoutSecs/5; i++ {
		time.Sleep(5 * time.Second)

		_, err := coreclient.Namespaces().Get(name, metav1.GetOptions{})
		if err == nil {
			log.Debugf("Waiting for namespace %s to delete...", name)
			continue
		} else if apierrors.IsNotFound(err) {
			log.Infof("Deleted namespace %s", name)
			return nil
		} else {
			return err
		}
	}

	errorMsg := fmt.Sprintf("Timed out trying to delete namespace %s", name)
	log.Error(errorMsg)
	return errors.New(errorMsg)
}

func (ls *LabScheduler) createNamespace(req *LabScheduleRequest) (*corev1.Namespace, error) {

	coreclient, err := corev1client.NewForConfig(ls.Config)
	if err != nil {
		panic(err)
	}

	nsName := fmt.Sprintf("%d-%s-ns", req.LabDef.LabID, req.Session)

	log.Infof("Creating namespace: %s", req.Session)

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsName,
			Labels: map[string]string{
				"labId":          fmt.Sprintf("%d", req.LabDef.LabID),
				"sessionId":      req.Session,
				"syringeManaged": "yes",
			},
			Namespace: nsName,
		},
	}

	result, err := coreclient.Namespaces().Create(namespace)
	if err == nil {
		log.Infof("Created namespace: %s", result.ObjectMeta.Name)
	} else if apierrors.IsAlreadyExists(err) {
		log.Warnf("Namespace %s already exists.", nsName)

		// In this case we are returning what we tried to create. This means that when this lab is cleaned up,
		// syringe will delete the pod that already existed.
		return namespace, err
	} else {
		log.Errorf("Problem creating namespace %s: %s", nsName, err)
		return nil, err
	}
	return result, err
}
