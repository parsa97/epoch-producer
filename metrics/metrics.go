package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var ProducedMessageCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "produced_message_total",
		Help: "How many message produced in which partition.",
	},
	[]string{"partition", "topic"},
)

func Exporter() error {
	listernPort, metricsPath := exporterConfig()
	prometheus.MustRegister(ProducedMessageCounter)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`<html>
<head><title>Producer Exporter</title></head>
<body>
<p><a href= %s >Metrics</a></p>
</body>
</html>
`, metricsPath)))
	})
	http.Handle(metricsPath, promhttp.Handler())
	err := http.ListenAndServe(":"+listernPort, nil)
	if err != nil {
		log.Error(nil)
	}
	return err
}
