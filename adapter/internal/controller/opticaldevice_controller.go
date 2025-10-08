package controller

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	hardwarev1alpha1 "ran.example.com/o-ran-adapter/internal/api/v1alpha1"
)

var (
	opticaldeviceInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "opticaldevice_info",
			Help: "Information about the OpticalDevice.",
		},
		[]string{"hostname", "bandwidth"},
	)
	opticaldeviceReconcileTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "opticaldevice_reconcile_total",
			Help: "Total number of successful reconciliations.",
		},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(opticaldeviceInfo, opticaldeviceReconcileTotal)
}

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

	var opticalDevice hardwarev1alpha1.OpticalDevice
	if err := r.Get(ctx, req.NamespacedName, &opticalDevice); err != nil {
		logger.Info("OpticalDevice resource not found. Ignoring since object must be deleted.")
		// When a resource is deleted, we should also clean up its associated metric.
		// Note: This is a best-effort cleanup. If the controller is down, the metric might persist.
		// A more robust solution might involve a finalizer.
		opticaldeviceInfo.DeletePartialMatch(prometheus.Labels{"hostname": req.Name})
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling OpticalDevice",
		"hostname", opticalDevice.Spec.ControllerConfig.Hostname,
		"port", opticalDevice.Spec.ControllerConfig.Port,
		"bandwidth", opticalDevice.Spec.Parameters.Bandwidth,
		"laserPower", opticalDevice.Spec.Parameters.LaserPower,
	)

	// (Simulation Logic)
	time.Sleep(50 * time.Millisecond)

	// Update status
	meta.SetStatusCondition(&opticalDevice.Status.Conditions, metav1.Condition{
		Type:    "Ready",
		Status:  metav1.ConditionTrue,
		Reason:  "ReconciliationSuccess",
		Message: "Device configured successfully (simulated)",
	})
	opticalDevice.Status.ObservedBandwidth = opticalDevice.Spec.Parameters.Bandwidth
	opticalDevice.Status.ObservedLaserPower = opticalDevice.Spec.Parameters.LaserPower
	opticalDevice.Status.LastSyncTime = metav1.Now()

	if err := r.Status().Update(ctx, &opticalDevice); err != nil {
		logger.Error(err, "Failed to update OpticalDevice status")
		return ctrl.Result{}, err
	}

	// Update Prometheus metrics
	opticaldeviceInfo.WithLabelValues(
		opticalDevice.Spec.ControllerConfig.Hostname,
		opticalDevice.Spec.Parameters.Bandwidth,
	).Set(1)
	opticaldeviceReconcileTotal.Inc()

	logger.Info("Successfully reconciled OpticalDevice")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpticalDeviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hardwarev1alpha1.OpticalDevice{}).
		Complete(r)
}
