package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PGDatabaseSpec defines the desired state of PGDatabase
type PGDatabaseSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=3
	Instances int32 `json:"instances,omitempty"`

	// Storage size, e.g., "1Gi"
	StorageSize string `json:"storageSize"`

	// Postgres version (Major version, e.g., "15")
	// +kubebuilder:default="15"
	Version string `json:"version,omitempty"`
}

// PGDatabaseStatus defines the observed state of PGDatabase
type PGDatabaseStatus struct {
	// Status of the underlying CNPG Cluster
	Phase string `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PGDatabase is the Schema for the pgdatabases API
type PGDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PGDatabaseSpec   `json:"spec,omitempty"`
	Status PGDatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PGDatabaseList contains a list of PGDatabase
type PGDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PGDatabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PGDatabase{}, &PGDatabaseList{})
}
