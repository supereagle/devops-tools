package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Jenkins cleaner"
	app.Usage = "Clean slaves with specified labels for Jenkins Masters"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Config file",
		},
	}

	app.Action = func(c *cli.Context) {
		// Read config file
		filePath := c.String("config")
		cfg, err := ReadConfig(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, server := range cfg.JenkinsServers {
			fmt.Printf("Start to clean slaves with labels %v for Jenkins %s\n", cfg.NodeLabels, server)
			cleaner, err := NewJenkinsCleaner(server, cfg.JenkinsUser, cfg.JenkinsPassword)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = cleaner.CleanSlaves(cfg.NodeLabels...)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}

	app.Run(os.Args)
}

type JenkinsCleaner struct {
	jenkins *gojenkins.Jenkins
}

func NewJenkinsCleaner(jenkinsServer, jenkinsUser, jenkinsPassword string) (*JenkinsCleaner, error) {
	jenkinsServer = strings.TrimSpace(jenkinsServer)
	if len(jenkinsServer) == 0 {
		return nil, fmt.Errorf("The Jenkins server url must not be empty")
	}

	// Create the Jenkins Instance
	jenkins, err := gojenkins.CreateJenkins(jenkinsServer, jenkinsUser, jenkinsPassword).Init()
	if err != nil {
		return nil, fmt.Errorf("Fail to create Jenkins instance as %s", err.Error())
	}

	cleaner := &JenkinsCleaner{jenkins}

	return cleaner, nil
}

func (jc *JenkinsCleaner) CleanSlaves(labelStrs ...string) error {
	if len(labelStrs) == 0 {
		return fmt.Errorf("No need to clean as no labels")
	}

	for _, labelStr := range labelStrs {
		label, err := jc.jenkins.GetLabel(labelStr)
		if err != nil {
			return fmt.Errorf("Fail to get label as %s", err.Error())
		}

		deletedNodes := []string{}
		for _, ln := range label.Raw.Nodes {
			if ln.Class == "org.csanchez.jenkins.plugins.kubernetes.KubernetesSlave" {
				nodeName := ln.NodeName
				node, err := jc.jenkins.GetNode(nodeName)
				if err != nil {
					return fmt.Errorf("Fail to get node %s as %s", nodeName, err.Error())
				}

				idle, err := node.IsIdle()
				if err != nil {
					return fmt.Errorf("Fail to judge whether node %s is idle as %s", nodeName, err.Error())
				}

				if idle {
					deleted, err := node.Delete()
					if err != nil {
						return fmt.Errorf("Fail to delete node %s as %s", nodeName, err.Error())
					}

					if !deleted {
						return fmt.Errorf("Fail to delete node %s", nodeName)
					}
					deletedNodes = append(deletedNodes, nodeName)
				}
			}
		}
		fmt.Printf("Succeed to delete the nodes for label %s: %v\n", labelStr, deletedNodes)
	}

	return nil
}
