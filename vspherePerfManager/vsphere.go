package vspherePerfManager

import (
	"github.com/vmware/govmomi"
	"net/url"
	"strings"
	"context"
)

type VspherePerfManager struct {
	Config       Config
	client       *govmomi.Client
	metricsInfo  []metricInfo
	objects      map[string]map[string]ManagedObject
}

func (v *VspherePerfManager) Init() (error) {
	err := v.connect(v.Config.Vcenter)
	if err != nil {
		return err
	}
	v.metricsInfo, err = v.getMetricsInfo()
	if err != nil {
		return err
	}

	return v.managedObjects()
}

func (v *VspherePerfManager) connect(c Vcenter) error {
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

func (v *VspherePerfManager) Get(entityType PmSupportedEntities) (map[string]ManagedObject, error) {
	return v.fetchMetrics(string(entityType))
}

func (v *VspherePerfManager) fetchMetrics(ObjectType string) (map[string]ManagedObject, error) {
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