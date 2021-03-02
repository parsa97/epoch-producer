package metrics

import (
	"os"
	"testing"
)

func TestExporterConfig(t *testing.T) {
	valueListenPort := "8090"
	valueMetricsPath := "/test"
	os.Setenv("LISTEN_PORT", valueListenPort)
	os.Setenv("METRICS_PATH", valueMetricsPath)
	listenPort, metrricsPath := exporterConfig()
	if listenPort != valueListenPort || metrricsPath != valueMetricsPath {
		t.Error()
	}
}
