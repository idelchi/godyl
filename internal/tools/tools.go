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
					Mode: Extract,
				},
			}
		} else {
			*t = Tools{
				Tool{
					Name: cfg,
					Mode: Extract,
				},
			}
		}
		return nil
	}

	file, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}
