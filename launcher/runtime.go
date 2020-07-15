package launcher

import (
	"github.com/dfuse-io/bstream"
	dmeshClient "github.com/dfuse-io/dmesh/client"
)

// TODO: TransformableBlock Interface is not placed in bstream yet, since we want
// to ensure that we really narrowed down it's behavior, and naming convention
type BlockTransformer interface {
	TransformInPlace(blk *bstream.Block) error
}

type RuntimeModules struct {
	SearchDmeshClient dmeshClient.SearchClient
	Launcher          *Launcher
	BlockFilter       BlockTransformer
}
