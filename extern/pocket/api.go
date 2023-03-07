package pocket

import (
	"encoding/json"
	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/utils"
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
