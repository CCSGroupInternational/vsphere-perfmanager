package config

import (
	"time"
)

type VspherePerfManagerConfig struct {
	Vcenter Vcenter
	QueryInterval time.Duration
	Metrics map[EntitiesType][]MetricDef
}

type Vcenter struct {
	Username string
	Password string
	Host     string
	Insecure bool
}

type MetricDef struct {
	Metric   string
	Instance string
}

type EntitiesType string

const (
	VMs  EntitiesType = "VirtualMachines"
	Hosts EntitiesType = "Hostsystems"
)