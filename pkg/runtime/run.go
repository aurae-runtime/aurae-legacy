/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package runtime

import (
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/container"
	"github.com/weaveworks/libgitops/pkg/runtime"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RunContainer(name string) error {
	//return nil
	// func ExecuteFirecracker(vm *api.VM, fcIfaces firecracker.NetworkInterfaces) (err error) {

	//client := firecracker.NewClient("/run/firecracker.socket", logrus.NewEntry(logrus.New()), true)
	vm := &ignite.VM{}
	vm.SetName(name)
	vm.SetImage(&ignite.Image{
		TypeMeta: runtime.TypeMeta{
			TypeMeta: v1.TypeMeta{
				Kind:       "",
				APIVersion: "",
			},
		},
		ObjectMeta: runtime.ObjectMeta{
			Name:        name,
			Labels:      nil,
			Annotations: nil,
		},
		Spec: ignite.ImageSpec{},
	})
	netinterfaces := []firecracker.NetworkInterface{
		{
			StaticConfiguration: &firecracker.StaticNetworkConfiguration{
				MacAddress:  "",
				HostDevName: "aurae",
				//IPConfiguration: &firecracker.IPConfiguration{
				//	IPAddr: net.IPNet{
				//		IP:   net.IP{},
				//		Mask: nil,
				//	},
				//	Gateway:     nil,
				//	Nameservers: nil,
				//	IfName:      "",
				//},
			},
			AllowMMDS: false,
		},
	}
	return container.ExecuteFirecracker(vm, netinterfaces)
}
