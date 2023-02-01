package logic

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func (c *Client) RunCmdPluginList() error {
	plugins := c.Plugins()

	if len(plugins) > 0 {
		log.Println("The following plugins are installed:")
	} else {
		log.Println("No plugins are installed.")
	}

	for _, pluginPath := range plugins {
		//nolint:forbidigo
		fmt.Printf("- %s\n", filepath.Base(pluginPath))
	}

	return nil
}

func (c *Client) RunCmdPluginListAvailable() error {
	plugins := c.PluginsAvailable()

	if len(plugins) > 0 {
		log.Println("The following plugins are available online:")
	} else {
		log.Println("No plugins are available online.")
	}

	for plugin := range plugins {
		//nolint:forbidigo
		fmt.Printf("- %s\n", filepath.Base(plugin))
	}

	return nil
}

func (c *Client) RunCmdPluginInstall(args []string) error {
	// TODO: implement.
	for _, plugin := range args {
		log.Printf("Installing plugin %s...", plugin)
	}

	log.Print("...plugin installed.")

	return nil
}

func (c *Client) RunCmdPluginUpdate(args []string) error {
	// TODO: implement.
	for _, plugin := range args {
		log.Printf("Updating plugin %s...", plugin)
	}

	log.Print("...plugins updated.")

	return nil
}
