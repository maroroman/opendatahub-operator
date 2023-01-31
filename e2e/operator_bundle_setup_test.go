package e2e

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	kfdefv1 "github.com/opendatahub-io/opendatahub-operator/apis/kfdef.apps.kubeflow.org/v1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sclient "k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlruntime "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	testNamespace string
	scheme        = runtime.NewScheme()
)

// Holds information specific to individual tests
type testContext struct {
	// Rest config
	cfg *rest.Config
	// client for k8s resources
	kubeClient *k8sclient.Clientset
	// custom client for managing custom resources
	customClient client.Client
	// namespace for running the tests
	testNamespace string
	// context for accessing resources
	ctx   context.Context
	kfDef kfdefv1.KfDef
}

func NewTestContext() (*testContext, error) {

	// GetConfig(): If KUBECONFIG env variable is set, it is used to create
	// the client, else the inClusterConfig() is used.
	// Lastly if none of the them are set, it uses  $HOME/.kube/config to create the client.
	config, err := ctrlruntime.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating the config object %v", err)
	}

	kc, err := k8sclient.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Kubernetes client")
	}

	// custom client to manages resources like KfDef, Route etc
	custClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize custom client")
	}

	return &testContext{
		cfg:           config,
		kubeClient:    kc,
		customClient:  custClient,
		testNamespace: testNamespace,
		ctx:           context.TODO(),
		kfDef:         setupCoreKfDef(),
	}, nil
}

// THis function sets up the testing suit for the operator bundle
func TestE2EBundle(t *testing.T) {

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kfdefv1.AddToScheme(scheme))
	// individual test suites after the operator is running
	// if !t.Run("validate controllers", testNotebookControllerValidation) {
	// 	return
	// }
	// Run create and delete tests for all the test notebooks
	t.Run("create", creationTestSuite)
	// if !skipDeletion {
	// 	t.Run("delete", deletionTestSuite)
	// }
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.StringVar(&testNamespace, "namespace",
		"opendatahub", "Custom namespace where the operator is deployed")
	flag.Parse()
	os.Exit(m.Run())
}

func setupCoreKfDef() kfdefv1.KfDef {
	result := kfdefv1.KfDef{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KfDef",
			APIVersion: "kfdef.apps.kubeflow.org/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "operator-e2e",
			Namespace: testNamespace,
		},
		Spec: kfdefv1.KfDefSpec{
			Version: "master",
			Applications: []kfdefv1.Application{
				{
					Name: "odh-common",
					KustomizeConfig: &kfdefv1.KustomizeConfig{
						RepoRef: &kfdefv1.RepoRef{
							Name: "manifests",
							Path: "odh-common",
						},
					},
				},
				{
					Name: "odh-dashboard",
					KustomizeConfig: &kfdefv1.KustomizeConfig{
						RepoRef: &kfdefv1.RepoRef{
							Name: "manifests",
							Path: "odh-dashboard",
						},
					},
				},
				{
					Name: "odh-notebook-controller",
					KustomizeConfig: &kfdefv1.KustomizeConfig{
						RepoRef: &kfdefv1.RepoRef{
							Name: "manifests",
							Path: "odh-notebook-controller",
						},
					},
				},
				{
					Name: "notebook-images",
					KustomizeConfig: &kfdefv1.KustomizeConfig{
						RepoRef: &kfdefv1.RepoRef{
							Name: "manifests",
							Path: "notebook-images",
						},
					},
				},
			},
			Repos: []kfdefv1.Repo{
				{
					Name: "manifests",
					URI:  "https://github.com/opendatahub-io/odh-manifests/tarball/master", // TODO: Make this a variable
				},
			},
		},
	}

	return result
}
