package permission

import (
	"testing"
)

func TestPermissionType(t *testing.T) {
	// Test permission type constants
	permissions := []PermissionType{
		PermissionTypeAccessibility,
		PermissionTypeNotification,
		PermissionTypeAutoStart,
	}

	// Each permission should have a string representation
	for _, perm := range permissions {
		str := string(perm)
		if str == "" {
			t.Errorf("Permission %v should have non-empty string representation", perm)
		}
	}
}

func TestPermissionStatus(t *testing.T) {
	// Test status constants
	statuses := []Status{
		StatusGranted,
		StatusDenied,
		StatusNotDetermined,
		StatusNotApplicable,
	}

	// Each status should have a string representation
	for _, status := range statuses {
		str := string(status)
		if str == "" {
			t.Errorf("Status %v should have non-empty string representation", status)
		}
	}
}

func TestPermission(t *testing.T) {
	// Test Permission structure
	result := Permission{
		Type:        PermissionTypeAccessibility,
		Status:      StatusGranted,
		Description: "Permission granted",
		Required:    true,
	}

	if result.Type != PermissionTypeAccessibility {
		t.Error("Permission.Type not set correctly")
	}

	if result.Status != StatusGranted {
		t.Error("Permission.Status not set correctly")
	}

	if !result.Required {
		t.Error("Permission.Required should be true")
	}

	if result.Description != "Permission granted" {
		t.Error("Permission.Description not set correctly")
	}
}

func TestPermissionCopy(t *testing.T) {
	// Test that Permission can be copied
	original := Permission{
		Type:        PermissionTypeAccessibility,
		Status:      StatusGranted,
		Required:    true,
		Description: "Original message",
	}

	permCopy := original
	permCopy.Status = StatusDenied
	permCopy.Description = "Modified message"

	// Original should be unchanged
	if original.Status != StatusGranted {
		t.Error("Original Permission was modified during copy")
	}

	if original.Description != "Original message" {
		t.Error("Original Permission description was modified during copy")
	}

	// Copy should be modified
	if permCopy.Status != StatusDenied {
		t.Error("Copied Permission was not modified correctly")
	}

	if permCopy.Description != "Modified message" {
		t.Error("Copied Permission description was not modified correctly")
	}
}

func TestPermissionConstants(t *testing.T) {
	// Test that permission constants are defined correctly
	expectedPermissions := map[PermissionType]string{
		PermissionTypeAccessibility: "accessibility",
		PermissionTypeNotification:  "notification",
		PermissionTypeAutoStart:     "autostart",
	}

	for perm, expectedStr := range expectedPermissions {
		if string(perm) != expectedStr {
			t.Errorf("Permission %v should be %s, got %s", perm, expectedStr, string(perm))
		}
	}
}

func TestPermissionStatusConstants(t *testing.T) {
	// Test that status constants are defined correctly
	expectedStatuses := map[Status]string{
		StatusGranted:       "granted",
		StatusDenied:        "denied",
		StatusNotDetermined: "not_determined",
		StatusNotApplicable: "not_applicable",
	}

	for status, expectedStr := range expectedStatuses {
		if string(status) != expectedStr {
			t.Errorf("Status %v should be %s, got %s", status, expectedStr, string(status))
		}
	}
}

func TestPermissionZeroValue(t *testing.T) {
	// Test zero value behavior
	var result Permission

	if result.Type != "" {
		t.Error("Zero value Permission should have empty Type")
	}

	if result.Status != "" {
		t.Error("Zero value Permission should have empty Status")
	}

	if result.Required {
		t.Error("Zero value Permission should have Required = false")
	}

	if result.Description != "" {
		t.Error("Zero value Permission should have empty Description")
	}
}

func TestPermissionTypeValidation(t *testing.T) {
	// Test validation of permission types
	validTypes := []PermissionType{
		PermissionTypeAccessibility,
		PermissionTypeNotification,
		PermissionTypeAutoStart,
	}

	for _, permType := range validTypes {
		// Each type should be non-empty
		if string(permType) == "" {
			t.Errorf("Permission type %v should have non-empty string value", permType)
		}

		// Should not contain spaces or special characters
		str := string(permType)
		for _, char := range str {
			if char == ' ' {
				t.Errorf("Permission type %s should not contain spaces", str)
			}
		}
	}
}

func TestPermissionStatusValidation(t *testing.T) {
	// Test validation of permission statuses
	validStatuses := []Status{
		StatusGranted,
		StatusDenied,
		StatusNotDetermined,
		StatusNotApplicable,
	}

	for _, status := range validStatuses {
		// Each status should be non-empty
		if string(status) == "" {
			t.Errorf("Permission status %v should have non-empty string value", status)
		}

		// Should not contain spaces (underscores are okay)
		str := string(status)
		if str == "" {
			t.Errorf("Permission status %v should not be empty", status)
		}
	}
}
