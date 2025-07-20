package action

import (
	"context"
	"sort"

	"github.com/SphereStacking/silentcast/internal/action/app"
	"github.com/SphereStacking/silentcast/internal/action/script"
	"github.com/SphereStacking/silentcast/internal/action/url"
	"github.com/SphereStacking/silentcast/internal/elevated"
	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/errors"
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

// UpdateActions updates the grimoire with new actions
func (m *Manager) UpdateActions(grimoire map[string]config.ActionConfig) {
	m.grimoire = grimoire
}

// Execute executes an action by spell name
func (m *Manager) Execute(ctx context.Context, spellName string) error {
	action, exists := m.grimoire[spellName]
	if !exists {
		// Get available spells for context
		availableSpells := make([]string, 0, len(m.grimoire))
		for spell := range m.grimoire {
			availableSpells = append(availableSpells, spell)
		}
		sort.Strings(availableSpells)
		
		return errors.New(errors.ErrorTypeConfig, "spell not found").
			WithContext("spell_name", spellName).
			WithContext("available_spells", availableSpells).
			WithContext("error_type", "spell_not_found").
			WithContext("suggested_action", "check spellbook.yml configuration")
	}

	executor, err := m.createExecutor(&action)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to create executor", err).
			WithContext("spell_name", spellName).
			WithContext("action_type", action.Type)
	}

	if err := executor.Execute(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to execute spell", err).
			WithContext("spell_name", spellName).
			WithContext("action_type", action.Type)
	}

	return nil
}

// createExecutor creates an executor based on action type
func (m *Manager) createExecutor(action *config.ActionConfig) (Executor, error) {
	var executor Executor
	
	switch action.Type {
	case "app":
		executor = app.NewAppExecutor(action)
	case "script":
		executor = script.NewScriptExecutor(action)
	case "url":
		executor = url.NewURLExecutor(action)
	default:
		return nil, errors.New(errors.ErrorTypeConfig, "unknown action type").
			WithContext("action_type", action.Type).
			WithContext("valid_types", []string{"app", "script", "url"}).
			WithContext("suggested_action", "check action type in spellbook.yml")
	}
	
	// Wrap with elevated executor if admin is required
	if action.Admin {
		executor = elevated.NewElevatedExecutor(executor, true)
	}
	
	return executor, nil
}
