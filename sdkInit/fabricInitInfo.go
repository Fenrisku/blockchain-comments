package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

type InitInfo struct {
	ChannelID      string
	ChannelConfig  string
	OrgAdmin       string
	Org1Name       string
	Org2Name       string
	OrdererOrgName string
	OrgResMgmt     *resmgmt.Client
	OrgResMgmt2    *resmgmt.Client

	ChaincodeID     string
	ChaincodeGoPath string
	ChaincodePath   string
	UserName        string
}
