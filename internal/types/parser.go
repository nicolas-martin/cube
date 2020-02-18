package types

import (
	"encoding/json"
	"log"
)

// ParseInt parses int from json message
func ParseInt(j json.RawMessage) int {
	var v int
	if err := json.Unmarshal(j, &v); err != nil {
		log.Panicf("failed to parse int from %q: %v", j, err)
	}
	return v
}
