package service

import (
	"testing"
)

func TestServiceStatus(t *testing.T) {
	// Test zero value
	status := ServiceStatus{}
	
	if status.Installed {
		t.Error("Default ServiceStatus should have Installed = false")
	}
	
	if status.Running {
		t.Error("Default ServiceStatus should have Running = false")
	}
	
	if status.StartType != "" {
		t.Error("Default ServiceStatus should have empty StartType")
	}
	
	if status.Message != "" {
		t.Error("Default ServiceStatus should have empty Message")
	}
}

func TestServiceStatusFields(t *testing.T) {
	// Test setting all fields
	status := ServiceStatus{
		Installed: true,
		Running:   true,
		StartType: "auto",
		Message:   "Service is running normally",
	}
	
	if !status.Installed {
		t.Error("ServiceStatus.Installed should be true")
	}
	
	if !status.Running {
		t.Error("ServiceStatus.Running should be true")
	}
	
	if status.StartType != "auto" {
		t.Errorf("ServiceStatus.StartType = %s, want auto", status.StartType)
	}
	
	if status.Message != "Service is running normally" {
		t.Errorf("ServiceStatus.Message = %s, want 'Service is running normally'", status.Message)
	}
}

func TestServiceStatusStartTypes(t *testing.T) {
	// Test different start types
	startTypes := []string{"auto", "manual", "disabled", ""}
	
	for _, startType := range startTypes {
		status := ServiceStatus{
			StartType: startType,
		}
		
		if status.StartType != startType {
			t.Errorf("StartType not preserved: got %s, want %s", status.StartType, startType)
		}
	}
}

func TestServiceStatusCopy(t *testing.T) {
	// Test that ServiceStatus can be copied by value
	original := ServiceStatus{
		Installed: true,
		Running:   false,
		StartType: "manual",
		Message:   "Service stopped",
	}
	
	copy := original
	copy.Running = true
	copy.Message = "Service started"
	
	// Original should be unchanged
	if original.Running {
		t.Error("Original ServiceStatus was modified during copy")
	}
	
	if original.Message != "Service stopped" {
		t.Error("Original ServiceStatus message was modified during copy")
	}
	
	// Copy should be modified
	if !copy.Running {
		t.Error("Copied ServiceStatus was not modified correctly")
	}
	
	if copy.Message != "Service started" {
		t.Error("Copied ServiceStatus message was not modified correctly")
	}
}