package template

import (
	"devman/internal/user"
	"fmt"
	"io/ioutil"
)

func GetTemplate(name string) string {
	realUser := user.GetName()

	path := "/home/" + realUser + "/.config/devman/" + name + ".stub"

	content, _ := ioutil.ReadFile(path)

	return string(content)
}

func WpHtaccess() string {
	return `
# BEGIN WordPress
# The directives (lines) between "BEGIN WordPress" and "END WordPress" are
# dynamically generated, and should only be modified via WordPress filters.
# Any changes to the directives between these markers will be overwritten.
<IfModule mod_rewrite.c>
RewriteEngine On
RewriteBase /
RewriteRule ^index\.php$ - [L]
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule . /index.php [L]
</IfModule>
# END WordPress
`
}

func ApacheConfig(domain, docRoot string) string {
	return fmt.Sprintf(`
    <VirtualHost *:80>
        ServerName %s

        ServerAdmin webmaster@localhost
        DocumentRoot /var/www/html/%s

        <Directory /var/www/html/%s/>
        Options FollowSymLinks
        AllowOverride All
        Order allow,deny
        allow from all
        </Directory>

		# Available loglevels: trace8, ..., trace1, debug, info, notice, warn,
        # error, crit, alert, emerg.
        # It is also possible to configure the loglevel for particular
        # modules, e.g.
        #LogLevel info ssl:warn

        ErrorLog ${APACHE_LOG_DIR}/error.log
        CustomLog ${APACHE_LOG_DIR}/access.log combined

        <DirectoryMatch "^/.*/\.git/">
                Order deny,allow
                Deny from all
        </DirectoryMatch>
    </VirtualHost>
    <VirtualHost *:443>
		ServerName %s

		ServerAdmin webmaster@localhost
        DocumentRoot /var/www/html/%s

		ErrorLog ${APACHE_LOG_DIR}/error.log
        CustomLog ${APACHE_LOG_DIR}/access.log combined

		SSLEngine on
        SSLCertificateFile /etc/ssl/certs/ssl-cert-snakeoil.pem
        SSLCertificateKeyFile /etc/ssl/private/ssl-cert-snakeoil.key

		<FilesMatch \"\.(cgi|shtml|phtml|php)$\">
                       SSLOptions +StdEnvVars
        </FilesMatch>

		<Directory /var/www/html/%s/>
        Options FollowSymLinks
        AllowOverride All
        Order allow,deny
        allow from all
                   SSLOptions +StdEnvVars
        </Directory>

		BrowserMatch \"MSIE [2-6]\" \
                   nokeepalive ssl-unclean-shutdown \
                   downgrade-1.0 force-response-1.0
        BrowserMatch \"MSIE [17-9]\" ssl-unclean-shutdown

        <DirectoryMatch "^/.*/\.git/">
                Order deny,allow
                Deny from all
        </DirectoryMatch>
    </VirtualHost>
`, domain, docRoot, domain, domain, docRoot, domain)
}

func NginxConfig(domain, docRoot, phpVersion string) string {
	return fmt.Sprintf(`
	server {
		listen 80;
		listen [::]:80;

		# SSL configuration
		#
		listen 443 ssl;
		listen [::]:443 ssl;

		include snippets/snakeoil.conf;

		root /var/www/html/%s;

		# Add index.php to the list if you are using PHP
		index index.php index.html index.htm index.nginx-debian.html;

		server_name %s;

		location / {
				# First attempt to serve request as file, then
				# as directory, then fall back to displaying a 404.
				try_files $uri $uri/ /index.php?$query_string;
		}

		# pass PHP scripts to FastCGI server
		#
		location ~ .php$ {
				include snippets/fastcgi-php.conf;
				# fastcgi_pass 127.0.0.1:9000;
				fastcgi_pass unix:/var/run/php/php%s-fpm.sock;
		}

		location ~ /\.git {
		  deny all;
		}
	}
`, docRoot, domain, phpVersion)
}
