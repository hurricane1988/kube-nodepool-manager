/*
Copyright 2025 codeFuthure.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodePoolSpec defines the desired state of NodePool.
type NodePoolSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NodeSelector is used to dynamically select nodes for the pool
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Architecture specifies the node architecture (e.g., amd64, arm64)
	Architecture string `json:"architecture,omitempty"`

	// OperatingSystem specifies the operating system of the nodes (e.g., linux, windows)
	OperatingSystem string `json:"operatingSystem,omitempty"`

	// OSImage specifies the OS image (e.g., Ubuntu 20.04)
	OSImage string `json:"osImage,omitempty"`

	// KernelVersion specifies the version of the node's kernel (e.g., 5.4.0-42-generic)
	KernelVersion string `json:"kernelVersion,omitempty"`

	// KubeletVersion specifies the version of kubelet installed on the node
	KubeletVersion string `json:"kubeletVersion,omitempty"`

	// CpuVendor specifies the CPU vendor (e.g., Intel, AMD)
	CPUVendor string `json:"cpuVendor,omitempty"`

	// Nodes is an optional list of node names in the pool
	Nodes []string `json:"nodes,omitempty"`
}

// NodePoolCapacity defines the overall resource capacity of the NodePool
type NodePoolCapacity struct {
	CPU              string `json:"cpu"`
	EphemeralStorage string `json:"ephemeral-storage"`
	Hugepages1Gi     string `json:"hugepages-1Gi"`
	Hugepages2Mi     string `json:"hugepages-2Mi"`
	Memory           string `json:"memory"`
	Pods             string `json:"pods"`
}

// NodeDetail contains the information about nodes with specific properties
type NodeDetail struct {
	Architecture    string `json:"architecture,omitempty"`
	OperatingSystem string `json:"operatingSystem,omitempty"`
	OSImage         string `json:"osImage,omitempty"`
	KernelVersion   string `json:"kernelVersion,omitempty"`
	KubeletVersion  string `json:"kubeletVersion,omitempty"`
	CPUVendor       string `json:"cpuVendor,omitempty"`
	Count           int    `json:"count"`
}

// NodePoolStatus defines the observed state of NodePool.
type NodePoolStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// AvailableNodes contains the list of nodes in the pool
	AvailableNodes []string `json:"availableNodes,omitempty"`

	// NodeDetails stores nodes categorized by architecture, OS, etc.
	NodeDetails map[string]NodeDetail `json:"nodeDetails,omitempty"`

	// Capacity defines the resource capacity of the NodePool
	Capacity NodePoolCapacity `json:"capacity,omitempty"`

	// Nodes stores the current state of each node (Schedulable, Unschedulable)
	Nodes map[string]string `json:"nodes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NodePool is the Schema for the nodepools API.
type NodePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodePoolSpec   `json:"spec,omitempty"`
	Status NodePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodePoolList contains a list of NodePool.
type NodePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodePool{}, &NodePoolList{})
}
