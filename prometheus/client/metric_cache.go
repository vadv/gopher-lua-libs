package prometheus_client

import (
	"sync"
)

var metricCache = newPrometheusMetricCache()

type promMetricCache struct {
	lock  *sync.Mutex
	cache map[string]*luaMetric
}

func newPrometheusMetricCache() *promMetricCache {
	return &promMetricCache{
		lock:  &sync.Mutex{},
		cache: make(map[string]*luaMetric, 0),
	}
}

func (c *promMetricCache) get(key string) (*luaMetric, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	m, ok := c.cache[key]
	return m, ok
}

func (c *promMetricCache) set(key string, m *luaMetric) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[key] = m
}
