// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// privtangle is a cluster of SUI nodes started for testing purposes.
package privtangle_sui

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/iotaledger/wasp/packages/testutil/privtangle_sui/miniclient"
	"github.com/iotaledger/wasp/packages/testutil/privtangle_sui/privtangledefaults"
	"github.com/iotaledger/wasp/packages/util"
)

type LogFunc func(format string, args ...interface{})

type PrivTangle struct {
	SnapshotInit string
	ConfigFile   string
	BaseDir      string
	BasePort     int
	NodeCount    int
	NodeCommands []*exec.Cmd
	ctx          context.Context
	logfunc      LogFunc
}

func Start(ctx context.Context, baseDir string, basePort, nodeCount int, logfunc LogFunc) *PrivTangle {
	pt := PrivTangle{

		SnapshotInit: "snapshot.init",
		ConfigFile:   "config.json",
		BaseDir:      baseDir,
		BasePort:     basePort,
		NodeCount:    nodeCount,
		NodeCommands: make([]*exec.Cmd, nodeCount),
		ctx:          ctx,
		logfunc:      logfunc,
	}

	pt.logf("Starting in baseDir=%s with basePort=%d, nodeCount=%d ...", baseDir, basePort, nodeCount)

	if err := os.MkdirAll(pt.BaseDir, 0o755); err != nil {
		panic(fmt.Errorf("unable to create dir %v: %w", pt.BaseDir, err))
	}

	pt.StartServers(true)

	return &pt
}

func (pt *PrivTangle) StartServers(deleteExisting bool) {
	ts := time.Now()
	pt.logf("Starting all SUI nodes...")

	for i := range pt.NodeCommands {
		pt.startNode(i, deleteExisting)
	}

	pt.logf("Starting all SUI nodes... done! took: %v", time.Since(ts).Truncate(time.Millisecond))

	pt.waitAllHealthy(5 * time.Minute)

	pt.logf("Privtangle started successfully")
}

func (pt *PrivTangle) startNode(i int, deleteExisting bool) {
	nodePath := filepath.Join(pt.BaseDir, "sui-node")

	if deleteExisting {
		if err := os.RemoveAll(nodePath); err != nil {
			panic(fmt.Errorf("unable to delete dir %v: %w", nodePath, err))
		}
		if err := os.MkdirAll(nodePath, 0o755); err != nil {
			panic(fmt.Errorf("unable to create dir %v: %w", nodePath, err))
		}
	}

	args := []string{
		// Setting a config path here will crash the node as it does not create required config files itself.
		// We either need to keep working with the args alone, or establish config templates.
		//fmt.Sprintf("--config-dir=%s", nodePath),
		fmt.Sprintf("--epoch-duration-ms=%d", 60000),
		fmt.Sprintf("--fullnode-rpc-port=%d", pt.NodePortRestAPI(i)),
		fmt.Sprintf("--faucet-port=%d", pt.NodePortFaucet(i)),
		fmt.Sprintf("--graphql-port=%d", pt.NodePortGraphQL(i)),
		// Indexer requires a running postgres DB
		//"--with-indexer",
		//fmt.Sprintf("--indexer-rpc-port=%d", pt.NodePortIndexer(i)),
	}

	testValidatorCmd := exec.CommandContext(pt.ctx, "/home/luke/dev/iota_sui/target/debug/sui-test-validator", args...)

	// kill SUI cmd if the go test process is killed
	util.TerminateCmdWhenTestStops(testValidatorCmd)

	testValidatorCmd.Env = os.Environ()
	testValidatorCmd.Dir = nodePath

	pt.NodeCommands[i] = testValidatorCmd

	writeOutputToFiles(nodePath, testValidatorCmd)

	if err := testValidatorCmd.Start(); err != nil {
		panic(fmt.Errorf("cannot start SUI node[%d]: %w", i, err))
	}
}

func (pt *PrivTangle) Stop() {
	pt.logf("Stopping...")

	for i, c := range pt.NodeCommands {
		if err := c.Process.Signal(syscall.SIGTERM); err != nil {
			panic(fmt.Errorf("unable to send INT signal to SUI node [%d]: %w", i, err))
		}
	}

	for i, c := range pt.NodeCommands {
		if err := c.Wait(); err != nil {
			var errCode *exec.ExitError
			ok := errors.As(err, &errCode)

			if ok && strings.Contains(errCode.Error(), "terminated") {
				pt.logf("Stopping... Done")
				return
			}

			panic(fmt.Errorf("SUI node [%d] failed: %s", i, c.ProcessState.String()))
		}
	}
}

func (pt *PrivTangle) nodeClient(i int) *miniclient.MiniClient {
	return miniclient.NewMiniClient(fmt.Sprintf("http://localhost:%d", pt.NodePortRestAPI(i)))
}

func (pt *PrivTangle) waitAllHealthy(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(pt.ctx, timeout)
	defer cancel()

	ts := time.Now()
	pt.logf("Waiting for all SUI nodes to become healthy...")

	for {
		if ctx.Err() != nil {
			panic("nodes didn't become healthy in time")
		}

		allOK := true
		for i := range pt.NodeCommands {
			res, err := pt.nodeClient(i).GetLatestSuiSystemState(pt.ctx)
			if err != nil || res.Result.PendingActiveValidatorsSize != "0" {
				pt.logf("Waiting healthy... node #%d not ready yet. time waiting: %v", i, time.Since(ts).Truncate(time.Millisecond))
				allOK = false
			}
		}
		if allOK {
			pt.logf("Waiting for all SUI nodes to become healthy... done! took: %v", time.Since(ts).Truncate(time.Millisecond))
			return
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (pt *PrivTangle) NodePortRestAPI(i int) int {
	return pt.BasePort + i*100 + privtangledefaults.NodePortOffsetRestAPI
}

func (pt *PrivTangle) NodePortPeering(i int) int {
	return pt.BasePort + i*100 + privtangledefaults.NodePortOffsetPeering
}

func (pt *PrivTangle) NodePortFaucet(i int) int {
	return pt.BasePort + i*100 + privtangledefaults.NodePortOffsetFaucet
}

func (pt *PrivTangle) NodePortGraphQL(i int) int {
	return pt.BasePort + i*100 + privtangledefaults.NodePortOffsetGraphQL
}

func (pt *PrivTangle) NodePortIndexer(i int) int {
	return pt.BasePort + i*100 + privtangledefaults.NodePortOffsetIndexer
}

func (pt *PrivTangle) logf(msg string, args ...interface{}) {
	if pt.logfunc != nil {
		pt.logfunc("SUI Cluster: "+msg, args...)
	}
}

func writeOutputToFiles(path string, cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Errorf("unable to get stdout for SUI [path: %s]: %w", path, err))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(fmt.Errorf("unable to get stdout for SUI[path: %s]: %w", path, err))
	}
	outFilePath := filepath.Join(path, "stdout.log")
	outFile, err := os.Create(outFilePath)
	if err != nil {
		panic(err)
	}
	errFilePath := filepath.Join(path, "stderr.log")
	errFile, err := os.Create(errFilePath)
	if err != nil {
		panic(err)
	}
	go scanLog(
		stderr,
		func(line string) {
			_, err := errFile.WriteString(fmt.Sprintln(line))
			if err != nil {
				panic(fmt.Errorf("error writing to file %s: %w", errFilePath, err))
			}
		},
	)
	go scanLog(
		stdout,
		func(line string) {
			_, err := outFile.WriteString(fmt.Sprintln(line))
			if err != nil {
				panic(fmt.Errorf("error writing to file %s: %w", outFilePath, err))
			}
		},
	)
}

func scanLog(reader io.Reader, hooks ...func(string)) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		for _, hook := range hooks {
			hook(line)
		}
	}
}
