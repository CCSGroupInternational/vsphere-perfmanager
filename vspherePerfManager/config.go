package vspherePerfManager

type Config struct {
	Vcenter Vcenter
	Samples int32
	Metrics map[PmSupportedEntities][]MetricDef
	Data    map[string][]string
}

type Vcenter struct {
	Username string
	Password string
	Host     string
	Insecure bool
}

type MetricDef struct {
	Metric   []string
	Instance []string
	Entities []string
}

type PmSupportedEntities string

const (
	VMs           PmSupportedEntities = "VirtualMachine"
	Hosts         PmSupportedEntities = "HostSystem"
	ResourcePools PmSupportedEntities = "ResourcePool"
	Datastores    PmSupportedEntities = "Datastore"
	Clusters      PmSupportedEntities = "ClusterComputeResource"
	Vapp          PmSupportedEntities = "VirtualApp"
	Datacenter    PmSupportedEntities = "Datacenter"
)

var ALL = []string{"*"}
