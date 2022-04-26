package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Importer func() *string

type Pair struct {
	K string
	V float64
}

type Parser func(*string) []Pair

type redisCollector struct {
	f  Importer
	p  Parser
	ns string
}

func NewCollector(f Importer, p Parser, ns string) prometheus.Collector {
	return &redisCollector{
		f:  f,
		p:  p,
		ns: ns,
	}
}

func (c *redisCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *redisCollector) Collect(ch chan<- prometheus.Metric) {
	info := c.f()
	if info == nil {
		return
	}
	for _, p := range c.p(info) {
		d := prometheus.NewDesc(fmt.Sprintf("%s_%s", c.ns, p.K), p.K, nil, nil)
		ch <- prometheus.MustNewConstMetric(d, prometheus.GaugeValue, p.V)
	}
}
