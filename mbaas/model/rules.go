package model

// BlinkyType distinguishes between single-LED and multiple-LED (i.e.
// one group of LEDs per letter) blinkies.
type BlinkyType string

// Allowed values for the blinky type field.
const (
	SingleLED BlinkyType = "single"
	MultiLED  BlinkyType = "multi"
)

// Job represents a single blinky creation job.
type Job struct {
	Text              string     `json:"text"`
	Type              BlinkyType `json:"type"`
	LEDForwardVoltage float32    `json:"led_forward_voltage_V"`
	LEDForwardCurrent float32    `json:"led_forward_current_mA"`
	LEDGroups         []int      `json:"led_groups",omitempty`
	BlinkRate         int        `json:"blink_rate_ms"`
}
