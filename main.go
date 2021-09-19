package main

// @title Comments-Blockchain API
// @version 1.0
// @description This is a server for comments data of blockchain system.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host chen-v1:8000
// @BasePath /

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"net/http"

	"github.com/blockchain.com/comments/sdkInit"
	"github.com/blockchain.com/comments/service"
	"github.com/blockchain.com/comments/web"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//fabric版本1.4
const (
	configFile  = "config.yaml"
	initialized = false
	ComCC       = "comcc"
	chainflag   = false //是否首次启动
)

// 首次启动时创建通道及安装实例化链码
func chanCreate(sdk *fabsdk.FabricSDK, initInfo *sdkInit.InitInfo) *channel.Client {
	err := sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}

	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	return channelClient
}

//程序中断后再次启动
func chanExc(sdk *fabsdk.FabricSDK, initInfo *sdkInit.InitInfo) *channel.Client {
	clientChannelContext := sdk.ChannelContext(initInfo.ChannelID, fabsdk.WithUser(initInfo.UserName), fabsdk.WithOrg(initInfo.Org1Name))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Errorf("创建应用通道客户端失败: %v", err)
	}
	// fmt.Println(channelClient)
	return channelClient
}

//数据修改
func comMod(total int, com []service.Comment, serviceSetup service.ServiceSetup) {
	for i := 0; i < total; i++ {
		// com[i].OJID = com[i].OID.ID
		// fmt.Println(com[i])
		msg, err := serviceSetup.ModifyCom(com[i])
		// fmt.Println(com[i])
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("信息发布成功, 交易编号为: " + msg)
		}
	}
}

//读取课程
func readClass(url string) service.DataClass {
	var class service.DataClass
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	err = json.Unmarshal([]byte(body), &class)
	if err != nil {
		fmt.Println(err)
	}
	return class
}

//读取评论信息
func readCom(url string) []service.Comment {
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var Dcom service.DataCom
	var comc []service.Comment
	err := json.Unmarshal([]byte(body), &Dcom)
	if err != nil {
		fmt.Println(err)
	}
	comc = Dcom.Data
	return comc
}

//查询统计课程列表
func queryAllClass(serviceSetup service.ServiceSetup) (int, []string) {
	var classlist []string
	// classmap := make(map[int]string)
	Allcount := 0
	var com []service.Comment
	result, err := serviceSetup.FindAll()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		json.Unmarshal(result, &com)
		// fmt.Println("查询所有信息成功：")
		classlist = append(classlist, com[0].ClassName)
		// classmap[com[0].Cid] = com[0].ClassName
		for i := 0; i < len(com); i++ {
			flag := 0
			for j := 0; j < len(classlist); j++ {
				if classlist[j] == com[i].ClassName {
					flag = 1
				}
			}
			if flag != 1 {
				classlist = append(classlist, com[i].ClassName)
				// classmap[com[i].Cid] = com[i].ClassName
			}
		}
		Allcount = len(com)
		// fmt.Printf("共有: %d条\n", len(com))
		// fmt.Println("classlist:", classlist)
	}
	// fmt.Println(com)
	// fmt.Println(classmap)
	return Allcount, classlist
}

//按分数查询统计课程列表
func queryClassByScore(score int, serviceSetup service.ServiceSetup) (int, []string) {
	var classlist []string
	Allcount := 0
	result, err := serviceSetup.FindTotalScore(strconv.Itoa(score))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var com []service.Comment
		json.Unmarshal(result, &com)
		fmt.Println("查询所有信息成功：")
		classlist = append(classlist, com[0].ClassName)
		for i := 0; i < len(com); i++ {
			flag := 0
			for j := 0; j < len(classlist); j++ {
				if classlist[j] == com[i].ClassName {
					flag = 1
				}
			}
			if flag != 1 {
				classlist = append(classlist, com[i].ClassName)
			}
		}
		Allcount = len(com)
		// fmt.Printf("共有: %d条\n", len(com))
		// fmt.Println("classlist:", classlist)
	}
	return Allcount, classlist
}

//统计总共的评分情况
func queryTotalScore(serviceSetup service.ServiceSetup, c chan service.BlockCount) {
	var scorecount service.BlockCount
	scorecount.Name = "总计"
	for j := 0; j < 5; j++ {
		score := strconv.Itoa(j + 1)
		result, _ := serviceSetup.FindTotalScore(score)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		var com []service.Comment
		json.Unmarshal(result, &com)

		if j == 4 {
			scorecount.Score.Five = len(com)
		} else if j == 3 {
			scorecount.Score.Four = len(com)
		} else if j == 2 {
			scorecount.Score.Three = len(com)
		} else if j == 1 {
			scorecount.Score.Two = len(com)
		} else if j == 0 {
			scorecount.Score.One = len(com)
		}
		// }
	}
	c <- scorecount
}

func queryScore(classlist []string, serviceSetup service.ServiceSetup, c chan []service.BlockCount) {
	//统计各店的评分情况
	var ScoreGroup []service.BlockCount
	comlen := 0
	for i := 0; i < len(classlist); i++ {
		var scount service.BlockCount
		scount.Name = classlist[i]
		for j := 0; j < 5; j++ {
			score := strconv.Itoa(j + 1)
			result, _ := serviceSetup.FindCommentScore(classlist[i], score)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			var com []service.Comment
			json.Unmarshal(result, &com)
			if j == 4 {
				scount.Score.Five = len(com)
			} else if j == 3 {
				scount.Score.Four = len(com)
			} else if j == 2 {
				scount.Score.Three = len(com)
			} else if j == 1 {
				scount.Score.Two = len(com)
			} else if j == 0 {
				scount.Score.One = len(com)
			}
			comlen += len(com)
			// }
		}
		ScoreGroup = append(ScoreGroup, scount)
	}
	c <- ScoreGroup
}

func traceQuery(Allcount int, serviceSetup service.ServiceSetup, c chan []service.Comment) {
	//朔源数据
	// var Trace []service.Trace
	var TraceCom []service.Comment
	for i := 0; i < Allcount; i++ {
		count := strconv.Itoa(i + 1)
		result, err := serviceSetup.FindComInfoByID(count)
		if err != nil {
			fmt.Println(err.Error())
		}
		var comment service.Comment
		json.Unmarshal(result, &comment)
		if err == nil {
			// data := service.Trace{Com: comment, Index: count}
			// Trace = append(Trace, data)
			TraceCom = append(TraceCom, comment)
		}
	}
	c <- TraceCom
}

func classSort(classmap map[int]string, serviceSetup service.ServiceSetup) map[string][]service.Comment {
	mapcom := make(map[string][]service.Comment)
	//存储各店名的评论
	for key, _ := range classmap {
		key := strconv.Itoa(key)
		result, err := serviceSetup.FindByClass(key)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("结果：", mapcom)
		var com []service.Comment
		json.Unmarshal(result, &com)
		mapcom[key] = com
	}
	return mapcom
}

func traceSort(tracecom []service.Comment, serviceSetup service.ServiceSetup) map[int][]service.Comment {
	mapcom := make(map[int][]service.Comment)
	for _, item := range tracecom {
		mapcom[item.Cid] = append(mapcom[item.Cid], item)
	}
	return mapcom
}

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID:     "mychannel",
		ChannelConfig: "/home/star/go/src/github.com/blockchain.com/comments/fixtures/network-base/channel-artifacts/channel.tx",

		OrgAdmin:       "Admin",
		Org1Name:       "Org1",
		Org2Name:       "Org2",
		OrdererOrgName: "orderer.example.com",

		ChaincodeID:     ComCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/blockchain.com/commentsV2_edu/chaincode/",
		UserName:        "User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	//===========================================//
	//初始化
	var channelClient *channel.Client
	if chainflag {
		channelClient = chanCreate(sdk, initInfo)
	} else {
		channelClient = chanExc(sdk, initInfo)
	}

	//===========================================//

	serviceSetup := service.ServiceSetup{
		ChaincodeID: ComCC,
		Client:      channelClient,
	}

	//模拟数据读取
	url := "http://docker.hbuter.com:8012/course/?size=100"
	class := readClass(url)
	var com []service.Comment
	for _, item := range class.Data {
		url2 := "http://docker.hbuter.com:8012/course/" + strconv.Itoa(item.Cid) + "/comments?size=5000"
		comc := readCom(url2)
		com = append(com, comc...)
	}
	classmap := make(map[int]string)
	agencymap := make(map[int]string)
	for _, item := range class.Data {
		classmap[item.Cid] = item.ClassName
		agencymap[item.Cid] = item.AgencyName
	}

	//启动区块链程序写入模拟数据
	writeflag := false
	if writeflag {
		total := len(com)
		//数据上链
		for i := 0; i < total; i++ {
			// com[i].OJID = com[i].OID.ID
			com[i].ID = strconv.Itoa(i + 1)
			com[i].ClassName = classmap[com[i].Cid]
			com[i].AgencyName = agencymap[com[i].Cid]
			// fmt.Println(com[i])
			msg, err := serviceSetup.SaveCom(com[i])
			// fmt.Println(com[i])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("信息发布成功, 交易编号为: " + msg)
			}
		}
		//数据更新,模拟评论修改
		var com2 = make([]service.Comment, len(com))
		copy(com2, com)
		for idx, _ := range com2 {
			com2[idx].CommentDetail = "追加评价(模拟评价修改)"
		}
		// fmt.Println(com2)
		comMod(total, com2, serviceSetup)
		comMod(total, com, serviceSetup)
	}

	//链上数据查询及统计分类=============================//

	c1 := make(chan []service.Comment)
	c2 := make(chan service.BlockCount)
	c3 := make(chan []service.BlockCount)
	defer close(c1)
	defer close(c2)
	defer close(c3)

	//课程数量及名称查询
	Allcount, classlist := queryAllClass(serviceSetup)
	//统计总共的评分情况
	go queryTotalScore(serviceSetup, c2)
	ScoreTotal := <-c2

	//统计各店的评分情况
	go queryScore(classlist, serviceSetup, c3)
	ScoreGroup := <-c3

	//朔源数据
	go traceQuery(Allcount, serviceSetup, c1)
	TraceCom := <-c1

	//溯源查询结果分类
	mapcom := traceSort(TraceCom, serviceSetup)
	//查询结果按课程名分类

	//===========================================//
	//启动后端网路服务
	defer web.WebStart(ScoreTotal, ScoreGroup, TraceCom, class.Data, mapcom)
}
