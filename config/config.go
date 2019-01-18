package config

import (
	"time"
	"github.com/vmware/govmomi/vim25/types"
)

type VspherePerfManagerConfig struct {
	Vcenter Vcenter
	QueryInterval time.Duration
	Metrics map[EntitiesType][]MetricDef
	Properties []types.PropertySpec
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

type EntitiesType string

const (
	VMs        EntitiesType = "VirtualMachine"
	Hosts      EntitiesType = "HostSystem"
)

var ALL = []string{"*"}
