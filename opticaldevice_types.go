/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUTHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpticalDeviceSpec 定義了光通訊設備的期望狀態
type OpticalDeviceSpec struct {
	// Hostname 是設備的管理主機名或 IP 位址
	Hostname string `json:"hostname"`
	// Port 是設備的管理埠號
	Port int `json:"port"`
	// Bandwidth 是期望的頻寬設定, e.g., "10Gbps"
	Bandwidth string `json:"bandwidth"`
	// LaserPower 是期望的雷射功率, e.g., "14dBm"
	LaserPower string `json:"laserPower"`
}

// OpticalDeviceStatus 定義了光通訊設備的觀察到的狀態
type OpticalDeviceStatus struct {
	// Conditions 儲存了資源的當前狀態條件列表
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ObservedBandwidth 是從硬體回讀的實際頻寬
	// +optional
	ObservedBandwidth string `json:"observedBandwidth,omitempty"`

	// LastUpdateTime 是最後一次成功同步的時間
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="HOSTNAME",type="string",JSONPath=".spec.hostname"
//+kubebuilder:printcolumn:name="BANDWIDTH",type="string",JSONPath=".spec.bandwidth"
//+kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"

// OpticalDevice 是光通訊設備的 Schema
type OpticalDevice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpticalDeviceSpec   `json:"spec,omitempty"`
	Status OpticalDeviceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpticalDeviceList 包含一個 OpticalDevice 的列表
type OpticalDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpticalDevice `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpticalDevice{}, &OpticalDeviceList{})
}
