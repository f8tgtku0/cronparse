package schedule

import (
	"encoding/json"
	"fmt"
	"os"
)

// fileSchema is the JSON shape used for persistence.
type fileSchema struct {
	Entries []struct {
		Name string `json:"name"`
		Raw  string `json:"expression"`
	} `json:"schedules"`
}

// LoadJSON reads a JSON file and populates the Schedule.
// The file must contain an object with a "schedules" array of
// {"name": "...", "expression": "..."} objects.
func (s *Schedule) LoadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("loadJSON: %w", err)
	}
	var schema fileSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("loadJSON: %w", err)
	}
	for _, e := range schema.Entries {
		if err := s.Add(e.Name, e.Raw); err != nil {
			return err
		}
	}
	return nil
}

// SaveJSON writes the current entries to path as JSON.
func (s *Schedule) SaveJSON(path string) error {
	var schema fileSchema
	for _, e := range s.entries {
		schema.Entries = append(schema.Entries, struct {
			Name string `json:"name"`
			Raw  string `json:"expression"`
		}{Name: e.Name, Raw: e.Raw})
	}
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("saveJSON: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
