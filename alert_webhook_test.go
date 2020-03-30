package alertfunc

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAlertFunc(t *testing.T) {
	rawTestJSON1 := `{
		"incident": {
		  "incident_id": "f2e08c333dc64cb09f75eaab355393bz",
		  "resource_id": "i-4a266a2d",
		  "resource_name": "webserver-85",
		  "state": "open",
		  "started_at": 1385085727,
		  "ended_at": null,
		  "policy_name": "Webserver Health",
		  "condition_name": "CPU usage",
		  "url": "https://console.cloud.google.com/monitoring/alerting/incidents?project=PROJECT_ID",
		  "summary": "CPU for webserver-85 is above the threshold of 1% with a value of 28.5%"
		},
		"version": "1.1"
	}`

	tests := []struct {
		authToken string
		body      string
		want      string
	}{
		{authToken: "valid", body: rawTestJSON1, want: "OK"},
		{authToken: "invalid", body: "{}", want: "Unauthorized\n"},
	}

	os.Setenv("AUTH_TOKEN", "valid")
	for _, test := range tests {
		url := "/?auth_token=" + test.authToken
		req := httptest.NewRequest("GET", url, strings.NewReader(test.body))
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		AlertFunc(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("AlertFunc(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}
