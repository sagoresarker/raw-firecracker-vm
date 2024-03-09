package instance

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

func Int64(v int64) *int64 {
	return &v
}

func printSeconds() {
	var secondsPassed int

	for {
		fmt.Printf("%d seconds\n", secondsPassed)
		secondsPassed++
		time.Sleep(1 * time.Second)
	}
}

func LaunchVM() {
	// Define the paths to the kernel image and the root drive image.
	kernelImagePath := "/run/media/sagoresarker/Study/SWE-Poridhi/kubernetes/cloned-repo/blog-posts-src/202012-firecracker_cloud_image_automation/images/bionic/bionic.vmlinux"
	rootDriveImagePath := "/run/media/sagoresarker/Study/SWE-Poridhi/kubernetes/cloned-repo/blog-posts-src/202012-firecracker_cloud_image_automation/images/bionic/bionic.rootfs"

	drive := models.Drive{
		DriveID:      firecracker.String("1"),
		PathOnHost:   &rootDriveImagePath,
		IsRootDevice: firecracker.Bool(true),
		IsReadOnly:   firecracker.Bool(false),
	}

	cfg := firecracker.Config{
		SocketPath:      "/tmp/firecracker.sock",
		KernelImagePath: kernelImagePath,
		InitrdPath:      "/run/media/sagoresarker/Study/SWE-Poridhi/kubernetes/cloned-repo/blog-posts-src/202012-firecracker_cloud_image_automation/images/bionic/bionic.initrd",
		// KernelArgs:      "console=ttyS0 reboot=k panic=1 pci=off nomodset=1 systemd.journald.forward_to_console systemd.log_level=err log_buf_len=1M random.trust_cpu=on random.random.trust_cpu=on",
		KernelArgs: "console=ttyS0 reboot=k panic=1 pci=off",
		Drives:     []models.Drive{drive},
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(1),
			MemSizeMib: firecracker.Int64(512),
		},
	}

	ctx := context.Background()
	cmd := firecracker.VMCommandBuilder{}.
		WithBin("firecracker").
		WithSocketPath(cfg.SocketPath).
		Build(ctx)

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		log.Fatalf("Failed to create machine: %s", err)
	}

	if err := m.Start(ctx); err != nil {
		log.Fatalf("Failed to start machine: %s", err)
	}

	go printSeconds()

	// Wait for the VM to start up.
	time.Sleep(100 * time.Second)

	if err := m.StopVMM(); err != nil {
		log.Fatalf("Failed to stop VM: %s", err)
	}

	// Wait for the VM to finish.
	if err := m.Wait(ctx); err != nil {
		log.Fatalf("Wait returned an error %s", err)
	}
}
