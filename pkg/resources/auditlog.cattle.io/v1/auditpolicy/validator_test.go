package auditpolicy

import (
	"testing"

	auditlogv1 "github.com/rancher/rancher/pkg/apis/auditlog.cattle.io/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestAdmitterValidateFields(t *testing.T) {
	type testCase struct {
		Name     string
		Policy   *auditlogv1.AuditPolicy
		Expected error
	}

	cases := []testCase{
		{
			Name: "filter action allow is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Filters: []auditlogv1.Filter{
						{
							Action: auditlogv1.FilterActionAllow,
						},
					},
				},
			},
		},
		{
			Name: "filter action deny is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Filters: []auditlogv1.Filter{
						{
							Action: auditlogv1.FilterActionDeny,
						},
					},
				},
			},
		},
		{
			Name: "invalid filter action is invalid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Filters: []auditlogv1.Filter{
						{
							Action: "you shall not pass",
						},
					},
				},
			},
			Expected: field.NotSupported(field.NewPath("auditpolicy", "spec", "filters").Index(0), "you shall not pass", []string{string(auditlogv1.FilterActionAllow), string(auditlogv1.FilterActionDeny)}),
		},
		{
			Name: "empty filter action is invalid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Filters: []auditlogv1.Filter{
						{},
					},
				},
			},
			Expected: field.NotSupported(field.NewPath("auditpolicy", "spec", "filters").Index(0), "", []string{string(auditlogv1.FilterActionAllow), string(auditlogv1.FilterActionDeny)}),
		},
		{
			Name: "valid filter request uri regex is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Filters: []auditlogv1.Filter{
						{
							Action:     auditlogv1.FilterActionAllow,
							RequestURI: "/some/endoint/.*",
						},
					},
				},
			},
		},
		{
			Name: "invalid filter request uri regex is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Filters: []auditlogv1.Filter{
						{
							Action:     auditlogv1.FilterActionAllow,
							RequestURI: "*",
						},
					},
				},
			},
			Expected: field.Invalid(field.NewPath("auditpolicy", "spec", "filters").Index(0), "*", "error parsing regexp: missing argument to repetition operator: `*`"),
		},

		{
			Name: "valid header regex is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					AdditionalRedactions: []auditlogv1.Redaction{
						{
							Headers: []string{
								".*",
							},
						},
					},
				},
			},
		},
		{
			Name: "invalid header regex is invalid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					AdditionalRedactions: []auditlogv1.Redaction{
						{
							Headers: []string{
								"*",
							},
						},
					},
				},
			},
			Expected: field.Invalid(field.NewPath("auditpolicy", "spec", "additionalRedactions").Index(0).Child("headers").Index(0), "*", "error parsing regexp: missing argument to repetition operator: `*`"),
		},
		{
			Name: "valid jsonpath is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					AdditionalRedactions: []auditlogv1.Redaction{
						{
							Paths: []string{
								"$..*",
							},
						},
					},
				},
			},
		},
		{
			Name: "invalid jsonpath is invalid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					AdditionalRedactions: []auditlogv1.Redaction{
						{
							Paths: []string{
								"..*",
							},
						},
					},
				},
			},
			Expected: field.Invalid(field.NewPath("auditpolicy", "spec", "additionalRedactions").Index(0).Child("paths").Index(0), "..*", "paths must begin with the root object identifier: '$'"),
		},

		{
			Name: "verbosity level 0 is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Verbosity: auditlogv1.LogVerbosity{
						Level: 0,
					},
				},
			},
		},
		{
			Name: "verbosity level 3 is valid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Verbosity: auditlogv1.LogVerbosity{
						Level: 3,
					},
				},
			},
		},
		{
			Name: "verbosity level -1 is invalid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Verbosity: auditlogv1.LogVerbosity{
						Level: -1,
					},
				},
			},
			Expected: field.Invalid(field.NewPath("auditpolicy", "spec", "verbosity", "level"), -1, ".spec.verbosity.level must be >= 0 or <= 3"),
		},
		{
			Name: "verbosity level 4 is invalid",
			Policy: &auditlogv1.AuditPolicy{
				Spec: auditlogv1.AuditPolicySpec{
					Verbosity: auditlogv1.LogVerbosity{
						Level: 4,
					},
				},
			},
			Expected: field.Invalid(field.NewPath("auditpolicy", "spec", "verbosity", "level"), 4, ".spec.verbosity.level must be >= 0 or <= 3"),
		},
	}

	a := admitter{}
	path := field.NewPath("auditpolicy", "spec")

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			err := a.validateFields(c.Policy, path)

			if c.Expected == nil && err != nil {
				assert.Failf(t, "received unexpected error '%s'", err.Error())
			} else if c.Expected != nil && err == nil {
				assert.Failf(t, "expected to receive err '%s'", c.Expected.Error())
			} else if c.Expected != nil && err != nil {
				assert.EqualError(t, err, c.Expected.Error())
			}
		})
	}
}
