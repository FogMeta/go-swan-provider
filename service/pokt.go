package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/filswan/go-swan-lib/client/swan"
	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/utils"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"strconv"
	"swan-provider/common"
	"swan-provider/config"
	"swan-provider/extern/docker"
	mydc "swan-provider/extern/docker"
	api "swan-provider/extern/pocket"
	"swan-provider/logs"
	"swan-provider/models"
	"time"
)

type Status struct {
	Version     string
	Address     string
	Height      string
	Balance     string
	Award       string
	Jailed      string
	JailedBlock string
	JailedUntil string
}

type PoktService struct {
	PoktApiUrl           string
	PoktAccessToken      string
	PoktScanInterval     time.Duration
	ApiHeartbeatInterval time.Duration
	PoktServerApiUrl     string
	PoktServerApiPort    int
	PoktNetType          string
	PoktNetSeed          string

	dkImage    string
	dkName     string
	dkConfPath string

	dkCli *mydc.DockerCli

	PoktAddress string
	AlarmBlc    big.Int
	CurStatus   Status
}

var myPoktSvr *PoktService

func GetMyPoktService() *PoktService {
	if myPoktSvr == nil {
		confPokt := config.GetConfig().Pokt

		myPoktSvr = &PoktService{
			PoktApiUrl:           confPokt.PoktApiUrl,
			PoktAccessToken:      confPokt.PoktAccessToken,
			PoktScanInterval:     confPokt.PoktScanInterval,
			ApiHeartbeatInterval: config.GetConfig().Main.SwanApiHeartbeatInterval,
			PoktServerApiUrl:     confPokt.PoktServerApiUrl,
			PoktServerApiPort:    confPokt.PoktServerApiPort,
			PoktNetType:          confPokt.PoktNetType,
			PoktNetSeed:          confPokt.PoktNetSeed,
			dkImage:              confPokt.PoktDockerImage,
			dkName:               confPokt.PoktDockerName,
			dkConfPath:           confPokt.PoktConfigPath,
			CurStatus:            Status{},
		}
		myPoktSvr.dkCli = docker.GetMyCli(myPoktSvr.dkImage, myPoktSvr.dkName, myPoktSvr.dkConfPath)

		logs.GetLog().Debugf("New myPoktSvr :%+v ", *myPoktSvr)

		return myPoktSvr
	}
	return myPoktSvr
}

func (psvc *PoktService) GetCli() *mydc.DockerCli {
	if psvc.dkCli == nil {
		psvc.dkCli = docker.GetMyCli(psvc.dkImage, psvc.dkName, psvc.dkConfPath)
		logs.GetLog().Infof("GetCli New Docker Cli")
	}
	return psvc.dkCli
}
func (psvc *PoktService) StartPoktContainer(op []string) {

	cli := psvc.dkCli
	if !cli.PoktCtnExist() {

		logs.GetLog().Debug("Init Pocket Container ... ")
		fs := flag.NewFlagSet("Start", flag.ExitOnError)
		passwd := fs.String("passwd", "", "password for create account")
		err := fs.Parse(op[1:])
		if *passwd == "" || err != nil {
			printPoktUsage()
			panic("need password for create account.")
			return
		}

		pass := *passwd
		logs.GetLog().Debug("POCKET_CORE_PASSPHRASE=", pass)
		env := []string{"POCKET_CORE_KEY=", "POCKET_CORE_PASSPHRASE=" + pass}

		accCmd := []string{"pocket", "accounts", "create"}
		cli.PoktCtnPullAndCreate(accCmd, env, true)
		cli.PoktCtnStart()

		for {
			if cli.PoktCtnExist() {
				logs.GetLog().Info("Wait for Creating Account...")
				cli.PoktCtnList()
				time.Sleep(time.Second * 3)
				continue
			}
			break
		}
		logs.GetLog().Debug("Init Creating Account Over")

		runCmd := []string{
			"pocket",
			"start",
			"--seeds=" + psvc.PoktNetSeed,
			"--" + psvc.PoktNetType}

		if "simulate" == psvc.PoktNetType {
			runCmd = []string{
				"pocket",
				"start",
				"--simulateRelay"}

			logs.GetLog().Debug("Pocket Start Simulate Relay")
		}

		cli.PoktCtnCreateRun(runCmd, env, false)

	}

	if !cli.PoktCtnStart() {
		logs.GetLog().Error("Pocket Start FALSE")
	}
}

func (psvc *PoktService) StartScan() {
	url := psvc.PoktApiUrl
	height, err := api.PoktApiGetHeight(url)
	if err != nil {
		logs.GetLog().Error(err)
	}
	logs.GetLog().Info("Pokt Get Current Height=", height)

}

func (psvc *PoktService) SendPoktHeartbeatRequest(swanClient *swan.SwanClient) {

	params := ""
	confPokt := config.GetPoktConfig().Pokt
	selfUrl := utils.UrlJoin(confPokt.PoktServerApiUrl, API_POCKET_V1)

	apiUrl := utils.UrlJoin(selfUrl, "status")
	response, err := web.HttpGetNoToken(apiUrl, params)
	if err != nil {
		fmt.Printf("Heartbeat Get Pocket Status err: %s \n", err)
		return
	}

	res := &models.StatusResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Heartbeat Parse Response (%s) err: %s \n", response, err)
		return
	}
	// Swan Server Is Not Ready!
	//stat := swan.PocketHeartbeatOnlineParams{
	//	Address:     res.Data.Address,
	//	Version:     res.Data.Version,
	//	Height:      res.Data.Height,
	//	Balance:     res.Data.Balance,
	//	Award:       res.Data.Award,
	//	Jailed:      res.Data.Jailed,
	//	JailedBlock: res.Data.JailedBlock,
	//	JailedUntil: res.Data.JailedUntil,
	//}
	//err = swanClient.SendPoktHeartbeatRequest(stat)
	//if err != nil {
	//	fmt.Printf("Heartbeat Send err: %s \n", err)
	//	return
	//}

	{
		title := color.New(color.FgGreen).Sprintf("%s", "Version")
		value := color.New(color.FgYellow).Sprintf("%s", res.Data.Version)
		fmt.Printf("%s\t\t: %s\n", title, value)

		title = color.New(color.FgGreen).Sprintf("%s", "Height")
		value = color.New(color.FgYellow).Sprintf("%d", res.Data.Height)
		fmt.Printf("%s\t\t: %s\n", title, value)

		title = color.New(color.FgGreen).Sprintf("%s", "Address")
		value = color.New(color.FgYellow).Sprintf("%s", res.Data.Address)
		fmt.Printf("%s\t\t: %s\n", title, value)

		title = color.New(color.FgGreen).Sprintf("%s", "Balance")
		value = color.New(color.FgYellow).Sprintf("%d", res.Data.Balance)
		fmt.Printf("%s\t\t: %s\n", title, value)

		title = color.New(color.FgGreen).Sprintf("%s", "Jailed")
		value = color.New(color.FgRed).Sprintf("%t", res.Data.Jailed)
		fmt.Printf("%s\t\t: %s\n", title, value)

		title = color.New(color.FgGreen).Sprintf("%s", "JailedBlock")
		value = color.New(color.FgRed).Sprintf("%d", res.Data.JailedBlock)
		fmt.Printf("%s\t: %s\n", title, value)

		title = color.New(color.FgGreen).Sprintf("%s", "JailedUntil")
		value = color.New(color.FgRed).Sprintf("%s", res.Data.JailedUntil)
		fmt.Printf("%s\t: %s\n", title, value)
	}

	return
}

func HttpGetPoktVersion(c *gin.Context) {
	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecVersion()
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(cmdOut))
}

func HttpGetPoktCurHeight(c *gin.Context) {
	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecHeight()
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(cmdOut))
}

func HttpGetPoktNodeAddr(c *gin.Context) {
	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecNodeAddress()
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(cmdOut))
}

func HttpGetPoktStatus(c *gin.Context) {
	poktSvr := GetMyPoktService()
	data := &models.StatusData{}

	versionData, err := poktSvr.GetCli().PoktCtnExecVersion()
	if err != nil {
		logs.GetLog().Error(err)
	} else {
		data.Version = versionData.Version
	}

	heightData, err := poktSvr.GetCli().PoktCtnExecHeight()
	if err != nil {
		logs.GetLog().Error(err)
	} else {
		data.Height = heightData.Height
	}

	address, err := poktSvr.GetCli().PoktCtnExecNodeAddress()
	if err != nil {
		logs.GetLog().Error(err)
	} else {
		data.Address = address
	}

	balanceData, err := poktSvr.GetCli().PoktCtnExecBalance(address)
	if err != nil {
		logs.GetLog().Error(err)
	} else {
		data.Balance = balanceData.Balance
	}

	nodeData, err := poktSvr.GetCli().PoktCtnExecNode(address)
	if err != nil {
		logs.GetLog().Error(err)
	} else {
		data.Jailed = nodeData.Jailed
	}

	signData, err := poktSvr.GetCli().PoktCtnExecSignInfo(address)
	if err != nil || len(signData) == 0 {
		logs.GetLog().Error(err)
	} else {
		signInfo := signData[0]
		data.JailedUntil = signInfo.JailedUntil
		data.JailedBlock = signInfo.JailedBlocksCounter
	}

	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

func HttpGetPoktBalance(c *gin.Context) {
	var params models.BalancePoktParams
	err := c.BindJSON(&params)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	cmdOut, err := poktSvr.GetCli().PoktCtnExecBalance(params.Address)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	logs.GetLog().Debug("pocket query balance result:", cmdOut)

	data := &models.BalanceCmdData{
		Height:  params.Height,
		Address: params.Address,
		Balance: strconv.FormatUint(cmdOut.Balance, 10)}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

func HttpGetPoktThreshold(c *gin.Context) {
	var params models.ThresholdParams
	err := c.BindJSON(&params)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &models.ThresholdData{
		Address:   params.Address,
		Threshold: params.Threshold,
		Active:    true,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

///////////////////////////////////////////////////////////////////////////////

func HttpSetPoktValidator(c *gin.Context) {
	var params models.ValidatorParams
	err := c.BindJSON(&params)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	result, err := poktSvr.GetCli().PoktCtnExecSetValidator(
		params.Address,
		params.Passwd,
	)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &models.ValidatorData{
		Result: result,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

///////////////////////////////////////////////////////////////////////////////

func HttpSetPoktCustodial(c *gin.Context) {
	var params models.CustodialParams
	err := c.BindJSON(&params)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	result, err := poktSvr.GetCli().PoktCtnExecCustodial(
		params.Address,
		params.Amount,
		params.RelayChainIDs,
		params.ServiceURI,
		params.NetworkID,
		params.Fee,
		params.IsBefore,
		params.Passwd,
	)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &models.CustodialData{
		Result: result,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}

///////////////////////////////////////////////////////////////////////////////

func HttpSetPoktNonCustodial(c *gin.Context) {
	var params models.NonCustodialParams
	err := c.BindJSON(&params)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusOK, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	poktSvr := GetMyPoktService()
	result, err := poktSvr.GetCli().PoktCtnExecNonCustodial(
		params.PubKey,
		params.OutputAddr,
		params.Amount,
		params.RelayChainIDs,
		params.ServiceURI,
		params.NetworkID,
		params.Fee,
		params.IsBefore,
		params.Passwd,
	)
	if err != nil {
		logs.GetLog().Error(err)
		c.JSON(http.StatusInternalServerError, common.CreateErrorResponse("-1", err.Error()))
		return
	}

	data := &models.NonCustodialData{
		Result: result,
	}
	c.JSON(http.StatusOK, common.CreateSuccessResponse(data))
}
