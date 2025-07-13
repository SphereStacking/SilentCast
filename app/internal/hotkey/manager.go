//go:build !nogohook
// +build !nogohook

package hotkey

import (
	"fmt"
	"strings"
	"sync"
	"time"

	hook "github.com/robotn/gohook"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/pkg/logger"
)


// Modifier mask constants from gohook's C header (iohook.h)
// These are not exposed as Go constants in the gohook package
const (
	// Left modifiers
	maskShiftL = 1 << 0
	maskCtrlL  = 1 << 1
	maskMetaL  = 1 << 2
	maskAltL   = 1 << 3
	
	// Right modifiers
	maskShiftR = 1 << 4
	maskCtrlR  = 1 << 5
	maskMetaR  = 1 << 6
	maskAltR   = 1 << 7
	
	// Combined masks (left or right)
	maskShift = maskShiftL | maskShiftR
	maskCtrl  = maskCtrlL | maskCtrlR
	maskMeta  = maskMetaL | maskMetaR
	maskAlt   = maskAltL | maskAltR
)

// DefaultManager implements the Manager interface using gohook
type DefaultManager struct {
	mu        sync.RWMutex
	parser    *Parser
	validator *Validator
	handler   Handler
	running   bool
	stopChan  chan struct{}
	eventChan chan hook.Event
	keyMapper KeyMapper

	// Prefix key configuration
	prefixKey       KeySequence
	prefixTimeout   time.Duration
	sequenceTimeout time.Duration

	// Current state
	prefixActive    bool
	prefixTime      time.Time
	currentSequence []Key

	// Registered sequences
	sequences map[string]string // normalized sequence -> spell name
}

// NewManager creates a new hotkey manager
func NewManager(cfg *config.HotkeyConfig) (*DefaultManager, error) {
	parser := NewParser()

	// Parse prefix key
	prefixSeq, err := parser.Parse(cfg.Prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prefix key '%s': %w", cfg.Prefix, err)
	}

	return &DefaultManager{
		parser:          parser,
		validator:       NewValidator(),
		stopChan:        make(chan struct{}),
		eventChan:       make(chan hook.Event, 100),
		keyMapper:       GetKeyMapper(),
		prefixKey:       prefixSeq,
		prefixTimeout:   cfg.Timeout.ToDuration(),
		sequenceTimeout: cfg.SequenceTimeout.ToDuration(),
		sequences:       make(map[string]string),
		currentSequence: make([]Key, 0),
	}, nil
}

// Start begins listening for hotkeys
func (m *DefaultManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("hotkey manager already running")
	}

	logger.Debug("Starting hotkey manager...")
	logger.Debug("Prefix key: %s", m.prefixKey.String())
	logger.Debug("Registered sequences: %d", len(m.sequences))

	m.running = true
	m.stopChan = make(chan struct{})

	// Start event processing goroutine
	go m.processEvents()

	// Start hook event collection
	go m.collectEvents()

	logger.Info("✅ Hotkey manager started with prefix: %s", m.prefixKey.String())
	return nil
}

// Stop stops listening for hotkeys
func (m *DefaultManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return fmt.Errorf("hotkey manager not running")
	}

	m.running = false
	close(m.stopChan)

	// Stop gohook
	hook.End()

	return nil
}

// Register registers a hotkey sequence with a spell name
func (m *DefaultManager) Register(sequence string, spellName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate the sequence
	if err := m.validator.Register(sequence, spellName); err != nil {
		return err
	}

	// Parse and normalize the sequence
	keySeq, err := m.parser.Parse(sequence)
	if err != nil {
		return err
	}

	normalized := keySeq.String()
	m.sequences[normalized] = spellName

	return nil
}

// Unregister removes a hotkey registration
func (m *DefaultManager) Unregister(sequence string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Parse and normalize the sequence
	keySeq, err := m.parser.Parse(sequence)
	if err != nil {
		return err
	}

	normalized := keySeq.String()
	delete(m.sequences, normalized)
	m.validator.Unregister(sequence)

	return nil
}

// SetHandler sets the event handler
func (m *DefaultManager) SetHandler(handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handler = handler
}

// IsRunning returns whether the manager is actively listening
func (m *DefaultManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// collectEvents collects keyboard events from gohook
func (m *DefaultManager) collectEvents() {
	logger.Debug("Starting gohook event collection...")
	evChan := hook.Start()
	defer hook.End()
	
	logger.Info("🎮 Hotkey listener started successfully")

	for {
		select {
		case <-m.stopChan:
			logger.Debug("Stopping event collection")
			return
		case ev := <-evChan:
			// Log all key events for debugging
			if ev.Kind == hook.KeyDown {
				// Convert to readable key name for logging
				keyName := m.getKeyName(ev)
				logger.Debug("⌨️  Key pressed: %s (keycode=%d, rawcode=%d, mask=%d)", keyName, ev.Keycode, ev.Rawcode, ev.Mask)
				
				// Also log modifier states
				if ev.Mask != 0 {
					var mods []string
					if ev.Mask&maskCtrl != 0 { mods = append(mods, "Ctrl") }
					if ev.Mask&maskShift != 0 { mods = append(mods, "Shift") }
					if ev.Mask&maskAlt != 0 { mods = append(mods, "Alt") }
					if len(mods) > 0 {
						logger.Debug("   Active modifiers: %s", strings.Join(mods, "+"))
					}
				}
				
				select {
				case m.eventChan <- ev:
				default:
					// Drop event if channel is full
					logger.Warn("Event channel full, dropping key event")
				}
			}
		}
	}
}

// processEvents processes keyboard events
func (m *DefaultManager) processEvents() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return

		case ev := <-m.eventChan:
			m.handleKeyEvent(ev)

		case <-ticker.C:
			// Check for timeouts
			m.checkTimeouts()
		}
	}
}

// handleKeyEvent handles a single key event
func (m *DefaultManager) handleKeyEvent(ev hook.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Convert gohook event to our Key type
	key := m.convertEvent(ev)
	if key == nil {
		return
	}

	// Check if this is the prefix key
	if !m.prefixActive && m.isPrefix(key) {
		m.prefixActive = true
		m.prefixTime = time.Now()
		m.currentSequence = []Key{}
		logger.Info("🔵 Prefix key detected - waiting for command...")
		return
	}

	// If prefix is active, build sequence
	if m.prefixActive {
		m.currentSequence = append(m.currentSequence, *key)

		// Check if we have a match
		currentSeq := KeySequence{Keys: m.currentSequence}
		normalized := currentSeq.String()

		// Check for exact match
		if spellName, exists := m.sequences[normalized]; exists {
			// We have a match! Execute the handler
			logger.Info("✅ Executed: %s", spellName)
			if m.handler != nil {
				event := Event{
					Sequence:  currentSeq,
					SpellName: spellName,
					Timestamp: time.Now(),
				}

				// Execute handler in goroutine to not block
				go m.handler.Handle(event)
			}

			// Reset state
			m.resetState()
			return
		}

		// Check if this could be a prefix of any registered sequence
		isPossiblePrefix := false
		for seq := range m.sequences {
			if len(normalized) < len(seq) && strings.HasPrefix(seq, normalized) {
				isPossiblePrefix = true
				break
			}
		}

		// If not a possible prefix, reset
		if !isPossiblePrefix {
			logger.Info("❌ Unknown command: %s", normalized)
			m.resetState()
		}
	}
}

// convertEvent converts a gohook event to our Key type
func (m *DefaultManager) convertEvent(ev hook.Event) *Key {
	keyName := ""
	modifiers := []string{}
	
	// Check if this is a modifier key itself
	if modName, isModifier := m.keyMapper.IsModifierKey(ev.Rawcode); isModifier {
		keyName = modName
		modifiers = append(modifiers, modName)
	} else {
		// Regular key - check for active modifiers
		if ev.Mask&maskCtrl != 0 {
			modifiers = append(modifiers, "ctrl")
		}
		if ev.Mask&maskShift != 0 {
			modifiers = append(modifiers, "shift")
		}
		if ev.Mask&maskAlt != 0 {
			modifiers = append(modifiers, "alt")
		}
		
		// Try platform-specific mapping first
		keyName = m.keyMapper.GetKeyNameFromRawcode(ev.Rawcode)
		
		// Fallback to gohook's built-in conversion
		if keyName == "" {
			keyName = hook.RawcodetoKeychar(ev.Rawcode)
		}
	}
	
	if keyName == "" {
		return nil
	}

	return &Key{
		Code:      ev.Keycode, // Use keycode instead of rawcode
		Modifiers: modifiers,
		Name:      strings.ToLower(keyName),
	}
}

// isPrefix checks if a key matches the prefix key
func (m *DefaultManager) isPrefix(key *Key) bool {
	if len(m.prefixKey.Keys) != 1 {
		return false // Only support single key prefix for now
	}

	prefixKey := m.prefixKey.Keys[0]

	// Check key name
	if key.Name != prefixKey.Name {
		return false
	}

	// Check modifiers
	if len(key.Modifiers) != len(prefixKey.Modifiers) {
		return false
	}

	for i, mod := range key.Modifiers {
		if mod != prefixKey.Modifiers[i] {
			return false
		}
	}

	return true
}

// checkTimeouts checks for prefix and sequence timeouts
func (m *DefaultManager) checkTimeouts() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.prefixActive {
		return
	}

	now := time.Now()

	// Check prefix timeout
	if now.Sub(m.prefixTime) > m.prefixTimeout {
		logger.Info("⏱️  Timeout - command cancelled")
		m.resetState()
		return
	}

	// Check sequence timeout (if we have started a sequence)
	if len(m.currentSequence) > 0 && now.Sub(m.prefixTime) > m.sequenceTimeout {
		logger.Info("⏱️  Timeout - command cancelled")
		m.resetState()
	}
}

// resetState resets the current hotkey state
func (m *DefaultManager) resetState() {
	m.prefixActive = false
	m.currentSequence = []Key{}
}

// getKeyName converts a gohook event to a readable key name
func (m *DefaultManager) getKeyName(ev hook.Event) string {
	// Try to convert to our Key type first
	key := m.convertEvent(ev)
	if key != nil {
		return key.String()
	}
	
	// Fallback to raw keycode
	return fmt.Sprintf("Unknown(%d)", ev.Keycode)
}

