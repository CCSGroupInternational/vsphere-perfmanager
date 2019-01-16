package vspherePerfManager

import (
	"github.com/vmware/govmomi/vim25/types"
	"time"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/CCSGroupInternational/vsphere-perfmanager/config"
	"strings"
	"context"
	"github.com/vmware/govmomi/vim25/mo"
	u "github.com/ahl5esoft/golang-underscore"
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

func getStartEndTime(interval time.Duration) (time.Time, time.Time){
	endTime := time.Now().Add(time.Duration(-1) * time.Second)
	startTime := endTime.Add(-interval -1  * time.Second)
	return startTime, endTime
}

func (v *VspherePerfManager) getAvailablePerfMetrics(entity types.ManagedObjectReference) []types.PerfQuerySpec {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime, endTime := getStartEndTime(v.config.QueryInterval)

	perfReq := types.QueryAvailablePerfMetric{
		This:       *v.client.ServiceContent.PerfManager,
		Entity:     entity,
		BeginTime:  &startTime,
		EndTime:    &endTime,
		IntervalId: int32(20),
	}

	perfRes, _ := methods.QueryAvailablePerfMetric(ctx, v.client.RoundTripper, &perfReq )

	return []types.PerfQuerySpec{{
		Entity:     entity,
		StartTime:  &startTime,
		EndTime:    &endTime,
		MetricId:   perfRes.Returnval,
		IntervalId: int32(20),
	}}
}

func (v *VspherePerfManager) getMetricsFromConfig(entity types.ManagedObjectReference) []types.PerfQuerySpec {

	startTime, endTime := getStartEndTime(v.config.QueryInterval)

	var metricsIds []types.PerfMetricId

	for _, metric := range v.config.Metrics[config.EntitiesType(entity.Type)] {
		metricInfo := u.WhereBy(v.metricsInfo, map[string]interface{}{
			"Counter": strings.Split(metric.Metric, ".")[1],
			"Group":   strings.Split(metric.Metric, ".")[0],
			"Rollup":  strings.Split(metric.Metric, ".")[2],
		}).([]metricInfo)

		metricsIds = append(metricsIds, types.PerfMetricId{
			CounterId: metricInfo[0].Key,
			Instance:  metric.Instance,
		})
	}
	return []types.PerfQuerySpec{{
		Entity:     entity,
		StartTime:  &startTime,
		EndTime:    &endTime,
		MetricId:   metricsIds,
		IntervalId: int32(20),
	}}
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