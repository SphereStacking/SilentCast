package action

import (
	"context"
	"fmt"

	"github.com/SphereStacking/silentcast/internal/config"
)

// Manager manages action execution
type Manager struct {
	grimoire map[string]config.ActionConfig
}

// NewManager creates a new action manager
func NewManager(grimoire map[string]config.ActionConfig) *Manager {
	return &Manager{
		grimoire: grimoire,
	}
}

// Execute executes an action by spell name
func (m *Manager) Execute(ctx context.Context, spellName string) error {
	action, exists := m.grimoire[spellName]
	if !exists {
		return fmt.Errorf("spell '%s' not found in grimoire", spellName)
	}

	executor, err := m.createExecutor(&action)
	if err != nil {
		return fmt.Errorf("failed to create executor for spell '%s': %w", spellName, err)
	}

	if err := executor.Execute(ctx); err != nil {
		return fmt.Errorf("failed to execute spell '%s': %w", spellName, err)
	}

	return nil
}

// createExecutor creates an executor based on action type
func (m *Manager) createExecutor(action *config.ActionConfig) (Executor, error) {
	switch action.Type {
	case "app":
		return NewAppExecutor(action), nil
	case "script":
		return NewScriptExecutor(action), nil
	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Type)
	}
}
