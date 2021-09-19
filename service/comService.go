package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) SaveCom(com Comment) (string, error) {
	// fmt.Printf(req)
	// fmt.Printf("no error:1")
	eventID := "eventAddCom"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)

	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将com对象序列化成为字节数组
	// fmt.Printf("no error:2")
	b, err := json.Marshal(com)
	if err != nil {
		return "", fmt.Errorf("指定的com对象序列化时发生错误")
	}
	// fmt.Printf("no error:3")
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addCom", Args: [][]byte{b, []byte(eventID)}}
	// fmt.Printf("no error:a")
	// fmt.Printf(req)
	respone, err := t.Client.Execute(req)
	// fmt.Printf("no error:b")
	if err != nil {
		// fmt.Printf("this is fault")
		return "", err
	}
	// fmt.Printf("no error:4")
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	fmt.Printf(string(respone.Payload))
	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) FindComInfoByID(ID string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryComInfoByID", Args: [][]byte{[]byte(ID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// func (t *ServiceSetup) FindBlockInfoByID(ID string) (string, error) {
func (t *ServiceSetup) FindBlockInfoByID(ID string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getHistory", Args: [][]byte{[]byte(ID)}}
	respone, err := t.Client.Query(req)
	// if err != nil {
	// 	return "", err
	// }

	// return string(respone.Payload), nil

	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) FindHistory(ID string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getHistory", Args: [][]byte{[]byte(ID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return []byte(respone.TransactionID), nil
}

func (t *ServiceSetup) FindComByName(className, usrName string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryComByName", Args: [][]byte{[]byte(className), []byte(usrName)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) FindCommentScore(className, commentScore string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryComByScore", Args: [][]byte{[]byte(className), []byte(commentScore)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) FindTotalScore(commentScore string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryTotalByScore", Args: [][]byte{[]byte(commentScore)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) FindByClass(className string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryByClass", Args: [][]byte{[]byte(className)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) FindAll() ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryAll"}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) ModifyCom(com Comment) (string, error) {

	eventID := "eventModifyCom"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将com对象序列化成为字节数组
	b, err := json.Marshal(com)
	if err != nil {
		return "", fmt.Errorf("指定的com对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateCom", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}
