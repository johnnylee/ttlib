package main

import (
	"fmt"
	"github.com/johnnylee/ttlib"
	"log"
	"os"
	"strings"
)

func printUsage() {
	fmt.Println("")
	fmt.Printf("Usage: %v <directory> <listen_addr> <public_addr>\n",
		os.Args[0])
	fmt.Println("")
	fmt.Println("directory:")
	fmt.Println("    The directory in which to store the configuration.")
	fmt.Println("")
	fmt.Println("listen_addr:")
	fmt.Println("    The address to listen on.")
	fmt.Println("")
	fmt.Println("public_addr:")
	fmt.Println("    The public address that clients will connect to.")
	fmt.Println("")
}

func main() {
	var err error
	var host string
	var sc *ttlib.ServerConfig

	if len(os.Args) != 4 {
		printUsage()
		return
	}

	// Read in command line arguments.
	baseDir := os.Args[1]
	listenAddr := os.Args[2]
	publicAddr := os.Args[3]

	// Create base directory.
	fmt.Println("Creating base directory...")
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		log.Fatalln("Error:", err)
	}

	// Create the clients directory.
	fmt.Println("Creating clients directory...")
	if err = os.MkdirAll(ttlib.ClientPwdDir(baseDir), 0700); err != nil {
		log.Fatalln("Error:", err)
	}

	// Create the client-config directory.
	fmt.Println("Creating client configuration directory...")
	if err = os.MkdirAll(ttlib.ClientConfigDir(baseDir), 0700); err != nil {
		log.Fatalln("Error:", err)
	}

	// Create and save the server configuration.
	fmt.Println("Creating and saving configuration file...")
	sc = new(ttlib.ServerConfig)
	sc.ListenAddr = listenAddr
	sc.PublicAddr = publicAddr
	if err = sc.Save(baseDir); err != nil {
		log.Fatalln("Error:", err)
	}

	// Generate key and self-signed certificate.
	host = strings.Split(publicAddr, ":")[0]
	if err = ttlib.GenerateKeyAndCert(baseDir, host); err != nil {
		log.Fatalln("Error:", err)
	}

	return
}
