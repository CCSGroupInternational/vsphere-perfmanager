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
	Metrics   []string
	Instances []string
	Entities  []string
}

type PmSupportedEntities string

const (
	VMs               PmSupportedEntities = "VirtualMachine"
	Hosts             PmSupportedEntities = "HostSystem"
	ResourcePools     PmSupportedEntities = "ResourcePool"
	Datastores        PmSupportedEntities = "Datastore"
	Clusters          PmSupportedEntities = "ClusterComputeResource"
	Vapps             PmSupportedEntities = "VirtualApp"
	Datacenters       PmSupportedEntities = "Datacenter"
	Folders                               = "Folder"
	DatastoreClusters                     = "StoragePod"
)

var ALL = []string{"*"}
