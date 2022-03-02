package launcher

import (
	"github.com/streamingfast/bstream"
)

type Runtime struct {
	AbsDataDir string
	Launcher   *Launcher
	Tracker    *bstream.Tracker

	ProtocolSpecificModules map[string]interface{}
}
