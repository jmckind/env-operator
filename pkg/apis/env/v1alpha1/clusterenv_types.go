package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

// ClusterEnvSpec defines the desired state of ClusterEnv
// +k8s:openapi-gen=true
type ClusterEnvSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// ClusterEnvStatus defines the observed state of ClusterEnv
// +k8s:openapi-gen=true
type ClusterEnvStatus struct {
	Product string `json:"product"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterEnv is the Schema for the clusterenvs API
// +k8s:openapi-gen=true
type ClusterEnv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterEnvSpec   `json:"spec,omitempty"`
	Status ClusterEnvStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterEnvList contains a list of ClusterEnv
type ClusterEnvList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterEnv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterEnv{}, &ClusterEnvList{})
}
