package vagrant

import (
	"strings"
)

// UpCommand specifies options and output from vagrant up.
type UpCommand struct {
	BaseCommand
	UpResponse

	// Enable or disable provisioning (default: enabled)
	Provisioning bool

	// Enabled provisioners by type or name (default: blank which means they're
	// all enable or disabled depending on the Provisioning flag)
	Provisioners []string

	// Destroy on error (default: true)
	DestroyOnError bool

	// Enable parallel execution if the provider supports it (default: false)
	Parallel bool

	// Provider to use (default: blank which means vagrant will use the default
	// provider)
	Provider string

	// Install the provider if it isn't installed, if possible (default: false)
	InstallProvider bool
}

// Run vagrant up. After setting options as appropriate, you must call Run()
// or Start() followed by Wait() to execute. Output will be in VMInfo or Error.
func (client *VagrantClient) Up() *UpCommand {
	return &UpCommand{
		BaseCommand:     newBaseCommand(client),
		UpResponse:      newUpResponse(),
		Provisioning:    true,
		DestroyOnError:  true,
		InstallProvider: true,
	}
}

func (cmd *UpCommand) buildArguments() []string {
	args := []string{}
	if !cmd.Provisioning {
		args = append(args, "--no-provision")
	}
	if cmd.Provisioners != nil && len(cmd.Provisioners) > 0 {
		args = append(args, "--provision-with", strings.Join(cmd.Provisioners, ","))
	}
	if !cmd.DestroyOnError {
		args = append(args, "--no-destroy-on-error")
	}
	if cmd.Parallel {
		args = append(args, "--parallel")
	}
	if len(cmd.Provider) > 0 {
		args = append(args, "--provider", cmd.Provider)
	}
	if !cmd.InstallProvider {
		args = append(args, "--no-install-provider")
	}
	return args
}

func (cmd *UpCommand) init() error {
	args := cmd.buildArguments()
	return cmd.BaseCommand.init(&cmd.UpResponse, "up", args...)
}

// Run the command
func (cmd *UpCommand) Run() error {
	if err := cmd.init(); err != nil {
		return err
	}
	return cmd.BaseCommand.Run()
}

// Start the command. You must call Wait() to complete execution.
func (cmd *UpCommand) Start() error {
	if err := cmd.init(); err != nil {
		return err
	}
	return cmd.BaseCommand.Start()
}
