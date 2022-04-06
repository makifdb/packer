package packer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/JustinTimperio/osinfo"
)

var (
	apk     = Manager{Name: "apk", InstallArg: []string{"add"}, UpdateArg: []string{"update"}, RemoveArg: []string{"del"}}
	apt     = Manager{Name: "apt", InstallArg: []string{"-y", "install"}, UpdateArg: []string{"update"}, RemoveArg: []string{"remove"}}
	brew    = Manager{Name: "brew", InstallArg: []string{"install"}, UpdateArg: []string{"update"}, RemoveArg: []string{"remove"}}
	dnf     = Manager{Name: "dnf", InstallArg: []string{"install"}, UpdateArg: []string{"upgrade"}, RemoveArg: []string{"erase"}}
	flatpak = Manager{Name: "flatpak", InstallArg: []string{"install"}, UpdateArg: []string{"update"}, RemoveArg: []string{"uninstall"}}
	snap    = Manager{Name: "snap", InstallArg: []string{"install"}, UpdateArg: []string{"upgrade"}, RemoveArg: []string{"remove"}}
	pacman  = Manager{Name: "pacman", InstallArg: []string{"-S", "--noconfirm"}, UpdateArg: []string{"--noconfirm", "-Syuu"}, RemoveArg: []string{"--noconfirm", "-Rscn"}}
	paru    = Manager{Name: "paru", InstallArg: []string{"-S"}, UpdateArg: []string{"-Syuu"}, RemoveArg: []string{"-R"}}
	yay     = Manager{Name: "yay", InstallArg: []string{"-S"}, UpdateArg: []string{"-Syuu"}, RemoveArg: []string{"-R"}}
	zypper  = Manager{Name: "zypper", InstallArg: []string{"-n", "install"}, UpdateArg: []string{"update"}, RemoveArg: []string{"-n", "remove"}}
)

type Manager struct {
	Name       string
	InstallArg []string
	UpdateArg  []string
	RemoveArg  []string
}

func DetectManager() (Manager, error) {
	switch opsystem := osinfo.GetVersion().Runtime; opsystem {
	default:
		// windows, freebsd, plan9 ...
		return Manager{}, fmt.Errorf("%s is not supported", opsystem)
	case "darwin":
		return brew, nil
	case "linux":
		switch distro := osinfo.GetVersion().Linux.Distro; distro {

		case "arch":
			if Check("pacman") {
				return pacman, nil
			} else if Check("yay") {
				return yay, nil
			} else if Check("paru") {
				return paru, nil
			}

		case "alpine":
			if Check("apk") {
				return apk, nil
			}

		case "fedora":
			if Check("dnf") {
				return dnf, nil
			}

		case "opensuse":
			if Check("zypper") {
				return zypper, nil
			}

		case "debian":
			if Check("apt") {
				return apt, nil
			} else if Check("snap") {
				return snap, nil
			}

		default:
			if Check("snap") {
				return snap, nil
			} else if Check("flatpak") {
				return flatpak, nil
			} else {
				return Manager{}, fmt.Errorf("No package manager found")
			}

		}
		return Manager{}, fmt.Errorf("No package manager found")
	}
}

func Check(packageName string) bool {
	_, err := exec.LookPath(packageName)
	if err != nil {
		return false
	} else {
		return true
	}
}

func Install(packageName string) error {
	pkgManager, err := DetectManager()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch ln := cap(pkgManager.InstallArg); ln {
	default:
		cmd = exec.Command("sudo", pkgManager.Name, pkgManager.InstallArg[0], packageName)
	case 2:
		cmd = exec.Command("sudo", pkgManager.Name, pkgManager.InstallArg[0], pkgManager.InstallArg[1], packageName)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Printf("Install failed with %s\n", err)
	}
	return err
}

func Remove(packageName string) error {
	pkgManager, err := DetectManager()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch ln := cap(pkgManager.InstallArg); ln {
	default:
		cmd = exec.Command("sudo", pkgManager.Name, pkgManager.RemoveArg[0], packageName)
	case 2:
		cmd = exec.Command("sudo", pkgManager.Name, pkgManager.RemoveArg[0], pkgManager.RemoveArg[1], packageName)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Printf("Remove failed with %s\n", err)
	}

	return err
}

func Update() error {
	pkgManager, err := DetectManager()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch ln := cap(pkgManager.InstallArg); ln {
	default:
		cmd = exec.Command("sudo", pkgManager.Name, pkgManager.UpdateArg[0])
	case 2:
		cmd = exec.Command("sudo", pkgManager.Name, pkgManager.UpdateArg[0], pkgManager.UpdateArg[1])
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Printf("Update failed with %s\n", err)
	}

	return err
}
