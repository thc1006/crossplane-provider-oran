package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hardwarev1alpha1 "ran.example.com/o-ran-adapter/internal/api/v1alpha1"
)

var _ = Describe("OpticalDevice Controller", func() {
	const (
		DeviceName      = "test-device"
		DeviceNamespace = "default"
		DeviceHostname  = "laser-a1.site1"
		DeviceBandwidth = "100Gbps"
	)

	ctx := context.Background()

	Context("When reconciling a resource", func() {
		It("should successfully update the status", func() {
			By("creating a new OpticalDevice resource")
			device := &hardwarev1alpha1.OpticalDevice{
				ObjectMeta: metav1.ObjectMeta{
					Name:      DeviceName,
					Namespace: DeviceNamespace,
				},
				Spec: hardwarev1alpha1.OpticalDeviceSpec{
					Location: "site1",
					ControllerConfig: hardwarev1alpha1.ControllerConfigSpec{
						Hostname: DeviceHostname,
						Port:     830,
					},
					Parameters: hardwarev1alpha1.OpticalParametersSpec{
						Bandwidth:  DeviceBandwidth,
						LaserPower: "15dBm",
					},
				},
			}
			Expect(k8sClient.Create(ctx, device)).To(Succeed())

			reconciler := &OpticalDeviceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      DeviceName,
					Namespace: DeviceNamespace,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			By("verifying the updated status")
			updatedDevice := &hardwarev1alpha1.OpticalDevice{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: DeviceName, Namespace: DeviceNamespace}, updatedDevice)
				if err != nil {
					return false
				}
				// Check if the status field we care about is updated
				return updatedDevice.Status.ObservedBandwidth == DeviceBandwidth
			}, time.Second*10, time.Millisecond*250).Should(BeTrue())

			By("verifying the 'Ready' condition")
			Expect(updatedDevice.Status.Conditions).NotTo(BeEmpty())
			readyCondition := updatedDevice.Status.Conditions[0]
			Expect(readyCondition.Type).To(Equal("Ready"))
			Expect(readyCondition.Status).To(Equal(metav1.ConditionTrue))
			Expect(readyCondition.Reason).To(Equal("ReconciliationSuccess"))
		})
	})
})
