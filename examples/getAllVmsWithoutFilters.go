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
		Data: map[string][]string{
			string(config.VMs): {"runtime.host"},
			string(config.Hosts): {},
		},
		Metrics: map[config.PmSupportedEntities][]config.MetricDef{
			config.VMs: {
				config.MetricDef{
					Metric:   "net.packets*",
					Entities: []string{"dropbox"},
					Instance: []string{"vmnic\\d"},
				},
			},
		},

	}

	vspherePerfManager, err := pm.Init(&vspherePmConfig)

	vms, err := vspherePerfManager.Get(config.VMs)

	if err != nil {
		fmt.Println("Error Getting Vms Metrics\n", err)
	}

	for _, vm := range vms {
		//fmt.Println("VM Name: " + vm.GetProperty("name").(string))
		host := vspherePerfManager.GetProperty(vm, "runtime.host").(pm.ManagedObject)
		fmt.Println(vspherePerfManager.GetProperty(host, "name").(string))
		//fmt.Println("Host ID :" + vspherePerfManager.GetProperty("runtime.host", *vspherePerfManager).(pm.managedObject).GetProperty("name", *vspherePerfManager).(string))
		//for _, metric := range vm.Metrics {
		//	fmt.Println( "Metric Info: " + metric.Info.Metric )
		//	fmt.Println( "Metric Instance: " + metric.Value.Instance)
		//	fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		//}
	}

}


