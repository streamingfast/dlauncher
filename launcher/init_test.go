package launcher

import "github.com/streamingfast/logging"

var zlog, _ = logging.PackageLogger("launcher", "github.com/streamingfast/dlauncher/launcher")

func init() {
	logging.TestingOverride()
}
