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

package launcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blendle/zapdriver"
	zapbox "github.com/dfuse-io/dlauncher/zap-box"
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var UserLog = zapbox.NewCLILogger(zap.NewNop())
var zlog *zap.Logger

func init() {
	logging.Register("github.com/dfuse-io/dlauncher/launcher", &zlog)
}

type LoggingOptions struct {
	WorkingDir    string // the folder where the data will be stored, in our case will be used to store the logger
	Verbosity     int    // verbosity level
	LogFormat     string // specifies the log format
	LogToFile     bool   // specifies if we should store the logs on disk
	LogListenAddr string // address that listens to change the logs
}

func SetupLogger(opts *LoggingOptions) {
	verbosity := opts.Verbosity
	logformat := opts.LogFormat

	// TODO: The logger expect that the dataDir already exists...

	var logFileWriter zapcore.WriteSyncer
	if opts.LogToFile {
		logFileWriter = createLogFileWriter(opts.WorkingDir)
	}
	logStdoutWriter := zapcore.Lock(os.Stdout)

	commonLogger := createLogger("common", []zapcore.Level{zap.WarnLevel, zap.WarnLevel, zap.InfoLevel, zap.DebugLevel}, verbosity, logFileWriter, logStdoutWriter, logformat)
	logging.Set(commonLogger)

	for _, appDef := range AppRegistry {
		logging.Set(createLogger(appDef.ID, appDef.Logger.Levels, verbosity, logFileWriter, logStdoutWriter, logformat), appDef.Logger.Regex)
	}

	logging.Set(createLogger("bstream", []zapcore.Level{zap.WarnLevel, zap.WarnLevel, zap.InfoLevel, zap.DebugLevel}, verbosity, logFileWriter, logStdoutWriter, logformat), "github.com/dfuse-io/bstream.*")

	userLog := UserLog.LoggerReference()
	*userLog = createLogger("dfuse", []zapcore.Level{zap.InfoLevel, zap.InfoLevel, zap.DebugLevel}, verbosity, logFileWriter, logStdoutWriter, logformat)

	// Fine-grain customization
	//
	// Note that `zapbox.WithLevel` used below does not work in all circumstances! See
	// https://github.com/uber-go/zap/issues/581#issuecomment-600641485 for details.

	if value := os.Getenv("WARN"); value != "" {
		changeLoggersLevel(value, zap.WarnLevel)
	}

	if value := os.Getenv("INFO"); value != "" {
		changeLoggersLevel(value, zap.InfoLevel)
	}

	if value := os.Getenv("DEBUG"); value != "" {
		changeLoggersLevel(value, zap.DebugLevel)
	}

	// Hijack standard Golang `log` and redirect it to our common logger
	zap.RedirectStdLogAt(commonLogger, zap.DebugLevel)

	if opts.LogListenAddr != "" {
		go func() {
			UserLog.Debug("starting atomic level switcher", zap.String("listen_addr", opts.LogListenAddr))
			if err := http.ListenAndServe(opts.LogListenAddr, http.HandlerFunc(handleHTTPLogChange)); err != nil {
				UserLog.Warn("failed starting atomic level switcher", zap.Error(err), zap.String("listen_addr", opts.LogListenAddr))
			}
		}()
	}

}

type logChangeReq struct {
	Inputs string `json:"inputs"`
	Level  string `json:"level"`
}

func handleHTTPLogChange(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read body: %s", err), 400)
		return
	}

	// Unmarshal
	var in logChangeReq
	err = json.Unmarshal(b, &in)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot unmarshal JSON body: %s", err), 400)
		return
	}

	if in.Inputs == "" {
		http.Error(w, fmt.Sprintf("inputs not defined, should be comma-separated list of words or a regular expressions: %s", err), 400)
		return
	}

	switch strings.ToLower(in.Level) {
	case "warn", "warning":
		changeLoggersLevel(in.Inputs, zap.WarnLevel)
	case "info":
		changeLoggersLevel(in.Inputs, zap.InfoLevel)
	case "debug":
		changeLoggersLevel(in.Inputs, zap.DebugLevel)
	default:
		http.Error(w, fmt.Sprintf("invalid value for 'level': %s", in.Level), 400)
		return
	}

	w.Write([]byte("ok"))
}

var appToAtomicLevel = map[string]zap.AtomicLevel{}
var appToAtomicLevelLock sync.Mutex

func createLogger(appID string, levels []zapcore.Level, verbosity int, fileSyncer zapcore.WriteSyncer, consoleSyncer zapcore.WriteSyncer, format string) *zap.Logger {

	// It's ok for concurrent use here, we assume all logger are created in a single goroutine
	appToAtomicLevel[appID] = zap.NewAtomicLevelAt(appLoggerLevel(levels, verbosity))
	opts := []zap.Option{zap.AddCaller()}

	var consoleCore zapcore.Core
	switch format {
	case "stackdriver":
		opts = append(opts, zapdriver.WrapCore(zapdriver.ReportAllErrors(true), zapdriver.ServiceName(appID)))
		encoderConfig := zapdriver.NewProductionEncoderConfig()
		consoleCore = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleSyncer, appToAtomicLevel[appID])
	default:
		consoleCore = zapcore.NewCore(zapbox.NewEncoder(verbosity), consoleSyncer, appToAtomicLevel[appID])
	}

	if fileSyncer == nil {
		return zap.New(consoleCore, opts...).Named(appID)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileSyncer, zap.InfoLevel)
	teeCore := zapcore.NewTee(consoleCore, fileCore)

	return zap.New(teeCore, opts...).Named(appID)

}

func changeLoggersLevel(inputs string, level zapcore.Level) {
	for _, input := range strings.Split(inputs, ",") {
		normalizeInput := strings.Trim(input, " ")
		if normalizeInput == "bstream" || normalizeInput == "dfuse" || AppRegistry[normalizeInput] != nil {
			changeAppLogLevel(normalizeInput, level)
		} else {
			// Assumes it's a regex, we use the unnormalized input, just in case it had some spaces
			logging.Extend(overrideLoggerLevel(level), input)
		}
	}
}

// At some point, we will want to control the level from the server directly. It will
// be possible to use this method to achieve that. However, it might be required to be
// moved to `dfuse` package directly, so it's available to be used by the `gRPC` server
// in dashboard. To be determined once the issue is tackled.
func changeAppLogLevel(appID string, level zapcore.Level) {
	appToAtomicLevelLock.Lock()
	defer appToAtomicLevelLock.Unlock()

	atomicLevel, found := appToAtomicLevel[appID]
	if found {
		atomicLevel.SetLevel(level)
	}
}

func overrideLoggerLevel(level zapcore.Level) logging.LoggerExtender {
	return func(current *zap.Logger) *zap.Logger {
		return current.WithOptions(zapbox.WithLevel(level))
	}
}

func appLoggerLevel(levels []zapcore.Level, verbosity int) zapcore.Level {
	severityIndex := verbosity
	if severityIndex > len(levels)-1 {
		severityIndex = len(levels) - 1
	}

	return levels[severityIndex]
}

func createLogFileWriter(dataDir string) zapcore.WriteSyncer {
	_ = os.Mkdir(dataDir, 0755)

	logFile := filepath.Join(dataDir, "dfuse.log.json")
	writer, _, err := zap.Open(logFile)
	if err != nil {
		tempLogFile := filepath.Join(os.TempDir(), "dfuse.log.json")
		fmt.Printf("Unable to use %q for logging purposes, trying with %q instead (error: %s)\n", logFile, tempLogFile, err)
		writer, _, err := zap.Open(logFile)
		if err != nil {
			fmt.Printf("Unable to use %q for logging purposes, logs won't be saved in a log file and will be printed to console only (error: %s)\n", tempLogFile, err)
			return writer
		}
	}

	// Might return `nil`, which is handled by logging
	return writer
}
