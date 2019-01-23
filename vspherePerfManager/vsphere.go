package vspherePerfManager

import (
	"github.com/vmware/govmomi"
	"net/url"
	"strings"
	"context"
	"github.com/CCSGroupInternational/vsphere-perfmanager/config"
)

type VspherePerfManager struct {
	client       *govmomi.Client
	metricsInfo  []metricInfo
	config       *config.VspherePerfManagerConfig
	objects      map[config.PmSupportedEntities]map[string]ManagedObject
}

func Init(c *config.VspherePerfManagerConfig) (*VspherePerfManager, error) {
	vspherePerfManager := VspherePerfManager{}
	err := vspherePerfManager.connect(c.Vcenter)
	if err != nil {
		return nil, err
	}
	vspherePerfManager.config = c
	vspherePerfManager.metricsInfo, err = vspherePerfManager.getMetricsInfo()
	if err != nil {
		return nil, err
	}

	err = vspherePerfManager.managedObjects()
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

func (v *VspherePerfManager) Get(entityType config.PmSupportedEntities) (map[string]ManagedObject, error) {
	return v.getMetrics(entityType)
}

func (v *VspherePerfManager) getMetrics(ObjectType config.PmSupportedEntities) (map[string]ManagedObject, error) {
	var err error

	entities := v.objects[ObjectType]

	for key, vm := range entities {
		entities[key], err = v.query(vm)
		//
		if err != nil {
			return nil, err
		}
	}
	return entities, nil
}