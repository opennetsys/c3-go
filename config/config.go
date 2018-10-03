package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/c3systems/c3-go/common/dirutil"
	"github.com/c3systems/c3-go/common/ipnsutil"
)

var fileperm = os.FileMode(0644)

// NOTE: properties must be uppercase (exported) to save as TOML
type config struct {
	Port            int    `toml:"port"`
	DataDir         string `toml:"dataDir"`
	PrivateKeyPath  string `toml:"privateKey"`
	Peer            string `toml:"peer"`
	BlockDifficulty int    `toml:"blockDifficulty"`
	configDir       string `toml:"-"` // NOTE: don't save to TOML
	configFilename  string `toml:"-"` // NOTE: don't save to TOML
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
			Port:            DefaultServerPort,
			configDir:       DefaultConfigDirectory,
			configFilename:  DefaultConfigFilename,
			DataDir:         DefaultStoreDirectory,
			PrivateKeyPath:  DefaultConfigDirectory + "/" + DefaultPrivateKeyFilename,
			Peer:            "",
			BlockDifficulty: DefaultBlockDifficulty,
		},
	}
	if err := cnf.setupConfig(); err != nil {
		log.Fatal(err)
	}

	return cnf
}

// NewFromFile ...
func NewFromFile(filepath string) *Config {
	if filepath == "" {
		log.Fatal("filepath is required")
	}

	paths := strings.Split(filepath, "/")
	filename := strings.Join(paths[len(paths)-1:len(paths)], "")
	filedir := "./"
	if len(paths) > 1 {
		filedir = strings.Join(paths[0:len(paths)-2], "/")
	}

	cnf := &Config{
		config: &config{
			Port:            DefaultServerPort,
			configDir:       filedir,
			configFilename:  filename,
			DataDir:         DefaultStoreDirectory,
			PrivateKeyPath:  DefaultConfigDirectory + "/" + DefaultPrivateKeyFilename,
			Peer:            "",
			BlockDifficulty: DefaultBlockDifficulty,
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

// NodeURI ...
func (cnf *Config) NodeURI() string {
	return fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", cnf.Port())
}

// DataDir ...
func (cnf *Config) DataDir() string {
	return dirutil.NormalizePath(cnf.config.DataDir)
}

// PrivateKeyPath ...
func (cnf *Config) PrivateKeyPath() string {
	return dirutil.NormalizePath(cnf.config.PrivateKeyPath)
}

// PrivateKeyIPNS ...
func (cnf *Config) PrivateKeyIPNS() string {
	id, err := ipnsutil.PEMToIPNS(cnf.PrivateKeyPath(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

// Peer ...
func (cnf *Config) Peer() string {
	return cnf.config.Peer
}

// BlockDifficulty ...
func (cnf *Config) BlockDifficulty() int {
	return cnf.config.BlockDifficulty
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
	return nil
}

func (cnf *Config) configDirPath() string {
	return dirutil.NormalizePath(cnf.config.configDir)
}

func (cnf *Config) configPath() string {
	return fmt.Sprintf("%v%v", cnf.configDirPath(), "/"+cnf.config.configFilename)
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
