package alertfunc

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAlertFuncSystem(t *testing.T) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

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
		{authToken: os.Getenv("AUTH_TOKEN"), body: rawTestJSON1, want: "OK"},
		{authToken: "invalid", body: "{}", want: "Unauthorized\n"},
	}

	for _, test := range tests {
		urlString := os.Getenv("BASE_URL") + "?auth_token=" + test.authToken
		testURL, err := url.Parse(urlString)
		if err != nil {
			t.Fatalf("url.Parse(%q): %v", urlString, err)
		}

		req := &http.Request{
			Method: http.MethodPost,
			Body:   ioutil.NopCloser(strings.NewReader(test.body)),
			URL:    testURL,
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("HelloHTTP http.Get: %v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("HelloHTTP ioutil.ReadAll: %v", err)
		}
		if got := string(body); got != test.want {
			t.Errorf("HelloHTTP(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}
