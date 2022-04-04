# packer	

Packer is a simple package management tool for Go

<img align="right" width="328" alt="icon" src="https://user-images.githubusercontent.com/31243845/161521421-ca0328fd-9395-47c2-8f0d-b348a89c09db.png">


| Operation Systems   | Package Managers |
|---------------------|------------------|
| Ubuntu              | apk 		 	 | 
| Debian              | apt 		     | 
| MXLinux             | brew         	 | 
| Mint                | dnf        		 | 
| Kali                | flatpak          | 
| ParrotOS            | snap             |
| OpenSUSE Leap       | pacman           | 
| OpenSUSE TumbleWeed | paru             |
| OpenSUSE SLES       | yay              |
| Arch                | zypper           |
| Manjaro             |                  |
| Alpine              |                  |
| Fedora              |                  |
| RHEL                |                  |
| CentOS              |                  |
| Oracle              |                  |
| MacOS               |                  |



## Example Usage

1. Check package installation

```go
func main() {
	p:= packer.Check("curl")
	fmt.Println(p)
}
// output: true
```

2. Install package

```go
func main() {
	packer.Install("curl")
}
```

3. Remove package

```go
func main() {
	packer.Remove("curl")
}
```


4. Update system

```go
func main() {	
	packer.Update()
}
```

5. Dedect Package Manager

```go
func main() {	
	mngr, _ := DedectManager()
	fmt.Println(mngr.Name)
}
// output: yay

```
