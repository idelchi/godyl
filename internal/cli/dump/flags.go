package dump

func (cmd *Command) Flags() {
	cmd.Command.Flags().StringP("format", "f", "yaml", "Output format (json or yaml)")
}
