package vspherePerfManager

import (
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/CCSGroupInternational/vsphere-perfmanager/config"
	"strings"
	"context"
	"github.com/vmware/govmomi/vim25/mo"
	u "github.com/ahl5esoft/golang-underscore"
	"time"
)

type Metric struct {
	Info  metricInfo
	Value metricValue
}

type metricInfo struct {
	Counter   string
	Group     string
	Rollup    string
	StatsType string
	UnitInfo  string
	Key       int32
}

type metricValue struct {
	Value    int64
	Instance string
}

func (v *VspherePerfManager) getAvailablePerfMetrics(entity types.ManagedObjectReference, startTime time.Time, endTime time.Time) []types.PerfQuerySpec {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	perfReq := types.QueryAvailablePerfMetric{
		This:       *v.client.ServiceContent.PerfManager,
		Entity:     entity,
		BeginTime:  &startTime,
		EndTime:    &endTime,
		IntervalId: int32(20),
	}

	perfRes, _ := methods.QueryAvailablePerfMetric(ctx, v.client.RoundTripper, &perfReq )

	return createPerfQuerySpec(entity, startTime, endTime, perfRes.Returnval)
}

func (v *VspherePerfManager) getMetricsFromConfig(entity types.ManagedObjectReference, startTime time.Time, endTime time.Time) []types.PerfQuerySpec {

	var availableMetrics []types.PerfQuerySpec

	if hasMetricsWithAllInstances(v.config.Metrics[config.EntitiesType(entity.Type)]) {
		availableMetrics = v.getAvailablePerfMetrics(entity, startTime, endTime)
	}

	var metricsIds []types.PerfMetricId

	for _, metricDef := range v.config.Metrics[config.EntitiesType(entity.Type)] {
		info := u.WhereBy(v.metricsInfo, map[string]interface{}{
			"Counter": strings.Split(metricDef.Metric, ".")[1],
			"Group":   strings.Split(metricDef.Metric, ".")[0],
			"Rollup":  strings.Split(metricDef.Metric, ".")[2],
		})

		if info == nil {
			continue
		}

		if metricDef.Instance[0] == config.AllInstances[0] {

			availableMetricInstances := u.WhereBy(availableMetrics[0].MetricId, map[string]interface{} {
				"CounterId": info.([]metricInfo)[0].Key,
			})

			if availableMetricInstances != nil {
				for _, metricInstance := range availableMetricInstances.([]types.PerfMetricId) {
					metricsIds = append(metricsIds, types.PerfMetricId{
						CounterId: info.([]metricInfo)[0].Key,
						Instance:  metricInstance.Instance,
					})
				}
			}

		} else {
			for _, instance := range metricDef.Instance {
				metricsIds = append(metricsIds, types.PerfMetricId{
					CounterId: info.([]metricInfo)[0].Key,
					Instance:  instance,
				})
			}
		}

	}

	return createPerfQuerySpec(entity, startTime, endTime, metricsIds)
}

func (v *VspherePerfManager) getMetricsInfo() ([]metricInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var perfmanager mo.PerformanceManager
	err := v.client.RetrieveOne(ctx, *v.client.ServiceContent.PerfManager, nil, &perfmanager)
	if err != nil {
		return nil, err
	}

	var metrics []metricInfo

	for _, metric := range perfmanager.PerfCounter {
		metrics = append(metrics, metricInfo{
			Counter   : metric.NameInfo.GetElementDescription().Key,
			Group     : metric.GroupInfo.GetElementDescription().Key,
			Rollup    : string(metric.RollupType),
			StatsType : string(metric.StatsType),
			UnitInfo  : metric.UnitInfo.GetElementDescription().Key,
			Key       : metric.Key,
		})
	}

	return metrics, nil

}

func hasMetricsWithAllInstances(metrics []config.MetricDef) bool {
	metricDefAllInstances := u.Where(metrics, func(metricDef config.MetricDef, i int) bool {
		return metricDef.Instance[0] == config.AllInstances[0]
	})

	if metricDefAllInstances != nil {
		return true
	}
	return false

}

func createPerfQuerySpec(entity types.ManagedObjectReference, startTime time.Time, endTime time.Time, metricsIds []types.PerfMetricId) []types.PerfQuerySpec {
	return []types.PerfQuerySpec{{
		Entity:     entity,
		StartTime:  &startTime,
		EndTime:    &endTime,
		MetricId:   metricsIds,
		IntervalId: int32(20),
	}}

}