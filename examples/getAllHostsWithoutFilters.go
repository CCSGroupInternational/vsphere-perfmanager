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
			QueryInterval: time.Duration(20) * time.Second,
			Data: map[string][]string{
				string(pm.Hosts): {"parent"},
				string(pm.Clusters): {},
			},
		},
	}


	err = vspherePm.Init()

	if err != nil {
		fmt.Println("Error on Initializing Vsphere Performance Manager\n", err)
	}

	hosts, err := vspherePm.Get(pm.Hosts)

	if err != nil {
		fmt.Println("Error Getting Hosts Metrics\n", err)
	}


	for _, host := range hosts {
		fmt.Println("Host Name: " + vspherePm.GetProperty(host, "name").(string))
		fmt.Println("Cluster Name: " + vspherePm.GetProperty(vspherePm.GetProperty(host, "parent").(pm.ManagedObject),"name").(string))
		for _, metric := range host.Metrics {
			fmt.Println( "Metric : " + metric.Info.Metric )
			fmt.Println( "Metric Instance: " + metric.Value.Instance)
			fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		}
	}

}


