package main

import (
	"fmt"

	"github.com/pkg/errors"
)

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
type configuration struct {
	// AWS access key
	AWSAccessKeyID string

	// AWS secret key
	AWSSecretAccessKey string

	// AWS region with "us-east-1" as default
	AWSRegion string

	// disable plugin
	disabled bool
}

// Clone deep copies the configuration. Your implementation may only require a shallow copy if
// your configuration has no reference types.
func (c *configuration) Clone() *configuration {
	return &configuration{
		AWSAccessKeyID:     c.AWSAccessKeyID,
		AWSSecretAccessKey: c.AWSSecretAccessKey,
		AWSRegion:          c.AWSRegion,
		disabled:           c.disabled,
	}
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

func (p *Plugin) diffConfiguration(newConfiguration *configuration) {
	oldConfiguration := p.getConfiguration()
	configurationDiff := make(map[string]interface{})

	if newConfiguration.AWSAccessKeyID != oldConfiguration.AWSAccessKeyID {
		configurationDiff["aws_access_key_id"] = newConfiguration.AWSAccessKeyID
	}
	if newConfiguration.AWSSecretAccessKey != oldConfiguration.AWSSecretAccessKey {
		configurationDiff["aws_secret_access_key"] = newConfiguration.AWSSecretAccessKey
	}
	if newConfiguration.AWSRegion != oldConfiguration.AWSRegion {
		configurationDiff["aws_region"] = newConfiguration.AWSRegion
	}
	if newConfiguration.disabled != oldConfiguration.disabled {
		configurationDiff["disabled"] = newConfiguration.disabled
	}
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	configuration := p.getConfiguration().Clone()

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	p.diffConfiguration(configuration)

	p.setConfiguration(configuration)

	return nil
}

// setEnabled wraps setConfiguration to configure if the plugin is enabled.
func (p *Plugin) setEnabled(enabled bool) {
	var configuration = p.getConfiguration().Clone()
	configuration.disabled = !enabled

	p.setConfiguration(configuration)
}

// IsValid validates plugin configuration
func (p *Plugin) IsValid() error {
	configuration := p.getConfiguration()
	if configuration.AWSAccessKeyID == "" {
		return fmt.Errorf("Must have AWS Access Key ID")
	}

	if configuration.AWSSecretAccessKey == "" {
		return fmt.Errorf("Must have AWS Secret Access Key")
	}

	if configuration.AWSRegion == "" {
		configuration.AWSRegion = "us-east-1"
	}

	return nil
}
