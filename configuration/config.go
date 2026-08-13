package configuration

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	ID         string `yaml:"id"`
	Role       string `yaml:"role"`
	Listen     string `yaml:"listen"`
	GRPCListen string `yaml:"grpc_listen"`
	MainAddr   string `yaml:"main_addr,omitempty"`
}

type Config struct {
	Servers           []Server `yaml:"server"`
	ReplicationFactor int      `yaml:"replication_factor"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if len(cfg.Servers) == 0 {
		return nil, fmt.Errorf("config %s: no servers defined", path)
	}
	if cfg.ReplicationFactor < 1 {
		return nil, fmt.Errorf("config %s: replication_factor must be >= 1", path)
	}
	if len(cfg.Servers) < cfg.ReplicationFactor {
		return nil, fmt.Errorf("config %s: replication_factor %d exceeds number of servers %d", path, cfg.ReplicationFactor, len(cfg.Servers))
	}

	return &cfg, nil
}

func (c *Config) FindByID(id string) (*Server, error) {
	for i := range c.Servers {
		if c.Servers[i].ID == id {
			return &c.Servers[i], nil
		}
	}
	return nil, fmt.Errorf("server %q not found in config", id)
}

func (c *Config) ByRole(role string) []Server {
	var out []Server
	for i := range c.Servers {
		if c.Servers[i].Role == role {
			out = append(out, c.Servers[i])
		}
	}
	return out
}
