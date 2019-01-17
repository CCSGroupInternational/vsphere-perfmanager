package vspherePerfManager

import (
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/CCSGroupInternational/vsphere-perfmanager/config"
	"context"
	"github.com/vmware/govmomi/vim25/mo"
	u "github.com/ahl5esoft/golang-underscore"
	"time"
	"regexp"
)

type Metric struct {
	Info  metricInfo
	Value metricValue
}

type metricInfo struct {
	Metric    string
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

		metrics := getMetricsInfoFromConfig(v.metricsInfo, metricDef)

		for _, info := range metrics {

			if len(metricDef.Instance) == 0 || metricDef.Instance[0] == config.AllInstances[0] {
				metricsIds = setAllInstancesToMetrics(availableMetrics[0].MetricId, info, metricsIds)
				continue
			}

			for _, instance := range metricDef.Instance {
				metricsIds = append(metricsIds, types.PerfMetricId{
					CounterId: info.Key,
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
			Metric   :  metric.GroupInfo.GetElementDescription().Key + "." + metric.NameInfo.GetElementDescription().Key + "." + string(metric.RollupType),
			StatsType : string(metric.StatsType),
			UnitInfo  : metric.UnitInfo.GetElementDescription().Key,
			Key       : metric.Key,
		})
	}

	return metrics, nil

}

func hasMetricsWithAllInstances(metrics []config.MetricDef) bool {
	metricDefAllInstances := u.Where(metrics, func(metricDef config.MetricDef, i int) bool {
		if len(metricDef.Instance) == 0 {
			return true
		}
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

func setAllInstancesToMetrics(availableMetrics []types.PerfMetricId, metricInfo metricInfo, metricsIds []types.PerfMetricId) []types.PerfMetricId {
	availableMetricInstances := u.WhereBy(availableMetrics, map[string]interface{}{
		"CounterId": metricInfo.Key,
	})

	if availableMetricInstances != nil {
		for _, metricInstance := range availableMetricInstances.([]types.PerfMetricId) {
			metricsIds = append(metricsIds, types.PerfMetricId{
				CounterId: metricInfo.Key,
				Instance:  metricInstance.Instance,
			})
		}
	}

	return metricsIds
}

func getMetricsInfoFromConfig(metricsInfo []metricInfo, metricDef config.MetricDef) []metricInfo {

	metrics := u.Where(metricsInfo, func(metric metricInfo, i int) bool {
		re := regexp.MustCompile(metricDef.Metric)
		return re.MatchString(metric.Metric)
	})

	if metrics == nil {
		return []metricInfo{}
	}
	return metrics.([]metricInfo)
}