package tools

func (cmd *Command) Flags() {
	cmd.Command.Flags().BoolP("full", "f", false, "Show the tools in full syntax")
	cmd.Command.Flags().StringSliceP("tags", "t", []string{""}, "Tags to filter tools by. Prefix with '!' to exclude")
}
