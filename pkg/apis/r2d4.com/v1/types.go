package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const FooResourcePlural = "foos"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooSpec   `json:"spec"`
	Status FooStatus `json:"status,omitempty"`
	Test   FooStatus
}

type FooSpec struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type FooStatus struct {
	State   FooState `json:"state,omitempty"`
	Message string   `json:"message,omitempty"`
}

type FooState string

const (
	A FooState = "A"
	B FooState = "B"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FooList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}
