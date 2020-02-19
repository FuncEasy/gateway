package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FunctionSpec defines the desired state of Function
type FunctionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Function         string            `json:"function" binding:"required"`
	Identifier       string            `json:"identifier" binding:"required"`
	Version          string            `json:"version" binding:"required"`
	Runtime          string            `json:"runtime" binding:"required"`
	Deps             string            `json:"deps,omitempty"`
	Handler          string            `json:"handler" binding:"required"`
	ContentType      string            `json:"contentType" binding:"required"`
	Timeout          string            `json:"timeout" binding:"required"`
	Size             *int32            `json:"size" binding:"required"`
	ExposedPort      int32             `json:"exposedPorts,omitempty"`
	ExternalService  map[string]string `json:"externalService,omitempty"`
	DataSource       string            `json:"dataSource,omitempty"`
	DataServiceToken string            `json:"dataServiceToken,omitempty"`
}

// FunctionStatus defines the observed state of Function
type FunctionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodsStatus []PodsStatus `json:"podStatus"`
}

type PodsStatus struct {
	PodName               string                   `json:"podName"`
	PodPhase              corev1.PodPhase          `json:"podPhase"`
	InitContainerStatuses []corev1.ContainerStatus `json:"initContainerStatuses"`
	ContainerStatuses     []corev1.ContainerStatus `json:"containerStatuses"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionSpec   `json:"spec,omitempty"`
	Status FunctionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FunctionList contains a list of Function
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Function `json:"items"`
}
