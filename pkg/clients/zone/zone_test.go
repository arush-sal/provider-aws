package zone

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/provider-aws/apis/network/v1alpha3"
)

func TestIsErrorNoSuchHostedZone(t *testing.T) {
	tests := map[string]struct {
		err  error
		want bool
	}{
		"validError": {
			err:  awserr.New(route53.ErrCodeNoSuchHostedZone, "The specified hosted zone does not exist.", errors.New(route53.ErrCodeNoSuchHostedZone)),
			want: true,
		},
		"invalidAwsError": {
			err:  awserr.New(route53.ErrCodeHostedZoneNotFound, "The specified HostedZone can't be found.", errors.New(route53.ErrCodeHostedZoneNotFound)),
			want: false,
		},
		"randomError": {
			err:  errors.New("the specified hosted zone does not exist"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			if got := IsErrorNoSuchHostedZone(tt.err); got != tt.want {
				t.Errorf("IsErrorNoSuchHostedZone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		status *v1alpha3.ZoneObservation
		op     *route53.HostedZone
	}
	tests := map[string]struct {
		args
		want *v1alpha3.ZoneObservation
	}{
		"ValidInput": {
			args: args{
				status: &v1alpha3.ZoneObservation{},
				op: &route53.HostedZone{
					Id:                     aws.String("/hostedzone/XXXXXXXXXXXXXXXXXXX"),
					ResourceRecordSetCount: aws.Int64(2),
				},
			},
			want: &v1alpha3.ZoneObservation{
				ID:                  "/hostedzone/XXXXXXXXXXXXXXXXXXX",
				ResourceRecordCount: 2,
			},
		},
		"EmpytInput": {
			args: args{
				status: &v1alpha3.ZoneObservation{
					ID:                  "/hostedzone/XXXXXXXXXXXXXXXXXXX",
					ResourceRecordCount: 2,
				},
				op: &route53.HostedZone{},
			},
			want: &v1alpha3.ZoneObservation{},
		},
		"UpdateInput": {
			args: args{
				status: &v1alpha3.ZoneObservation{
					ID:                  "/hostedzone/XXXXXXXXXXXXXXXXXXX",
					ResourceRecordCount: 2,
				},
				op: &route53.HostedZone{
					Id:                     aws.String("/hostedzone/XXXXXXXXXXXXXXXXXXX"),
					ResourceRecordSetCount: aws.Int64(20),
				},
			},
			want: &v1alpha3.ZoneObservation{
				ID:                  "/hostedzone/XXXXXXXXXXXXXXXXXXX",
				ResourceRecordCount: 20,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			Update(tt.args.status, tt.args.op)
			if diff := cmp.Diff(tt.want, tt.args.status, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}
