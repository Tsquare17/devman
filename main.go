package main

import (
	"flag"
	"github.com/tsquare17/devman/internal/commands"
	"github.com/tsquare17/devman/internal/output"
	"github.com/tsquare17/devman/internal/setup"
	"os"
)

const version = "0.1.0"

func main() {
	output.Info("DevMan " + version)

	setup.Init()

	var help bool
	const helpUsage = "Print this help message."

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage + " short-hand")

	var versionInput bool
	const versionUsage = "Show the version."
	flag.BoolVar(&versionInput, "version", false, versionUsage)
	flag.BoolVar(&versionInput, "v", false, versionUsage)

	var newSiteInput string
	const newSiteInputUsage = "Enter the domain of the site to be created."

	flag.StringVar(&newSiteInput, "new", "", newSiteInputUsage)
	flag.StringVar(&newSiteInput, "n", "", newSiteInputUsage + " short-hand")

	var removeSiteInput string
	const removeSiteUsage = "Enter the domain of the site to be removed."
	flag.StringVar(&removeSiteInput, "remove", "", removeSiteUsage)
	flag.StringVar(&removeSiteInput, "rm", "", removeSiteUsage + " short-hand")

	var gitHooksInput string
	var gitHooksUsage = "Command"
	flag.StringVar(&gitHooksInput, "git", "", gitHooksUsage)

	flag.Parse()

	if help == true {
		flag.Usage()
		os.Exit(0)
	}

	executed := false
	if newSiteInput != "" {
		commands.NewSite(newSiteInput)
		executed = true
	}

	if removeSiteInput != "" {
		commands.RemoveSite(removeSiteInput)
		executed = true
	}

	if gitHooksInput != "" {
		if gitHooksInput == "gitdeny" {
			commands.AddGitDenyTag()
			executed = true
		}
	}

	if executed == false {
		flag.Usage()
	}
}
