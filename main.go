package main

import (
	"devman/internal/commands"
	"devman/internal/output"
	"flag"
	"os"
)

const version = "0.1.0"

func main() {
	output.Info("DevMan " + version)

	var help bool
	const helpUsage = "Print this help message."

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage + " short-hand")

	var newSiteInput string
	const newSiteInputUsage = "Enter the domain"

	flag.StringVar(&newSiteInput, "new", "", newSiteInputUsage)
	flag.StringVar(&newSiteInput, "n", "", newSiteInputUsage + " short-hand")

	flag.Parse()

	if help == true {
		flag.Usage()
		os.Exit(0)
	}

	if newSiteInput != "" {
		commands.NewSite(newSiteInput)
	}
}
