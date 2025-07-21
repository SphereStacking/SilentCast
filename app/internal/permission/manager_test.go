package permission

import (
	"context"
	"strings"
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

	// Test unknown permission type
	unknownType := PermissionType("unknown")
	instructions := manager.GetInstructions(unknownType)
	if instructions == "" {
		t.Errorf("Expected some instructions for unknown type")
	}
}

func TestManager_Request(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test requesting permissions (these may not grant actual permissions in test environment)
	permTypes := []PermissionType{
		PermissionTypeAccessibility,
		PermissionTypeNotification,
		PermissionTypeAutoStart,
	}

	ctx := context.Background()
	for _, permType := range permTypes {
		// Request shouldn't panic or crash
		reqErr := manager.Request(ctx, permType)
		// We don't assert error == nil because requesting permissions
		// may fail in test environments or when already granted
		if reqErr != nil {
			t.Logf("Request for %s returned error (expected in test): %v", permType, reqErr)
		}
	}

	// Test unknown permission type - stub manager doesn't validate types
	unknownType := PermissionType("unknown")
	err = manager.Request(ctx, unknownType)
	// Stub manager always returns error, so this is expected
	if err == nil {
		t.Logf("Request for unknown type succeeded (platform-dependent behavior)")
	}
}

func TestPermissionStatus_String(t *testing.T) {
	statuses := []Status{
		StatusGranted,
		StatusDenied,
		StatusNotDetermined,
		StatusNotApplicable,
	}

	for _, status := range statuses {
		str := status.String()
		if str == "" {
			t.Errorf("Status %v has empty string representation", status)
		}
	}

	// Test unknown status
	unknownStatus := Status("unknown")
	str := unknownStatus.String()
	if str != "unknown" {
		t.Errorf("Expected 'unknown', got '%s'", str)
	}
}

func TestPermissionType_String(t *testing.T) {
	types := []PermissionType{
		PermissionTypeAccessibility,
		PermissionTypeNotification, 
		PermissionTypeAutoStart,
	}

	for _, permType := range types {
		str := permType.String()
		if str == "" {
			t.Errorf("PermissionType %v has empty string representation", permType)
		}
	}

	// Test unknown type
	unknownType := PermissionType("unknown")
	str := unknownType.String()
	if str != "unknown" {
		t.Errorf("Expected 'unknown', got '%s'", str)
	}
}

func TestPermission_IsRequired(t *testing.T) {
	tests := []struct {
		name     string
		perm     Permission
		expected bool
	}{
		{
			name: "required and denied",
			perm: Permission{
				Type:     PermissionTypeAccessibility,
				Status:   StatusDenied,
				Required: true,
			},
			expected: true,
		},
		{
			name: "required and not determined",
			perm: Permission{
				Type:     PermissionTypeAccessibility,
				Status:   StatusNotDetermined,
				Required: true,
			},
			expected: true,
		},
		{
			name: "required but granted",
			perm: Permission{
				Type:     PermissionTypeAccessibility,
				Status:   StatusGranted,
				Required: true,
			},
			expected: false,
		},
		{
			name: "not required",
			perm: Permission{
				Type:     PermissionTypeAutoStart,
				Status:   StatusDenied,
				Required: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.perm.IsRequired()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestManager_OpenSettings(t *testing.T) {
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
		// OpenSettings shouldn't panic or crash
		openErr := manager.OpenSettings(permType)
		// We don't assert error == nil because opening settings
		// may fail in test environments or headless systems
		if openErr != nil {
			t.Logf("OpenSettings for %s returned error (expected in test): %v", permType, openErr)
		}
	}

	// Test unknown permission type - stub manager doesn't validate types  
	unknownType := PermissionType("unknown")
	err = manager.OpenSettings(unknownType)
	// Stub manager always returns error, so this is expected
	if err == nil {
		t.Logf("OpenSettings for unknown type succeeded (platform-dependent behavior)")
	}
}

func TestLinuxManager_AllPermissionTypes(t *testing.T) {
	// This test improves coverage by testing all permission types on Linux
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	ctx := context.Background()
	
	// Test with all defined permission types plus some custom ones
	allTypes := []PermissionType{
		PermissionTypeAccessibility,
		PermissionTypeNotification,
		PermissionTypeAutoStart,
		PermissionType("custom1"),
		PermissionType("custom2"),
	}
	
	for _, permType := range allTypes {
		// Linux manager returns true only for known types
		supported := manager.IsSupported(permType)
		
		// Check if this is a known type
		isKnownType := permType == PermissionTypeAccessibility || 
					   permType == PermissionTypeNotification || 
					   permType == PermissionTypeAutoStart
		
		if isKnownType && !supported {
			t.Errorf("IsSupported(%v) = false, want true for known type on Linux", permType)
		} else if !isKnownType && supported {
			t.Errorf("IsSupported(%v) = true, want false for unknown type on Linux", permType)
		}
		
		// Instructions should vary based on permission type
		instructions := manager.GetInstructions(permType)
		
		// Accessibility has special message, others have generic message
		if permType == PermissionTypeAccessibility {
			expectedMsg := "No special permissions needed on Linux for global hotkeys."
			if instructions != expectedMsg {
				t.Errorf("GetInstructions(%v) = %q, want %q", permType, instructions, expectedMsg)
			}
		} else {
			expectedMsg := "This permission type is not applicable to Linux."
			if instructions != expectedMsg {
				t.Errorf("GetInstructions(%v) = %q, want %q", permType, instructions, expectedMsg)
			}
		}
		
		// All instructions should mention Linux somehow
		if !strings.Contains(instructions, "Linux") || !strings.Contains(instructions, "not applicable") {
			// Special case for accessibility which mentions Linux differently
			if permType != PermissionTypeAccessibility || !strings.Contains(instructions, "Linux") {
				t.Errorf("GetInstructions(%v) should mention Linux or 'not applicable', got: %s", permType, instructions)
			}
		}
	}
	
	// Check should return consistent results
	perms1, err := manager.Check(ctx)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	
	perms2, err := manager.Check(ctx)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	
	if len(perms1) != len(perms2) {
		t.Errorf("Check() returned inconsistent results: %d vs %d permissions", len(perms1), len(perms2))
	}
}

func TestManager_ContextCancellation(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	// Check should still work with cancelled context on Linux
	perms, err := manager.Check(ctx)
	if err != nil {
		t.Logf("Check() with cancelled context error = %v (platform-dependent)", err)
	} else {
		t.Logf("Check() returned %d permissions with cancelled context", len(perms))
	}
	
	// Request should still work
	err = manager.Request(ctx, PermissionTypeAccessibility)
	if err != nil {
		t.Logf("Request() with cancelled context error = %v (platform-dependent)", err)
	}
}
