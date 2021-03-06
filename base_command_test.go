package vagrant

import (
	"io"
	"testing"
)

// Use this method in tests to cleanup a BaseCommand that was init()'d but not
// executed. init() creates a goroutine that needs to be cleaned up.
func (b *BaseCommand) cleanup() {
	if b.cmd == nil {
		return
	}

	// closing Stdout will cause the goroutine to exit
	b.cmd.Stdout.(io.Closer).Close()

	// wait for the goroutine to exit - should be immediate
	b.readers.Wait()

	// cleanup
	b.cmd = nil
	b.readers = nil
}

func TestBaseCommand_init(t *testing.T) {
	client := newMockVagrantClient()
	handler := newMockOutputHandler()

	t.Run("base", func(t *testing.T) {
		cmd := newBaseCommand(client)
		cmd.Env = []string{"ENV1=value1"}
		cmd.init(handler, "test", "arg1")
		defer cmd.cleanup()

		if cmd.cmd.Path != client.executable {
			t.Errorf("Expected cmd path to be %v; got %v", client.executable, cmd.cmd.Path)
		}

		argsLength := len(cmd.cmd.Args)
		if argsLength < 6 {
			t.Errorf("Expected len(args) to be at least 6; got %v", argsLength)
		} else {
			if cmd.cmd.Args[0] != client.executable {
				t.Errorf("Expected args[0] to be %v; got %v", client.executable, cmd.cmd.Args[0])
			}
			if cmd.cmd.Args[argsLength-3] != "test" {
				t.Errorf("Expected args[1] to be 'test'; got %v", cmd.cmd.Args[argsLength-3])
			}
			if cmd.cmd.Args[argsLength-2] != "--machine-readable" {
				t.Errorf("Expected args[1] to be '--machine-readable'; got %v", cmd.cmd.Args[argsLength-2])
			}
			if cmd.cmd.Args[argsLength-1] != "arg1" {
				t.Errorf("Expected args[2] to be 'arg1'; got %v", cmd.cmd.Args[argsLength-1])
			}
		}

		if cmd.cmd.Env == nil || len(cmd.cmd.Env) == 0 {
			t.Errorf("Expected cmd.Env to be set")
		} else if cmd.cmd.Env[len(cmd.cmd.Env)-1] != "ENV1=value1" {
			t.Errorf("Expected env[-1] to be 'ENV1=value1'; got %v", cmd.cmd.Env[len(cmd.cmd.Env)-1])
		}

		if cmd.cmd.Dir != client.VagrantfileDir {
			t.Errorf("Expected Dir to be %v; got %v", client.VagrantfileDir, cmd.cmd.Dir)
		}
	})

	t.Run("alreadyStarted", func(t *testing.T) {
		cmd := newBaseCommand(client)
		cmd.init(handler, "test")
		defer cmd.cleanup()

		if err := cmd.init(handler, "test"); err == nil {
			t.Errorf("Expected init to return an error on second call.")
		}
	})

	t.Run("additionalArguments", func(t *testing.T) {
		cmd := newBaseCommand(client)
		cmd.AdditionalArgs = []string{"--myarg"}
		cmd.init(handler, "test")
		defer cmd.cleanup()

		argsLength := len(cmd.cmd.Args)
		if cmd.cmd.Args[argsLength-1] != "--myarg" {
			t.Errorf("Expected last arg to be '--myarg'; got %v", cmd.cmd.Args[argsLength-1])
		}
	})
}

func TestBaseCommand_Run(t *testing.T) {
	client := newMockVagrantClient()
	handler := newMockOutputHandler()
	cmd := newBaseCommand(client)
	cmd.init(handler, "test", "arg1")

	if err := cmd.Run(); err != nil {
		t.Fatalf("Error: %v", err)
	}

	if cmd.Process == nil {
		t.Errorf("Expected Process to be set.")
	}

	if cmd.ProcessState == nil {
		t.Errorf("Expected ProcessState to be set.")
	}

	if handler.subcommand != "test" {
		t.Errorf("Expected subcommand 'test'; got %v", handler.subcommand)
	}
	if !handler.machineReadable {
		t.Errorf("Expected machine-readable")
	}

	if len(handler.args) != 1 {
		t.Fatalf("Expected 1 arg; got %v", len(handler.args))
	}
	if handler.args[0] != "arg1" {
		t.Errorf("Expected arg 1 to be 'arg1'; got %v", handler.args[0])
	}
}
