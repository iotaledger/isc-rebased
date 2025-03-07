package chain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/webapi/apierrors"
	"github.com/iotaledger/wasp/packages/webapi/controllers/controllerutils"
	"github.com/iotaledger/wasp/packages/webapi/interfaces"
	"github.com/iotaledger/wasp/packages/webapi/models"
	"github.com/iotaledger/wasp/packages/webapi/params"
	"github.com/iotaledger/wasp/packages/webapi/services"
)

func (c *Controller) getCommitteeInfo(e echo.Context) error {
	controllerutils.SetOperation(e, "get_committee_info")
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService)
	if err != nil {
		return err
	}

	chain, err := c.chainService.GetChainInfoByChainID(chainID, "")
	if err != nil {
		return apierrors.ChainNotFoundError()
	}

	chainNodeInfo, err := c.committeeService.GetCommitteeInfo(chainID)
	if err != nil {
		if errors.Is(err, services.ErrNotInCommittee) {
			return e.JSON(http.StatusOK, models.CommitteeInfoResponse{})
		}
		return err
	}

	chainInfo := models.CommitteeInfoResponse{
		ChainID:        chainID.String(),
		Active:         chain.IsActive,
		StateAddress:   chainNodeInfo.Address.String(),
		CommitteeNodes: models.MapCommitteeNodes(chainNodeInfo.CommitteeNodes),
		AccessNodes:    models.MapCommitteeNodes(chainNodeInfo.AccessNodes),
		CandidateNodes: models.MapCommitteeNodes(chainNodeInfo.CandidateNodes),
	}

	return e.JSON(http.StatusOK, chainInfo)
}

func (c *Controller) getChainInfo(e echo.Context) error {
	controllerutils.SetOperation(e, "get_chain_info")
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService)
	if err != nil {
		return err
	}

	chainInfo, err := c.chainService.GetChainInfoByChainID(chainID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if errors.Is(err, interfaces.ErrChainNotFound) {
		return e.NoContent(http.StatusNotFound)
	} else if err != nil {
		return err
	}

	evmChainID := uint16(0)
	if chainInfo.IsActive {
		evmChainID, err = c.chainService.GetEVMChainID(chainID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
		if err != nil {
			return err
		}
	}

	chainInfoResponse := models.MapChainInfoResponse(chainInfo, evmChainID)

	return e.JSON(http.StatusOK, chainInfoResponse)
}

func (c *Controller) getChainList(e echo.Context) error {
	chainIDs, err := c.chainService.GetAllChainIDs()
	c.log.Infof("After allChainIDS %v", err)
	if err != nil {
		return err
	}

	chainList := make([]models.ChainInfoResponse, 0)

	for _, chainID := range chainIDs {
		chainInfo, err := c.chainService.GetChainInfoByChainID(chainID, "")
		c.log.Infof("getchaininfo %v", err)

		if errors.Is(err, interfaces.ErrChainNotFound) {
			// TODO: Validate this logic here. Is it possible to still get more chain info?
			chainList = append(chainList, models.ChainInfoResponse{
				IsActive: false,
				ChainID:  chainID.String(),
			})
			continue
		} else if err != nil {
			return err
		}

		evmChainID := uint16(0)
		if chainInfo.IsActive {
			evmChainID, err = c.chainService.GetEVMChainID(chainID, "")
			c.log.Infof("getevmchainid %v", err)

			if err != nil {
				return err
			}
		}

		chainInfoResponse := models.MapChainInfoResponse(chainInfo, evmChainID)
		c.log.Infof("mapchaininfo %v", err)

		chainList = append(chainList, chainInfoResponse)
	}

	return e.JSON(http.StatusOK, chainList)
}

func (c *Controller) getState(e echo.Context) error {
	controllerutils.SetOperation(e, "get_state")
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService)
	if err != nil {
		return err
	}

	stateKey, err := cryptolib.DecodeHex(e.Param(params.ParamStateKey))
	if err != nil {
		return apierrors.InvalidPropertyError(params.ParamStateKey, err)
	}

	state, err := c.chainService.GetState(chainID, stateKey)
	if err != nil {
		panic(err)
	}

	response := models.StateResponse{
		State: hexutil.Encode(state),
	}

	return e.JSON(http.StatusOK, response)
}

var dumpAccountsMutex = sync.Mutex{}

func (c *Controller) dumpAccounts(e echo.Context) error {
	chainID, err := controllerutils.ChainIDFromParams(e, c.chainService)
	if err != nil {
		return err
	}
	ch := lo.Must(c.chainService.GetChainByID(chainID))

	if !dumpAccountsMutex.TryLock() {
		return e.String(http.StatusLocked, "account dump in progress")
	}

	go func() {
		defer dumpAccountsMutex.Unlock()
		chainState := lo.Must(ch.LatestState(chain.ActiveOrCommittedState))
		blockIndex := chainState.BlockIndex()
		stateRoot := chainState.TrieRoot()
		filename := fmt.Sprintf("block_%d_stateroot_%s.json", blockIndex, stateRoot.String())

		err := os.MkdirAll(filepath.Join(c.accountDumpsPath, chainID.String()), os.ModePerm)
		if err != nil {
			c.log.Errorf("dumpAccounts - Creating dir failed: %s", err.Error())
			return
		}
		f, err := os.Create(filepath.Join(c.accountDumpsPath, chainID.String(), filename))
		if err != nil {
			c.log.Errorf("dumpAccounts - Creating account dump file failed: %s", err.Error())
			return
		}
		_, err = f.WriteString("{")
		if err != nil {
			c.log.Errorf("dumpAccounts - writing to account dump file failed: %s", err.Error())
			return
		}
		sa := accounts.NewStateReaderFromChainState(allmigrations.DefaultScheme.LatestSchemaVersion(), chainState)

		// because we don't know when the last account will be, we save each account string and write it in the next iteration
		// this way we can remove the trailing comma, thus getting a valid JSON
		prevString := ""

		sa.AllAccountsAsDict().ForEach(func(key kv.Key, value []byte) bool {
			if prevString != "" {
				_, err2 := f.WriteString(prevString)
				if err2 != nil {
					c.log.Errorf("dumpAccounts - writing to account dump file failed: %s", err2.Error())
					return false
				}
			}
			accKey := kv.Key(key)
			agentID := lo.Must(accounts.AgentIDFromKey(accKey, ch.ID()))
			accountAssets := sa.GetAccountFungibleTokens(agentID, chainID)
			assetsJSON, err2 := json.Marshal(accountAssets)
			if err2 != nil {
				c.log.Errorf("dumpAccounts - generating JSON for account %s assets failed%s", agentID.String(), err2.Error())
				return false
			}
			prevString = fmt.Sprintf("%q:%s,", agentID.String(), string(assetsJSON))
			return true
		})
		// delete last ',' for a valid json
		prevString = prevString[:len(prevString)-1]
		_, err = f.WriteString(fmt.Sprintf("%s}\n", prevString))
		if err != nil {
			c.log.Errorf("dumpAccounts - writing to account dump file failed: %s", err.Error())
		}
	}()

	return e.NoContent(http.StatusAccepted)
}
