package main

type StatusCommand struct {
	BaseCommand
	StatusResponse
}

// Run vagrant status. After setting options as appropriate, you must call
// Run() or Start() followed by Wait() to execute. Errors will be recorded in
// Error.
func (client *VagrantClient) Status() *StatusCommand {
	return &StatusCommand{
		BaseCommand:    newBaseCommand(client),
		StatusResponse: newStatusResponse(),
	}
}

func (cmd *StatusCommand) init() error {
	return cmd.BaseCommand.init(&cmd.StatusResponse, "status")
}

// Run the command
func (cmd *StatusCommand) Run() error {
	if err := cmd.init(); err != nil {
		return err
	}
	return cmd.BaseCommand.Run()
}

// Start the command. You must call Wait() to complete execution.
func (cmd *StatusCommand) Start() error {
	if err := cmd.init(); err != nil {
		return err
	}
	return cmd.BaseCommand.Start()
}