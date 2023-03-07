package pocket

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"time"
)

type PoktConfig struct {
	PoktLogLevel          string        `toml:"pokt_log_level"`
	PoktApiUrl            string        `toml:"pokt_api_url"`
	PoktAccessToken       string        `toml:"pokt_access_token"`
	PoktAddress           string        `toml:"pokt_address"`
	PoktDockerImage       string        `toml:"pokt_docker_image"`
	PoktDockerName        string        `toml:"pokt_docker_name"`
	PoktConfigPath        string        `toml:"pokt_path"`
	PoktScanInterval      time.Duration `toml:"pokt_scan_interval"`
	PoktHeartbeatInterval time.Duration `toml:"pokt_heartbeat_interval"`
	PoktServerApiUrl      string        `toml:"pokt_server_api_url"`
	PoktServerApiPort     int           `toml:"pokt_server_api_port"`
	PoktNetworkType       string        `toml:"pokt_network_type"`
}

var config *PoktConfig

func GetConfig() PoktConfig {
	if config == nil {
		initConfig()
	}
	return *config
}

func initConfig() {
	configPath := os.Getenv("SWAN_PATH")
	if configPath == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			GetLog().Fatal("Cannot get home directory.")
		}

		configPath = filepath.Join(homedir, ".swan/")
	}

	initPoktConfig(filepath.Join(configPath, "provider/config-pokt.toml"))

}

func initPoktConfig(configFile string) {
	GetLog().Debug("Your pokt config file is:", configFile)

	if metaData, err := toml.DecodeFile(configFile, &config); err != nil {
		GetLog().Fatal("error:", err)
	} else {
		if !requiredPoktAreGiven(metaData) {
			GetLog().Fatal("required fields not given")
		}
	}
}

func requiredPoktAreGiven(metaData toml.MetaData) bool {
	requiredFields := [][]string{
		{"pokt", "log_level"},
		{"pokt", "pokt_api_url"},
		{"pokt", "pokt_access_token"},
		{"pokt", "pokt_docker_image"},
		{"pokt", "pokt_docker_name"},
		{"pokt", "pokt_path"},
		{"pokt", "pokt_scan_interval"},
		{"pokt", "pokt_heartbeat_interval"},
		{"pokt", "pokt_server_api_url"},
		{"pokt", "pokt_server_api_port"},
		{"pokt", "pokt_network_type"},
	}

	for _, v := range requiredFields {
		if !metaData.IsDefined(v...) {
			GetLog().Fatal("required conf fields ", v)
		}
	}

	return true
}
