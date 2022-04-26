package parser

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/vontikov/prom-redis/internal/collector"
	"github.com/vontikov/prom-redis/internal/logging"
)

type Pair = collector.Pair

type converter func(v string) float64

var logger = logging.NewLogger("parser")

var converters = map[string]converter{
	"redis_mode": func(v string) float64 {
		if v == "standalone" {
			return 0.0
		}
		return 1.0
	},
	"process_supervised": func(v string) float64 {
		if v == "no" {
			return 0.0
		}
		return 1.0
	},
	"role": func(v string) float64 {
		if v == "master" {
			return 1.0
		}
		return 0.0
	},
}

func Parse(in *string) []Pair {
	var res []Pair

	scanner := bufio.NewScanner(strings.NewReader(*in))
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) == 0 || strings.HasPrefix(s, "#") {
			continue
		}
		arr := strings.Split(s, ":")
		k := arr[0]
		v := arr[1]

		if c, ok := converters[k]; ok {
			res = append(res, Pair{K: k, V: c(v)})
			continue
		}

		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			res = append(res, Pair{K: k, V: f})
			continue
		}

		logger.Trace("could not parse", "key", k, "value", v)
	}

	return res
}
