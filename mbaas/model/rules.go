package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
)

// BlinkyType distinguishes between single-LED and multiple-LED (i.e.
// one group of LEDs per letter) blinkies.
type BlinkyType string

// Allowed values for the blinky type field.
const (
	SingleLED         BlinkyType = "single"
	SingleLEDGroup    BlinkyType = "group"
	LEDGroupPerLetter BlinkyType = "multi"
)

// Rules represents a single blinky creation job.
type Rules struct {
	Type              BlinkyType        `json:"type"`
	LEDForwardVoltage float32           `json:"led_forward_voltage_V"`
	LEDForwardCurrent float32           `json:"led_forward_current_mA"`
	LEDGroups         []int             `json:"led_groups",omitempty`
	BlinkRate         int               `json:"blink_rate_ms"`
	TransistorDrivers bool              `json:"transistor_drivers"`
	Footprints        map[string]string `json:"footprints"`
}

// ProjectName generates a hash-based project name from the text of a
// blinky and its rules.
func (r *Rules) ProjectName(text string) ([]byte, string, string, error) {
	projName := strings.ReplaceAll(strings.ToLower(text), " ", "-")
	unformattedJSONRules, err := json.Marshal(r)
	if err != nil {
		return nil, "", "", err
	}
	var out bytes.Buffer
	json.Indent(&out, unformattedJSONRules, "", "  ")
	jsonRules := out.Bytes()
	h := fnv.New32a()
	h.Write(jsonRules)
	fullProjName := fmt.Sprintf("%s-%0x", projName, h.Sum32())
	return jsonRules, projName, fullProjName, nil
}
