/*
Copyright 2019 The Crossplane Authors.

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
)

// +kubebuilder:object:root=true

// Zone is a managed resource that represents an AWS Route53 Hosted Zone.
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.atProvider.Id"
// +kubebuilder:printcolumn:name="RRs",type="integer",JSONPath=".status.atProvider.ResourceRecordSetCount"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type Zone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZoneSpec   `json:"spec"`
	Status ZoneStatus `json:"status,omitempty"`
}

// ZoneSpec defines the desired state of an AWS Route53 Hosted Zone.
type ZoneSpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  ZoneParameters `json:"forProvider"`
}

// ZoneStatus represents the observed state of a Zone.
type ZoneStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	AtProvider                     ZoneObservation `json:"atProvider"`
}

// ZoneParameters define the desired state of an AWS Route53 Hosted Zone.
type ZoneParameters struct {
	// The name of the domain. Specify a fully qualified domain name, for example,
	// www.example.com. The trailing dot is optional; Amazon Route 53 assumes that
	// the domain name is fully qualified. This means that Route 53 treats www.example.com
	// (without a trailing dot) and www.example.com. (with a trailing dot) as identical.
	//
	// If you're creating a public hosted zone, this is the name you have registered
	// with your DNS registrar. If your domain name is registered with a registrar
	// other than Route 53, change the name servers for your domain to the set of
	// NameServers that CreateHostedZone returns in DelegationSet.
	// +immutable
	// +kubebuilder:validation:Required
	Name *string `json:"name"`

	// CallerReference is an unique string that identifies the request and that
	// allows failed Zone create requests to be retried without the risk of
	// executing the operation twice.
	// +immutable
	// +optional
	CallerReference *string `json:"callerref,omitempty"`

	// Any comments that you want to include about the hosted zone.
	// +optional
	Comment *string `json:"comment,omitempty"`

	// A value that indicates whether this is a private hosted zone.
	// +immutable
	// +optional
	PrivateZone *bool `json:"privateZone,omitempty"`

	// (Private hosted zones only) The ID of an Amazon VPC that you're
	// associating with this hosted zone. You can specify only one Amazon VPC
	// when you create a private hosted zone.
	// +immutable
	// +optional
	VPCID *string `json:"VPCId,omitempty"`

	// (Private hosted zones only) The region that an Amazon VPC was created in.
	// +immutable
	// +optional
	VPCRegion *string `json:"VPCRegion,omitempty"`

	// (Private hosted zones only) VPCIdRef references a VPC to retrieves its VPC Id
	// +optional
	VPCIDRef *runtimev1alpha1.Reference `json:"VPCIdRef,omitempty"`

	// VPCIdSelector selects a reference to a VPC to retrieves its VPC Id
	// +optional
	VPCIDSelector *runtimev1alpha1.Selector `json:"VPCIdSelector,omitempty"`
}

// ZoneObservation keeps the state for the external resource.
type ZoneObservation struct {
	// ID is the unique URL representing the new hosted zone.
	ID string `json:"Id"`

	// The number of resource record sets in the hosted zone.
	ResourceRecordCount int64 `json:"ResourceRecordSetCount"`
}

// DelegationSet describes the name servers for this hosted zone.
type DelegationSet struct {
	// NameServers contains a list of the authoritative name servers for a hosted zone.
	NameServers []string `json:"NameServers"`
}

// +kubebuilder:object:root=true

// ZoneList contains a list of Zone.
type ZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Zone `json:"items"`
}
