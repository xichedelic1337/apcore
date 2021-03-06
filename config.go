// apcore is a server framework for implementing an ActivityPub application.
// Copyright (C) 2019 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package apcore

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/ini.v1"
)

const (
	postgresDB = "postgres"
)

// Overall configuration file structure
type config struct {
	ServerConfig      serverConfig      `ini:"server" comment:"HTTP server configuration"`
	OAuthConfig       oAuthConfig       `ini:"oauth" comment:"OAuth 2 configuration"`
	DatabaseConfig    databaseConfig    `ini:"database" comment:"Database configuration"`
	ActivityPubConfig activityPubConfig `ini:"activitypub" comment:"ActivityPub configuration"`
}

func defaultConfig(dbkind string) (c *config, err error) {
	var dbc databaseConfig
	dbc, err = defaultDatabaseConfig(dbkind)
	if err != nil {
		return
	}
	c = &config{
		ServerConfig:      defaultServerConfig(),
		OAuthConfig:       defaultOAuthConfig(),
		DatabaseConfig:    dbc,
		ActivityPubConfig: defaultActivityPubConfig(),
	}
	return
}

// Configuration section specifically for the HTTP server.
type serverConfig struct {
	Host                        string `ini:"sr_host" comment:"(required) Host with TLD for this instance (basically, the fully qualified domain or subdomain); ignored in debug mode"`
	CertFile                    string `init:"sr_cert_file" comment:"(required) Path to the certificate file used to establish TLS connections for HTTPS"`
	KeyFile                     string `init:"sr_cert_file" comment:"(required) Path to the private key file used to establish TLS connections for HTTPS"`
	CookieAuthKeyFile           string `ini:"sr_cookie_auth_key_file" comment:"(required) Path to private key file used for cookie authentication"`
	CookieEncryptionKeyFile     string `ini:"sr_cookie_encryption_key_file" comment:"Path to private key file used for cookie encryption"`
	CookieMaxAge                int    `ini:"sr_cookie_max_age" comment:"(default: 86400 seconds) Number of seconds a cookie is valid; 0 indicates no Max-Age (browser-dependent, usually session-only); negative value is invalid"`
	CookieSessionName           string `ini:"sr_cookie_session_name" comment:"(required) Cookie session name to use for the application"`
	HttpsReadTimeoutSeconds     int    `ini:"sr_https_read_timeout_seconds" comment:"Timeout in seconds for incoming HTTPS requests; a zero or unset value does not timeout"`
	HttpsWriteTimeoutSeconds    int    `ini:"sr_https_write_timeout_seconds" comment:"Timeout in seconds for outgoing HTTPS responses; a zero or unset value does not timeout"`
	RedirectReadTimeoutSeconds  int    `ini:"sr_redirect_read_timeout_seconds" comment:"Timeout in seconds for incoming HTTP requests, which will be redirected to HTTPS; a zero or unset value does not timeout"`
	RedirectWriteTimeoutSeconds int    `ini:"sr_redirect_write_timeout_seconds" comment:"Timeout in seconds for outgoing HTTP redirect-to-HTTPS responses; a zero or unset value does not timeout"`
	StaticRootDirectory         string `ini:"sr_static_root_directory" comment:"(required) Root directory for serving static content, such as ECMAScript, CSS, favicon; !!!Warning: Everything in this directory will be served and accessible!!!"`
	SaltSize                    int    `ini:"sr_salt_size" comment:"(default: 32) The size of salts to use with passwords when hashing, anything smaller than 16 will be treated as 16"`
	BCryptStrength              int    `ini:"sr_bcrypt_strength" comment:"(default: 10) The hashing cost to use with the bcrypt hashing algorithm, between 4 and 31; the higher the cost, the slower the hash comparisons for passwords will take for attackers and regular users alike"`
	RSAKeySize                  int    `ini:"sr_rsa_private_key_size" comment:"(default: 1024) The size of the RSA private key for a user; values less than 1024 are forbidden"`
}

func defaultServerConfig() serverConfig {
	return serverConfig{
		CookieMaxAge:   86400,
		SaltSize:       32,
		BCryptStrength: bcrypt.DefaultCost,
		RSAKeySize:     1024,
	}
}

type oAuthConfig struct {
	AccessTokenExpiry  int `ini:"oauth_access_token_expiry" comment:"(default: 3600 seconds) Duration in seconds until an access token expires; zero or negative values are invalid."`
	RefreshTokenExpiry int `ini:"oauth_refresh_token_expiry" comment:"(default: 7200 seconds) Duration in seconds until a refresh token expires; zero or negative values are invalid."`
}

func defaultOAuthConfig() oAuthConfig {
	return oAuthConfig{
		AccessTokenExpiry:  3600,
		RefreshTokenExpiry: 7200,
	}
}

// Configuration section specifically for the database.
type databaseConfig struct {
	DatabaseKind              string         `ini:"db_database_kind" comment:"(required) Only \"postgres\" supported"`
	ConnMaxLifetimeSeconds    int            `ini:"db_conn_max_lifetime_seconds" comment:"(default: indefinite) Maximum lifetime of a connection in seconds; a value of zero or unset value means indefinite"`
	MaxOpenConns              int            `ini:"db_max_open_conns" comment:"(default: infinite) Maximum number of open connections to the database; a value of zero or unset value means infinite"`
	MaxIdleConns              int            `ini:"db_max_idle_conns" comment:"(default: 2) Maximum number of idle connections in the connection pool to the database; a value of zero maintains no idle connections; a value greater than max_open_conns is reduced to be equal to max_open_conns"`
	DefaultCollectionPageSize int            `ini:"db_default_collection_page_size" comment:"(default: 10) The default collection page size when fetching a page of an ActivityStreams collection"`
	PostgresConfig            postgresConfig `ini:"db_postgres,omitempty" comment:"Only needed if database_kind is postgres, and values are based on the github.com/lib/pq driver"`
}

func defaultDatabaseConfig(dbkind string) (d databaseConfig, err error) {
	d = databaseConfig{
		DatabaseKind: dbkind,
		// This default is implicit in Go but could change, so here we
		// make it explicit instead
		MaxIdleConns: 2,
		// This default is arbitrarily chosen
		DefaultCollectionPageSize: 10,
	}
	if dbkind != postgresDB {
		err = fmt.Errorf("unsupported database kind: %s", dbkind)
		return
	}
	d.PostgresConfig = defaultPostgresConfig()
	return
}

// Configuration section specifically for ActivityPub.
type activityPubConfig struct {
	ClockTimezone                    string               `ini:"ap_clock_timezone" comment:"(default: UTC) Timezone for ActivityPub related operations: unset and \"UTC\" are UTC, \"Local\" is local server time, otherwise use IANA Time Zone database values"`
	OutboundRateLimitQPS             float64              `ini:"ap_outbound_rate_limit_qps" comment:"(default: 10) Global outbound rate limit for delivery of federated messages under steady state conditions; a negative value or value of zero is invalid"`
	OutboundRateLimitBurst           int                  `ini:"ap_outbound_rate_limit_burst" comment:"(default: 50) Global outbound burst tolerance for delivery of federated messages; a negative value or value of zero is invalid"`
	HttpSignaturesConfig             httpSignaturesConfig `ini:"ap_http_signatures" comment:"HTTP Signatures configuration"`
	MaxInboxForwardingRecursionDepth int                  `ini:"ap_max_inbox_forwarding_recursion_depth" comment:"(default: 50) The maximum recursion depth to use when determining whether to do inbox forwarding, which if triggered ensures older thread participants are able to receive messages; zero means no limit (only used if the application has S2S enabled)"`
	MaxDeliveryRecursionDepth        int                  `ini:"ap_max_delivery_recursion_depth" comment:"(default: 50) The maximum depth to search for peers to deliver due to inbox forwarding, which ensures messages received by this server are propagated to them and no \"ghost reply\" problems occur; zero means no limit (only used if the application has S2S enabled)"`
}

func defaultActivityPubConfig() activityPubConfig {
	return activityPubConfig{
		ClockTimezone:                    "UTC",
		OutboundRateLimitQPS:             10,
		OutboundRateLimitBurst:           50,
		HttpSignaturesConfig:             defaultHttpSignaturesConfig(),
		MaxInboxForwardingRecursionDepth: 50,
		MaxDeliveryRecursionDepth:        50,
	}
}

// Configuration for HTTP Signatures.
type httpSignaturesConfig struct {
	Algorithms      []string `ini:"http_sig_algorithms" comment:"(default: \"sha256,sha512\") Comma-separated list of algorithms used by the go-fed/httpsig library to sign outgoing HTTP signatures; the first algorithm in this list will be the one used to verify other peers' HTTP signatures"`
	DigestAlgorithm string   `ini:"http_sig_digest_algorithm" comment:"(default: \"SHA-256\") RFC ???? algorithm for use in signing header Digests"` // TODO: Find the Digest header RFC for reference
	GetHeaders      []string `ini:"http_sig_get_headers" comment:"(default: \"(request-target),Date,Digest\") Comma-separated list of HTTP headers to sign in GET requests; must contain \"(request-target)\", \"Date\", and \"Digest\""`
	PostHeaders     []string `ini:"http_sig_post_headers" comment:"(default: \"(request-target),Date,Digest\") Comma-separated list of HTTP headers to sign in POST requests; must contain \"(request-target)\", \"Date\", and \"Digest\""`
}

func defaultHttpSignaturesConfig() httpSignaturesConfig {
	return httpSignaturesConfig{
		Algorithms:      []string{"sha256", "sha512"},
		DigestAlgorithm: "SHA-256",
		GetHeaders:      []string{"(request-target)", "Date", "Digest"},
		PostHeaders:     []string{"(request-target)", "Date", "Digest"},
	}
}

// Configuration section specifically for Postgres databases.
type postgresConfig struct {
	DatabaseName            string `ini:"pg_db_name" comment:"(required) Database name"`
	UserName                string `ini:"pg_user" comment:"(required) User to connect as (any password will be prompted)"`
	Host                    string `ini:"pg_host" comment:"(default: localhost) The Postgres host to connect to"`
	Port                    int    `ini:"pg_port" comment:"(default: 5432) The port to connect to"`
	SSLMode                 string `ini:"pg_ssl_mode" comment:"(default: require) SSL mode to use when connecting (options are: \"disable\", \"require\", \"verify-ca\", \"verify-full\")"`
	FallbackApplicationName string `ini:"pg_fallback_application_name" comment:"An application_name to fall back to if one is not provided"`
	ConnectTimeout          int    `ini:"pg_connect_timeout" comment:"(default: indefinite) Maximum wait when connecting to a database, zero or unset means indefinite"`
	SSLCert                 string `ini:"pg_ssl_cert" comment:"PEM-encoded certificate file location"`
	SSLKey                  string `ini:"pg_ssl_key" comment:"PEM-encoded private key file location"`
	SSLRootCert             string `ini:"pg_ssl_root_cert" comment:"PEM-encoded root certificate file location"`
	Schema                  string `ini:"pg_schema" comment:"Postgres schema prefix to use"`
}

func defaultPostgresConfig() postgresConfig {
	return postgresConfig{}
}

func loadConfigFile(filename string, a Application, debug bool) (c *config, err error) {
	InfoLogger.Infof("Loading config file: %s", filename)
	var cfg *ini.File
	cfg, err = ini.Load(filename)
	if err != nil {
		return
	}
	c = &config{}
	err = cfg.MapTo(c)
	if err != nil {
		return
	}
	appCfg := a.NewConfiguration()
	err = cfg.MapTo(appCfg)
	if err != nil {
		return
	}
	err = a.SetConfiguration(appCfg)
	if err != nil {
		return
	}
	if debug {
		c.ServerConfig.Host = "localhost"
	}
	return
}

func saveConfigFile(filename string, c *config, others ...interface{}) error {
	InfoLogger.Infof("Saving config file: %s", filename)
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, c)
	if err != nil {
		return err
	}
	for _, o := range others {
		err = ini.ReflectFrom(cfg, o)
		if err != nil {
			return err
		}
	}
	return cfg.SaveTo(filename)
}

func promptNewConfig(file string) (c *config, err error) {
	fmt.Println(clarkeSays(fmt.Sprintf(`
Welcome to the configuration guided flow!

Here we will visit common configuration choices. While not every option is asked
in the guided flow, you can always open the resulting configuration file to
change options. You can also change your answers to this flow. Note that in
order to take advantage of changed configuration values, the application will
need to be restarted.

Let's go!`)))

	var s string
	s, err = promptSelection(
		"Please choose the database you are using",
		postgresDB)
	if err != nil {
		return
	}
	c, err = defaultConfig(s)
	if err != nil {
		return
	}

	// Prompt for ServerConfig
	c.ServerConfig.Host, err = promptStringWithDefault(
		"Enter the host for this server (ignored in debug mode)",
		"example.com")
	if err != nil {
		return
	}
	c.ServerConfig.CertFile, err = promptString(
		"Enter the path to the file containing the certificate used in HTTPS connections")
	if err != nil {
		return
	}
	c.ServerConfig.KeyFile, err = promptString(
		"Enter the path to the file containing the private key for the certificate used in HTTPS connections")
	if err != nil {
		return
	}
	c.ServerConfig.StaticRootDirectory, err = promptStringWithDefault(
		"Enter the directory for serving static content (WARNING: Everything in it will be served)?",
		"static")
	if err != nil {
		return
	}
	var have bool
	if have, err = promptYN("Do you already have a file containing a cookie authentication private key?"); err != nil {
		return
	} else if have {
		c.ServerConfig.CookieAuthKeyFile, err = promptStringWithDefault(
			"Enter the existing file name for the cookie authentication private key",
			"cookie_authn.key")
		if err != nil {
			return
		}
	} else {
		c.ServerConfig.CookieAuthKeyFile, err = promptStringWithDefault(
			"Enter the new file name for the cookie authentication private key",
			"cookie_authn.key")
		if err != nil {
			return
		}
		err = createKeyFile(c.ServerConfig.CookieAuthKeyFile)
		if err != nil {
			return
		}
	}
	var want bool
	if have, err = promptYN("Do you already have a file containing a cookie encryption private key?"); err != nil {
		return
	} else if have {
		c.ServerConfig.CookieEncryptionKeyFile, err = promptStringWithDefault(
			"Enter the existing file name for the cookie encryption private key",
			"cookie_enc.key")
		if err != nil {
			return
		}
	} else if want, err = promptYN("Do you want to use a cookie encryption private key?"); err != nil {
		return
	} else if want {
		c.ServerConfig.CookieEncryptionKeyFile, err = promptStringWithDefault(
			"Enter the new file name for the cookie encryption private key",
			"cookie_enc.key")
		if err != nil {
			return
		}
		err = createKeyFile(c.ServerConfig.CookieEncryptionKeyFile)
		if err != nil {
			return
		}
	}
	c.ServerConfig.CookieSessionName, err = promptStringWithDefault(
		"Session name used to find cookies",
		"my_apcore_session_name")
	if err != nil {
		return
	}
	c.ServerConfig.HttpsReadTimeoutSeconds, err = promptIntWithDefault(
		"Enter the deadline (in seconds) for reading & writing HTTP & HTTPS requests. A value of zero means connections do not timeout",
		60)
	if err != nil {
		return
	}
	c.ServerConfig.HttpsWriteTimeoutSeconds = c.ServerConfig.HttpsReadTimeoutSeconds
	c.ServerConfig.RedirectReadTimeoutSeconds = c.ServerConfig.HttpsReadTimeoutSeconds
	c.ServerConfig.RedirectWriteTimeoutSeconds = c.ServerConfig.HttpsReadTimeoutSeconds

	// Prompt for ActivityPubConfig
	c.ActivityPubConfig.ClockTimezone, err = promptStringWithDefault(
		"Please enter an IANA Time Zone for the server, \"UTC\", or \"Local\"",
		"UTC")
	if err != nil {
		return
	}
	c.ActivityPubConfig.OutboundRateLimitQPS, err = promptFloat64WithDefault(
		"Please enter the steady-state rate limit for outbound ActivityPub QPS",
		10)
	if err != nil {
		return
	}
	c.ActivityPubConfig.OutboundRateLimitBurst, err = promptIntWithDefault(
		"Please enter the burst limit for outbound ActivityPub QPS",
		50)
	if err != nil {
		return
	}

	// Prompt for DatabaseConfig
	c.DatabaseConfig.ConnMaxLifetimeSeconds, err = promptIntWithDefault(
		"Enter the maximum lifetime (in seconds) for database connections. A value of zero means connections do not timeout",
		60)
	if err != nil {
		return
	}
	c.DatabaseConfig.MaxOpenConns, err = promptIntWithDefault(
		"Enter the maximum number of database connections allowed. A value of zero means infinite are permitted.",
		0)

	switch c.DatabaseConfig.DatabaseKind {
	case postgresDB:
		err = promptPostgresConfig(c)
	default:
		err = fmt.Errorf("unknown database kind: %s", c.DatabaseConfig.DatabaseKind)
	}
	return
}

func promptPostgresConfig(c *config) (err error) {
	fmt.Println("Prompting for Postgres database configuration options...")
	c.DatabaseConfig.PostgresConfig.DatabaseName, err = promptStringWithDefault(
		"Enter the postgres database name",
		"pgdb")
	if err != nil {
		return
	}
	c.DatabaseConfig.PostgresConfig.UserName, err = promptStringWithDefault(
		"Enter the postgres user name",
		"pguser")
	if err != nil {
		return
	}
	c.DatabaseConfig.PostgresConfig.Host, err = promptStringWithDefault(
		"Enter the postgres database host name",
		"localhost")
	if err != nil {
		return
	}
	c.DatabaseConfig.PostgresConfig.Port, err = promptIntWithDefault(
		"Enter the postgres database port",
		5432)
	if err != nil {
		return
	}
	c.DatabaseConfig.PostgresConfig.SSLMode, err = promptSelection(
		"Please choose a SSL mode (see https://www.postgresql.org/docs/current/libpq-ssl.html)",
		"disable",
		"require",
		"verify-ca",
		"verify-full")
	if err != nil {
		return
	}
	if mode := c.DatabaseConfig.PostgresConfig.SSLMode; mode == "require" || mode == "verify-ca" || mode == "verify-full" {
		fmt.Println(clarkeSays(fmt.Sprintf(`
Hey, Clarke the Cow here, I noticed you chose %q! Be sure to check your
configuration file for the %q, %q, and/or %q options to get SSL set up properly!
Toodlemoo~`,
			mode,
			"pg_ssl_cert",
			"pg_ssl_key",
			"pg_ssl_root_cert")))
	}
	return
}
