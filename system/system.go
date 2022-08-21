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

package system

import "encoding/json"

type Capability bool

// Aurae is the current state of the system.
type Aurae struct {

	// SocketComponents are currently adopted sockets by unique name.
	// SocketComponents can provide N capabilities to Aurae.
	SocketComponents map[string]Socket

	// ServiceComponents are compiled in services part of Aurae.
	// ServiceComponents can provide N capabilities to Aurae.
	ServiceComponents map[string]Service

	// CapRunVirtualMachine enables Aurae to run as a VM hypervisor.
	// A single Socket instance is required to provide the underlying
	// support for launching virtual machines.
	//
	// Example: Firecracker socket /var/run/firecracker.socket
	CapRunVirtualMachine Socket

	CapRunProcess Service

	CapRunContainer Socket
}

type AuraeSafe struct {
	SocketComponents     map[string]string `json:"SocketComponents"`
	CapRunVirtualMachine bool              `json:"CapRunVirtualMachine"`
	CapRunContainer      bool              `json:"CapRunContainer"`
	CapRunProcess        bool              `json:"CapRunProcess"`
}

// AuraeToSafe is the conversion logic to a read-only instance of Aurae (primarily used to transmitting status).
// Here we manually map the state of the system to the output we wish to expose to users.
func AuraeToSafe(a *Aurae) *AuraeSafe {
	safe := &AuraeSafe{
		SocketComponents: make(map[string]string),
	}
	for _, c := range a.SocketComponents {
		safe.SocketComponents[c.Name()] = c.Path()
	}
	if a.CapRunVirtualMachine != nil {
		safe.CapRunVirtualMachine = true
	}
	if a.CapRunContainer != nil {
		safe.CapRunVirtualMachine = true
	}
	if a.CapRunProcess != nil {
		safe.CapRunProcess = true
	}
	return safe
}

var a *Aurae

// AuraeInstance is a singleton for the main state of the system
func AuraeInstance() *Aurae {
	if a == nil {
		a = &Aurae{
			SocketComponents: make(map[string]Socket),
		}
	}
	return a
}

func (a *Aurae) Encapsulate() ([]byte, error) {
	return json.Marshal(AuraeToSafe(a))
}

func StringToAuraeSafe(e string) (*AuraeSafe, error) {
	return BytesToAuraeSafe([]byte(e))
}

func BytesToAuraeSafe(e []byte) (*AuraeSafe, error) {
	auraeInstance := &AuraeSafe{}
	err := json.Unmarshal(e, auraeInstance)
	return auraeInstance, err
}
