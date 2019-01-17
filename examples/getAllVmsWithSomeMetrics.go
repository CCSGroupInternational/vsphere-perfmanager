package main

import (
	"fmt"
	"strconv"
	"os"
	"github.com/CCSGroupInternational/vsphere-perfmanager/config"
	pm "github.com/CCSGroupInternational/vsphere-perfmanager/vspherePerfManager"
	"time"
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
		Metrics: map[config.EntitiesType][]config.MetricDef{
			config.VMs: {
				config.MetricDef{
					Metric:   "cpu.usage.average",
					Instance: []string{""},
				},
				config.MetricDef{
					Metric:   "cpu.usagemhz.average",
					Instance: config.AllInstances,
				},
			},
		},
	}

	vspherePerfManager, err := pm.Init(&vspherePmConfig)

	vms, err := vspherePerfManager.Vms()


	if err != nil {
		fmt.Println("Error Getting Vms Metrics\n", err)
	}

	for _, vm := range vms {
		fmt.Println("VM Name: " + vm.Properties[0].Val.(string))
		for _, metric := range vm.Metrics {
			fmt.Println( "Metric Info: " + metric.Info.Group + "." + metric.Info.Counter + "." + metric.Info.Rollup )
			fmt.Println( "Metric Instance: " + metric.Value.Instance)
			fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		}
	}

}


