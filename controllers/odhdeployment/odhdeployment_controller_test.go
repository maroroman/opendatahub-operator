package odhdeployment

import (
	//"github.com/go-logr/logr"
	//"github.com/golang-jwt/jwt/v4/request"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	//"sigs.k8s.io/kustomize/api/filters/namespace"
)



var _ = Describe("OdhdeploymentController", func() {
	Describe("Test odhdeployment reconciliation", func() {
		Context("with standard odhdeployment", func() {
			const (
				Name = "test-odhdeployment"
				Namespace = "test-namespace"
			)
			It("should reconcile with no error", func() {
				request := &ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: Namespace,
						Name: Name,
					},
				}
				_, err  := reconciler.Reconcile(ctx, *request)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
