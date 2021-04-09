package commands

import (
	"devman/internal/database"
	"devman/internal/environment"
	"devman/internal/output"
	"devman/internal/prompt"
	"devman/internal/template"
	"devman/internal/user"
	"devman/internal/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

		lastChar := documentRoot[len(documentRoot) -1:]
		if lastChar == "/" {
			documentRoot = documentRoot[:len(documentRoot) - 1]
		}
	}

	var dbPass = prompt.MySQLPasswordPrompt()

	var dbName = ""
	if isInstallWordPress && dbPass == "" {
		dbName = prompt.DatabaseName(strings.Replace(domain, ".", "_", 1))
	} else if dbPass != "" {
		dbName = strings.Replace(domain, ".", "_", 1)
	}

	var phpVersion string
	if environment.IsRunningNginx() {
		phpVersion = prompt.PhpVersion()
	}

	output.Info("Site URL: " + domain)
	output.Info("Site path: " + sitePath)

	if !isInstallWordPress && documentRoot != "" {
		output.Info("Document root: " + documentRoot)
	}

	var confirm = prompt.Confirm()

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
		output.Info("Creating database...")
		createdDb = database.CreateDatabase(dbName, dbPass)

		if !createdDb {
			output.Danger("Database creation failed. Does it already exist?")
		} else {
			output.Info("Database " + dbName + " created.")
		}
	}

	if isInstallWordPress {
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

		output.Info("Setting permissions...")
		err = utils.ChownR(sitePath, uid, gid)
		if err != nil {
			panic(err)
		}

		output.Info("Moving files... ")
		cmd := exec.Command("/bin/sh", "-c", "mv " + sitePath + "/wordpress/* " + sitePath + "/")
		_ = cmd.Run()
		_ = os.RemoveAll(sitePath + "/wordpress")
		_ = ioutil.WriteFile(sitePath+"/.htaccess", []byte(htaccess), 0644)
		_ = os.Chown(sitePath+"/.htaccess", uid, gid)

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

	enabledConfigLocation := "/etc/" + webserverSlug + "/sites-enabled/" + domain + ".conf"
	configLocation := "/etc/" + webserverSlug + "/sites-available/" + domain + ".conf"

	// Setup Nginx/Apache config.

	output.Info("Configuring " + webserver + "...")

	if utils.FileExists(enabledConfigLocation) {
		output.Warning(enabledConfigLocation + " already exists.")
	} else {
		output.Info("Creating configuration...")
		var config string
		if environment.IsRunningApache() {
			config = template.ApacheConfig(domain, documentRoot)
		} else if environment.IsRunningNginx() {
			config = template.NginxConfig(domain, documentRoot, phpVersion)
		}

		utils.WriteFile(configLocation, config)

		output.Info("Creating symbolic link " + configLocation + " /etc/" + webserverSlug + "/sites-enabled/")
		cmd := exec.Command("ln", "-s", configLocation, "/etc/" + webserverSlug + "/sites-enabled/")
		_ = cmd.Run()
	}

	if utils.FileExists("/var/www/html/" + domain) {
		output.Warning("/var/www/html/" + domain + " already exists.")
	} else {
		output.Info("Creating symbolic link " + sitePath + "/ /var/www/html/" + domain)
		cmd := exec.Command("ln", "-s", sitePath + "/", "/var/www/html/" + domain)
		_ = cmd.Run()
	}

	output.Info("Mapping DNS...")

	if utils.FileContains("/etc/hosts", "127.0.0.1 " + domain) {
		output.Warning("Entry already exists in /etc/hosts")
	} else {
		utils.AppendToFile("/etc/hosts", "127.0.0.1 " + domain)
	}

	output.Info("Restarting " + webserver)
	cmd := exec.Command("service", environment.WebServerProcessName(), "restart")
	_ = cmd.Run()

	output.Success("Finished!")
}
