package radio

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/theoutdoorclub/radio/shared"
)

type Config struct {
	Credentials CredentialsConfig `toml:"credentials"`
	Nodes       []NodeConfig      `toml:"nodes"`
}

type CredentialsConfig struct {
	Token string `toml:"token"`
}

type NodeConfig struct {
	Name     string `toml:"name"`
	Address  string `toml:"address"`
	Password string `toml:"password"`
	Secure   bool   `toml:"secure"`
}

func ParseConfig() (Config, error) {
	var conf Config
	_, err := toml.DecodeFile(filepath.Join(shared.CWD, "config.toml"), &conf)

	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
