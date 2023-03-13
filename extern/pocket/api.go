package pocket

import (
	"encoding/json"
	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/utils"
	ctypes "github.com/tendermint/tendermint/rpc/coretypes"
	"time"
)

func PoktApiGetHeight(url string) (uint64, error) {
	params := &HeightPoktParams{}
	apiUrl := utils.UrlJoin(url, "/query/height")
	response, err := web.HttpPostNoToken(apiUrl, params)
	if err != nil {
		GetLog().Error(err)
		return 0, err
	}

	poktRes := HeightData{}
	err = json.Unmarshal([]byte(response), &poktRes)
	if err != nil {
		GetLog().Error(err)
		return 0, err
	}

	return poktRes.Height, nil
}

type TmStatusResponse struct {
	JsonRpc string              `json:"jsonrpc"`
	Id      int                 `json:"id"`
	Result  ctypes.ResultStatus `json:"message,omitempty"`
}

func PoktApiGetSync() (bool, error) {
	url := "http://127.0.0.1:26657/status"
	response, err := web.HttpPostNoToken(url, "")
	if err != nil {
		GetLog().Error(err)
		return false, err
	}

	statusRes := TmStatusResponse{}
	err = json.Unmarshal([]byte(response), &statusRes)
	if err != nil {
		GetLog().Error(err)
		return false, err
	}

	diff := time.Now().UTC().Sub(statusRes.Result.SyncInfo.LatestBlockTime)
	seconds := int(diff.Seconds())
	GetLog().Info("Latest block was out ", seconds, " seconds ago.")
	return seconds < 60*30, nil
}
