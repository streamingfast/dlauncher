package launcher

type Runtime struct {
	AbsDataDir string
	Launcher   *Launcher

	ProtocolSpecificModules map[string]interface{}
}
