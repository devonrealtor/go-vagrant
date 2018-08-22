package main

type VersionCommand struct {
	BaseCommand
	VersionResponse
}

// Run vagrant version. After setting options as appropriate, you must call
// Run() or Start() followed by Wait() to execute. Errors will be recorded in
// Error.
func (client *VagrantClient) Version() *VersionCommand {
	return &VersionCommand{
		BaseCommand:     newBaseCommand(client),
		VersionResponse: newVersionResponse(),
	}
}

func (cmd *VersionCommand) init() error {
	return cmd.BaseCommand.init(&cmd.VersionResponse, "version")
}

// Run the command
func (cmd *VersionCommand) Run() error {
	if err := cmd.init(); err != nil {
		return err
	}
	return cmd.BaseCommand.Run()
}

// Start the command. You must call Wait() to complete execution.
func (cmd *VersionCommand) Start() error {
	if err := cmd.init(); err != nil {
		return err
	}
	return cmd.BaseCommand.Start()
}