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
					Instance: "*",
				},
				config.MetricDef{
					Metric:   "cpu.usagemhz.average",
					Instance: "",
				},
			},
		},
	}

	vspherePerfManager, err := pm.Init(&vspherePmConfig)

	vspherePerfManager.Vms()

}


