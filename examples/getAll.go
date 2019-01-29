package main

import (
	"strconv"
	"os"
	"fmt"
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
			Samples: 6,
			Data: map[string][]string{
				string(pm.VMs):      {"runtime.host"},
				string(pm.Hosts):    {"parent"},
				pm.Clusters: {},
				string(pm.ResourcePools): {"parent", "vm"},
				string(pm.Datastores): {"summary.url"},
			},
			Metrics: map[pm.PmSupportedEntities][]pm.MetricDef{
				pm.Datastores: {
					pm.MetricDef{
						Entities: []string{"VMWareCP03"},
						Metric: []string{"disk.unshared.latest"},
					},
				},
				pm.Hosts: {
					pm.MetricDef{
						Metric: []string{"net.packets*"},
					},
				},
				pm.VMs: {
					pm.MetricDef{
						Metric:   []string{"net.packets*"},
						Entities: []string{"openshift01"},
						Instance: []string{"vmnic\\d"},
					},
				},
			},
		},
	}


	err = vspherePm.Init()

	if err != nil {
		fmt.Println("Error on Initializing Vsphere Performance Manager\n", err)
	}

	vms := vspherePm.Get(pm.VMs)

	for _, vm := range vms {
		fmt.Println("VM Name: " + vspherePm.GetProperty(vm, "name").(string))
		host := vspherePm.GetProperty(vm, "runtime.host").(pm.ManagedObject)
		fmt.Println("Host Name :" + vspherePm.GetProperty(host, "name").(string))
		fmt.Println("Cluster Name :" + vspherePm.GetProperty(vspherePm.GetProperty(host, "parent").(pm.ManagedObject), "name").(string))
		for _, metric := range vm.Metrics {
			fmt.Println("Metric : " + metric.Info.Metric)
			fmt.Println("Metric Instance: " + metric.Value.Instance)
			fmt.Println("Result: " + strconv.FormatInt(metric.Value.Value, 10))
		}
	}


	hosts := vspherePm.Get(pm.Hosts)

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

	dataStores := vspherePm.Get(pm.Datastores)

	for _, dataStore := range dataStores {
		fmt.Println("Datastore Name: " + vspherePm.GetProperty(dataStore, "name").(string))
		fmt.Println("Summary Url: " + vspherePm.GetProperty(dataStore, "summary.url").(string) )
		for _, metric := range dataStore.Metrics {
			fmt.Println( "Metric : " + metric.Info.Metric )
			var instance string
			if len(metric.Value.Instance) != 0 {
				if _, err := strconv.Atoi(metric.Value.Instance); err == nil {
					instance = vspherePm.GetProperty(vspherePm.GetObject(string(pm.VMs), "vm-" + metric.Value.Instance), "name").(string)
				} else {
					instance = metric.Value.Instance
				}
			} else {
				instance = metric.Value.Instance
			}
			fmt.Println("Metric Instance: " + instance)
			fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		}
	}

	resourcePools := vspherePm.Get(pm.ResourcePools)

	if err != nil {
		fmt.Println("Error Getting ResourcePool Metrics\n", err)
	}

	for _, resourcePool := range resourcePools {
		fmt.Println("Resource Pool: " + vspherePm.GetProperty(resourcePool, "name").(string))
		switch parentType := vspherePm.GetProperty(resourcePool, "parent").(pm.ManagedObject).Entity.Type; parentType {
		case string(pm.Clusters):
			fmt.Println("Cluster Name: " + vspherePm.GetProperty(vspherePm.GetProperty(resourcePool, "parent").(pm.ManagedObject), "name").(string))
		case string(pm.ResourcePools):
			fmt.Println("Cluster Name: " + vspherePm.GetProperty(vspherePm.GetProperty(vspherePm.GetProperty(resourcePool, "parent").(pm.ManagedObject),"parent").(pm.ManagedObject), "name").(string))
		}
		for _, metric := range resourcePool.Metrics {
			fmt.Println( "Metric : " + metric.Info.Metric )
			fmt.Println( "Metric Instance: " + metric.Value.Instance)
			fmt.Println( "Result: " + strconv.FormatInt(metric.Value.Value, 10) )
		}
	}

	clusters := vspherePm.Get(pm.Clusters)

	if err != nil {
		fmt.Println("Error Getting ResourcePool Metrics\n", err)
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
