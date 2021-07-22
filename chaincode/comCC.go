package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const DOC_TYPE = "comObj"

func PutCom(stub shim.ChaincodeStubInterface, com Comment) ([]byte, bool) {
	com.ObjectType = DOC_TYPE
	b, err := json.Marshal(com)
	if err != nil {
		return nil, false
	}
	// 保存com状态
	err = stub.PutState(com.ID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

func GetComInfo(stub shim.ChaincodeStubInterface, ID string) (Comment, bool) {
	var com Comment
	// 根据ID查询信息状态
	b, err := stub.GetState(ID)
	if err != nil {
		return com, false
	}

	if b == nil {
		return com, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &com)
	if err != nil {
		return com, false
	}

	// 返回结果
	return com, true
}

func getComByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n 查询部分", buffer.String())

	return buffer.Bytes(), nil

}

//使用CouchDB数据库
func (t *CommentChaincode) queryBySome(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	className := args[0]
	usrName := args[1]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"ClassName\":\"%s\", \"nick_name\":\"%s\"}}", DOC_TYPE, className, usrName)
	resultsIterator, err := stub.GetQueryResult(queryString)

	if err != nil {
		return shim.Error("Rich query failed")
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	buffer.WriteString(`{"result":[`)

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Fail")
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString(`]}`)
	fmt.Print("Query result: %s", buffer.String())

	return shim.Success(buffer.Bytes())
}

//使用CouchDB数据库查询评价分数
func (t *CommentChaincode) queryScoreCount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	className := args[0]
	// commentScore, _ := strconv.Atoi(args[1])
	commentScore := args[1]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"ClassName\":\"%s\", \"first_comment_score\":%s}}", DOC_TYPE, className, commentScore)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("Rich query failed")
	}
	defer resultsIterator.Close() //释放迭代器
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	buffer.WriteString(`{"result":[`)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next() //获取迭代器中的每一个值
		if err != nil {
			return shim.Error("Fail")
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value)) //将查询结果放入Buffer中
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString(`]}`)
	fmt.Print("Query result: %s", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (t *CommentChaincode) addCom(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	var com Comment
	err := json.Unmarshal([]byte(args[0]), &com)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}
	_, exist := GetComInfo(stub, com.ID)
	if exist {
		return shim.Error("要添加的评论序号已存在")
	}
	txid := stub.GetTxID()
	ts, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	var time = ts.String()
	s, err := stub.GetSignedProposal()
	if err != nil {
		return shim.Error(err.Error())
	}
	var signature = s.GetSignature()
	// var proposal = s.GetProposalBytes()
	bind, err := stub.GetBinding()
	if err != nil {
		return shim.Error(err.Error())
	}
	com.TxId = txid
	com.Timestamp = time
	com.Signature = signature
	com.Bind = bind
	_, bl := PutCom(stub, com)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

func (t *CommentChaincode) queryComByName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	className := args[0]
	usrName := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"comObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"ClassName\":\"%s\", \"nick_name\":\"%s\"}}", DOC_TYPE, className, usrName)

	// 查询数据
	result, err := getComByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据name查询信息时发生错误")
	}
	if result == nil || string(result) == "[]" {
		return shim.Error("根据指定的name没有查询到相关的信息")
	}
	return shim.Success(result)
}

func (t *CommentChaincode) queryComByScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	className := args[0]
	// commentScore, _ := strconv.Atoi(args[1])
	commentScore := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"comObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"ClassName\":\"%s\", \"first_comment_score\":%s}}", DOC_TYPE, className, commentScore)

	// 查询数据
	result, err := getComByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据score查询信息时发生错误")
	}
	if result == nil || string(result) == "[]" {
		return shim.Error("根据score没有查询到相关的信息")
	}
	return shim.Success(result)
}

func (t *CommentChaincode) queryTotalByScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	// commentScore, _ := strconv.Atoi(args[0])
	commentScore := args[0]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"comObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"first_comment_score\":%s}}", DOC_TYPE, commentScore)

	// 查询数据
	result, err := getComByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据score查询信息时发生错误")
	}
	if result == nil || string(result) == "[]" {
		return shim.Error("根据score没有查询到相关的信息")
	}
	return shim.Success(result)
}

func (t *CommentChaincode) queryByClass(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	cid := args[0]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"comObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"cid\":%s}}", DOC_TYPE, cid)

	// 查询数据
	result, err := getComByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据CID查询信息时发生错误")
	}
	if result == nil || string(result) == "[]" {
		return shim.Error("根据CID没有查询到相关的信息")
	}
	return shim.Success(result)
}

func (t *CommentChaincode) queryAll(stub shim.ChaincodeStubInterface) peer.Response {
	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"comObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", DOC_TYPE)

	// 查询数据
	result, err := getComByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据score查询信息时发生错误")
	}
	if result == nil || string(result) == "[]" {
		return shim.Error("根据score没有查询到相关的信息")
	}
	return shim.Success(result)
}

func (t *CommentChaincode) queryComInfoByID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	// 根据号码查询com状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据号码查询信息失败")
	}

	if b == nil {
		return shim.Error("根据号码没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var com Comment
	err = json.Unmarshal(b, &com)
	if err != nil {
		return shim.Error("反序列化com信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(com.ID)
	if err != nil {
		return shim.Error("根据指定的号码查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []HistoryItemN
	var hisCom Comment
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取com的历史变更数据失败")
		}

		var historyItem HistoryItemN
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisCom)

		if hisData.Value == nil {
			var empty Comment
			historyItem.Comment = empty
		} else {
			historyItem.Comment = hisCom
		}

		historys = append(historys, historyItem)

	}

	com.Historys = historys

	// 返回
	result, err := json.Marshal(com)
	if err != nil {
		return shim.Error("序列化com信息时发生错误")
	}
	return shim.Success(result)
}

func (t *CommentChaincode) updateCom(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var info Comment
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return shim.Error("反序列化com信息失败")
	}

	// 查询信息
	result, bl := GetComInfo(stub, info.ID)
	if !bl {
		return shim.Error("根据id号码查询信息时发生错误")
	}

	result.ID = info.ID
	result.UsrName = info.UsrName
	result.CommentTime = info.CommentTime
	result.CommentScore = info.CommentScore
	result.CommentDetail = info.CommentDetail
	result.StudyTime = info.StudyTime
	result.ClassName = info.ClassName
	result.Cid = info.Cid

	txid := stub.GetTxID()

	ts, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	var time = ts.String()

	s, err := stub.GetSignedProposal()
	if err != nil {
		return shim.Error(err.Error())
	}
	var signature = s.GetSignature()

	bind, err := stub.GetBinding()
	if err != nil {
		return shim.Error(err.Error())
	}

	result.TxId = txid
	result.Timestamp = time
	result.Signature = signature
	result.Bind = bind

	_, bl = PutCom(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

func (t *CommentChaincode) getHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	// 查询com状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据ID查询信息失败")
	}

	if b == nil {
		return shim.Error("根据ID没有查询到相关的信息")
	}

	result := stub.GetTxID()

	return shim.Success([]byte(result))
}
