package config

import "github.com/spf13/cobra"

func Flags(cmd *cobra.Command) {
	cmd.Flags().BoolP("full", "f", false, "Dump all configuration values, including unset ones")
}
