package e2e

import (
	"fmt"
	"testing"

	kfdefv1 "github.com/opendatahub-io/opendatahub-operator/apis/kfdef.apps.kubeflow.org/v1"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func creationTestSuite(t *testing.T) {
	testCtx, err := NewTestContext()
	require.NoError(t, err)
	t.Run(testCtx.kfDef.ObjectMeta.Name, func(t *testing.T) {
		t.Run("Creation of Notebook instance", func(t *testing.T) {
			err = testCtx.testKfDefCreation()
			require.NoError(t, err, "error creating KfDef ")
		})
	})
}

func (testCtx *testContext) testKfDefCreation() error {
	testKfDef := testCtx.kfDef
	kfDefLookup := types.NamespacedName{Name: testKfDef.ObjectMeta.Name, Namespace: testKfDef.ObjectMeta.Namespace}
	createdKfDef := kfdefv1.KfDef{}

	err := testCtx.customClient.Get(testCtx.ctx, kfDefLookup, &createdKfDef)

	if err != nil {
		if errors.IsNotFound(err) {
			createErr := testCtx.customClient.Create(testCtx.ctx, &testKfDef)
			if createErr != nil {
				return fmt.Errorf("Error creating KfDef %v: %v", testKfDef.ObjectMeta.Name, createErr)
			}
		} else {
			return fmt.Errorf("Error Getting KfDef %v: %v", testKfDef.ObjectMeta.Name, err)
		}
	}

	return nil
}

// func (tc *testContext) testNotebookCreation(nbContext notebookContext) error {

// 	testNotebook := &nbv1.Notebook{
// 		ObjectMeta: *nbContext.nbObjectMeta,
// 		Spec:       *nbContext.nbSpec,
// 	}

// 	// Create test Notebook resource if not already created
// 	notebookLookupKey := types.NamespacedName{Name: testNotebook.Name, Namespace: testNotebook.Namespace}
// 	createdNotebook := nbv1.Notebook{}

// 	err := tc.customClient.Get(tc.ctx, notebookLookupKey, &createdNotebook)
// 	if err != nil {
// 		if errors.IsNotFound(err) {
// 			nberr := wait.Poll(tc.resourceRetryInterval, tc.resourceCreationTimeout, func() (done bool, err error) {
// 				creationErr := tc.customClient.Create(tc.ctx, testNotebook)
// 				if creationErr != nil {
// 					log.Printf("Error creating Notebook resource %v: %v, trying again",
// 						testNotebook.Name, creationErr)
// 					return false, nil
// 				} else {
// 					return true, nil
// 				}
// 			})
// 			if nberr != nil {
// 				return fmt.Errorf("error creating test Notebook %s: %v", testNotebook.Name, nberr)
// 			}
// 		} else {
// 			return fmt.Errorf("error getting test Notebook %s: %v", testNotebook.Name, err)
// 		}
// 	}
// 	return nil
// }

// func (tc *testContext) testNotebookValidation(nbMeta *metav1.ObjectMeta) error {
// 	// Verify StatefulSet is running
// 	err := wait.Poll(tc.resourceRetryInterval, tc.resourceCreationTimeout, func() (done bool, err error) {
// 		notebookStatefulSet, err1 := tc.kubeClient.AppsV1().StatefulSets(tc.testNamespace).Get(tc.ctx,
// 			nbMeta.Name, metav1.GetOptions{})

// 		if err1 != nil {
// 			if errors.IsNotFound(err1) {
// 				return false, nil
// 			} else {
// 				log.Printf("Failed to get %s statefulset", nbMeta.Name)
// 				return false, err1
// 			}
// 		}
// 		if notebookStatefulSet.Status.AvailableReplicas == 1 &&
// 			notebookStatefulSet.Status.ReadyReplicas == 1 {
// 			return true, nil
// 		}
// 		return false, nil
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error validating notebook StatefulSet: %s", err)
// 	}
// 	return nil
// }
