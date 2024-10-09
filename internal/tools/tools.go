package tools

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Tools []Tool

func (t *Tools) Load(cfg string) (err error) {
	if !strings.HasSuffix(cfg, ".yml") && !strings.HasSuffix(cfg, ".yaml") {
		if strings.HasPrefix(cfg, "http") {
			*t = Tools{
				Tool{
					Name: filepath.Base(cfg),
					Path: cfg,
				},
			}
		} else {
			*t = Tools{
				Tool{
					Name: cfg,
				},
			}
		}
		return nil
	}

	file, err := os.Open(cfg)
	if err != nil {
		return err
	}

	dec := yaml.NewDecoder(file)
	dec.KnownFields(true)
	if err := dec.Decode(t); err != nil {
		return err
	}

	return nil
}
