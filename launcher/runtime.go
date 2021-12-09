package launcher

import (
	"github.com/streamingfast/bstream"
	dmeshClient "github.com/streamingfast/dmesh/client"
	pbblockmeta "github.com/streamingfast/pbgo/sf/blockmeta/v1"
)

// TODO: TransformableBlock Interface is not placed in bstream yet, since we want
// to ensure that we really narrowed down it's behavior, and naming convention
type BlockTransformer interface {
	TransformInPlace(blk *bstream.Block) error
}

type Runtime struct {
	SearchDmeshClient dmeshClient.SearchClient
	Launcher          *Launcher
	BlockFilter       BlockTransformer
	BlockMeta         pbblockmeta.BlockIDClient
	Tracker           *bstream.Tracker
	AbsDataDir        string
}
