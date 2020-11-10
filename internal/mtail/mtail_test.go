// Copyright 2011 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package mtail

import (
	"expvar"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/google/mtail/internal/metrics"
	"github.com/google/mtail/internal/testutil"
	"github.com/google/mtail/internal/watcher"
)

const testProgram = "/$/ { }\n"

func startMtailServer(t *testing.T, options ...func(*Server) error) *Server {
	expvar.Get("lines_total").(*expvar.Int).Set(0)
	expvar.Get("log_count").(*expvar.Int).Set(0)
	expvar.Get("log_rotations_total").(*expvar.Map).Init()
	expvar.Get("prog_loads_total").(*expvar.Map).Init()

	w, err := watcher.NewLogWatcher(0, true)
	if err != nil {
		t.Errorf("Couodn't make a log watcher: %s", err)
	}
	m, err := New(metrics.NewStore(), w, options...)
	testutil.FatalIfErr(t, err)
	if pErr := m.l.CompileAndRun("test", strings.NewReader(testProgram)); pErr != nil {
		t.Errorf("Couldn't compile program: %s", pErr)
	}

	if err := m.StartTailing(); err != nil {
		t.Errorf("StartTailing failed: %s", err)
	}
	return m
}

// func TestHandleLogDeletes(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode")
// 	}
//	workdir, rmWorkdir := testutil.TestTempDir(t)
//	defer rmWorkdir()
// 	// touch log file
// 	logFilepath := path.Join(workdir, "log")
// 	logFile, err := os.Create(logFilepath)
// 	if err != nil {
// 		t.Errorf("could not touch log file: %s", err)
// 	}
// 	defer logFile.Close()
// 	m := startMtailServer(t, LogPathPatterns(logFilepath))
// 	defer m.Close(true)

// 	if err = os.Remove(logFilepath); err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := "0"
// 	check := func() (bool, error) {
// 		if expvar.Get("log_count").String() != expected {
// 			return false, nil
// 		}
// 		return true, nil
// 	}
// 	ok, err := testutil.DoOrTimeout(check, 100*time.Millisecond, 10*time.Millisecond)
// 	if err != nil {
// 		buf := make([]byte, 1<<16)
// 		count := runtime.Stack(buf, true)
// 		t.Log("Timed out: Dumping goroutine stack")
// 		t.Log(string(buf[:count]))
// 		t.Fatal(err)
// 	}
// 	if !ok {
// 		t.Errorf("Log count not decreased\n\texpected: %s\n\treceived %s", expected, expvar.Get("log_count").String())
// 	}
// }

func TestBuildInfo(t *testing.T) {
	buildInfo := BuildInfo{
		Branch:   "foo",
		Version:  "bar",
		Revision: "baz",
	}

	buildInfoWant := fmt.Sprintf(
		"mtail version bar git revision baz go version %s go arch %s go os %s",
		runtime.Version(),
		runtime.GOARCH,
		runtime.GOOS,
	)
	buildInfoGot := buildInfo.String()

	if buildInfoWant != buildInfoGot {
		t.Errorf("Unexpected build info string, want: %q, got: %q", buildInfoWant, buildInfoGot)
	}
}
