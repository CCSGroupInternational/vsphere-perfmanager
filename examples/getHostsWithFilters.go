package main

import (
	"fmt"
	"strconv"
	"os"
	"github.com/CCSGroupInternational/vsphere-perfmanager/config"
	pm "github.com/CCSGroupInternational/vsphere-perfmanager/vspherePerfManager"
	"time"
	"github.com/vmware/govmomi/vim25/types"
)

func main() {
	insecure, err := strconv.ParseBool(os.Getenv("VSPHERE_INSECURE"))

	if err != nil {
		fmt.Println("Error to convert VSPHERE_INSECURE env var to bool type\n", err)
	}

	vspherePmConfig := config.VspherePerfManagerConfig{
		Vcenter: config.Vcenter {
			Username : os.Getenv("VSPHERE_USER"),
			Password : os.Getenv("VSPHERE_PASSWORD"),
			Host     : os.Getenv("VSPHERE_HOST"),
			Insecure : insecure,
		},
		QueryInterval: time.Duration(20) * time.Second,
		Metrics: map[config.PmSupportedEntities][]config.MetricDef{
			config.Hosts: {
				config.MetricDef{
					Metric:   "cpu.usage.average",
					Entities: config.ALL,
					Instance: []string{"^0$"},
				},
				config.MetricDef{
					Metric:   "datastore.datastoreVMObservedLatency",
					Instance: []string{"5a5e17e3-e66424d4-7eb5-7ca23e8e9504"},
				},
			},
		},
		Properties: []types.PropertySpec{{
			Type: string(config.Hosts),
			PathSet: []string{"parent"},
		}},
	}

	vspherePerfManager, err := pm.Init(&vspherePmConfig)

	hosts, err := vspherePerfManager.Hosts()

	if err != nil {
		fmt.Println("Error Getting Hosts Metrics\n", err)
	}

	for _, host := range hosts {
		fmt.Println("Host Name: " + host.GetProperty("name").(string))
		fmt.Println("Cluster ID: " + host.GetProperty("parent").(types.ManagedObjectReference).Value)
		for _, metric := range host.Metrics {
			fmt.Println( "Metric : " + metric.Info.Metric )
			fmt.Println( "Metric Instance: " + metric.Value.Instance)
			fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		}
	}

}


