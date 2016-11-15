/*
   Copyright (c) 2016 VMware, Inc. All Rights Reserved.
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package config

import (
	"fmt"
	"os"
	"strings"
)

// ConfigLoader is the interface to load configurations
type ConfigLoader interface {

	// Load will load configuration from different source into a string map, the values in the map will be parsed in to configurations.
	Load() (map[string]string, error)
}

// EnvConfigLoader loads the config from env vars.
type EnvConfigLoader struct {
	keys []string
}

// Load ...
func (ec *EnvConfigLoader) Load() (map[string]string, error) {
	m := make(map[string]string)
	for _, k := range ec.keys {
		m[k] = os.Getenv(k)
	}
	return m, nil
}

// ConfigParser
type ConfigParser interface {
	//Parse parse the input raw map into a config map
	Parse(raw map[string]string, config map[string]interface{}) error
}

type Config struct {
	config map[string]interface{}
	loader ConfigLoader
	parser ConfigParser
}

func (conf *Config) load() error {
	rawMap, err := conf.loader.Load()
	if err != nil {
		return err
	}
	err = conf.parser.Parse(rawMap, conf.config)
	return err
}

type MySQLSetting struct {
	Database string
	User     string
	Password string
	Host     string
	Port     string
}

type SQLiteSetting struct {
	FilePath string
}

type commonParser struct{}

// Parse parses the db settings, veryfy_remote_cert, ext_endpoint, token_endpoint
func (cp *commonParser) Parse(raw map[string]string, config map[string]interface{}) error {
	db := strings.ToLower(raw["DATABASE"])
	if db == "mysql" || db == "" {
		db = "mysql"
		mySQLDB := raw["MYSQL_DATABASE"]
		if len(mySQLDB) == 0 {
			mySQLDB = "registry"
		}
		setting := MySQLSetting{
			mySQLDB,
			raw["MYSQL_USR"],
			raw["MYSQL_PWD"],
			raw["MYSQL_HOST"],
			raw["MYSQL_PORT"],
		}
		config["mysql"] = setting
	} else if db == "sqlite" {
		f := raw["SQLITE_FILE"]
		if len(f) == 0 {
			f = "registry.db"
		}
		setting := SQLiteSetting{
			f,
		}
		config["sqlite"] = setting
	} else {
		return fmt.Errorf("Invalid DB: %s", db)
	}
	config["database"] = db

	//By default it's true
	config["verify_remote_cert"] = raw["VERIFY_REMOTE_CERT"] != "off"

	config["ext_endpoint"] = raw["EXT_ENDPOINT"]
	config["token_endpoint"] = raw["TOKEN_ENDPOINT"]
	config["log_level"] = raw["LOG_LEVEL"]
	return nil
}

var commonConfig *Config

func init() {
	commonKeys := []string{"DATABASE", "MYSQL_DATABASE", "MYSQL_USR", "MYSQL_PWD", "MYSQL_HOST", "MYSQL_PORT", "SQLITE_FILE", "VERIFY_REMOTE_CERT", "EXT_ENDPOINT", "TOKEN_ENDPOINT", "LOG_LEVEL"}
	commonConfig = &Config{
		config: make(map[string]interface{}),
		loader: &EnvConfigLoader{keys: commonKeys},
		parser: &commonParser{},
	}
	if err := commonConfig.load(); err != nil {
		panic(err)
	}
}

// Reload will reload the configuration.
func Reload() error {
	return commonConfig.load()
}

// Database returns the DB type in configuration.
func Database() string {
	return commonConfig.config["database"].(string)
}

// MySQL returns the mysql setting in configuration.
func MySQL() MySQLSetting {
	return commonConfig.config["mysql"].(MySQLSetting)
}

// SQLite returns the SQLite setting
func SQLite() SQLiteSetting {
	return commonConfig.config["sqlite"].(SQLiteSetting)
}

// VerifyRemoteCert returns bool value.
func VerifyRemoteCert() bool {
	return commonConfig.config["verify_remote_cert"].(bool)
}

// ExtEndpoint ...
func ExtEndpoint() string {
	return commonConfig.config["ext_endpoint"].(string)
}

// TokenEndpoint returns the endpoint string of token service, which can be accessed by internal service of Harbor.
func TokenEndpoint() string {
	return commonConfig.config["token_endpoint"].(string)
}

func LogLevel() string {
	return commonConfig.config["log_level"].(string)
}
