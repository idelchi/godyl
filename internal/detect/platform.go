package detect

import (
	"github.com/idelchi/godyl/internal/detect/platform"
)

// Platform encapsulates system-specific characteristics and capabilities.
type Platform struct {
	OS           platform.OS
	Library      platform.Library
	Distribution platform.Distribution
	Extension    platform.Extension
	Architecture platform.Architecture
}

// Parse extracts platform information from a string identifier.
// Attempts to determine OS, architecture, and library details from the input.
func (p *Platform) ParseFrom(name string) {
	p.OS.ParseFrom(name)
	p.Architecture.ParseFrom(name)
	p.Library.ParseFrom(name)
	p.Extension.ParseFrom(p.OS)
}

// Parse extracts platform information from the configured fields.
func (p *Platform) Parse() error {
	if err := p.OS.Parse(); err != nil {
		return err
	}

	if err := p.Architecture.Parse(); err != nil {
		return err
	}

	if err := p.Library.Parse(); err != nil {
		p.Library = p.Library.Default(p.OS, p.Distribution)
	}

	p.Extension.ParseFrom(p.OS)

	return nil
}

// Merge merges unset fields from other into p and reports whether anything changed.
func (p *Platform) Merge(other Platform) (changed bool) {
	if p.OS.IsNil() {
		p.OS, changed = other.OS, true
	}

	if p.Architecture.IsNil() {
		p.Architecture, changed = other.Architecture, true
	}

	if p.Library.IsNil() {
		p.Library, changed = other.Library, true
	}

	if p.Extension.IsNil() {
		p.Extension, changed = other.Extension, true
	}

	if p.Distribution.IsNil() {
		p.Distribution, changed = other.Distribution, true
	}

	return changed
}

// ToMap converts the Platform configuration into a map for templating.
// Includes derived values like architecture type, version, and capability flags.
func (p Platform) ToMap() map[string]any {
	platformMap := make(map[string]any)
	platformMap["OS"] = p.OS.String()
	platformMap["ARCH"] = p.Architecture.Type()
	platformMap["ARCH_VERSION"] = p.Architecture.Version()
	platformMap["ARCH_LONG"] = p.Architecture.String()
	platformMap["IS_ARM"] = p.Architecture.IsARM()
	platformMap["IS_X86"] = p.Architecture.IsX86()
	platformMap["LIBRARY"] = p.Library.String()
	platformMap["EXTENSION"] = p.Extension.String()
	platformMap["DISTRIBUTION"] = p.Distribution.String()

	return platformMap
}
