package main

import (
	"errors"
	"strings"
)

type VMInfo struct {
	// Name of the VM set by vagrant (ex: mydir_default_1534347044260_6006)
	Name string

	// The VM provider (ex: virtualbox)
	Provider string
}

type UpResponse struct {
	// Info about all of the VMs constructed by the Vagrantfile. The map keys are
	// vagrant VM names (ex: default) and the values are VMInfo's.
	VMInfo map[string]*VMInfo

	// If set, there was an error while running vagrant up
	Error error
}

func newUpResponse() UpResponse {
	return UpResponse{
		VMInfo: make(map[string]*VMInfo),
	}
}

func (resp *UpResponse) handleOutput(target, key string, message []string) {
	// Only interested in the following output:
	// * target: X, key: metadata, message: [provider, Y]
	// * target: X, key: ui, message: [_, "target: Setting the name of the VM: Y"]
	// * key: error-exit, message: X
	if target != "" && len(message) == 2 {
		info, exists := resp.VMInfo[target]
		if !exists {
			info = &VMInfo{}
			resp.VMInfo[target] = info
		}

		if key == "metadata" && message[0] == "provider" {
			info.Provider = message[1]
		} else if key == "ui" && strings.Contains(message[1], "Setting the name of the VM:") {
			idx := strings.LastIndex(message[1], ":")
			info.Name = strings.TrimSpace(message[1][idx+1:])
		}
	} else if key == "error-exit" {
		resp.Error = errors.New(strings.Join(message, ", "))
	}
}