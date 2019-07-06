package prometheus_client

import (
	"sort"

	"github.com/prometheus/client_golang/prometheus"
	lua "github.com/yuin/gopher-lua"
)

// convert lua table to promMetricConfig
func luaTableToMetricConfig(config *lua.LTable, L *lua.LState) *promMetricConfig {
	result := &promMetricConfig{}
	config.ForEach(func(k lua.LValue, v lua.LValue) {
		switch k.String() {
		case `namespace`:
			result.namespace = v.String()
		case `subsystem`:
			result.subsystem = v.String()
		case `name`:
			result.name = v.String()
		case `help`:
			result.help = v.String()
		case `labels`:
			tbl, ok := v.(*lua.LTable)
			if !ok {
				L.ArgError(1, "labels must be string")
			}
			result.labels = luaTableToSlice(tbl)
		}
	})
	return result
}

// convert lua table to sorted []string
func luaTableToSlice(tbl *lua.LTable) []string {
	result := make([]string, 0)
	tbl.ForEach(func(k lua.LValue, v lua.LValue) {
		result = append(result, v.String())
	})
	sort.Strings(result)
	return result
}

//convert lua table to prometheus.Label
func luaTableToPrometheusLabels(tbl *lua.LTable) prometheus.Labels {
	result := make(map[string]string, 0)
	tbl.ForEach(func(k lua.LValue, v lua.LValue) {
		result[k.String()] = v.String()
	})
	return result
}
