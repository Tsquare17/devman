package prompt

import (
	"devman/internal/utils"
	"github.com/manifoldco/promptui"
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
	return getPrompt("Enter root MySQL password (if creating database)")
}

func IsInstallWordPress() bool {
	result := getPrompt("Install WordPress? (Y/n) [n]")

	if result == "y" || result == "Y" {
		return true
	}

	return false
}

func DatabaseName(defaultName string) string {
	return getPrompt("Enter database name [" + defaultName + "]")
}

func Confirm() bool {
	result := getPrompt("Confirm (Y/n) [Y]")

	if result == "n" || result == "N" {
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
