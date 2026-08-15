package config

import "errors"

var AppConfig *Config

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Log      LoggingConfig
}

func Load() (*Config, error) {
	cfg := &Config{
		Database: loadDatabaseConfig(),
		Server:   loadServerConfig(),
		Log:      loadLoggingConfig(),
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	AppConfig = cfg
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return errors.New("SERVER_PORT is required")
	}
	if c.Server.Mode == "" {
		return errors.New("SERVER_MODE is required")
	}
	if c.Server.Env == "" {
		return errors.New("SERVER_ENV is required")
	}
	if c.Database.Host == "" {
		return errors.New("DB_HOST is required")
	}
	if c.Database.Port == "" {
		return errors.New("DB_PORT is required")
	}
	if c.Database.User == "" {
		return errors.New("DB_USER is required")
	}
	if c.Database.Password == "" {
		return errors.New("DB_PASSWORD is required")
	}
	if c.Database.Name == "" {
		return errors.New("DB_NAME is required")
	}
	if c.Log.Level == "" {
		return errors.New("LOG_LEVEL is required")
	}
	if c.Log.Format == "" {
		return errors.New("LOG_FORMAT is required")
	}
	if c.Server.JwtSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	return nil
}
