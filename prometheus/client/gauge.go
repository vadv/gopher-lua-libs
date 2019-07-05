package prometheus_client

import (
	"fmt"

	prometheus "github.com/prometheus/client_golang/prometheus"
	lua "github.com/yuin/gopher-lua"
)

var regGaugeUserData = make(map[string]*lua.LUserData, 0)

type luaGauge struct {
	prometheus.Gauge
}

func checkGauge(L *lua.LState, n int) *luaGauge {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaGauge); ok {
		return v
	}
	L.ArgError(n, "prometheus_client_gauge_ud expected")
	return nil
}

// Gauge
// prometheus.gauge(config) return lua (user data, error)
// config table:
//   {
//     namespace="node_scout",
//     subsystem="nf_conntrack",
//     name="insert_failed",
//     help="insert_failed from nf_conntrack",
//   }
func Gauge(L *lua.LState) int {

	config := L.CheckTable(1)

	namespace, subsystem, name, help := "", "", "", ""
	config.ForEach(func(k lua.LValue, v lua.LValue) {
		switch k.String() {
		case `namespace`:
			namespace = v.String()
		case `subsystem`:
			subsystem = v.String()
		case `name`:
			name = v.String()
		case `help`:
			help = v.String()
		}
	})

	fullName := fmt.Sprintf("%s_%s_%s", namespace, subsystem, name)
	if ud, ok := regGaugeUserData[fullName]; ok {
		L.Push(ud)
		return 1
	}

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})

	if err := prometheus.Register(gauge); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = &luaGauge{gauge}
	L.SetMetatable(ud, L.GetTypeMetatable("prometheus_client_gauge_ud"))
	L.Push(ud)
	regGaugeUserData[fullName] = ud
	return 1
}

// GaugeSet lua prometheus_client_gauge:set(value)
func GaugeSet(L *lua.LState) int {
	gauge := checkGauge(L, 1)
	value := L.CheckNumber(2)
	gauge.Set(float64(value))
	return 0
}

// GaugeAdd lua prometheus_client_gauge:set(value)
func GaugeAdd(L *lua.LState) int {
	gauge := checkGauge(L, 1)
	value := L.CheckNumber(2)
	gauge.Add(float64(value))
	return 0
}
