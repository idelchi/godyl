package download

func (cmd *Command) Flags() {
	cmd.Command.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Command.Flags().String("os", "", "Operating system to install the tools for")
	cmd.Command.Flags().String("arch", "", "Architecture to install the tools for")

	cmd.Command.Flags().String("source", "github", "Source from which to install the tools (github, gitlab, url, go, none)")
	cmd.Command.Flags().StringSlice("hints", []string{""}, "Hints to use for tool resolution")
	cmd.Command.Flags().String("version", "", "Version of the tool to install. Empty means latest. Not so useful when downloading multiple tools.")
}
