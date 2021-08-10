// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dashboard

import (
	"time"

	"github.com/streamingfast/shutter"
	"github.com/streamingfast/dlauncher/launcher"
	"github.com/streamingfast/dlauncher/metrics"
	dmeshCli "github.com/streamingfast/dmesh/client"
	"go.uber.org/zap"
)

type Config struct {
	HTTPListenAddr      string
	GRPCListenAddr      string
	NodeManagerAPIAddr  string
	DmeshServiceVersion string
	MetricsHTTPAddr     string

	// Dashboard configuration payload
	Title              string
	BlockExplorerName  string
	HeadBlockNumberApp string
}

type Modules struct {
	Launcher    *launcher.Launcher
	DmeshClient dmeshCli.SearchClient
}

type App struct {
	*shutter.Shutter
	config   *Config
	launcher *launcher.Launcher
	Ready    chan interface{}
	ready    bool
	modules  *Modules
}

func New(config *Config, modules *Modules) *App {
	zlog.Info("new dashboard app", zap.Reflect("config", config))
	return &App{
		Shutter: shutter.New(),
		config:  config,
		Ready:   make(chan interface{}),
		modules: modules,
	}
}

func (a *App) Run() error {
	// Launch MetricManager
	mgr := metrics.NewManager(a.config.MetricsHTTPAddr+"/metrics", []string{"head_block_time_drift", "head_block_number"}, 5*time.Second, launcher.GetMetricAppMeta())
	go mgr.Launch()

	s := newServer(a.config, a.modules, mgr)

	a.OnTerminating(s.Shutdown)

	go func() {
		a.Shutdown(s.Launch())
	}()

	close(a.Ready)
	a.ready = true

	return nil
}

func (a *App) OnReady(f func()) {
	<-a.Ready
	f()
}

func (a *App) IsReady() bool {
	return a.ready
}
