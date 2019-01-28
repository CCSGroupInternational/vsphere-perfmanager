package main

import (
	"fmt"
	"strconv"
	"os"
	pm "github.com/CCSGroupInternational/vsphere-perfmanager/vspherePerfManager"
	"time"
)

func main() {
	insecure, err := strconv.ParseBool(os.Getenv("VSPHERE_INSECURE"))

	if err != nil {
		fmt.Println("Error to convert VSPHERE_INSECURE env var to bool type\n", err)
	}

	vspherePm := pm.VspherePerfManager{
		Config: pm.Config {
			Vcenter: pm.Vcenter {
				Username : os.Getenv("VSPHERE_USER"),
				Password : os.Getenv("VSPHERE_PASSWORD"),
				Host     : os.Getenv("VSPHERE_HOST"),
				Insecure : insecure,
			},
			Samples: time.Duration(20) * time.Second,
			Data: map[string][]string{
				string(pm.VMs): {"runtime.host"},
				string(pm.Hosts): {},
			},
		},
	}

	err = vspherePm.Init()

	if err != nil {
		fmt.Println("Error on Initializing Vsphere Performance Manager\n", err)
	}

	vms, err := vspherePm.Get(pm.VMs)

	if err != nil {
		fmt.Println("Error Getting Vms Metrics\n", err)
	}

	for _, vm := range vms {
		fmt.Println("VM Name: " + vspherePm.GetProperty(vm, "name").(string))
		fmt.Println("Host Name :" + vspherePm.GetProperty(vspherePm.GetProperty(vm,"runtime.host").(pm.ManagedObject), "name").(string))
		for _, metric := range vm.Metrics {
			fmt.Println( "Metric : " + metric.Info.Metric )
			fmt.Println( "Metric Instance: " + metric.Value.Instance)
			fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		}
	}

}


