package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	red     = color.New(color.FgRed).Add(color.Bold)
	blue    = color.New(color.FgBlue).Add(color.Bold)
	magenta = color.New(color.FgHiMagenta)
)

func usage() {
	// Top
	fmt.Printf("\r\n")

	// PRIMARY COMMANDS
	yellow.Printf("	Primary Commands:\r\n\r\n")

	// find Command
	fmt.Print("	$ ")
	cyan.Print("find ")
	magenta.Print("<directory/path> <file>")
	fmt.Printf("  ->  Recursively searches target directory for file\n")
	fmt.Printf("\r\n")

	// dns Command
	fmt.Print("	$ ")
	cyan.Print("dns")
	fmt.Printf("  ->  Checks outbound DNS using internal and Cloudflare DNS Servers\n")
	fmt.Printf("\r\n")

	// goget Command
	fmt.Print("	$ ")
	cyan.Print("goget ")
	magenta.Print("https://hostname.com/path/to/file.ext newfile.ext")
	fmt.Printf(" -> Downloads remote file\n")
	fmt.Printf("\r\n")

	// OS COMMANDS
	yellow.Printf("	OS Commands:\n\n")

	// exit Command
	fmt.Print("	$ ")
	cyan.Print("exit")
	fmt.Printf("	-> Exits the program\n")
	fmt.Printf("\r\n")

	// cd Command
	fmt.Print("	$ ")
	cyan.Print("cd ")
	magenta.Print("<directory/path>")
	fmt.Printf("  ->  Changes directory to target directory/path\n")
	fmt.Printf("\r\n")

	// pwd Command
	fmt.Print("	$ ")
	cyan.Print("pwd")
	fmt.Printf("  ->  Prints working directory\n")
	fmt.Printf("\r\n")

	// ls Command
	fmt.Print("	$ ")
	cyan.Print("ls")
	fmt.Print("  ->  List all files in current directory\n")
	fmt.Printf("\r\n")

	// clear Command
	fmt.Print("	$ ")
	cyan.Print("clear")
	fmt.Printf("  ->  Clears the screen\n")
	fmt.Printf("\r\n")

	// id Command
	fmt.Print("	$ ")
	cyan.Print("id")
	fmt.Printf("  ->  Prints the hostname and current user\n")
	fmt.Printf("\r\n")

	// ip Command
	fmt.Print("	$ ")
	cyan.Print("ip")
	fmt.Printf("  ->  Prints the Network Interfaces\n")
	fmt.Printf("\r\n")

	// cmd Command
	fmt.Print("	$ ")
	cyan.Print("cmd ")
	magenta.Print("<Command> <Subcommand>")
	fmt.Printf("  ->  Allows you to run OS commands\n")
	fmt.Printf("\r\n")

	// FILE COMMANDS
	yellow.Printf("	File Commands:\n\n")

	// fileinfo Command
	fmt.Print("	$ ")
	cyan.Print("fileinfo ")
	magenta.Print("<file>")
	fmt.Printf("  ->  Prints information about a file\n")
	fmt.Printf("\r\n")

	// mkfile Command
	fmt.Print("	$ ")
	cyan.Print("mkfile ")
	magenta.Print("<file>")
	fmt.Printf("  ->  Creates empty file\n")
	fmt.Printf("\r\n")

	// apfile Command
	fmt.Print("	$ ")
	cyan.Print("apfile ")
	magenta.Print("<content> <file>")
	fmt.Printf("  ->  Add or append content at the end of text file\n")
	fmt.Printf("\r\n")

	// delfile Command
	fmt.Print("	$ ")
	cyan.Print("delfile ")
	magenta.Print("<file>")
	fmt.Printf("  ->  Deletes file\n")
	fmt.Printf("\r\n")

	// cp Command
	fmt.Print("	$ ")
	cyan.Print("cp ")
	magenta.Print("<source> <new>")
	fmt.Printf("  ->  Copies a file\n")
	fmt.Printf("\r\n")

	// mkdir Command
	fmt.Print("	$ ")
	cyan.Print("mkdir ")
	magenta.Print("<directory>")
	fmt.Printf("  ->  Creates directory\n")
	fmt.Printf("\r\n")

	// deldir Command
	fmt.Print("	$ ")
	cyan.Print("deldir ")
	magenta.Print("<directory>")
	fmt.Printf("  ->  Deletes directory\n")
	fmt.Printf("\r\n")
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func homeDir() {
	usrDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	os.Chdir(usrDir)
}

func command(cmdStr string) error {
	cmdStr = strings.TrimSuffix(cmdStr, "\n")
	argStr := strings.Fields(cmdStr)

	if len(argStr) == 0 {
		fmt.Printf("\r")
		return nil
	}

	if argStr[0] == "help" || argStr[0] == "h" {
		usage()
		return nil
	}

	if argStr[0] == "exit" {
		clearScreen()
		os.Exit(0)
	}

	if argStr[0] == "find" {
		regexes := []*regexp.Regexp{
			regexp.MustCompile(`(?i)` + argStr[1]),
		}

		//loading := spinner.New(spinner.CharSets[69], 100*time.Millisecond)

		fmt.Print("\r\n")
		yellow.Print("File(s) Found:\r\n")
		if err := filepath.Walk("/", func(path string, f os.FileInfo, err error) error {
			for _, r := range regexes {
				if r.MatchString(path) {
					fmt.Printf("%s", filepath.Dir(path))
					fmt.Print("\\")
					fmt.Printf(filepath.Base(path))
					fmt.Printf("\n")
				}
			}
			return nil
		}); err != nil {
			log.Fatalln(err)
		}
		fmt.Print("\r\n")
		yellow.Print("Done!\n")
	}
	if argStr[0] == "cd" {
		chDir := os.Chdir(argStr[1])
		cuDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if chDir == nil {
			green.Print(cuDir)
			fmt.Printf("\r\n")
			return chDir
		}
		fmt.Printf("Not a directory!\r\n")
	}
	if argStr[0] == "pwd" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		green.Print(dir)
		fmt.Printf("\r\n")
	}
	if argStr[0] == "ls" || argStr[0] == "dir" {
		list, err := ioutil.ReadDir(".")
		if err != nil {
			panic(err)
		}

		for _, file := range list {
			green.Println(file.Name())
		}
		fmt.Printf("\r")
		return nil
	}
	if argStr[0] == "id" || argStr[0] == "whoami" {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Print("\r\n")
		green.Print(usr.Username)
		fmt.Printf("\r\n\r\n")
		return nil
	}
	if argStr[0] == "ip" {
		ifaces, err := net.Interfaces()
		if err != nil {
			fmt.Print(err)
			return nil
		}
		fmt.Print("\n")
		yellow.Printf("Network Interfaces:\n")
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				fmt.Print(err)
				continue
			}
			for _, a := range addrs {
				switch v := a.(type) {
				case *net.IPAddr:
					green.Printf("%v: %s (%s)\n", i.Name, v, v.IP.DefaultMask())

				case *net.IPNet:
					green.Printf("%v: %s [%v/%v]\n", i.Name, v, v.IP, v.Mask)
				}
			}
		}
		fmt.Print("\n")
	}

	// File-Folder Commands
	if argStr[0] == "fileinfo" {
		fileInfo, err := os.Stat(argStr[1])
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("\r\n")
		fmt.Println("File Name: ", fileInfo.Name())
		fmt.Println("File Size: ", fileInfo.Size())
		fmt.Println("Permissions: ", fileInfo.Mode())
		fmt.Println("Last Modified: ", fileInfo.ModTime())
	}

	if argStr[0] == "mkfile" {
		mkFile, err := os.Create(argStr[1])
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(mkFile)
		mkFile.Close()
	}

	if argStr[0] == "apfile" {
		content := argStr[1]
		fileName := argStr[2]

		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		defer f.Close()
		fmt.Fprintf(f, "%s\n", content)
	}

	if argStr[0] == "delfile" {
		dFile := os.Remove(argStr[1])
		if dFile != nil {
			log.Fatalln(dFile)
		}
	}

	if argStr[0] == "cp" {
		sourceFile, err := os.Open(argStr[1])
		if err != nil {
			log.Fatalln(err)
		}
		defer sourceFile.Close()

		// Create New File
		newFile, err := os.Create(argStr[2])
		if err != nil {
			log.Fatalln(err)
		}
		defer newFile.Close()

		bytesCopied, err := io.Copy(newFile, sourceFile)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Copied %d bytes.", bytesCopied)
	}

	if argStr[0] == "mkdir" {
		_, err := os.Stat(argStr[1])
		if os.IsNotExist(err) {
			mDir := os.MkdirAll(argStr[1], 0755)
			if mDir != nil {
				log.Fatalln(mDir)
			}
		}
	}

	if argStr[0] == "deldir" {
		dFolder := os.RemoveAll(argStr[1])
		if dFolder != nil {
			log.Fatalln(dFolder)
		}
	}

	if argStr[0] == "dns" {
		_, err := net.Dial("tcp", argStr[1]+":443")
		if err == nil {
			blue.Println("Internal DNS Successful!")
		} else {
			red.Print("Internal DNS Error!")
		}
		fmt.Printf("\r\n")

		var (
			dnsResolverIP        = "1.1.1.1:53"
			dnsResolverProto     = "udp"
			dnsResolverTimeoutMS = 5000
		)

		dialer := &net.Dialer{
			Resolver: &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						Timeout: time.Duration(dnsResolverTimeoutMS) * time.Millisecond,
					}
					return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
				},
			},
		}

		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}

		http.DefaultTransport.(*http.Transport).DialContext = dialContext
		httpClient := &http.Client{}

		resp, err := httpClient.Get("https://" + argStr[1])
		if err != nil {
			log.Fatalln(err)
			red.Print("Cloudflare DNS Error!")
		}
		blue.Print("Cloudflare DNS Successful!")
		resp.Body.Close()
		fmt.Printf("\r\n")
		return nil
	}

	if argStr[0] == "goget" {
		err := DownloadFile(argStr[2], argStr[1])
		if err != nil {
			panic(err)
		}

		yellow.Printf("Download Finished!\n")
	}

	if argStr[0] == "clear" || argStr[0] == "cls" {
		cmd := exec.Command("clear")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
		return nil
	}
	if argStr[0] == "cmd" {
		cmd := exec.Command(argStr[1], argStr[2:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
		return nil
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	homeDir()
	clearScreen()

	yellow.Print("Type help or h to see a list of commands\n\n")

	for {
		fmt.Print("GoMechanic~$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		command(cmdString)
	}
}
