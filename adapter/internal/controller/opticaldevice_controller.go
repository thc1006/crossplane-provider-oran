package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	hardwarev1alpha1 "ran.example.com/o-ran-adapter/api/v1alpha1"
)

// OpticalDeviceReconciler reconciles a OpticalDevice object
type OpticalDeviceReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=hardware.ran.example.com,resources=opticaldevices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hardware.ran.example.com,resources=opticaldevices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hardware.ran.example.com,resources=opticaldevices/finalizers,verbs=update

func (r *OpticalDeviceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciliation loop started")

	// 1. Fetch the OpticalDevice instance
	var opticalDevice hardwarev1alpha1.OpticalDevice
	if err := r.Get(ctx, req.NamespacedName, &opticalDevice); err != nil {
		logger.Info("OpticalDevice resource not found. Ignoring since object must be deleted.")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// --- Placeholder for a real hardware client ---
	// In a real implementation, you would initialize a client here to communicate
	// with the local controller's API (e.g., an HTTP client for a REST API).
	// hardwareClient := NewHardwareClient(opticalDevice.Spec.ControllerConfig)
	// ---------------------------------------------

	logger.Info("Processing OpticalDevice",
		"hostname", opticalDevice.Spec.ControllerConfig.Hostname,
		"bandwidth", opticalDevice.Spec.Parameters.Bandwidth,
		"laserPower", opticalDevice.Spec.Parameters.LaserPower,
	)

	// 2. (SIMULATION) Pretend to configure the hardware
	// This is where you would call the actual hardware API.
	// err := hardwareClient.Configure(opticalDevice.Spec.Parameters)
	// For now, we just log and simulate success.
	logger.Info("SIMULATING: Sending configuration to hardware via NI-VISA API wrapper",
		"target", fmt.Sprintf("%s:%d", opticalDevice.Spec.ControllerConfig.Hostname, opticalDevice.Spec.ControllerConfig.Port),
	)
	time.Sleep(1 * time.Second) // Simulate network latency

	// 3. Update the status of the OpticalDevice resource
	// This is crucial for closing the control loop.
	opticalDevice.Status.ObservedBandwidth = opticalDevice.Spec.Parameters.Bandwidth
	opticalDevice.Status.ObservedLaserPower = opticalDevice.Spec.Parameters.LaserPower
	opticalDevice.Status.LastSyncTime = metav1.Now()

	// Here you would set conditions based on the result of the hardware interaction
	// setStatusCondition(&opticalDevice, "Ready", metav1.ConditionTrue, "Configured", "Successfully configured hardware")
	// setStatusCondition(&opticalDevice, "Synced", metav1.ConditionTrue, "Reconciled", "Reconciliation successful")

	if err := r.Status().Update(ctx, &opticalDevice); err != nil {
		logger.Error(err, "Failed to update OpticalDevice status")
		r.Recorder.Event(&opticalDevice, "Warning", "StatusUpdateFailed", "Could not update resource status")
		return ctrl.Result{}, err
	}

	r.Recorder.Event(&opticalDevice, "Normal", "Configured", "Hardware configuration applied successfully (simulated)")
	logger.Info("Reconciliation loop finished successfully")

	// Requeue after a certain period to periodically check hardware status
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpticalDeviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hardwarev1alpha1.OpticalDevice{}).
		Complete(r)
}
