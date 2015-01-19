package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"fmt"
	"github.com/johnnylee/ttlib"
	"log"
	"os"
)

func printUsage() {
	fmt.Println("")
	fmt.Printf("Usage: %v <directory> <user_name>\n",
		os.Args[0])
	fmt.Println("")
	fmt.Println("directory:")
	fmt.Println("    The directory in which to store the configuration.")
	fmt.Println("")
	fmt.Println("user_name:")
	fmt.Println("    The user's unique username.")
	fmt.Println("")
}

func main() {
	if len(os.Args) != 3 {
		printUsage()
		return
	}

	baseDir := os.Args[1]
	user := os.Args[2]

	// Load server configuration.
	fmt.Println("Loading server configuration...")
	serverConfig, err := ttlib.LoadServerConfig(baseDir)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	// Create and write client configuration.
	fmt.Println("Creating client configuration file...")
	cc := new(ttlib.ClientConfig)
	cc.Host = serverConfig.PublicAddr
	cc.User = user
	cc.Pwd = make([]byte, 48)

	if _, err = rand.Read(cc.Pwd); err != nil {
		log.Fatalln("Error:", err)
	}

	_, cc.CaCert, err = ttlib.LoadKeyAndCert(baseDir)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	if err = cc.Save(ttlib.ClientConfigFile(baseDir, user)); err != nil {
		log.Fatalln("Error:", err)
	}

	// Create and write the client password file.
	fmt.Println("Creating password file for client...")
	pf := new(ttlib.PwdFile)
	pf.PwdHash, err = bcrypt.GenerateFromPassword(cc.Pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	if err = pf.Save(baseDir, user); err != nil {
		log.Fatalln("Error:", err)
	}
}
