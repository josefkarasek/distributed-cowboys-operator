/*
Copyright 2022.

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

// ShootoutSpec defines the desired state of Shootout
type ShootoutSpec struct {
	// Cowboys is a raw JSON representation of cowboys taking part in this shootout
	Cowboys string `json:"cowboys,omitempty"`
}

// ShootoutStatus defines the observed state of Shootout
type ShootoutStatus struct {
	// Winner is the winner of this Shootout
	Winner string `json:"winner,omitempty"`
	// Error represents non-recoverable error (ie. invalid user input)
	Error string `json:"error,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Shootout is the Schema for the shootouts API
type Shootout struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShootoutSpec   `json:"spec,omitempty"`
	Status ShootoutStatus `json:"status,omitempty"`
}

// func (s Shootout) GenerateName(cowboyName string) string {
// 	return fmt.Sprintf("%s-%s", s.Name, cowboyName)
// }

//+kubebuilder:object:root=true

// ShootoutList contains a list of Shootout
type ShootoutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Shootout `json:"items"`
}

type Cowboy struct {
	Name   string `json:"name,omitempty"`
	Health int64  `json:"health,omitempty"`
	Damage int64  `json:"damage,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Shootout{}, &ShootoutList{})
}
