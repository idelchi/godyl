package core

import (
	"fmt"

	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/defaults"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Handler manages the common operations for CLI commands including configuration resolution,
// logging, and embedded resource management.
// Handler manages configuration resolution, logging, and embedded resources for CLI operations.
type Handler struct {
	config   root.Config
	logger   *logger.Logger
	embedded Embedded
}

// NewHandler creates a new Handler instance with the provided configuration and embedded resources.
// NewHandler creates a new Handler instance with the provided configuration and embedded resources.
func NewHandler(cfg root.Config, embedded Embedded) *Handler {
	return &Handler{
		config:   cfg,
		embedded: embedded,
	}
}

// Resolve loads and processes the configuration for tools, including defaults, inheritance, and platform-specific
// settings.
// Resolve processes the configuration hierarchy by loading defaults, applying inheritance,
// and merging platform-specific settings for all tools.
func (c *Handler) Resolve(defaultFile file.File, tools *tools.Tools) (err error) {
	/* Steps:

	1. Load the defaults file
	2. Create a defaults map to hold the defaults file content as map[string]*Tool
	3. Resolve the inheritance scheme of the defaults
	4. Retrieve the common configuration as a tool, force setting
	5. Merge all defaults with this configuration, to replace any missing values
	6. Get all the flags that were explicitly set by the user as a tool
	7. Merge all defaults with this configuration, to force replace any values set
	8. Assign default inheritance to all tools
	9. Resolve the inheritance scheme of all the tools. Important is to do the final
	   merge with the defaults as a "UnmarshalYAML", to have the custom unmarshalling mechanisms kick in.
	10. Finally, ensure that no nil points are left in the tools
	11. Now we can merge the platform settings
	12. Finally, we can merge the platform settings
	*/
	// Continue with setting up the defaults
	defaultMap := c.embedded.Defaults

	// Attempt to load the defaults file
	if defaultFile != "" && defaultFile.Exists() {
		if defaultMap, err = defaultFile.Read(); err != nil {
			return fmt.Errorf("reading defaults file: %w", err)
		}
	}

	// Create a defaults map to hold the defaults file content as map[string]*Tool
	defs, err := defaults.NewDefaultsFromBytes(defaultMap)
	if err != nil {
		return fmt.Errorf("loading defaults: %w", err)
	}

	// Resolve the inheritance scheme of the defaults
	if err := defs.ResolveInheritance(); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Retrieve the common configuration as a tool, force setting
	// to get all the default values from the current flags
	toolFromFlag := c.config.ToTool(true)

	// Merge all defaults with this configuration, to replace any missing values
	if err := defs.MergeWith(toolFromFlag); err != nil {
		return fmt.Errorf("merging defaults with configuration: %w", err)
	}

	// Get all the flags that were explicitly set by the user as a tool
	toolFromFlag = c.config.ToTool(false)

	// Merge all defaults with this configuration, to force replace any values set
	if err := defs.MergeFrom(toolFromFlag); err != nil {
		return fmt.Errorf("merging defaults with forced configuration: %w", err)
	}

	// Assign default inheritance to all tools
	tools.DefaultInheritance(c.config.Inherit)

	// We can now resolve the inheritance scheme of all the tools.
	if err := tools.ResolveInheritance(defs); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Finally, ensure that no nil points are left in the tools
	if err := tools.ResolveNilPointers(); err != nil {
		return fmt.Errorf("resolving nil pointers: %w", err)
	}

	// Now we can merge the platform settings
	if err := tools.MergePlatform(); err != nil {
		return fmt.Errorf("merging platform: %w", err)
	}

	return nil
}

// Logger returns the logger instance associated with the handler.
func (c *Handler) Logger() *logger.Logger {
	return c.logger
}

// SetupLogger creates and configures a logger with the specified log level.
func (c *Handler) SetupLogger(level string) (err error) {
	c.logger, err = SetupLogger(level)

	return err
}

// SetupLogger creates a new logger instance with the specified log level.
func SetupLogger(level string) (*logger.Logger, error) {
	// Retrieve log level and create a logger instance
	lvl, err := logger.LevelString(level)
	if err != nil {
		return nil, fmt.Errorf("parsing log level: %w", err)
	}

	l, err := logger.New(lvl)
	if err != nil {
		return nil, fmt.Errorf("creating logger: %w", err)
	}

	return l, nil
}
