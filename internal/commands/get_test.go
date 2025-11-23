package commands

import (
	"strings"
	"testing"
)

func TestGetCmdUsage(t *testing.T) {
	cmd := NewGetCmd()

	// Verify command basic properties
	if cmd.Use != "get <recording-id|latest>" {
		t.Errorf("expected Use to be 'get <recording-id|latest>', got %q", cmd.Use)
	}

	if !strings.Contains(cmd.Long, "latest") {
		t.Error("expected Long description to mention 'latest'")
	}

	// Verify flags exist
	jsonFlag := cmd.Flags().Lookup("json")
	if jsonFlag == nil {
		t.Error("expected --json flag to exist")
	}

	clubFlag := cmd.Flags().Lookup("club")
	if clubFlag == nil {
		t.Error("expected --club flag to exist")
	}
}
