package cloudwatch

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	lua "github.com/yuin/gopher-lua"
)

// GetMetricData lua clw:get_metric_data({start_time=1, end_time=1, queries={}}) result, error
func GetMetricData(L *lua.LState) int {
	clw := checkluaClW(L, 1)
	config := L.CheckTable(2)

	startTime := time.Now().Add(-10 * time.Minute)
	endTime := time.Now()
	var err error
	queries := make([]*cloudwatch.MetricDataQuery, 0)

	// parse parameters
	config.ForEach(func(k lua.LValue, v lua.LValue) {
		switch k.String() {
		case `start_time`, `end_time`:
			if value, ok := v.(lua.LNumber); !ok {
				err = fmt.Errorf("time must be number")
				return
			} else {
				if k.String() == `start_time` {
					startTime = time.Unix(int64(value), 0)
				} else {
					endTime = time.Unix(int64(value), 0)
				}
			}
		case `queries`:
			if value, ok := v.(*lua.LTable); !ok {
				err = fmt.Errorf("queries must be table")
				return
			} else {
				value.ForEach(func(k lua.LValue, v lua.LValue) {
					name := k.String()
					tbl, ok := v.(*lua.LTable)
					if !ok {
						err = fmt.Errorf("query must be table")
						return
					}
					query, errParse := parseGetMetricDataQuery(tbl, name)
					if errParse != nil {
						err = errParse
						return
					}
					queries = append(queries, query)
				})
			}
		}
	})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	clwResult, err := clw.cloudWatchClient.GetMetricData(&cloudwatch.GetMetricDataInput{
		StartTime:         &startTime,
		EndTime:           &endTime,
		MetricDataQueries: queries,
	})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	for _, metric := range clwResult.MetricDataResults {
		luaRow := L.NewTable()

		var id string
		if metric.Id == nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("unknown id in response"))
			return 2
		}
		id = *metric.Id

		for i := 0; i < len(metric.Timestamps); i++ {
			tt := metric.Timestamps[i]
			value := metric.Values[i]
			if tt == nil || value == nil {
				continue
			}
			ts := *tt
			luaRow.RawSet(lua.LNumber(ts.Unix()), lua.LNumber(*value))
		}

		result.RawSetString(id, luaRow)
	}
	L.Push(result)
	return 1
}

func parseGetMetricDataQuery(tbl *lua.LTable, id string) (*cloudwatch.MetricDataQuery, error) {
	var metric, namespace, dimensionName, dimensionValue string
	stat := `Average`
	period := int64(60)
	var err error
	tbl.ForEach(func(k lua.LValue, v lua.LValue) {
		switch k.String() {
		case `namespace`:
			namespace = v.String()
		case `metric`:
			metric = v.String()
		case `dimension_name`:
			dimensionName = v.String()
		case `dimension_value`:
			dimensionValue = v.String()
		case `stat`:
			stat = v.String()
		case `period`:
			if value, ok := v.(lua.LNumber); !ok {
				err = fmt.Errorf("period must be number")
				return
			} else {
				period = int64(value)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	query := &cloudwatch.MetricDataQuery{
		Id: &id,
		MetricStat: &cloudwatch.MetricStat{
			Metric: &cloudwatch.Metric{
				Namespace:  &namespace,
				MetricName: &metric,
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  &dimensionName,
						Value: &dimensionValue,
					},
				},
			},
			Period: &period,
			Stat:   &stat,
		},
	}
	return query, nil
}
