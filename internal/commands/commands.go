package commands

import (
	"github.com/tsquare17/devman/internal/database"
	"github.com/tsquare17/devman/internal/environment"
	"github.com/tsquare17/devman/internal/output"
	"github.com/tsquare17/devman/internal/prompt"
	"github.com/tsquare17/devman/internal/template"
	"github.com/tsquare17/devman/internal/user"
	"github.com/tsquare17/devman/internal/utils"
	"io/ioutil"
	"os"
	"os/exec"
)

func NewSite(domain string) {
	if user.IsSuper() != true {
		output.Danger("You must run this command with sudo.")
		os.Exit(0)
	}

	var webserver string
	var webserverSlug string
	if environment.IsRunningApache() {
		webserver = "Apache"
		webserverSlug = "apache2"
	} else if environment.IsRunningNginx() {
		webserver = "Nginx"
		webserverSlug = "nginx"
	} else if environment.IsRunningCaddy() {
		webserver = "Caddy"
		webserverSlug = "caddy"
	} else {
		output.Danger("Could not detect web server...")
		os.Exit(0)
	}

	output.Info("Detected " + webserver + " web server.")

	var sitePath = prompt.SitePath()
	lastChar := sitePath[len(sitePath) -1:]
	if lastChar == "/" {
		sitePath = sitePath[:len(sitePath) - 1]
	}

	var isInstallWordPress = prompt.IsInstallWordPress()
	var documentRoot = ""

	if !isInstallWordPress {
		documentRoot = prompt.SiteDocumentRoot()

		if documentRoot != "" {
			lastChar := documentRoot[len(documentRoot) -1:]
			if lastChar == "/" {
				documentRoot = documentRoot[:len(documentRoot) - 1]
			}
		}
	}

	var dbPass = prompt.MySQLPasswordPrompt()

	var dbName = ""
	if isInstallWordPress && dbPass == "" {
		dbName = prompt.DatabaseName(utils.Slugify(domain))
	} else if dbPass != "" {
		dbName = utils.Slugify(domain)
	}

	var phpVersion string
	phpVersion = prompt.PhpVersion()

	output.Info("Site URL: " + domain)
	output.Info("Site path: " + sitePath)

	if !isInstallWordPress && documentRoot != "" {
		output.Info("Document root: " + documentRoot)
	}

	var confirm = prompt.Confirm(true)

	if !confirm {
		output.Info("Exiting.")
		os.Exit(0)
	}

	output.Info("Running...")

	if documentRoot == "" {
		documentRoot = domain
	} else {
		documentRoot = domain + "/" + documentRoot
	}

	if !utils.FileExists(sitePath) {
		output.Danger("Site path does not exist.")
		os.Exit(0)
	}

	uid, gid := user.GetUidGid(sitePath)

	var createdDb = false
	if dbName != "" && dbPass != "" {
		output.Info("Creating database: " + dbName)
		createdDb = database.CreateDatabase(dbName, dbPass)

		if !createdDb {
			output.Danger("Database creation failed. Does it already exist?")
		} else {
			output.Info("Database created")
		}
	}

	if isInstallWordPress {
		// TODO: Check for user defined templates in the template method itself.
		htaccess := template.GetTemplate("wp-htaccess")
		if htaccess == "" {
			htaccess = template.WpHtaccess()
		}

		output.Info("Downloading WordPress archive...")
		err := utils.DownloadFile("/tmp/latest.tar.gz", "https://wordpress.org/latest.tar.gz")
		if err != nil {
			if createdDb {
				output.Danger("Failed to download WordPress archive. Rolling back...")
				database.DeleteDatabase(dbName, dbPass)
			} else {
				output.Danger("Failed to download WordPress archive. Exiting.")
			}

			os.Exit(0)
		}

		wpPathReader, err := os.Open("/tmp/latest.tar.gz")
		if err != nil {
			output.Danger("Failed to open archive...")
		}

		output.Info("Unpacking archive...")
		err = utils.UnTar(sitePath, wpPathReader)
		if err != nil {
			output.Danger("Failed to unpack archive...")
		}

		output.Info("Moving files... ")
		cmd := exec.Command("/bin/sh", "-c", "mv " + sitePath + "/wordpress/* " + sitePath + "/")
		_ = cmd.Run()
		_ = os.RemoveAll(sitePath + "/wordpress")
		_ = ioutil.WriteFile(sitePath+"/.htaccess", []byte(htaccess), 0664)

		err = os.Chown(sitePath+"/.htaccess", uid, gid)
		if err != nil {
			output.Warning(err.Error())
		}

		if dbName != "" {
			cmd = exec.Command("/bin/sh", "-c", "sed -i \"s/database_name_here/" + dbName + "/g\" " + sitePath + "/wp-config-sample.php")
			_ = cmd.Run()

			cmd = exec.Command("/bin/sh", "-c", "sed -i \"s/username_here/root/g\" " + sitePath + "/wp-config-sample.php")
			_ = cmd.Run()
		}

		if dbPass != "" {
			cmd = exec.Command("/bin/sh", "-c", "sed -i \"s/password_here/" + dbPass + "/g\" " + sitePath + "/wp-config-sample.php")
			_ = cmd.Run()
		}

		if dbName != "" || dbPass != "" {
			_ = os.Rename(sitePath+"/wp-config-sample.php", sitePath+"/wp-config.php")
		}
	}

	output.Info("Setting permissions on: " + sitePath)
	err := utils.ChownR(sitePath, uid, gid)
	if err != nil {
		output.Warning(err.Error())
	}

	enabledConfigLocation := "/etc/" + webserverSlug + "/sites-enabled/" + domain + ".conf"

	var configLocation string
	if environment.IsRunningCaddy() {
		configLocation = "/etc/caddy/Caddyfile"
	} else {
		configLocation = "/etc/" + webserverSlug + "/sites-available/" + domain + ".conf"
	}

	// Setup web server config.

	output.Info("Configuring " + webserver + "...")

	if utils.FileExists(enabledConfigLocation) && !environment.IsRunningCaddy() {
		output.Warning(enabledConfigLocation + " already exists.")
	} else {
		output.Info("Creating configuration...")
		var config string
		if environment.IsRunningApache() {
			config = template.ApacheConfig(domain, documentRoot, phpVersion)
		} else if environment.IsRunningNginx() {
			config = template.NginxConfig(domain, documentRoot, phpVersion)
		} else if environment.IsRunningCaddy() {
			config = template.CaddyConfig(domain, documentRoot, phpVersion)
		}

		if !environment.IsRunningCaddy() {
			utils.WriteFile(configLocation, config)

			output.Info("Creating symbolic link " + configLocation + " /etc/" + webserverSlug + "/sites-enabled/")
			cmd := exec.Command("ln", "-s", configLocation, "/etc/" + webserverSlug + "/sites-enabled/")
			_ = cmd.Run()
		} else {
			if !utils.FileContains(configLocation, domain) {
				utils.AppendFile(configLocation, config)

				output.Info("Added " + domain + " to Caddyfile")
			}
		}
	}

	if utils.FileExists("/var/www/html/" + domain) {
		output.Warning("/var/www/html/" + domain + " already exists.")
	} else {
		output.Info("Creating symbolic link " + sitePath + "/ /var/www/html/" + domain)
		cmd := exec.Command("ln", "-s", sitePath + "/", "/var/www/html/" + domain)
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}

	output.Info("Mapping DNS...")

	if utils.FileContains("/etc/hosts", "127.0.0.1 " + domain) {
		output.Warning("Entry already exists in /etc/hosts")
	} else {
		utils.AppendToFile("/etc/hosts", "\n127.0.0.1 " + domain)
	}

	output.Info("Restarting " + webserver)
	cmd := exec.Command("service", environment.WebServerProcessName(), "restart")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	output.Success("Finished!")
}

func RemoveSite(domain string) {
	if user.IsSuper() != true {
		output.Danger("You must run this command with sudo.")
		os.Exit(0)
	}

	output.Warning("Are you sure you want to remove " + domain)

	var confirm = prompt.Confirm(true)

	if !confirm {
		output.Info("Exiting.")
		os.Exit(0)
	}


	var webserverSlug string
	if environment.IsRunningApache() {
		webserverSlug = "apache2"
	} else if environment.IsRunningNginx() {
		webserverSlug = "nginx"
	} else {
		output.Danger("Could not detect web server...")
	}

	output.Info("Removing hosts entry...")

	if utils.FileContains("/etc/hosts", "127.0.0.1 " + domain) {
		output.Warning("/etc/hosts does not contain entry for " + domain)
	} else {
		utils.RemoveFromFile("/etc/hosts", "127.0.0.1 " + domain)
	}

	output.Info("Removing site symlink...")

	if utils.FileExists("/var/www/html/" + domain) {
		utils.RemoveFile("/var/www/html/" + domain)
	} else {
		output.Warning("Site directory symlink wasn't found at /var/www/html/" + domain)
	}

	output.Info("Removing enabled config file...")

	enabledConfig := "/etc/" + webserverSlug + "/sites-enabled/" + domain + ".conf"
	if utils.FileExists(enabledConfig) {
		utils.RemoveFile(enabledConfig)
	} else {
		output.Warning("Configuration not found at " + enabledConfig)
	}

	output.Info("Removing available config file...")

	availableConfig := "/etc/" + webserverSlug + "/sites-available/" + domain + ".conf"
	if utils.FileExists(availableConfig) {
		utils.RemoveFile(availableConfig)
	} else {
		output.Warning("Configuration not found at " + availableConfig)
	}

	output.Success("Finished!")
}
