package hotkey

import (
	"fmt"
	"strings"
)

// Validator validates key sequences and checks for conflicts
type Validator struct {
	parser     *Parser
	registered map[string]string // sequence -> spell name mapping
}

// NewValidator creates a new key sequence validator
func NewValidator() *Validator {
	return &Validator{
		parser:     NewParser(),
		registered: make(map[string]string),
	}
}

// Validate validates a key sequence and checks for conflicts
func (v *Validator) Validate(sequence string, spellName string) error {
	// Parse the sequence
	keySeq, err := v.parser.Parse(sequence)
	if err != nil {
		return err
	}
	
	// Check for empty sequence
	if len(keySeq.Keys) == 0 {
		return ValidationError{
			Sequence: sequence,
			Message:  "sequence must contain at least one key",
		}
	}
	
	// Check for exact duplicates
	normalized := v.normalize(sequence)
	if existing, exists := v.registered[normalized]; exists {
		if existing != spellName {
			return ValidationError{
				Sequence: sequence,
				Message:  fmt.Sprintf("sequence already registered for spell '%s'", existing),
			}
		}
	}
	
	// Check for prefix conflicts (e.g., "g" conflicts with "g,s")
	if err := v.checkPrefixConflicts(normalized, spellName); err != nil {
		return err
	}
	
	return nil
}

// Register marks a sequence as registered
func (v *Validator) Register(sequence string, spellName string) error {
	if err := v.Validate(sequence, spellName); err != nil {
		return err
	}
	
	normalized := v.normalize(sequence)
	v.registered[normalized] = spellName
	return nil
}

// Unregister removes a sequence registration
func (v *Validator) Unregister(sequence string) {
	normalized := v.normalize(sequence)
	delete(v.registered, normalized)
}

// Clear removes all registrations
func (v *Validator) Clear() {
	v.registered = make(map[string]string)
}

// GetRegistered returns all registered sequences
func (v *Validator) GetRegistered() map[string]string {
	result := make(map[string]string)
	for seq, spell := range v.registered {
		result[seq] = spell
	}
	return result
}

// normalize normalizes a key sequence for comparison
func (v *Validator) normalize(sequence string) string {
	// Parse and re-stringify to normalize
	keySeq, err := v.parser.Parse(sequence)
	if err != nil {
		// Return as-is if parsing fails
		return strings.ToLower(sequence)
	}
	return keySeq.String()
}

// checkPrefixConflicts checks for prefix conflicts
func (v *Validator) checkPrefixConflicts(normalized string, spellName string) error {
	parts := strings.Split(normalized, ",")
	
	// Check if this sequence is a prefix of any existing sequence
	for registered, existingSpell := range v.registered {
		if registered == normalized {
			continue // Skip exact match (already checked)
		}
		
		registeredParts := strings.Split(registered, ",")
		
		// Check if normalized is a prefix of registered
		if len(parts) < len(registeredParts) {
			isPrefix := true
			for i := 0; i < len(parts); i++ {
				if parts[i] != registeredParts[i] {
					isPrefix = false
					break
				}
			}
			if isPrefix {
				return ValidationError{
					Sequence: normalized,
					Message:  fmt.Sprintf("sequence conflicts with longer sequence '%s' (spell: %s)", registered, existingSpell),
				}
			}
		}
		
		// Check if registered is a prefix of normalized
		if len(registeredParts) < len(parts) {
			isPrefix := true
			for i := 0; i < len(registeredParts); i++ {
				if registeredParts[i] != parts[i] {
					isPrefix = false
					break
				}
			}
			if isPrefix {
				return ValidationError{
					Sequence: normalized,
					Message:  fmt.Sprintf("sequence conflicts with shorter sequence '%s' (spell: %s)", registered, existingSpell),
				}
			}
		}
	}
	
	return nil
}