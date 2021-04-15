package prompt

import (
	"github.com/manifoldco/promptui"
	"github.com/tsquare17/devman/internal/utils"
)

func SitePath() string {
	result := getPrompt("Enter site path (if different from current path)")

	if result != "" {
		return result
	}

	return utils.GetPwd()
}

func SiteDocumentRoot() string {
	return getPrompt("Enter document root (relative to the site path)")
}

func MySQLPasswordPrompt() string {
	return getPrompt("Enter root MySQL password (if creating mysql)")
}

func IsInstallWordPress() bool {
	result := getPrompt("Install WordPress? (Y/n) [n]")

	if result == "y" || result == "Y" {
		return true
	}

	return false
}

func DatabaseName(defaultName string) string {
	return getPrompt("Enter mysql name [" + defaultName + "]")
}

func PhpVersion() string {
	var phpVersion = ""
	for phpVersion == "" {
		phpVersionTmp := getPrompt("Enter PHP version")

		if utils.GetCommandOutput("bash", "-c", "command -v php-fpm" + phpVersionTmp) != "" {
			phpVersion = phpVersionTmp
		}
	}

	return phpVersion
}

func Confirm(defaultY bool) bool {
	var confirmation string
	if defaultY {
		confirmation = "Confirm (Y/n) [Y]"
	} else {
		confirmation = "Confirm (Y/n) [n]"
	}

	result := getPrompt(confirmation)

	if defaultY == true && result == "n" || result == "N" {
		return false
	}

	if defaultY == false && result != "y" || result != "Y" {
		return false
	}

	return true
}

func getPrompt(message string) string {
	prompt := promptui.Prompt {
		Label: message,
	}

	result, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	return result
}
