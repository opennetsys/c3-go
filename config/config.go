package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/c3systems/c3-go/common/dirutil"
)

var fileperm = os.FileMode(0644)

// NOTE: properties must be uppercase (exported) to save as TOML
type config struct {
	Port           int    `toml:"port"`
	DataDir        string `toml:"dataDir"`
	PrivateKeyPath string `toml:"privateKey"`
}

// Config ...
type Config struct {
	config  *config
	muxlock sync.Mutex
}

// New ...
func New() *Config {
	cnf := &Config{
		config: &config{
			Port:           DefaultServerPort,
			DataDir:        DefaultStoreDirectory,
			PrivateKeyPath: DefaultConfigDirectory + "/priv.pem",
		},
	}
	if err := cnf.setupConfig(); err != nil {
		log.Fatal(err)
	}

	return cnf
}

// Port ...
func (cnf *Config) Port() int {
	return cnf.config.Port
}

// DataDir ...
func (cnf *Config) DataDir() string {
	return cnf.config.DataDir
}

// PrivateKeyPath ...
func (cnf *Config) PrivateKeyPath() string {
	return cnf.config.PrivateKeyPath
}

func (cnf *Config) setupConfig() error {
	err := cnf.makeConfigDir()
	if err != nil {
		return err
	}
	err = cnf.makeConfigFile()
	if err != nil {
		return err
	}
	conf, err := cnf.parseConfig()
	if err != nil {
		return err
	}
	cnf.config = conf
	err = cnf.loadConfig()
	if err != nil {
		return err
	}
	return nil
}

func (cnf *Config) loadConfig() error {
	log.Println(cnf.config)

	return nil
}

func (cnf *Config) configDirPath() string {
	homedir := dirutil.UserHomeDir()
	return fmt.Sprintf("%s%s", homedir, "/.c3")
}

func (cnf *Config) configPath() string {
	return fmt.Sprintf("%v%v", cnf.configDirPath(), "/config.toml")
}

func (cnf *Config) makeConfigDir() error {
	path := cnf.configDirPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func (cnf *Config) makeConfigFile() error {
	path := cnf.configPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fo, err := os.Create(path)
		if err != nil {
			return err
		}
		defer fo.Close()
		b, err := cnf.configToToml()
		if err != nil {
			return err
		}
		if _, err := fo.Write(b); err != nil {
			return err
		}
	}

	return nil
}

func (cnf *Config) saveConfig() error {
	cnf.muxlock.Lock()
	defer cnf.muxlock.Unlock()
	path := cnf.configPath()
	if _, err := os.Stat(path); err == nil {
		b, err := cnf.configToToml()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, b, fileperm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cnf *Config) parseConfig() (*config, error) {
	var conf config
	path := cnf.configPath()
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (cnf *Config) configToToml() ([]byte, error) {
	var b bytes.Buffer
	encoder := toml.NewEncoder(&b)

	err := encoder.Encode(cnf.config)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
