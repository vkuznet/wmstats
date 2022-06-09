package main

// configuration module for wmstats
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Configuration stores configuration parameters
type Configuration struct {
	Port            int      `json:"port"`              // server port number
	StaticDir       string   `json:"staticdir"`         // location of static directory
	Base            string   `json:"base"`              // server base path
	Verbose         int      `json:"verbose"`           // verbosity level
	LogFile         string   `json:"log_file"`          // server log file (should ends with .log) or log area
	Hmac            string   `json:"hmac"`              // cmsweb hmac file location
	LimiterPeriod   string   `json:"limiter_rate"`      // limiter rate value
	LimiterHeader   string   `json:"limiter_header"`    // limiter header to use
	LimiterSkipList []string `json:"limiter_skip_list"` // limiter skip list
	MetricsPrefix   string   `json:"metrics_prefix"`    // metrics prefix used for prometheus
	CMSRole         string   `json:"cms_role"`          // cms role for write access
	CMSGroup        string   `json:"cms_group"`         // cms group for write access
	AccessURI       string   `json:"access_uri"`        // access URI, either URL or filename

	// server static parts
	Templates string `json:"templates"` // location of server templates
	Jscripts  string `json:"jscripts"`  // location of server JavaScript files
	Images    string `json:"images"`    // location of server images
	Styles    string `json:"styles"`    // location of server CSS styles

	// security parts
	ServerKey  string `json:"serverkey"`  // server key for https
	ServerCrt  string `json:"servercrt"`  // server certificate for https
	RootCA     string `json:"rootCA"`     // RootCA file
	CSRFKey    string `json:"csrfKey"`    // CSRF 32-byte-long-auth-key
	Production bool   `json:"production"` // production server or not
}

// Config represents global configuration object
var Config Configuration

// String returns string representation of Config
func (c *Configuration) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("ERROR: fail to marshal configuration", err)
		return ""
	}
	return string(data)
}

// ParseConfig parses given configuration file and initialize Config object
func ParseConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("unable to read config file", configFile, err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("unable to parse config file", configFile, err)
		return err
	}
	if Config.LimiterPeriod == "" {
		Config.LimiterPeriod = "100-S"
	}
	if Config.Templates == "" {
		Config.Templates = fmt.Sprintf("%s/templates", Config.StaticDir)
	}
	if Config.MetricsPrefix == "" {
		Config.MetricsPrefix = "wmstats"
	}
	return nil
}
