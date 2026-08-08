package pb

import (
	"testing"

	"github.com/pipe01/flydigictl/pkg/flydigi/config"
)

func TestLEDConversionUnknownMode(t *testing.T) {
	bean := &config.NewLedConfigBean{LedMode: config.LedModeBreathing}

	// Modes not modeled yet must convert without panicking.
	conf := ConvertLEDConfiguration(bean)
	if _, ok := conf.Leds.(*LedsConfiguration_Off); !ok {
		t.Errorf("unknown mode converted to %T, want off", conf.Leds)
	}

	(&LedsConfiguration{Leds: &LedsConfiguration_Off{}}).ApplyTo(bean)
	if bean.LedMode != config.LedModeOff {
		t.Errorf("after applying off, mode = %d, want %d", bean.LedMode, config.LedModeOff)
	}
}
