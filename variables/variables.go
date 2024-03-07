package variables

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	NUMBER_VMS         = 3
	VM_VCPUS           = 2
	VM_RAM_GB          = 4
	IMAGE_SIZE         = "8G"
	FIRECRACKER_BRIDGE = "fcbr0"
	VMS_NETWORK_PREFIX = "172.26.0"
	UBUNTU_VERSION     = "bionic"
	IMAGE_ROOTFS       = "images/" + UBUNTU_VERSION + "/" + UBUNTU_VERSION + ".rootfs"
	KERNEL_IMAGE       = "images/" + UBUNTU_VERSION + "/" + UBUNTU_VERSION + ".vmlinux"
	INITRD             = "images/" + UBUNTU_VERSION + "/" + UBUNTU_VERSION + ".initrd"
	KEYPAIR_DIR        = "keypairs"
	DEFAULT_KP         = "kp"
)

// GetEgressInterface returns the EGRESS_IFACE variable
func GetEgressInterface() string {
	// Get the EGRESS_IFACE using shell command
	output, err := exec.Command("sh", "-c", "ip route get 8.8.8.8 | grep uid | sed 's/.* dev \\([^ ]*\\) .*/\\1/'").Output()
	if err != nil {
		fmt.Println("Error getting EGRESS_IFACE:", err)
		return ""
	}

	EGRESS_IFACE := strings.TrimSpace(string(output))

	return EGRESS_IFACE
}
