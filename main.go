package main

import (
	networking "github.com/sagoresarker/raw-firecracker-vm/networking"

	downloadgenerateimage "github.com/sagoresarker/raw-firecracker-vm/download_generate_image"
)

func main() {
	networking.SetupNetworking()

	// downloads the image in images/variables.UBUNTU_VERSION/download/
	downloadgenerateimage.DownloadImage()
	// generate the image in images/variables.UBUNTU_VERSION/
	downloadgenerateimage.GenerateImage()
	// create the symlink
	downloadgenerateimage.CheckINITRD()
	// extract the vmlinux in images/variables.UBUNTU_VERSION/
	downloadgenerateimage.ExtractVmLinux()

}
