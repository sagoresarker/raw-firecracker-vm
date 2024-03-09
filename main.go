package main

import (
	vm "github.com/sagoresarker/raw-firecracker-vm/firecracker"
)

func main() {
	// networking.SetupNetworking()

	// // downloads the image in images/variables.UBUNTU_VERSION/download/
	// downloadgenerateimage.DownloadImage()
	// // generate the image in images/variables.UBUNTU_VERSION/
	// downloadgenerateimage.GenerateImage()
	// // create the symlink
	// downloadgenerateimage.CheckINITRD()
	// // extract the vmlinux in images/variables.UBUNTU_VERSION/
	// downloadgenerateimage.ExtractVmLinux()

	// // LaunchVM()
	vm.LaunchVM()

}
