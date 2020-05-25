package main

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

const minimumServerVersion = "5.22.0"

func (p *Plugin) checkServerVersion() error {
	serverVersion, err := semver.Parse(p.API.GetServerVersion())
	if err != nil {
		return errors.Wrap(err, "failed to parse server version")
	}

	r := semver.MustParseRange(">=" + minimumServerVersion)
	if !r(serverVersion) {
		return fmt.Errorf("This plugin requires Mattermost v%s or later", minimumServerVersion)
	}

	return nil
}

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account
func (p *Plugin) OnActivate() error {
	if err := p.checkServerVersion(); err != nil {
		return err
	}

	if ok, err := p.checkRequiredServerConfiguration(); err != nil {
		return errors.Wrap(err, "could not check required server configuration")
	} else if !ok {
		p.API.LogError("Server configuration is not compatible")
	}

	if err := p.IsValid(); err != nil {
		return err
	}

	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to use
// the API, and the plugin will be terminated shortly after this invocation.
//
// This demo implementation logs a message to the demo channel whenever the plugin is deactivated.
// func (p *Plugin) OnDeactivate() error {
// 	configuration := p.getConfiguration()

// 	return nil
// }

func (p *Plugin) checkRequiredServerConfiguration() (bool, error) {
	return p.Helpers.CheckRequiredServerConfiguration(manifest.RequiredConfig)
}
