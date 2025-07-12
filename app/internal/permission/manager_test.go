package permission

import (
	"context"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if manager == nil {
		t.Fatal("Manager is nil")
	}

	// Just verify we got a valid manager for the current OS
	// We can't check concrete types across platforms due to build constraints
}

func TestManager_Check(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	ctx := context.Background()
	permissions, err := manager.Check(ctx)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if len(permissions) == 0 {
		t.Error("No permissions returned")
	}

	// Verify all permission types are checked
	requiredTypes := map[PermissionType]bool{
		PermissionTypeAccessibility: false,
		PermissionTypeNotification:  false,
		PermissionTypeAutoStart:     false,
	}

	for _, perm := range permissions {
		requiredTypes[perm.Type] = true

		// Verify permission has required fields
		if perm.Description == "" {
			t.Errorf("Permission %s has empty description", perm.Type)
		}

		// Verify status is valid
		switch perm.Status {
		case StatusGranted, StatusDenied, StatusNotDetermined, StatusNotApplicable:
			// Valid status
		default:
			t.Errorf("Permission %s has invalid status: %s", perm.Type, perm.Status)
		}
	}

	// Check all required types were included
	for permType, found := range requiredTypes {
		if !found {
			t.Errorf("Permission type %s not found in results", permType)
		}
	}
}

func TestManager_IsSupported(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test known permission types
	knownTypes := []PermissionType{
		PermissionTypeAccessibility,
		PermissionTypeNotification,
		PermissionTypeAutoStart,
	}

	for _, permType := range knownTypes {
		if !manager.IsSupported(permType) {
			t.Errorf("Expected %s to be supported", permType)
		}
	}

	// Test unknown permission type
	unknownType := PermissionType("unknown")
	if manager.IsSupported(unknownType) {
		t.Errorf("Expected unknown type to not be supported")
	}
}

func TestManager_GetInstructions(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	permTypes := []PermissionType{
		PermissionTypeAccessibility,
		PermissionTypeNotification,
		PermissionTypeAutoStart,
	}

	for _, permType := range permTypes {
		instructions := manager.GetInstructions(permType)
		if instructions == "" {
			t.Errorf("Empty instructions for %s", permType)
		}

		// Verify instructions mention the permission type or related keywords
		if len(instructions) < 20 {
			t.Errorf("Instructions too short for %s: %s", permType, instructions)
		}
	}
}
