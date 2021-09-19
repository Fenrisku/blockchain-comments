package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type CommentChaincode struct {
}

func (t *CommentChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *CommentChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()
	fmt.Printf("invokerN")
	if fun == "addCom" {
		return t.addCom(stub, args)
	} else if fun == "queryComByName" {
		return t.queryComByName(stub, args)
	} else if fun == "queryComInfoByID" {
		return t.queryComInfoByID(stub, args)
	} else if fun == "updateCom" {
		return t.updateCom(stub, args)
	} else if fun == "getHistory" {
		return t.getHistory(stub, args)
	} else if fun == "queryComByScore" {
		return t.queryComByScore(stub, args)
	} else if fun == "queryTotalByScore" {
		return t.queryTotalByScore(stub, args)
	} else if fun == "queryByClass" {
		return t.queryByClass(stub, args)
	} else if fun == "queryAll" {
		return t.queryAll(stub)
	}
	return shim.Error("指定的函数名称错误")

}

func main() {
	err := shim.Start(new(CommentChaincode))
	if err != nil {
		fmt.Printf("启动CommentChaincode时发生错误: %s", err)
	}
}
