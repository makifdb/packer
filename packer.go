package packer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/JustinTimperio/osinfo"
)

type Manager struct {
	Name       string
	InstallArg string
	UpdateArg  string
	RemoveArg  string
}

var (
	apk     = Manager{"apk", "add", "update", "del"}
	apt     = Manager{"apt", "-y install", "update", "remove"}
	brew    = Manager{"brew", "install", "update", "remove"}
	dnf     = Manager{"dnf", "install", "upgrade", "erase"}
	flatpak = Manager{"flatpak", "install", "update", "uninstall"}
	snap    = Manager{"snap", "install", "upgrade", "remove"}
	pacman  = Manager{"pacman", "--noconfirm -S", "--noconfirm -Syuu", "--noconfirm -Rscn"}
	paru    = Manager{"paru", "-S", "-Syuu", "-R"}
	yay     = Manager{"yay", "-S", "-Syuu", "-R"}
	zypper  = Manager{"zypper", "-n install", "update", "-n remove"}
)

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
				return Manager{}, fmt.Errorf("no package manager found")
			}

		}
		return Manager{}, fmt.Errorf("no package manager found")
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
	mngr, err := DetectManager()
	if err != nil {
		return err
	}

	c := "sudo " + mngr.Name + " " + mngr.InstallArg + " " + packageName
	err = Command(c)
	return nil
}

func Remove(packageName string) error {
	mngr, err := DetectManager()
	if err != nil {
		return err
	}

	c := "sudo " + mngr.Name + " " + mngr.RemoveArg + " " + packageName
	err = Command(c)
	return nil
}

func Update() error {
	mngr, err := DetectManager()
	if err != nil {
		return err
	}

	c := "sudo " + mngr.Name + " " + mngr.UpdateArg
	err = Command(c)
	return nil
}

func Command(command string) error {
	args := strings.Fields(command)
	var cmd *exec.Cmd
	cmd = exec.Command(args[0], args[1:]...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		log.Printf("Command failed with %s\n", err)
	}
	return nil
}
