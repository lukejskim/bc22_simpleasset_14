// dev/simplaeasset/1.0/simpleasset.go
// echo > simpleasset.go

// 패키지 정의
package main

// 1. 외부 모듈 포함
import (
	"fmt"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)
// 설치 : 
// . ~/.profile
// cd ~/dev/simpleasset/1.0
// go mod init simpleasset
// go get -u "github.com/hyperledger/fabric/protos/peer"
// go get -u "github.com/hyperledger/fabric/core/chaincode/shim@v1.4"

// strconv.Atoi("305") => 305
// strconv.Itoa(305)   => "305"


// 2. 체인코드 클래스-구조체 정의 SimpleAsset
type SimpleAsset struct {
}

type Asset struct {
	Key 	string `json:key`
	Value 	string `json:value`
}

// 3. Init 함수
// *SimpleAsset : Init함수가 위 구조체 SimpleAsset의 함수이다. 
// 체인코드가 처음 배포될때 호출됨
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("init success"))

	// args := stub.GetStringArgs()
	// if len(args) != 2 {
	// 	return shim.Error("Incorrect arguments. Expecting a key and a value")
	// }

	// err := stub.PutState(args[0], []byte(args[1]))
	// if err != nil {
	// 	return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	// 	// return shim.Error("Failed to create asset: "+args[0])
	// }

	// return shim.Success(nil)
}

// 4. Invoke 함수
// 모든 체인코드 호출시
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	fn, args := stub.GetFunctionAndParameters()
	
	if fn == "set" {
		return t.Set(stub, args)
	} else if fn == "get" {
		return t.Get(stub, args)
	} else if fn == "del" {
		return t.Del(stub, args)
	} else if fn == "transfer" {
		return t.Transfer(stub, args)
	} else if fn == "transfer" {
		return t.Transfer(stub, args)
	} else if fn == "transfer2" {
		return t.Transfer2(stub, args)
	}

	return shim.Error("Not supported function name")
}


// 5. set 함수
func (t *SimpleAsset) Set(stub shim.ChaincodeStubInterface, args []string) peer.Response{
	
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and value")
	}
	// 오류체크 중복 키 검사 -> 덮어쓰기로 해결

	// 마샬링 할 구조체 생성
	asset := Asset{Key:args[0], Value:args[1]}

	assetAsByte, err := json.Marshal(asset)

	if err != nil {
		return shim.Error("Failed to marshal arguments:" + args[0] + " " + args[1])
	}

	err = stub.PutState(args[0], assetAsByte)
	if err != nil {
		return shim.Error("Failed to set asset:" + args[0])
	}

	return shim.Success(assetAsByte)
}


// 6. get 함수 
// stub : peer를 통해서 전달    cf. docker images
func (t *SimpleAsset) Get(stub shim.ChaincodeStubInterface, args []string) peer.Response{

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Filed to get asset: " + args[0] + " with error: " + err.Error())
	}
	if value == nil {
		return shim.Error("Asset not found: " + args[0])
	}

	return shim.Success([]byte(value))
}

// 6.1 del 함수  
// Hint : stub.DelState(key)
func (t *SimpleAsset) Del(stub shim.ChaincodeStubInterface, args []string) peer.Response{

	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Filed to get asset: " + args[0] + " with error: " + err.Error())
	}
	if value == nil {
		return shim.Error("Failed to delete Asset - Asset not found: " + args[0])
	}

	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error("Filed to del asset: " + args[0] + " with error: " + err.Error())
	}

	return shim.Success([]byte(args[0]))
}


// 6.2 transfer 함수  
// A -> B -> 100 송금
// A -= 100
// B += 100
func (t *SimpleAsset) Transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response{
	// 인자값 체크
	if len(args) != 3 {
		return shim.Error("Incorrect arguments. Expecting a from_key, to_key, amount")
	}
	// args[0] : from_key
	// args[1] : to_key
	// args[2] : amount

	// 보내는 이, 받는이의 자산을 조회 => from_asset, to_asset
	// 보내는 이 조회
	from_asset, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Filed to get asset: " + args[0] + " with error: " + err.Error())
	}
	if from_asset == nil {
		return shim.Error("Asset not found: " + args[0])
	}

	// 받는 이 조회
	to_asset, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Filed to get asset: " + args[0] + " with error: " + err.Error())
	}
	if to_asset == nil {
		return shim.Error("Asset not found: " + args[0])
	}

	// from_asset, to_asset => unmarshal
	from := Asset{}
	to := Asset{}

	json.Unmarshal(from_asset, &from)
	json.Unmarshal(to_asset, &to)

	// from_amount, to_amount, amount 정수형으로 변환
	from_amount, _ := strconv.Atoi(from.Value)
	to_amount, _ := strconv.Atoi(to.Value)
	amount, _ := strconv.Atoi(args[2])
	
	// 잔액 검증
	if (from_amount < amount) {
		return shim.Error("Not  enough asset value " + args[0])
	}

	// 송금된 결과값 계산
	from.Value = strconv.Itoa(from_amount-amount)
	to.Value = strconv.Itoa(to_amount+amount)

	// from, to를 marshal
	from_asset, _ = json.Marshal(from)
	to_asset, _ = json.Marshal(to)

	// PutState
	stub.PutState(args[0], from_asset)
	stub.PutState(args[1], to_asset)

	return shim.Success([]byte("transfer done!"))
}

func (t *SimpleAsset) Transfer2(stub shim.ChaincodeStubInterface, args []string) peer.Response{

	if len(args) != 3 {
		return shim.Error("Incorrect arguments. Expecting a key")
	}

	// args Check
	sender, receiver, amount := args[0], args[1], args[2]

	s_state, err := stub.GetState(sender)
	if err != nil {
		return shim.Error("Not found sender: " + sender + " with error: " + err.Error())
	}
	r_state, err := stub.GetState(receiver)
	if err != nil {
		return shim.Error("Not found receivr: " + receiver + " with error: " + err.Error())
	}


	// 3. 잔액변환 및 검증, 전송수행
	// 마샬링 할 구조체 생성
	s_asset := Asset{}
	r_asset := Asset{}
	json.Unmarshal(s_state, &s_asset)
	json.Unmarshal(r_state, &r_asset)

	// 금액(int) : 송금자, 수신자, 송금액 
	s_amount, _ := strconv.Atoi(s_asset.Value)
	r_amount, _ := strconv.Atoi(r_asset.Value)
	n_amount, _ := strconv.Atoi(amount)

	// 검증 : 송금액보다 적은금액 보유시 오류
	if(s_amount < n_amount) {
		return shim.Error("Not enough asset value: "+sender)
	}

	// 4. Marshal
	s_amount -= n_amount
	r_amount += n_amount
	s_asset.Value = strconv.Itoa(s_amount)
	r_asset.Value = strconv.Itoa(r_amount)

	s_state, _ = json.Marshal(s_asset)
	r_state, _ = json.Marshal(r_asset)

	// 5. PutState
	stub.PutState(sender, s_state)
	stub.PutState(receiver, r_state)

	return shim.Success([]byte("Success Transfer: " + sender + "=>" + receiver + " : " + amount))

}





// 7. main 함수
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode : %s", err)
	}
}

