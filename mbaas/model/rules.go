package model

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
	Type              BlinkyType `json:"type"`
	LEDForwardVoltage float32    `json:"led_forward_voltage_V"`
	LEDForwardCurrent float32    `json:"led_forward_current_mA"`
	LEDGroups         []int      `json:"led_groups",omitempty`
	BlinkRate         int        `json:"blink_rate_ms"`
	MOSFETDrivers     bool       `json:"mosfet_drivers"`
}
