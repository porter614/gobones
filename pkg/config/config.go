package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/porter614/gobones/pkg/watch"
	"github.com/porter614/logger"
)

type Configurable interface {
	Configure(interface{}) error
}

type CommonConfig struct {
	Log          logger.LogConfig `mapstructure:"log"`
	PollInterval int              `mapstructure:"poll-interval"`
}

type Config struct {
	Common CommonConfig `mapstructure:"common"`
}

type Configurator struct {
	Log  *logrus.Entry
	conf Config
}

// LoadConfig performs one-time viper setup then calls loadConfig,
// the package internal function that reads configuration
func (c *Configurator) LoadConfig(env, file string) error {

	// Get configuration from the environment
	viper.SetEnvPrefix(env)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Config file
	viper.SetConfigName(file)
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")

	return c.loadConfig()
}

// loadConfig merges configuration from the environment (with the specified
// prefix) and a config file, then marshals it into a Config struct instance
func (c *Configurator) loadConfig() error {
	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Unmarshal config
	if err := viper.Unmarshal(&c.conf); err != nil {
		return err
	}

	// Update all app subcomponent configs
	logger.Configure(&c.conf.Common.Log)
	/* Configure your other subcomponents here */

	return nil
}

// WatchConfig creates and starts a file watcher and then launches a separate
// goroutine to listen for changes to the config file
func (c *Configurator) WatchConfig() error {
	cfgch, errch, err := watch.WatchFile("config/config.json")
	if err != nil {
		return err
	}
	go c.watchConfig(cfgch, errch)
	return nil
}

// watchConfig blocks listening on the provided channel
func (c *Configurator) watchConfig(ch chan []byte, errch chan error) {
	for {
		select {
		case <-ch:
			c.Log.Debug("Config change detected")
			err := c.loadConfig()
			if err != nil {
				// TODO: Something...
				c.Log.Errorf("Failed to load config: %v", err)
				continue
			}
		case err := <-errch:
			c.Log.Errorf("Error in config file watcher: %v", err)
		}
	}
}
