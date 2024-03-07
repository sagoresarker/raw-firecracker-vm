package downloadgenerateimage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/sagoresarker/raw-firecracker-vm/variables"
)

func download(localPath, remoteURL string) error {
	fmt.Println("Downloading", remoteURL, "to", localPath)

	// Create the file
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(remoteURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded successfully")
	return nil
}

func downloadIfNotPresent(localPath, remoteURL string) error {
	// Check if file exists locally
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// If file does not exist, download it
		return download(localPath, remoteURL)
	}
	fmt.Println("File", localPath, "already exists, skipping download.")
	return nil
}

func DownloadImage() {
	// creating directiories
	directoryPath := fmt.Sprintf("images/%s/download", variables.UBUNTU_VERSION)

	// Check if directory already exists
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			fmt.Println("Error Creating Directories:", err)
			return
		}
		fmt.Println("Directory created successfully.")
	} else if err != nil {
		// Some other error occurred
		fmt.Println("Error:", err)
		return
	} else {
		// Directory already exists
		fmt.Println("Directory already exists.")
	}
	imageTar := fmt.Sprintf("%s-server-cloudimg-amd64-root.tar.xz", variables.UBUNTU_VERSION)
	localPath := filepath.Join("images", variables.UBUNTU_VERSION, "download", imageTar)
	remoteURL := fmt.Sprintf("https://cloud-images.ubuntu.com/%s/current/%s", variables.UBUNTU_VERSION, imageTar)
	err = downloadIfNotPresent(localPath, remoteURL)
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}
	fmt.Println("Download Image Completed Successfully!")
	kernel := fmt.Sprintf("%s-server-cloudimg-amd64-vmlinuz-generic", variables.UBUNTU_VERSION)
	localPath = filepath.Join("images", variables.UBUNTU_VERSION, "download", kernel)
	remoteURL = fmt.Sprintf("https://cloud-images.ubuntu.com/%s/current/unpacked/%s", variables.UBUNTU_VERSION, kernel)
	err = downloadIfNotPresent(localPath, remoteURL)
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}
	fmt.Println("Download Kernel Completed Successfully!")
	initTrd := fmt.Sprintf("%s-server-cloudimg-amd64-initrd-generic", variables.UBUNTU_VERSION)
	localPath = filepath.Join("images", variables.UBUNTU_VERSION, "download", initTrd)
	remoteURL = fmt.Sprintf("https://cloud-images.ubuntu.com/%s/current/unpacked/%s", variables.UBUNTU_VERSION, initTrd)
	err = downloadIfNotPresent(localPath, remoteURL)
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}
	fmt.Println("Download Inittrd Completed Successfully!")

}
func GenerateImage() {
	fmt.Println("Generating", variables.IMAGE_ROOTFS)

	// Create empty image file
	imageFile, err := os.Create(variables.IMAGE_ROOTFS)
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return
	}
	defer imageFile.Close()

	// Set the size of the image file
	err = imageFile.Truncate(variables.IMAGE_SIZE)
	if err != nil {
		fmt.Println("Error setting image file size:", err)
		return
	}

	// Create filesystem
	cmd := exec.Command("mkfs.ext4", "-F", variables.IMAGE_ROOTFS)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error creating filesystem:", err)
		fmt.Println("Command output:", string(out))
		return
	}

	// Mount image
	tmpPath := fmt.Sprintf("/tmp/.%d-%d", os.Getpid(), os.Getpid())
	err = os.MkdirAll(tmpPath, 0755)
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	defer os.RemoveAll(tmpPath)

	// Use sudo specifically for the mount command
	cmd = exec.Command("sudo", "mount", variables.IMAGE_ROOTFS, tmpPath)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error mounting image:", err)
		fmt.Println("Command output:", string(out))
		return
	}
	defer syscall.Unmount(tmpPath, 0)
	// Extract contents
	imageTar := fmt.Sprintf("%s-server-cloudimg-amd64-root.tar.xz", variables.UBUNTU_VERSION)
	cmd = exec.Command("sudo", "tar", "-xf", fmt.Sprintf("images/%s/download/%s", variables.UBUNTU_VERSION, imageTar), "--directory", tmpPath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error extracting contents:", err)
		return
	}

	fmt.Println("Image generated successfully.")
}
func CheckINITRD() {
	iniTrd := fmt.Sprintf("%s-server-cloudimg-amd64-initrd-generic", variables.UBUNTU_VERSION)
	// Check if the INITRD symlink already exists
	if _, err := os.Lstat(variables.INITRD); os.IsNotExist(err) {
		// Create a symbolic link to download/$initrd as INITRD
		err := os.Symlink("download/"+iniTrd, variables.INITRD)
		if err != nil {
			fmt.Println("Error creating symlink:", err)
			return
		}
		fmt.Println("Symlink created successfully.")
	} else if err != nil {
		fmt.Println("Error checking symlink existence:", err)
		return
	} else {
		fmt.Println("Symlink already exists.")
	}
}

func ExtractVmLinux() {
	kernel := fmt.Sprintf("%s-server-cloudimg-amd64-vmlinuz-generic", variables.UBUNTU_VERSION)
	err := extract(kernel)
	if err != nil {
		fmt.Println("Error extracting vmlinux:", err)
		return
	}
	fmt.Println("vmlinux extracted successfully.")
}
func extract(kernel string) error {
	fmt.Println("Extracting vmlinux to", variables.KERNEL_IMAGE)

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "extract-vmlinux")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Download extract-vmlinux script
	scriptURL := "https://raw.githubusercontent.com/torvalds/linux/master/scripts/extract-vmlinux"
	scriptPath := filepath.Join(tmpDir, "extract-vmlinux")
	if err := downloadFile(scriptPath, scriptURL); err != nil {
		return err
	}

	// Make the script executable
	if err := os.Chmod(scriptPath, 0755); err != nil {
		return err
	}

	// Execute extract-vmlinux script
	cmd := exec.Command(scriptPath, fmt.Sprintf("images/%s/download/%s", variables.UBUNTU_VERSION, kernel))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// Write output to KERNEL_IMAGE
	outFile, err := os.Create(variables.KERNEL_IMAGE)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := outFile.Write(out); err != nil {
		return err
	}

	return nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
