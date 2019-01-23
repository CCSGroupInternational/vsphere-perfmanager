package vspherePerfManager

import (
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/methods"
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

func (v *VspherePerfManager) getMetricsFromConfig(managedObject ManagedObject, startTime time.Time, endTime time.Time) []types.PerfQuerySpec {
	availableMetrics := v.getAvailablePerfMetrics(managedObject.Entity, startTime, endTime)

	var metricsIds []types.PerfMetricId

	for _, metricDef := range v.Config.Metrics[PmSupportedEntities(managedObject.Entity.Type)] {
		if checkEntity(metricDef, v.GetProperty(managedObject, "name").(string)) {
			metrics := getMetricsInfoFromConfig(v.metricsInfo, metricDef)

			for _, info := range metrics {
				metricsIds = setInstancesToMetrics(availableMetrics[0].MetricId, info, metricDef, metricsIds )
			}
		}
	}
	return createPerfQuerySpec(managedObject.Entity, startTime, endTime, metricsIds)
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

func setInstancesToMetrics(availableMetrics []types.PerfMetricId, metricInfo metricInfo, metricDef MetricDef, metricsIds []types.PerfMetricId) []types.PerfMetricId {

	availableMetricInstances := u.WhereBy(availableMetrics, map[string]interface{}{
		"CounterId": metricInfo.Key,
	})

	if availableMetricInstances != nil {
		for _, metricInstance := range availableMetricInstances.([]types.PerfMetricId) {
			if isToGetAllInstances(metricDef) {
				metricsIds = append(metricsIds, setMetricIds( metricInfo.Key, metricInstance.Instance))
			} else {
				for _, instance := range metricDef.Instance {
					re := regexp.MustCompile(instance)
					if re.MatchString(metricInstance.Instance) {
						metricsIds = append(metricsIds, setMetricIds( metricInfo.Key, metricInstance.Instance))
						continue
					}
				}

			}
		}
	}
	return metricsIds
}

func isToGetAllInstances(metricDef MetricDef) bool {
	return len(metricDef.Instance) == 0 || metricDef.Instance[0] == ALL[0]
}

func getMetricsInfoFromConfig(metricsInfo []metricInfo, metricDef MetricDef) []metricInfo {
	metrics := u.Where(metricsInfo, func(metric metricInfo, i int) bool {
		re := regexp.MustCompile(metricDef.Metric)
		return re.MatchString(metric.Metric)
	})

	if metrics == nil {
		return []metricInfo{}
	}
	return metrics.([]metricInfo)
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

func checkEntity(metricDef MetricDef, entityName string) bool {

	if len(metricDef.Entities) == 0 || metricDef.Entities[0] == ALL[0] {
		return true
	}

	for _, entity := range metricDef.Entities {
		re := regexp.MustCompile(entity)
		if re.MatchString(entityName) {
			return true
		}
	}
	return false
}

func setMetricIds(counterId int32, instance string) types.PerfMetricId {

	return types.PerfMetricId{
		CounterId: counterId,
		Instance:  instance,
	}

}