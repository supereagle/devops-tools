package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Config struct {
	JenkinsServers  []string `json:"jenkins_servers,omitempty"`
	JenkinsUser     string   `json:"jenkins_user,omitempty"`
	JenkinsPassword string   `json:"jenkins_password,omitempty"`
	NodeLabels      []string `json:"node_labels,omitempty"`
}

func ReadConfig(filePath string) (config *Config, err error) {
	filePath = strings.TrimSpace(filePath)
	if len(filePath) == 0 {
		return nil, fmt.Errorf("The config file path can not be empty")
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config = &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		err = fmt.Errorf("Fail to unmarshal config file to json object as %s", err.Error())
		return
	}

	return
}
