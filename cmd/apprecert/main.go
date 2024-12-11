package main

import (
	"flag"
	"log"

	"github.com/appremon/apprecert/cert"
	"github.com/appremon/apprecert/config"
	"github.com/appremon/apprecert/truststore"
)

func main() {
	// Define flags
	installFlag := flag.Bool("install", false, "Install the CA")
	uninstallFlag := flag.Bool("uninstall", false, "Uninstall the CA")
	helpFlag := flag.Bool("help", false, "Show usage information")

	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	// Initialize configuration
	cfg := config.Load()

	if *installFlag {
		if err := truststore.Install(cfg); err != nil {
			log.Fatalf("Failed to install CA: %v", err)
		}
		log.Println("CA installed successfully!")
		return
	}

	if *uninstallFlag {
		if err := truststore.Uninstall(cfg); err != nil {
			log.Fatalf("Failed to uninstall CA: %v", err)
		}
		log.Println("CA uninstalled successfully!")
		return
	}

	// Default action: generate certificate
	if len(flag.Args()) > 0 {
		if err := cert.Generate(cfg, flag.Args()); err != nil {
			log.Fatalf("Failed to generate certificate: %v", err)
		}
		log.Println("Certificate generated successfully!")
	} else {
		log.Println("No arguments provided. Use -help for usage information.")
	}
}

func printHelp() {
	log.Println("Usage of apprecert:")
	log.Println("  -install: Install the local CA.")
	log.Println("  -uninstall: Uninstall the local CA.")
	log.Println("  -help: Display usage information.")
}
