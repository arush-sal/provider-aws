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

package resourcerecordset

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/provider-aws/apis/network/v1alpha3"
	awsclients "github.com/crossplane/provider-aws/pkg/clients"
)

const (
	//RRSetNotFound is the error code that is returned if RRSet is present
	RRSetNotFound = "InvalidRRSetName.NotFound"
)

// Client defines ResourceRecordSet operations
type Client interface {
	ChangeResourceRecordSetsRequest(input *route53.ChangeResourceRecordSetsInput) route53.ChangeResourceRecordSetsRequest
	ListResourceRecordSetsRequest(input *route53.ListResourceRecordSetsInput) route53.ListResourceRecordSetsRequest
}

// NotFound will be raised when there is no ResourceRecordSet
type NotFound struct{}

func (err *NotFound) Error() string {
	return fmt.Sprint(RRSetNotFound)
}

// NewClient creates new AWS client with provided AWS Configuration/Credentials
func NewClient(config *aws.Config) Client {
	return route53.New(*config)
}

// GenerateChangeResourceRecordSetsInput prepares input for a ChangeResourceRecordSetsInput
func GenerateChangeResourceRecordSetsInput(p *v1alpha3.ResourceRecordSetParameters, action route53.ChangeAction) *route53.ChangeResourceRecordSetsInput {
	var ttl *int64
	if p.TTL == nil {
		num := int64(300)
		ttl = &num
	} else {
		ttl = p.TTL
	}

	resourceRecords := make([]route53.ResourceRecord, 0, len(p.Records))
	for _, r := range p.Records {
		record := r
		resourceRecords = append(resourceRecords, route53.ResourceRecord{
			Value: &record,
		})
	}

	return &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: p.ZoneID,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []route53.Change{
				{
					Action: action,
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            p.Name,
						Type:            route53.RRType(aws.StringValue(p.Type)),
						TTL:             ttl,
						ResourceRecords: resourceRecords,
					},
				},
			},
		},
	}
}

// IsUpToDate checks if object is up to date
func IsUpToDate(p v1alpha3.ResourceRecordSetParameters, rrset route53.ResourceRecordSet) (bool, error) {
	// check for the root "." found in DNS entries and add if not found
	if !strings.HasSuffix(*p.Name, ".") {
		p.Name = aws.String(fmt.Sprintf("%s.", *p.Name))
	}
	patch, err := CreatePatch(&rrset, &p)
	if err != nil {
		return false, err
	}
	return cmp.Equal(&v1alpha3.ResourceRecordSetParameters{}, patch), nil
}

// LateInitialize fills the empty fields in *v1alpha3.ResourceRecordSetParameters with
// the values seen in route53.ResourceRecordSet.
func LateInitialize(in *v1alpha3.ResourceRecordSetParameters, rrSet *route53.ResourceRecordSet) {
	if rrSet == nil {
		return
	}

	in.Name = awsclients.LateInitializeStringPtr(in.Name, rrSet.Name)

	rrType := string(rrSet.Type)
	in.Type = awsclients.LateInitializeStringPtr(in.Type, &rrType)

	in.TTL = awsclients.LateInitializeInt64Ptr(in.TTL, rrSet.TTL)
	if len(in.Records) == 0 && len(rrSet.ResourceRecords) != 0 {
		in.Records = make([]string, len(rrSet.ResourceRecords))
		for i, val := range rrSet.ResourceRecords {
			in.Records[i] = aws.StringValue(val.Value)
		}
	}
}

// CreatePatch creates a *v1beta1.ResourceRecordSetParameters that has only the changed
// values between the target *v1beta1.ResourceRecordSetParameters and the current
// *route53.ResourceRecordSet
func CreatePatch(in *route53.ResourceRecordSet, target *v1alpha3.ResourceRecordSetParameters) (*v1alpha3.ResourceRecordSetParameters, error) {
	currentParams := &v1alpha3.ResourceRecordSetParameters{}
	LateInitialize(currentParams, in)

	jsonPatch, err := awsclients.CreateJSONPatch(currentParams, target)

	if err != nil {
		return nil, err
	}
	patch := &v1alpha3.ResourceRecordSetParameters{}
	if err := json.Unmarshal(jsonPatch, patch); err != nil {
		return nil, err
	}
	patch.ZoneID = patch.Name
	return patch, nil
}

// GetResourceRecordSet returns recordSet if present or err
func GetResourceRecordSet(ctx context.Context, c Client, p v1alpha3.ResourceRecordSetParameters) (route53.ResourceRecordSet, error) {
	// check for the root "." found in DNS entries and add if not found
	if !strings.HasSuffix(*p.Name, ".") {
		p.Name = aws.String(fmt.Sprintf("%s.", *p.Name))
	}
	req := c.ListResourceRecordSetsRequest(&route53.ListResourceRecordSetsInput{
		HostedZoneId: p.ZoneID,
	})
	res, err := req.Send(ctx)
	if err != nil {
		return route53.ResourceRecordSet{}, err
	}

	for _, rrSet := range res.ResourceRecordSets {
		if aws.StringValue(rrSet.Name) == aws.StringValue(p.Name) && string(rrSet.Type) == aws.StringValue(p.Type) {
			return rrSet, nil
		}
	}

	return route53.ResourceRecordSet{}, &NotFound{}
}

// StringValueFromRRType returns the string value of aws ResourceRecordType
// passed in or "" if nil
func StringValueFromRRType(rrType route53.RRType) string {
	return string(rrType)
}
