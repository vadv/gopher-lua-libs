package prometheus_client

import (
	"github.com/prometheus/client_golang/prometheus"
	lua "github.com/yuin/gopher-lua"
)

type luaMetric struct {
	ud         *lua.LUserData
	isGauge    bool
	gauge      prometheus.Gauge
	gaugeVec   *prometheus.GaugeVec
	counter    prometheus.Counter
	counterVec *prometheus.CounterVec
	config     *promMetricConfig
}

func checkMetric(L *lua.LState, n int) *luaMetric {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaMetric); ok {
		return v
	}
	L.ArgError(n, "prometheus_client_metric_ud expected")
	return nil
}

func newMetric(L *lua.LState, isGauge bool) int {

	config := L.CheckTable(1)
	mConfig := luaTableToMetricConfig(config, L)

	fullKey := mConfig.getKey()
	if m, ok := metricCache.get(fullKey); ok {
		if m.config.equal(mConfig) {
			L.Push(m.ud)
			return 1
		}
		L.Push(lua.LNil)
		L.Push(lua.LString("already created with over config"))
		return 2
	}

	ud := L.NewUserData()
	metric := &luaMetric{config: mConfig}
	ud.Value = metric
	metric.ud = ud

	if isGauge {
		// is Gauge
		metric.isGauge = true
		if mConfig.hasLabels() {
			// is GaugeVec
			gaugeVec := prometheus.NewGaugeVec(mConfig.getGaugeOpts(), mConfig.labels)
			if err := prometheus.Register(gaugeVec); err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			metric.gaugeVec = gaugeVec
		} else {
			// is Gauge
			gauge := prometheus.NewGauge(mConfig.getGaugeOpts())
			if err := prometheus.Register(gauge); err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			metric.gauge = gauge
		}
	} else {
		// is Counter
		metric.isGauge = false
		if mConfig.hasLabels() {
			// is CounterVec
			counterVec := prometheus.NewCounterVec(mConfig.getCounterOpts(), mConfig.labels)
			if err := prometheus.Register(counterVec); err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			metric.counterVec = counterVec
		} else {
			// is Counter
			counter := prometheus.NewCounter(mConfig.getCounterOpts())
			if err := prometheus.Register(counter); err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			metric.counter = counter
		}
	}

	L.SetMetatable(ud, L.GetTypeMetatable("prometheus_client_metric_ud"))
	L.Push(ud)
	metricCache.set(fullKey, metric)
	return 1
}

// Gauge
// prometheus.gauge(config) return lua (user data, error)
// config table:
//   {
//     namespace="node_scout",
//     subsystem="nf_conntrack",
//     name="insert_failed",
//     help="insert_failed from nf_conntrack",
//     labels={"one", "two"}, -- optional
//   }
func Gauge(L *lua.LState) int {
	return newMetric(L, true)
}

// Counter
// prometheus.counter(config) return lua (user data, error)
// config table:
//   {
//     namespace="node_scout",
//     subsystem="nf_conntrack",
//     name="insert_failed",
//     help="insert_failed from nf_conntrack",
//     labels={"one", "two"}, -- optional
//   }
func Counter(L *lua.LState) int {
	return newMetric(L, false)
}

// Set lua prometheus_client_metric_ud:set(value)
func Set(L *lua.LState) int {
	metric := checkMetric(L, 1)
	if !metric.isGauge {
		L.ArgError(1, "unsupported operations for counter")
	}
	value := float64(L.CheckNumber(2))
	if metric.config.hasLabels() {
		labels := luaTableToPrometheusLabels(L.CheckTable(3))
		metric.gaugeVec.With(labels).Set(value)
	} else {
		metric.gauge.Set(value)
	}
	return 0
}

// Add lua prometheus_client_metric_ud:add(value)
func Add(L *lua.LState) int {
	metric := checkMetric(L, 1)
	value := float64(L.CheckNumber(2))
	if metric.isGauge {
		if metric.config.hasLabels() {
			labels := luaTableToPrometheusLabels(L.CheckTable(3))
			metric.gaugeVec.With(labels).Add(value)
		} else {
			metric.gauge.Add(value)
		}
	} else {
		if metric.config.hasLabels() {
			labels := luaTableToPrometheusLabels(L.CheckTable(3))
			metric.counterVec.With(labels).Add(value)
		} else {
			metric.counter.Add(value)
		}
	}
	return 0
}

// Inc lua prometheus_client_metric_ud:inc()
func Inc(L *lua.LState) int {
	metric := checkMetric(L, 1)
	if metric.isGauge {
		if metric.config.hasLabels() {
			labels := luaTableToPrometheusLabels(L.CheckTable(2))
			metric.gaugeVec.With(labels).Inc()
		} else {
			metric.gauge.Inc()
		}
	} else {
		if metric.config.hasLabels() {
			labels := luaTableToPrometheusLabels(L.CheckTable(2))
			metric.counterVec.With(labels).Inc()
		} else {
			metric.counter.Inc()
		}
	}
	return 0
}
