package service

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type DataClass struct {
	// Msg  string  `json:"msg"`
	Data []Class `json:"data"`
}

type Class struct {
	ObjectType               string `json:"docType"`
	ClassName                string `json:"name"`
	Cid                      int    `json:"cid"`
	AgencyName               string `json:"agency_name"`
	CommentGoodNum           int    `json:"comment_good_num"`
	Comment_medium_num       int    `json:"comment_medium_num"`
	Comment_bad_num          int    `json:"comment_bad_num"`
	Comment_all_num          int    `json:"comment_all_num"`
	Comment__good_percentage int    `json:"comment_good_percentage"`
}

type DataCom struct {
	// Msg  string  `json:"msg"`
	Data []Comment `json:"data"`
}

type Comment struct {
	ObjectType string `json:"docType"`
	Cid        int    `json:"cid"`
	UsrName    string `json:"nick_name"`
	Score      int    `json:"first_comment_score"`

	CommentDetail string `json:"first_comment"`
	CommentTime   int    `json:"first_comment_time"`
	StudyTime     int    `json:first_comment_study_time`

	ClassName  string
	AgencyName string
	// CommentScore string
	ID        string
	TxId      string
	Timestamp string
	Signature string
	Bind      []byte

	Historys []HistoryItemN // 当前com的历史记录
}

type HistoryItemN struct {
	TxId    string
	Comment Comment
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

type Trace struct {
	Com   Comment
	Index string
}

type TotalComments struct {
	Com  []Comment
	Name string
}

type BlockCount struct {
	Name  string
	Score ScoreCount
}

type ScoreCount struct {
	Five  int
	Four  int
	Three int
	Two   int
	One   int
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
