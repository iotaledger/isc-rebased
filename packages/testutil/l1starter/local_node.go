package l1starter

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/iotaledger/wasp/clients"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
)

var WaitUntilEffectsVisible = &iotaclient.WaitParams{
	Attempts:             10,
	DelayBetweenAttempts: 1 * time.Second,
}

type LocalIotaNode struct {
	ctx             context.Context
	config          Config
	iscPackageOwner iotasigner.Signer
	iscPackageID    iotago.PackageID
	container       testcontainers.Container
}

func NewLocalIotaNode(iscPackageOwner iotasigner.Signer) *LocalIotaNode {
	return &LocalIotaNode{
		iscPackageOwner: iscPackageOwner,
		config: Config{
			Host:   "http://localhost",
			Ports:  Ports{},
			Logger: Logger{},
		},
	}
}

func (in *LocalIotaNode) start(ctx context.Context) {
	in.ctx = ctx

	imagePlatform := "linux/amd64"
	if runtime.GOARCH == "arm64" {
		imagePlatform = "linux/arm64"
	}

	req := testcontainers.ContainerRequest{
		Image:         "iotaledger/iota-tools:v0.10.0-alpha",
		ImagePlatform: imagePlatform,
		ExposedPorts:  []string{"9000/tcp", "9123/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("9000/tcp"),
			wait.ForListeningPort("9123/tcp"),
		).WithDeadline(30 * time.Second),
		Cmd: []string{
			"iota",
			"start",
			"--force-regenesis",
			fmt.Sprintf("--epoch-duration-ms=%d", 60000),
			"--with-faucet",
			fmt.Sprintf("--faucet-amount=%d", iotaclient.SingleCoinFundsFromFaucetAmount),
		},
	}

	if runtime.GOOS == "linux" {
		req.Tmpfs = map[string]string{"/tmp": ""}
	}

	now := time.Now()

	in.logf("Starting LocalIotaNode...")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start container: %w", err))
	}

	in.container = container

	webAPIPort, err := container.MappedPort(ctx, "9000")
	if err != nil {
		container.Terminate(ctx)
		panic(fmt.Errorf("failed to get web API port: %w", err))
	}

	faucetPort, err := container.MappedPort(ctx, "9123")
	if err != nil {
		container.Terminate(ctx)
		panic(fmt.Errorf("failed to get faucet port: %w", err))
	}

	in.config.Ports.RPC = webAPIPort.Int()
	in.config.Ports.Faucet = faucetPort.Int()

	in.logf("Starting LocalIotaNode... done! took: %v", time.Since(now).Truncate(time.Millisecond))
	in.waitAllHealthy(5 * time.Minute)
	in.logf("Deploying ISC contracts...")

	packageID, err := in.L1Client().DeployISCContracts(ctx, ISCPackageOwner)
	if err != nil {
		panic(fmt.Errorf("isc contract deployment failed: %w", err))
	}

	in.iscPackageID = packageID

	in.logf("LocalIotaNode started successfully")
}

func (in *LocalIotaNode) stop() {
	in.logf("Stopping...")
	in.container.Terminate(context.Background(), testcontainers.StopTimeout(0))
	instance.Store(nil)
}

func (in *LocalIotaNode) ISCPackageID() iotago.PackageID {
	return in.iscPackageID
}

func (in *LocalIotaNode) APIURL() string {
	return fmt.Sprintf("%s:%d", in.config.Host, in.config.Ports.RPC)
}

func (in *LocalIotaNode) FaucetURL() string {
	return fmt.Sprintf("%s:%d/gas", in.config.Host, in.config.Ports.Faucet)
}

func (in *LocalIotaNode) L1Client() clients.L1Client {
	return clients.NewL1Client(clients.L1Config{
		APIURL:    in.APIURL(),
		FaucetURL: in.FaucetURL(),
	}, WaitUntilEffectsVisible)
}

func (in *LocalIotaNode) IsLocal() bool {
	return true
}

func (in *LocalIotaNode) waitAllHealthy(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(in.ctx, timeout)
	defer cancel()

	ts := time.Now()
	in.logf("Using temporary folder: %s", in.config.TempDir)
	in.logf("Waiting for all IOTA nodes to become healthy...")

	tryLoop := func(f func() bool) {
		for {
			if ctx.Err() != nil {
				panic("nodes didn't become healthy in time")
			}
			if f() {
				return
			}
			in.logf("Waiting until LocalIotaNode becomes ready. Time waiting: %v", time.Since(ts).Truncate(time.Millisecond))
			time.Sleep(500 * time.Millisecond)
		}
	}

	tryLoop(func() bool {
		res, err := in.L1Client().GetLatestIotaSystemState(in.ctx)
		if err != nil {
			in.logf("StatusLoop: err: %s", err)
		}
		if err != nil || res == nil {
			return false
		}
		if res.PendingActiveValidatorsSize.Uint64() != 0 {
			return false
		}
		return true
	})

	tryLoop(func() bool {
		err := iotaclient.RequestFundsFromFaucet(ctx, ISCPackageOwner.Address(), in.FaucetURL())
		if err != nil {
			in.logf("FaucetLoop: err: %s", err)
		}
		return err == nil
	})

	in.logf("Waiting until LocalIotaNode becomes ready... done! took: %v", time.Since(ts).Truncate(time.Millisecond))
}

func (in *LocalIotaNode) logf(msg string, args ...any) {
	if in.config.Logger != nil {
		in.config.Logger.Printf("Iota Node: "+msg+"\n", args...)
	}
}
