package config

type Dump struct {
	// Format to dump the configuration in
	Format string `validate:"oneof=json yaml"`
}
