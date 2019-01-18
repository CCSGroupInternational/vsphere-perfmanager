package vspherePerfManager

import (
	"github.com/vmware/govmomi/vim25/types"
	"context"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/find"
	u "github.com/ahl5esoft/golang-underscore"
)

type managedObject struct {
	Entity types.ManagedObjectReference
	Properties []types.DynamicProperty
	Metrics []Metric
}

func (v *VspherePerfManager) managedObjects(objectTypes []string) ([]types.ManagedObjectReference, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var viewManager mo.ViewManager
	err := v.client.RetrieveOne(ctx, *v.client.ServiceContent.ViewManager, nil, &viewManager)
	if err != nil {
		return nil, err
	}

	var mors []types.ManagedObjectReference

	datacenters, err := v.dataCenters()

	if err != nil {
		return nil, err
	}

	for _, datacenter := range datacenters {
		req := types.CreateContainerView{
			This: viewManager.Reference(),
			Container: datacenter,
			Type: objectTypes,
			Recursive: true,
		}

		res, err := methods.CreateContainerView(ctx, v.client.RoundTripper, &req)

		if err != nil {
			return nil, err
		}

		var containerView mo.ContainerView
		err = v.client.RetrieveOne(ctx, res.Returnval, nil, &containerView)
		if err != nil {
			return nil, err
		}
		mors = append(mors, containerView.View...)
	}

	return mors, nil

}

func (v *VspherePerfManager) getManagedObject(mors []types.ManagedObjectReference, propSets []types.PropertySpec) ([]managedObject, error) {
	var objectSet []types.ObjectSpec

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, mor := range mors {
		objectSet = append(objectSet, types.ObjectSpec{Obj: mor, Skip: types.NewBool(false)})
	}

	//retrieve properties
	propReq := types.RetrieveProperties{SpecSet: []types.PropertyFilterSpec{{ObjectSet: objectSet, PropSet: propSets}}}
	propRes, err := v.client.PropertyCollector().RetrieveProperties(ctx, propReq)

	if err != nil {
		return nil, err
	}

	var managedObjects []managedObject

	for _, objectContent := range propRes.Returnval {
		managedObjects = append(managedObjects, managedObject{
			Entity: objectContent.Obj,
			Properties: objectContent.PropSet,
		})
	}
	return managedObjects, nil
}

func (v *VspherePerfManager) dataCenters() ([]types.ManagedObjectReference, error) {

	var dataCenters []types.ManagedObjectReference

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finder := find.NewFinder(v.client.Client, true)

	dcs, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		return nil, err
	}

	for _, child := range dcs {
		dataCenters = append(dataCenters, child.Reference())
	}

	return dataCenters, nil

}

func (m *managedObject) GetProperty(property string) types.AnyType {
	props := u.Where(m.Properties, func(prop types.DynamicProperty, i int) bool {
		if prop.Name == property {
			return true
		}
		return false
	})

	if props == nil {
		return nil
	}

	return props.([]types.DynamicProperty)[0].Val
}

func getProperties(propertiesFromconfig []types.PropertySpec) []types.PropertySpec {

	properties := []types.PropertySpec{{
		Type   : "ManagedEntity",
		PathSet : []string{"name"},
	}}
	return append(properties, propertiesFromconfig...)
}