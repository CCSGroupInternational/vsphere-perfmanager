package vspherePerfManager

import (
	"context"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/methods"
	 u "github.com/ahl5esoft/golang-underscore"
	"time"
)

func (v *VspherePerfManager) query(managedObject *managedObject) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var metrics []types.PerfQuerySpec
	startTime, endTime := getStartEndTime(v.config.QueryInterval)

	if len(v.config.Metrics) != 0 {
		metrics = v.getMetricsFromConfig(managedObject.Entity, startTime, endTime )
	} else {
		metrics = v.getAvailablePerfMetrics(managedObject.Entity, startTime, endTime)
	}

	if len(metrics[0].MetricId) != 0 {
		perfQueryReq := types.QueryPerf{
			This: *v.client.ServiceContent.PerfManager,
			QuerySpec: metrics,
		}

		perfQueryRes, err := methods.QueryPerf(ctx, v.client.RoundTripper, &perfQueryReq )

		if err != nil {
			return err
		}

		if len(perfQueryRes.Returnval) == 0 {
			return nil
		}

		v.setMetrics(managedObject, perfQueryRes.Returnval)
	}

	return nil
}

func (v *VspherePerfManager) setMetrics(managedObject *managedObject, metrics []types.BasePerfEntityMetricBase) {
	for _, base := range metrics {
		pem := base.(*types.PerfEntityMetric)
		for _, baseSerie := range pem.Value {
			serie := baseSerie.(*types.PerfMetricIntSeries)

			info := u.WhereBy(v.metricsInfo, map[string]interface{}{
				"Key": serie.Id.CounterId,
			}).([]metricInfo)

			for _, value := range serie.Value {
				managedObject.Metrics = append(managedObject.Metrics, Metric{
					Info: info[0],
					Value: metricValue{
						Value: value,
						Instance: serie.Id.Instance,
					},
				})
			}
		}
	}
}

func getStartEndTime(interval time.Duration) (time.Time, time.Time){
	endTime := time.Now().Add(time.Duration(-1) * time.Second)
	startTime := endTime.Add(-interval -1  * time.Second)
	return startTime, endTime
}
