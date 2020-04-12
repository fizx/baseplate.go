package baseplate

import (
	"errors"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/reddit/baseplate.go/log"
)

// ServerConfig is a general purpose config for assembling a BaseplateServer
type ServerConfig struct {
	Addr string

	Timeout time.Duration

	Log struct {
		Level log.Level
	}

	Metrics struct {
		Namespace string
		Endpoint  string
	}

	Secrets struct {
		Path string
	}

	Sentry struct {
		DSN         string
		Environment string
		SampleRate  float64
	}

	Tracing struct {
		Namespace     string
		Endpoint      string
		RecordTimeout time.Duration
		SampleRate    float64
	}
}

// Server is the primary interface for baseplate servers.
type Server interface {
	Config() ServerConfig
	Serve() error
	Stop() error
}

// ParseServerConfig will populate a ServerConfig from a YAML file.
func ParseServerConfig(path string, cfg *ServerConfig) error {
	if path == "" {
		return errors.New("no config path given")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	log.Debugf("%#v", cfg)
	return nil
}
