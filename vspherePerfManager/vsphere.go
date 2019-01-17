package vspherePerfManager

import (
	"github.com/vmware/govmomi"
	"net/url"
	"strings"
	"context"
	"github.com/CCSGroupInternational/vsphere-perfmanager/config"
	"github.com/vmware/govmomi/vim25/types"
)

type VspherePerfManager struct {
	client       *govmomi.Client
	metricsInfo  []metricInfo
	config       *config.VspherePerfManagerConfig
}

func Init(c *config.VspherePerfManagerConfig) (*VspherePerfManager, error) {
	vspherePerfManager := VspherePerfManager{}
	err := vspherePerfManager.connect(c.Vcenter)
	vspherePerfManager.metricsInfo, err = vspherePerfManager.getMetricsInfo()
	vspherePerfManager.config = c
	if err != nil {
		return nil, err
	}
	return &vspherePerfManager, err
}

func (v *VspherePerfManager) connect(c config.Vcenter) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u, err := url.Parse(strings.Split(c.Host, "://")[0] + "://" +
		url.QueryEscape(c.Username) + ":" + url.QueryEscape(c.Password) + "@" +
		strings.Split(c.Host, "://")[1] + "/sdk")

	if err != nil {
		return err
	}

	client, err := govmomi.NewClient(ctx, u, c.Insecure)
	if err != nil {
		return err
	}

	v.client = client
	return nil
}

func (v *VspherePerfManager) Vms() ([]managedObject, error) {
	objects, err := v.managedObjects([]string{"VirtualMachine"})
	if err != nil {
		return nil, err
	}

	properties := []types.PropertySpec{{
		Type   : "ManagedEntity",
		PathSet : []string{"name"},
	}}

	vms, err := v.getManagedObject(objects, properties)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for key := range vms {
		err := v.query(&vms[key])

		if err != nil {
			return nil, err
		}
	}
	return vms, nil
}

func (v *VspherePerfManager) Hosts()  ([]managedObject, error) {

	objects, err := v.managedObjects([]string{"HostSystem"})
	if err != nil {
		return nil, err
	}

	properties := []types.PropertySpec{{
		Type   : "ManagedEntity",
		PathSet : []string{"name"},
	}}

	hosts, err := v.getManagedObject(objects, properties)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for key := range hosts {
		err := v.query(&hosts[key])

		if err != nil {
			return nil, err
		}
	}
	return hosts, nil
}