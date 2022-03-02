package launcher

import (
	"github.com/streamingfast/bstream"
)

// TODO: TransformableBlock Interface is not placed in bstream yet, since we want
// to ensure that we really narrowed down it's behavior, and naming convention
type BlockTransformer interface {
	TransformInPlace(blk *bstream.Block) error
}

type Runtime struct {
	AbsDataDir string
	Launcher   *Launcher
	Tracker    *bstream.Tracker

	ProtocolSpecificModules map[string]interface{}
}
