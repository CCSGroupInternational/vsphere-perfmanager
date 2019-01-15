package vsphere_perfmanager

import (
	"testing"
	"fmt"
	"os"
)

func TestConnect(t *testing.T) {
	fmt.Println(os.Getenv("VSPHERE_HOST"))
	fmt.Println(os.Getenv("VSPHERE_USER"))
	fmt.Println(os.Getenv("VSPHERE_PASSWORD"))
	fmt.Println(os.Getenv("VSPHERE_INSECURE"))
}
