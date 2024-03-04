package app

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/urfave/cli"
	"golang.org/x/sys/windows"
)

func Init() *cli.App {
	app := cli.NewApp()
	app.Name = "Local.Host DNS Manager"
	app.Usage = "Change and Reset system DNS settings"

	app.Commands = []cli.Command{
		{
			Name:    "change",
			Aliases: []string{"c"},
			Usage:   "Change DNS settings",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ipv4",
					Value: "127.0.0.1",
					Usage: "IP address to set as DNS server",
				},
				cli.StringFlag{
					Name:  "ipv6",
					Value: "::1",
					Usage: "IP address to set as DNS server",
				},
			},
			Action: changeDns,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List network interfaces",
			Action:  getInterfaces,
		},
		{
			Name:    "reset",
			Aliases: []string{"r"},
			Usage:   "Reset DNS settings",
			Action:  resetDNS,
		},
	}

	return app
}

func changeDns(c *cli.Context) {
	if !checkAdmin() {
		runMeElevated()
		return
	}

	ipv4 := c.String("ipv4")
	ipv6 := c.String("ipv6")

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range interfaces {
		fmt.Println("Setting IPv4 DNS for " + i.Name + " to " + ipv4)
		command := exec.Command("netsh", "interface", "ipv4", "set", "dns", "name="+i.Name, "static", ipv4)
		var out bytes.Buffer

		command.Stdout = &out
		cmdErr := command.Run()

		if cmdErr != nil {
			fmt.Println(cmdErr, out.String())
			return
		}

		fmt.Println("Setting IPv6 DNS for " + i.Name + " to " + ipv6)
		command = exec.Command("netsh", "interface", "ipv6", "set", "dns", "name="+i.Name, "static", ipv6)

		command.Stdout = &out
		cmdErr = command.Run()

		if cmdErr != nil {
			fmt.Println(cmdErr, out.String())
			return
		}
	}
}

func getInterfaces(c *cli.Context) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range interfaces {
		fmt.Println(i.Name)
	}
}

func resetDNS(c *cli.Context) {
	if !checkAdmin() {
		runMeElevated()
		return
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range interfaces {
		fmt.Println("Resetting DNS for", i.Name)
		command := exec.Command("netsh", "interface", "ipv4", "set", "dns", "name="+i.Name, "dhcp")
		var out bytes.Buffer

		command.Stdout = &out
		cmdErr := command.Run()

		if cmdErr != nil {
			fmt.Println(cmdErr, out.String())
		}
	}
}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func checkAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
