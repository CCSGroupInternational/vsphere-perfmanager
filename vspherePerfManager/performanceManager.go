package vspherePerfManager

import (
	"context"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/methods"
	 u "github.com/ahl5esoft/golang-underscore"
)

func (v *VspherePerfManager) query(managedObject *managedObject) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var metrics []types.PerfQuerySpec
	if len(v.config.Metrics) != 0 {
		metrics = v.getMetricsFromConfig(managedObject.Entity )
	} else {
		metrics = v.getAvailablePerfMetrics(managedObject.Entity)
	}

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
					Info: metricInfo{
						Counter   : info[0].Counter,
						Group     : info[0].Group,
						Rollup    : info[0].Rollup,
						StatsType : info[0].StatsType,
						UnitInfo  : info[0].UnitInfo,
						Key       : info[0].Key,
					},
					Value: metricValue{
						Value: value,
						Instance: serie.Id.Instance,
					},
				})
			}
		}
	}
}