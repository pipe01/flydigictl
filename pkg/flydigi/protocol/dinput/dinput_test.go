package dinput

import "testing"

func TestResolveUsbDataIgnoresStatePackets(t *testing.T) {
	p := &protocolDInput{}

	// A gamepad state packet whose sensor byte collides with the config
	// marker (235) must not be mistaken for a config chunk.
	packet := make([]byte, 32)
	packet[0] = 0x04
	packet[1] = 0xFE
	packet[15] = 235

	if msg, ok := p.resolveUsbData(packet); ok {
		t.Errorf("state packet resolved to %T, want ignored", msg)
	}
}
