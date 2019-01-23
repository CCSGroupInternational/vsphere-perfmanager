package vspherePerfManager

import (
	"time"
)

type Config struct {
	Vcenter Vcenter
	QueryInterval time.Duration
	Metrics map[PmSupportedEntities][]MetricDef
	Data map[string][]string
}

type Vcenter struct {
	Username string
	Password string
	Host     string
	Insecure bool
}

type MetricDef struct {
	Metric   string
	Instance []string
	Entities []string
}

type PmSupportedEntities string

const (
	VMs      PmSupportedEntities = "VirtualMachine"
	Hosts    PmSupportedEntities = "HostSystem"
	Clusters PmSupportedEntities = "ClusterComputeResource"
)

var ALL = []string{"*"}
