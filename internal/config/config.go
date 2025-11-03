package config

import (
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Keep int `yaml:"keep"`
    Prefix string `yaml:"prefix"`
    Remote string `yaml:"remote"`
    OlderThan string `yaml:"older-than"`
}

// LoadConfig attempts to read a YAML config from the given path. If path is empty,
// it looks for ~/.branch-pruner.yaml and ./.branch-pruner.yaml (in that order).
func LoadConfig(path string) (*Config, error) {
    if path != "" {
        return loadFromFile(path)
    }
    // try home
    home, err := os.UserHomeDir()
    if err == nil {
        p := home + string(os.PathSeparator) + ".branch-pruner.yaml"
        if _, err := os.Stat(p); err == nil {
            return loadFromFile(p)
        }
    }
    // try cwd
    if _, err := os.Stat(".branch-pruner.yaml"); err == nil {
        return loadFromFile(".branch-pruner.yaml")
    }
    // no config
    return &Config{}, nil
}

func loadFromFile(path string) (*Config, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var cfg Config
    if err := yaml.Unmarshal(b, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
