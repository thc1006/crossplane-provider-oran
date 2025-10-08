package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ControllerConfigSpec defines the connection details for the local controller.
type ControllerConfigSpec struct {
	// Hostname or IP address of the local controller PC.
	// +kubebuilder:validation:Required
	Hostname string `json:"hostname"`

	// Port of the API service (e.g., NI-VISA wrapper) on the local controller.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int `json:"port"`
}

// OpticalParametersSpec defines the desired configuration for the optical device.
type OpticalParametersSpec struct {
	// The desired bandwidth, e.g., "100Gbps".
	// +kubebuilder:validation:Required
	Bandwidth string `json:"bandwidth"`

	// The desired laser power, e.g., "14.5dBm".
	// +kubebuilder:validation:Required
	LaserPower string `json:"laserPower"`

	// Optional optical channel number.
	// +kubebuilder:validation:Optional
	Channel int `json:"channel,omitempty"`
}

// OpticalDeviceSpec defines the desired state of OpticalDevice
type OpticalDeviceSpec struct {
	// Location of the device.
	// +kubebuilder:validation:Required
	Location string `json:"location"`

	// Connection configuration for the local controller.
	// +kubebuilder:validation:Required
	ControllerConfig ControllerConfigSpec `json:"controllerConfig"`

	// Desired hardware parameters.
	// +kubebuilder:validation:Required
	Parameters OpticalParametersSpec `json:"parameters"`
}

// OpticalDeviceStatus defines the observed state of OpticalDevice
type OpticalDeviceStatus struct {
	// Represents the latest available observations of a OpticalDevice's state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditionsmetav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ObservedBandwidth reflects the last known bandwidth from the hardware.
	ObservedBandwidth string `json:"observedBandwidth,omitempty"`

	// ObservedLaserPower reflects the last known laser power from the hardware.
	ObservedLaserPower string `json:"observedLaserPower,omitempty"`

	// LastSyncTime is the timestamp of the last successful reconciliation.
	LastSyncTime metav1.Time `json:"lastSyncTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions.status"
//+kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions.status"
//+kubebuilder:printcolumn:name="BANDWIDTH",type="string",JSONPath=".spec.parameters.bandwidth"
//+kubebuilder:printcolumn:name="POWER",type="string",JSONPath=".spec.parameters.laserPower"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// OpticalDevice is the Schema for the opticaldevices API
type OpticalDevice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpticalDeviceSpec   `json:"spec,omitempty"`
	Status OpticalDeviceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpticalDeviceList contains a list of OpticalDevice
type OpticalDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           OpticalDevice `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpticalDevice{}, &OpticalDeviceList{})
}
