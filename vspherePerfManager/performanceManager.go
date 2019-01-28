package vspherePerfManager

import (
	"github.com/vmware/govmomi/vim25/types"
	"time"
	"github.com/vmware/govmomi/vim25/methods"
	"context"
	"github.com/vmware/govmomi"
)

func (v *VspherePerfManager) query(managedObject ManagedObject) (ManagedObject, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	summary, err := v.ProviderSummary(managedObject.Entity)
	if summary.RefreshRate == -1 {
		summary.RefreshRate = 300
	}

	if err != nil {
		return managedObject, err
	}

	startTime, err := getStartTime(v.Config.Samples, summary.RefreshRate, v.client )

	if err != nil {
		return managedObject, err
	}

	metrics := v.getAvailablePerfMetrics(managedObject.Entity, summary.RefreshRate, &startTime)
	metrics = v.filterWithConfig(metrics, managedObject)
	metricsSpec := createPerfQuerySpec(managedObject.Entity, metrics, summary.RefreshRate, &startTime)

	if len(metricsSpec[0].MetricId) != 0 {
		perfQueryReq := types.QueryPerf{
			This: *v.client.ServiceContent.PerfManager,
			QuerySpec: metricsSpec,
		}
		perfQueryRes, err := methods.QueryPerf(ctx, v.client.RoundTripper, &perfQueryReq )

		if err != nil {
			return managedObject, err
		}

		if len(perfQueryRes.Returnval) == 0 {
			return managedObject, err
		}

		v.setMetrics(&managedObject, perfQueryRes.Returnval)
	}

	return managedObject, nil
}

func (v *VspherePerfManager) setMetrics(managedObject *ManagedObject, metrics []types.BasePerfEntityMetricBase) {
	for _, base := range metrics {
		pem := base.(*types.PerfEntityMetric)

		for _, baseSerie := range pem.Value {
			serie := baseSerie.(*types.PerfMetricIntSeries)
			for _, value := range serie.Value {
				managedObject.Metrics = append(managedObject.Metrics, Metric{
					Info: v.metricsInfo[serie.Id.CounterId],
					Value: metricValue{
						Value: value,
						Instance: serie.Id.Instance,
					},
				})
			}
		}
	}
}

func getStartTime(samples int32, intervalId int32, client *govmomi.Client) (time.Time, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now, err := methods.GetCurrentTime(ctx, client)

	x := intervalId * -1 * samples
	return now.Add(time.Duration(x) * time.Second), err

}
