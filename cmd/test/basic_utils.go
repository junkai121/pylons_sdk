package inttest

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"

	"strings"

	"github.com/Pylons-tech/pylons_sdk/app"
	"github.com/cosmos/cosmos-sdk/x/auth"
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// CLIOptions is a struct to manage pylonscli options
type CLIOptions struct {
	CustomNode   string
	RestEndpoint string
	MaxWaitBlock int64
}

// CLIOpts is a variable to manage pylonscli options
var CLIOpts CLIOptions
var cliMux sync.Mutex

func init() {
	flag.StringVar(&CLIOpts.CustomNode, "node", "tcp://localhost:26657", "custom node url")
}

// GetMaxWaitBlock is a function to get configuration for maximum wait block, default 3
func GetMaxWaitBlock() int64 {
	if CLIOpts.MaxWaitBlock == 0 {
		return 3
	}
	return CLIOpts.MaxWaitBlock
}

// ReadFile is a utility function to read file
func ReadFile(fileURL string, t *testing.T) []byte {
	jsonFile, err := os.Open(fileURL)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

// GetAminoCdc is a utility function to get amino codec
func GetAminoCdc() *amino.Codec {
	return app.MakeCodec()
}

// KeyringBackendSetup is a utility function to setup keyring backend for pylonscli command
func KeyringBackendSetup(args []string) []string {
	if len(args) == 0 {
		return args
	}
	newArgs := append(args, "--keyring-backend", "test")
	switch args[0] {
	case "keys":
		return newArgs
	case "tx":
		if args[1] == "sign" {
			return newArgs
		}
		return args
	default:
		return args
	}
}

// NodeFlagSetup is a utility function to setup configured custom node
func NodeFlagSetup(args []string) []string {
	if len(CLIOpts.CustomNode) > 0 {
		if args[0] == "query" || args[0] == "tx" || args[0] == "status" {
			customNodes := strings.Split(CLIOpts.CustomNode, ",")
			randNodeIndex := rand.Intn(len(customNodes))
			randNode := customNodes[randNodeIndex]
			args = append(args, "--node", randNode)
		}
	}
	return args
}

// RunPylonsCli is a function to run pylonscli
func RunPylonsCli(args []string, stdinInput string) ([]byte, string, error) {
	args = NodeFlagSetup(args)
	args = KeyringBackendSetup(args)
	cliMux.Lock()
	cmd := exec.Command(path.Join(os.Getenv("GOPATH"), "/bin/pylonscli"), args...)
	cmd.Stdin = strings.NewReader(stdinInput)
	res, err := cmd.CombinedOutput()
	cliMux.Unlock()
	return res, fmt.Sprintf("cmd is \"pylonscli %s\", result is \"%s\"", strings.Join(args, " "), string(res)), err
}

// GetAccountAddr is a function to get account address from key
func GetAccountAddr(account string, t *testing.T) string {
	addrBytes, logstr, err := RunPylonsCli([]string{"keys", "show", account, "-a"}, "")
	addr := strings.Trim(string(addrBytes), "\n ")
	if t != nil && err != nil {
		t.Fatalf("error getting account address, account=%s, err=%+v, logstr=%s", account, err, logstr)
	}
	return addr
}

// GetAccountInfoFromAddr is a function to get account information from address
func GetAccountInfoFromAddr(addr string, t *testing.T) auth.BaseAccount {
	accBytes, logstr, err := RunPylonsCli([]string{"query", "account", addr}, "")
	if t != nil && err != nil {
		t.Fatalf("error getting account info addr=%s err=%+v, logstr=%s", addr, err, logstr)
	}
	var accInfo auth.BaseAccount
	err = GetAminoCdc().UnmarshalJSON(accBytes, &accInfo)
	t.MustNil(err)
	// t.Log("GetAccountInfo", accInfo)
	return accInfo
}

// GetAccountInfoFromName is a function to get account information from account key
func GetAccountInfoFromName(account string, t *testing.T) auth.BaseAccount {
	addr := GetAccountAddr(account, t)
	return GetAccountInfoFromAddr(addr, t)
}

// GetDaemonStatus is a function to get daemon status
func GetDaemonStatus() (*ctypes.ResultStatus, string, error) {
	var ds ctypes.ResultStatus

	dsBytes, logstr, err := RunPylonsCli([]string{"status"}, "")

	if err != nil {
		return nil, logstr, err
	}
	err = GetAminoCdc().UnmarshalJSON(dsBytes, &ds)

	if err != nil {
		return nil, logstr, err
	}
	return &ds, logstr, nil
}

// WaitForNextBlock is a function to wait until next block
func WaitForNextBlock() error {
	return WaitForBlockInterval(1)
}

// WaitForBlockInterval is a function to wait until block heights to flow
func WaitForBlockInterval(interval int64) error {
	ds, _, err := GetDaemonStatus()
	if err != nil {
		return err // couldn't get daemon status.
	}
	currentBlock := ds.SyncInfo.LatestBlockHeight

	counter := int64(1)
	for counter < 300*interval {
		ds, _, err = GetDaemonStatus()
		if err != nil {
			return err
		}
		if ds.SyncInfo.LatestBlockHeight >= currentBlock+interval {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		counter++
	}
	return errors.New("You are waiting too long time for interval")
}

// CleanFile is a function to remove file
func CleanFile(filePath string, t *testing.T) {
	err := os.Remove(filePath)
	ErrValidation(t, "error removing raw tx file json %+v", err)
}

// ErrValidation is a function to log output and finish test if error is available
func ErrValidation(t *testing.T, format string, err error) {
	if err != nil {
		t.Fatalf(format, err)
	}
}

// ErrValidationWithOutputLog is a function to log output and finish test if error is available
func ErrValidationWithOutputLog(t *testing.T, format string, bytes []byte, err error) {
	if err != nil {
		t.Fatalf(format, string(bytes), err)
	}
}
