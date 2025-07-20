//go:build !darwin && !windows && !linux
// +build !darwin,!windows,!linux

package permission

import (
	"context"
	"testing"
)

func TestStubManager_Check(t *testing.T) {
	manager := newStubManager()
	ctx := context.Background()
	
	perms, err := manager.Check(ctx)
	if err != nil {
		t.Errorf("Check() error = %v, want no error", err)
	}
	
	if len(perms) != 1 {
		t.Errorf("Check() returned %d permissions, want 1", len(perms))
	}
	
	if perms[0].Status != StatusNotApplicable {
		t.Errorf("Check() status = %v, want StatusNotApplicable", perms[0].Status)
	}
}

func TestStubManager_Request(t *testing.T) {
	manager := newStubManager()
	ctx := context.Background()
	
	err := manager.Request(ctx, PermissionTypeAccessibility)
	if err == nil {
		t.Error("Request() expected error but got nil")
	}
}

func TestStubManager_GetInstructions(t *testing.T) {
	manager := newStubManager()
	
	instructions := manager.GetInstructions(PermissionTypeAccessibility)
	if instructions == "" {
		t.Error("GetInstructions() returned empty string")
	}
}

func TestStubManager_OpenSettings(t *testing.T) {
	manager := newStubManager()
	
	err := manager.OpenSettings(PermissionTypeAccessibility)
	if err == nil {
		t.Error("OpenSettings() expected error but got nil")
	}
}

func TestStubManager_IsSupported(t *testing.T) {
	manager := newStubManager()
	
	supported := manager.IsSupported(PermissionTypeAccessibility)
	if supported {
		t.Error("IsSupported() returned true, want false")
	}
}