package config

import "github.com/spf13/cobra"

func flags(cmd *cobra.Command) {
	cmd.Flags().BoolP("all", "a", false, "Dump all configuration values, including unset ones")
}
