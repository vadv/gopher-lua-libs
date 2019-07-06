package prometheus_client

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type promMetricConfig struct {
	namespace string
	name      string
	subsystem string
	help      string
	labels    []string
}

func (m *promMetricConfig) hasLabels() bool {
	return len(m.labels) > 0
}

func (m *promMetricConfig) getKey() string {
	return fmt.Sprintf("%s_%s_%s", m.namespace, m.subsystem, m.name)
}

func (m *promMetricConfig) equal(m2 *promMetricConfig) bool {
	if m.getKey() != m2.getKey() {
		return false
	}
	if (len(m.labels) != 0) && (len(m2.labels) != 0) {
		// because labels sorted
		return strings.Join(m.labels, "") == strings.Join(m2.labels, "")
	}
	return false
}

func (m *promMetricConfig) getGaugeOpts() prometheus.GaugeOpts {
	return prometheus.GaugeOpts{
		Namespace: m.namespace,
		Subsystem: m.subsystem,
		Name:      m.name,
		Help:      m.help,
	}
}

func (m *promMetricConfig) getCounterOpts() prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Namespace: m.namespace,
		Subsystem: m.subsystem,
		Name:      m.name,
		Help:      m.help,
	}
}
