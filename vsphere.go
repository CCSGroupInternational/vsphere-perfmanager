package vsphere_perfmanager

import (
	"github.com/vmware/govmomi"
	"net/url"
	"strings"
	"context"
)

type Vsphere struct {
	Username string
	Password string
	Host     string
	Insecure bool
	client   *govmomi.Client
}

func (v *Vsphere) Connect() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u, err := url.Parse(strings.Split(v.Host, "://")[0] + "://" +
		url.QueryEscape(v.Username) + ":" + url.QueryEscape(v.Password) + "@" +
		strings.Split(v.Host, "://")[1] + "/sdk")

	if err != nil {
		return err
	}

	client, err := govmomi.NewClient(ctx, u, v.Insecure)
	if err != nil {
		return err
	}
	defer client.Logout(ctx)

	v.client = client
	return nil
}