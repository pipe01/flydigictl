package protocol

type cmd struct{}

func (cmd) command() {}

type CommandGetDongleVersion struct {
	cmd
}

type CommandGetDeviceInfo struct {
	cmd
}

type CommandReadConfig struct {
	cmd
	ConfigID byte
}

type CommandReadLEDConfig struct {
	cmd
	ConfigID byte
}

type CommandSendConfig struct {
	cmd
	Data     []byte
	ConfigID byte
}

type CommandSendLEDConfig struct {
	cmd
	Data     []byte
	ConfigID byte
}
