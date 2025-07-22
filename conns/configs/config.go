package configs

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"revonoir.com/jform/controllers/dto"
)

type Configuration struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Dbname   string `mapstructure:"dbname"`
		MaxConns int    `mapstructure:"max_conns"`
		Sslmode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
}

type ConfigManager struct {
	configName   string
	configLock   sync.RWMutex
	config       Configuration
	lastLoadTime time.Time
	cacheTTL     time.Duration
}

func NewConfigManager(configName string) (*ConfigManager, error) {
	configManager := ConfigManager{
		configName:   configName,
		lastLoadTime: time.Now().Add(-10 * time.Minute),
		cacheTTL:     5 * time.Minute,
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/APP/revonoir/jform/")
	viper.AddConfigPath("./resources/config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, dto.NewError(http.StatusInternalServerError, dto.ErrorTitleConfig, "failed to read config: "+err.Error())
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		slog.Info("config file changed", "event", in)
		if in.Op == fsnotify.Write {
			configManager.lastLoadTime = time.Now().Add(-10*time.Minute - configManager.cacheTTL) // Force the next read to read from file
		}
	})

	return &configManager, nil
}

/*
Get config from the configuration file
Data will be cached for 5 minutes for efficiency
*/
func (c *ConfigManager) GetConfig() (Configuration, error) {
	c.configLock.RLock()
	if time.Since(c.lastLoadTime) < c.cacheTTL {
		defer c.configLock.RUnlock()
		return c.config, nil
	}
	c.configLock.RUnlock()

	c.configLock.Lock()
	defer c.configLock.Unlock()

	// Check again in case another goroutine refreshed the config
	if time.Since(c.lastLoadTime) < c.cacheTTL {
		return c.config, nil
	}

	if err := viper.Unmarshal(&c.config); err != nil {
		slog.Error("failed to unmarshal config", "error", err)
		return c.config, dto.NewError(http.StatusInternalServerError, dto.ErrorTitleConfig, "failed to parse config: "+err.Error())
	}
	c.lastLoadTime = time.Now()

	return c.config, nil
}
