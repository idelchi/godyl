package status

// Status adds status-related command flags to the provided cobra command.
func (cmd *Command) Flags() {
	cmd.Command.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
}
