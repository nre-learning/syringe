package kubernetes

import (
	"testing"

	ot "github.com/opentracing/opentracing-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSyncSecret(t *testing.T) {
	span := ot.StartSpan("")
	defer span.Finish()

	k := createFakeKubernetesBackend()
	k.Config.SecretsNamespace = "prod"
	k.Config.PullCredName = "docker-pull-creds"

	_, err := k.Client.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: k.Config.SecretsNamespace,
			Labels: map[string]string{
				"name": k.Config.SecretsNamespace,
			},
		},
	})
	ok(t, err)
	_, err = k.Client.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testns",
			Labels: map[string]string{
				"name": "testns",
			},
		},
	})
	ok(t, err)

	_, err = k.Client.CoreV1().Secrets("prod").Create(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: k.Config.PullCredName,
		},
		Type: "kubernetes.io/dockerconfigjson",
		Data: map[string][]byte{
			".dockerconfigjson": {1, 2, 3},
		},
	})
	ok(t, err)

	err = k.syncSecret(span.Context(), k.Config.SecretsNamespace, "testns", k.Config.PullCredName)
	ok(t, err)

	syncedSecret, err := k.Client.CoreV1().Secrets("testns").Get(k.Config.PullCredName, metav1.GetOptions{})
	ok(t, err)

	assert(t, syncedSecret.Type == "kubernetes.io/dockerconfigjson", "")
}
