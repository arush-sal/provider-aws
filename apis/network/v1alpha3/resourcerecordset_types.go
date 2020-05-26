package v1alpha3

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSRecordType defines the valid DNS Record Types that can be used.
type DNSRecordType string

const (
	// DNSRecordTypeSOA represents DNS SOA record type.
	DNSRecordTypeSOA DNSRecordType = "SOA"

	// DNSRecordTypeA represents DNS A record type.
	DNSRecordTypeA DNSRecordType = "A"

	// DNSRecordTypeTXT represents DNS TXT record type.
	DNSRecordTypeTXT DNSRecordType = "TXT"

	// DNSRecordTypeNS represents DNS NS record type.
	DNSRecordTypeNS DNSRecordType = "NS"

	// DNSRecordTypeCNAME represents DNS CNAME record type.
	DNSRecordTypeCNAME DNSRecordType = "CNAME"

	// DNSRecordTypeMX represents DNS MX record type.
	DNSRecordTypeMX DNSRecordType = "MX"

	// DNSRecordTypeNAPTR represents DNS NAPTR record type.
	DNSRecordTypeNAPTR DNSRecordType = "NAPTR"

	// DNSRecordTypePTR represents DNS PTR record type.
	DNSRecordTypePTR DNSRecordType = "PTR"

	// DNSRecordTypeSRV represents DNS SRV record type.
	DNSRecordTypeSRV DNSRecordType = "SRV"

	// DNSRecordTypeSPF represents DNS SPF record type.
	DNSRecordTypeSPF DNSRecordType = "SPF"

	// DNSRecordTypeAAAA represents DNS AAAA record type.
	DNSRecordTypeAAAA DNSRecordType = "AAAA"

	// DNSRecordTypeCAA represents DNS CAA record type.
	DNSRecordTypeCAA DNSRecordType = "CAA"
)

// ChangeAction defines the valid actions that can be performed on a ResourceRecordSet.
type ChangeAction string

const (
	// ChangeActionCreate represents a Resource Record CREATE operation.
	ChangeActionCreate ChangeAction = "CREATE"

	// ChangeActionDelete represents a Resource Record DELETE operation.
	ChangeActionDelete ChangeAction = "DELETE"

	// ChangeActionUpsert represents a Resource Record UPSERT operation.
	ChangeActionUpsert ChangeAction = "UPSERT"
)

// +kubebuilder:object:root=true

// ResourceRecordSet is a managed resource that represents an AWS Route53 Resource Record.
// +kubebuilder:printcolumn:name="TYPE",type="string",JSONPath=".spec.forProvider.type"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type ResourceRecordSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceRecordSetSpec   `json:"spec"`
	Status ResourceRecordSetStatus `json:"status,omitempty"`
}

// ResourceRecordSetSpec defines the desired state of an AWS Route53 Resource Record.
type ResourceRecordSetSpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  ResourceRecordSetParameters `json:"forProvider"`
}

// ResourceRecordSetStatus represents the observed state of a ResourceRecordSet.
type ResourceRecordSetStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
}

// ResourceRecord holds the DNS value to be used for the record.
type ResourceRecord struct {
	// The current or new DNS record value, not to exceed 4,000 characters. In the
	// case of a DELETE action, if the current value does not match the actual value,
	// an error is returned.
	Value *string `json:"Value"`
}

// ResourceRecordSetParameters define the desired state of an AWS Route53 Resource Record.
type ResourceRecordSetParameters struct {
	// Name of the record that you want to create, update, or delete.
	// +kubebuilder:validation:Required
	Name *string `json:"name"`

	// Alias resource record sets only: Information about the AWS resource, such
	// as a CloudFront distribution a Elastic Load Balancer or an Amazon S3 bucket
	// that you want to route traffic to.
	// +optional
	Alias *bool `json:"alias,omitempty"`

	// +optional
	DNSName *string `json:"dnsName,omitempty"`

	// +optional
	EvaluateTargetHealth *bool `json:"evaluateHealthTarget,omitempty"`

	// Type represents the DNS record type
	// +kubebuilder:validation:Required
	Type *string `json:"type"`

	// The resource record cache time to live (TTL), in seconds.
	// +optional
	TTL *int64 `json:"ttl,omitempty"`

	// ResourceRecord holds the information about the resource records to act upon.
	// +optional
	Records []string `json:"records,omitempty"`

	// ZoneID of the HostedZone in which you want to CREATE, CHANGE, or DELETE the Resource Record.
	// +kubebuilder:validation:Required
	ZoneID *string `json:"zoneId,omitempty"`

	// ZoneIDRef references a Zone to retrieves its ZoneId
	// +optional
	ZoneIDRef *runtimev1alpha1.Reference `json:"zoneIdRef,omitempty"`

	// ZoneIDSelector selects a reference to a Zone to retrieves its ZoneID
	// +optional
	ZoneIDSelector *runtimev1alpha1.Selector `json:"zoneIdSelector,omitempty"`
}

// +kubebuilder:object:root=true

// ResourceRecordSetList contains a list of ResourceRecordSet.
type ResourceRecordSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ResourceRecordSet `json:"items"`
}
