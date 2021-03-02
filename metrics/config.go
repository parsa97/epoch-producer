package metrics

import "os"

func exporterConfig() (listenPort string, metricsPath string) {
	listenPort = GetEnv("LISTEN_PORT", "8080")
	metricsPath = GetEnv("METRICS_PATH", "/metrics")
	return listenPort, metricsPath
}

// GetEnv - Allows us to supply a fallback option if nothing specified
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
