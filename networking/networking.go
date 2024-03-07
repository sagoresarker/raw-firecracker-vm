package networking

import (
	"fmt"
	"os/exec"

	"github.com/sagoresarker/raw-firecracker-vm/variables"
	"github.com/vishvananda/netlink"
)

func SetupNetworking() {

	EGRESS_IFACE := variables.GetEgressInterface()

	// Create bridge interface
	bridge := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: variables.FIRECRACKER_BRIDGE,
		},
	}
	err := netlink.LinkAdd(bridge)
	if err != nil {
		fmt.Println("Error creating bridge:", err)
		return
	}

	// Set IP address for the bridge interface
	bridgeAddr, err := netlink.ParseAddr(variables.VMS_NETWORK_PREFIX + ".1/24")
	if err != nil {
		fmt.Println("Error parsing bridge IP address:", err)
		return
	}
	err = netlink.AddrAdd(bridge, bridgeAddr)
	if err != nil {
		fmt.Println("Error adding address to bridge:", err)
		return
	}

	// Bring up the bridge interface
	err = netlink.LinkSetUp(bridge)
	if err != nil {
		fmt.Println("Error bringing up bridge interface:", err)
		return
	}

	// Enable IP forwarding
	err = exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1").Run()
	if err != nil {
		fmt.Println("Error enabling IP forwarding:", err)
		return
	}

	// Setup iptables rules for NAT
	err = exec.Command("iptables", "--table", "nat", "--append", "POSTROUTING", "--out-interface", EGRESS_IFACE, "-j", "MASQUERADE").Run()
	if err != nil {
		fmt.Println("Error setting up iptables NAT rule:", err)
		return
	}

	err = exec.Command("iptables", "--insert", "FORWARD", "--in-interface", variables.FIRECRACKER_BRIDGE, "-j", "ACCEPT").Run()
	if err != nil {
		fmt.Println("Error setting up iptables FORWARD rule:", err)
		return
	}

	fmt.Println("Networking setup completed successfully.")
}
