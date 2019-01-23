package main

import (
	"strconv"
	"os"
	"fmt"
	"time"
	pm "github.com/CCSGroupInternational/vsphere-perfmanager/vspherePerfManager"
)

func main() {
	insecure, err := strconv.ParseBool(os.Getenv("VSPHERE_INSECURE"))

	if err != nil {
		fmt.Println("Error to convert VSPHERE_INSECURE env var to bool type\n", err)
	}

	vspherePm := pm.VspherePerfManager{
		Config: pm.Config{
			Vcenter: pm.Vcenter{
				Username: os.Getenv("VSPHERE_USER"),
				Password: os.Getenv("VSPHERE_PASSWORD"),
				Host:     os.Getenv("VSPHERE_HOST"),
				Insecure: insecure,
			},
			QueryInterval: time.Duration(20) * time.Second,
			Data: map[string][]string{
				string(pm.VMs):      {"runtime.host"},
				string(pm.Hosts):    {"parent"},
				pm.Clusters: {},
			},
		},
	}

	err = vspherePm.Init()

	if err != nil {
		fmt.Println("Error on Initializing Vsphere Performance Manager\n", err)
	}

	//vms, err := vspherePm.Get(pm.VMs)
	//
	//if err != nil {
	//	fmt.Println("Error Getting Vms Metrics\n", err)
	//}
	//
	//for _, vm := range vms {
	//	fmt.Println("VM Name: " + vspherePm.GetProperty(vm, "name").(string))
	//	host := vspherePm.GetProperty(vm, "runtime.host").(pm.ManagedObject)
	//	fmt.Println("Host Name :" + vspherePm.GetProperty(host, "name").(string))
	//	fmt.Println("Cluster Name :" + vspherePm.GetProperty(vspherePm.GetProperty(host, "parent").(pm.ManagedObject), "name").(string))
	//	for _, metric := range vm.Metrics {
	//		fmt.Println("Metric : " + metric.Info.Metric)
	//		fmt.Println("Metric Instance: " + metric.Value.Instance)
	//		fmt.Println("Result: " + strconv.FormatInt(metric.Value.Value, 10))
	//	}
	//}


	//hosts, err := vspherePm.Get(pm.Hosts)
	//
	//if err != nil {
	//	fmt.Println("Error Getting Hosts Metrics\n", err)
	//}
	//
	//for _, host := range hosts {
	//	fmt.Println("Host Name: " + vspherePm.GetProperty(host, "name").(string))
	//	fmt.Println("Cluster Name: " + vspherePm.GetProperty(vspherePm.GetProperty(host, "parent").(pm.ManagedObject),"name").(string))
	//	for _, metric := range host.Metrics {
	//		fmt.Println( "Metric : " + metric.Info.Metric )
	//		fmt.Println( "Metric Instance: " + metric.Value.Instance)
	//		fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
	//	}
	//}

	clusters, err := vspherePm.Get(pm.Clusters)

	if err != nil {
		fmt.Println("Error Getting Hosts Metrics\n", err)
	}


	for _, cluster := range clusters {
		fmt.Println("Cluster Name: " + vspherePm.GetProperty(cluster, "name").(string))
		for _, metric := range cluster.Metrics {
			fmt.Println( "Metric : " + metric.Info.Metric )
			fmt.Println( "Metric Instance: " + metric.Value.Instance)
			fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		}
	}

}
