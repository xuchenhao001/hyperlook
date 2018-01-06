package cadvisor

import (
	"strings"
	"fmt"

	"github.com/google/cadvisor/container"
)

type metricSetValue struct {
	container.MetricSet
}

func (ml *metricSetValue) String() string {
	var values []string
	for metric := range ml.MetricSet {
		values = append(values, string(metric))
	}
	return strings.Join(values, ",")
}

func (ml *metricSetValue) Set(value string) error {
	ml.MetricSet = container.MetricSet{}
	if value == "" {
		return nil
	}
	for _, metric := range strings.Split(value, ",") {
		if ignoreWhitelist.Has(container.MetricKind(metric)) {
			(*ml).Add(container.MetricKind(metric))
		} else {
			return fmt.Errorf("unsupported metric %q specified in disable_metrics", metric)
		}
	}
	return nil
}
