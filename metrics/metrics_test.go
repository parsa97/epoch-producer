package metrics

import (
	"os"
	"testing"
	"time"

	"github.com/lifeforms/httpcheck/httpcheck"
)

func TestExporter(t *testing.T) {
	os.Setenv("LISTEN_PORT", "8080")
	os.Setenv("METRICS_PATH", "/metrics")
	M := ProducedMessageCounter.WithLabelValues("0", "test")
	M.Inc()
	go Exporter()
	time.Sleep(2 * time.Second)
	httpcheck.Verbose = true
	httpcheck.RequestTimeout = 5
	httpcheck.ServerTimeout = 10

	manifest := httpcheck.Manifest{
		httpcheck.Server{
			Name: "exporter",
			Scenarios: []httpcheck.Scenario{
				httpcheck.Scenario{
					Name:  "http",
					Tests: httpcheck.Tests{httpcheck.Test{Url: "http://localhost:8080/metrics", Content: "produced_message_total"}}},
			},
		},
	}

	if manifest.Test() != nil {
		t.Error()
	}

}
