package card

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prometheus-msteams/prometheus-msteams/pkg/testutils"
)

const (
	testPromAlertFile   = "./testdata/prom_post_request.json"
	testSchemaContext   = "http://schema.org/extensions"
	testAlertTitle      = "Prometheus Alert (Firing)"
	testAlertSummary    = "Prometheus Test"
	testThemeColor      = "FFA500"
	testActivityTitle   = "[10.80.40.11 reported high memory usage with 23.28%.](http://docker.for.mac.host.internal:9093)"
	testMemorySummary   = "Server High Memory usage"
	testAlertname       = "alertname"
	testInstance        = "instance"
	testJob             = "job"
	testMonitor         = "master"
	testSeverity        = "severity"
	testLabelSummary    = "summary"
	testLabelMonitor    = "monitor"
	testSeverityWarning = "warning"
)

func Test_templatedCard_Convert(t *testing.T) {
	tests := []struct {
		name              string
		promAlertFile     string
		templateFile      string
		escapeUnderscores bool
		want              Office365ConnectorCard
		wantErr           bool
	}{
		{
			name:              "do not escape underscores",
			promAlertFile:     testPromAlertFile,
			templateFile:      "../../default-message-card.tmpl",
			escapeUnderscores: false,
			want: Office365ConnectorCard{
				Context:    testSchemaContext,
				Type:       messageCardType,
				Title:      testAlertTitle,
				Summary:    testAlertSummary,
				ThemeColor: testThemeColor,
				Sections: []Section{
					{
						ActivityTitle: testActivityTitle,
						Markdown:      true,
						Facts: []FactSection{
							{},
							{Name: testLabelSummary, Value: testMemorySummary},
							{Name: testAlertname, Value: `high_memory_load`},
							{Name: testInstance, Value: `instance-with-hyphen_and_underscore`},
							{Name: testJob, Value: `docker_nodes`},
							{Name: testLabelMonitor, Value: testMonitor},
							{Name: testSeverity, Value: testSeverityWarning},
						},
					},
				},
			},
		},
		{
			name:              "escape underscores",
			promAlertFile:     testPromAlertFile,
			templateFile:      "../../default-message-card.tmpl",
			escapeUnderscores: true,
			want: Office365ConnectorCard{
				Context:    testSchemaContext,
				Type:       messageCardType,
				Title:      testAlertTitle,
				Summary:    testAlertSummary,
				ThemeColor: testThemeColor,
				Sections: []Section{
					{
						ActivityTitle: testActivityTitle,
						Markdown:      true,
						Facts: []FactSection{
							{},
							{Name: testLabelSummary, Value: testMemorySummary},
							{Name: testAlertname, Value: `high\_memory\_load`},
							{Name: testInstance, Value: `instance-with-hyphen\_and\_underscore`},
							{Name: testJob, Value: `docker\_nodes`},
							{Name: testLabelMonitor, Value: testMonitor},
							{Name: testSeverity, Value: testSeverityWarning},
						},
					},
				},
			},
		},
		{
			name:              "action card",
			promAlertFile:     testPromAlertFile,
			templateFile:      "./testdata/action-message-card.tmpl",
			escapeUnderscores: true,
			want: Office365ConnectorCard{
				Context:    testSchemaContext,
				Type:       messageCardType,
				Title:      testAlertTitle,
				Summary:    testAlertSummary,
				ThemeColor: testThemeColor,
				Sections: []Section{
					{
						ActivityTitle: testActivityTitle,
						Markdown:      true,
						Facts: []FactSection{
							{},
							{Name: testLabelSummary, Value: testMemorySummary},
							{Name: testAlertname, Value: `high\_memory\_load`},
							{Name: testInstance, Value: `instance-with-hyphen\_and\_underscore`},
							{Name: testJob, Value: `docker\_nodes`},
							{Name: testLabelMonitor, Value: testMonitor},
							{Name: testSeverity, Value: testSeverityWarning},
						},
					},
				},
				PotentialAction: []Action{
					{
						"@context": string("http://schema.org"),
						"@type":    string("ViewAction"),
						"name":     string("Runbook"),
						"target":   []interface{}{string("https://github.com/bzon/prometheus-msteams")},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := ParseTemplateFile(tt.templateFile)
			if err != nil {
				t.Fatal(err)
			}

			a, err := testutils.ParseWebhookJSONFromFile(tt.promAlertFile)
			if err != nil {
				t.Fatal(err)
			}

			m := NewTemplatedCardCreator(tmpl, tt.escapeUnderscores)

			got, err := m.Convert(context.Background(), a)
			if (err != nil) != tt.wantErr {
				t.Errorf("templatedCard.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
