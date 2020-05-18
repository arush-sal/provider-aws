package zone

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/route53iface"
	"github.com/google/uuid"

	"github.com/crossplane/provider-aws/apis/network/v1alpha3"
)

// Client defines Route53 Client operations
type Client interface {
	CreateZoneRequest(cr *v1alpha3.Zone) route53.CreateHostedZoneRequest
	DeleteZoneRequest(id *string) route53.DeleteHostedZoneRequest
	GetZoneRequest(id string) route53.GetHostedZoneRequest
	UpdateZoneRequest(id, comment *string) route53.UpdateHostedZoneCommentRequest
}

type zoneClient struct {
	zone route53iface.ClientAPI
}

// NewClient creates new AWS Client with provided AWS Configurations/Credentials
func NewClient(config *aws.Config) Client {
	return &zoneClient{zone: route53.New(*config)}
}

// GetZoneRequest returns a route53 GetHostedZoneRequest using which a route53
// Hosted Zone can be fetched and checked for existence.
func (c *zoneClient) GetZoneRequest(id string) route53.GetHostedZoneRequest {
	return c.zone.GetHostedZoneRequest(&route53.GetHostedZoneInput{
		Id: &id,
	})
}

// CreateZoneRequest returns a route53 CreateHostedZoneRequest using which a route53
// Hosted Zone can be created.
func (c *zoneClient) CreateZoneRequest(cr *v1alpha3.Zone) route53.CreateHostedZoneRequest {
	if cr.Spec.ForProvider.CallerReference == nil {
		cr.Spec.ForProvider.CallerReference = getCallerRef()
	}

	reqInput := &route53.CreateHostedZoneInput{
		CallerReference: cr.Spec.ForProvider.CallerReference,
		Name:            cr.Spec.ForProvider.Name,
		HostedZoneConfig: &route53.HostedZoneConfig{
			PrivateZone: cr.Spec.ForProvider.PrivateZone,
			Comment:     cr.Spec.ForProvider.Comment,
		},
	}

	if *cr.Spec.ForProvider.PrivateZone {
		reqInput.HostedZoneConfig.PrivateZone = cr.Spec.ForProvider.PrivateZone

		if cr.Spec.ForProvider.VPCID != nil {
			reqInput.VPC = &route53.VPC{
				VPCId: cr.Spec.ForProvider.VPCID,
			}
		}

		if cr.Spec.ForProvider.VPCRegion != nil {
			reqInput.VPC.VPCRegion = route53.VPCRegion(*cr.Spec.ForProvider.VPCRegion)
		}
	}

	return c.zone.CreateHostedZoneRequest(reqInput)
}

// UpdateZoneRequest returns a route53 UpdateHostedZoneRequest using which a route53
// Hosted Zone can be updated.
func (c *zoneClient) UpdateZoneRequest(id, comment *string) route53.UpdateHostedZoneCommentRequest {
	return c.zone.UpdateHostedZoneCommentRequest(&route53.UpdateHostedZoneCommentInput{Comment: comment, Id: id})
}

// DeleteZoneRequest returns a route53 DeleteHostedZoneRequest using which a route53
// Hosted Zone can be deleted.
func (c *zoneClient) DeleteZoneRequest(id *string) route53.DeleteHostedZoneRequest {
	return c.zone.DeleteHostedZoneRequest(&route53.DeleteHostedZoneInput{
		Id: id,
	})
}

// IsErrorNoSuchHostedZone returns true if the error code indicates that the requested Zone was not found
func IsErrorNoSuchHostedZone(err error) bool {
	if zoneErr, ok := err.(awserr.Error); ok && zoneErr.Code() == route53.ErrCodeNoSuchHostedZone {
		return true
	}
	return false
}

// Update the status of the runtime object
func Update(status *v1alpha3.ZoneObservation, op *route53.HostedZone) {
	status.ID = aws.StringValue(op.Id)
	status.ResourceRecordCount = aws.Int64Value(op.ResourceRecordSetCount)
}

func getCallerRef() *string {
	return aws.String(uuid.New().String())
}
