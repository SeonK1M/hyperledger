package main

// import modules
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// define class
type SmartContract struct {
}

type Developer struct {
	ID                  string  `json:"id"`
	Avg                 float32 `json:"avg"`
	NumberOfProject     int     `json:"numproject"`
	CurrentState        int     `json:"state"`
	CurrentProject      string  `json:"pname"`
	CurrentProjectScore int     `json:"pscore"`
}

type Project struct {
	Name          string   `json:"name"`
	ProjectState  int      `json:"projectstate"`
	DeveloperList []string `json:"devlist"`
}

const (
	REGISTERED int = iota
	WORKING
	FINISHED
)

// TODO: 1. define project structure
// func init
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response { //모든 데이터를 초기화하기 위해 체인 코드 인스턴스화 중에 Init이 호출됨
	//when you want to instantiate or upgrade chaincode
	return shim.Success(nil)
}

// func invoke
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response { //트랜잭션 요청을 받을 때 호출되는 함수, (stub(인자 이름), shim.CSIf(자료형태))
	fn, args := stub.GetFunctionAndParameters()

	if fn == "registerUser" { //registerUser요청이면 registerUser 함수를 수행하라
		return s.registerUser(stub, args)
	} else if fn == "joinProject" {
		return s.joinProject(stub, args)
	} else if fn == "recordScore" {
		return s.recordScore(stub, args)
	} else if fn == "readUserInfoById" {
		return s.readUserInfoById(stub, args)
	} else if fn == "readProjectInfoByName" {
		return s.readProjectInfoByName(stub, args)
	} else if fn == "registerProject" {
		return s.registerProject(stub, args)
	} else if fn == "finalizeProject" {
		return s.finalizeProject(stub, args)
	}
	return shim.Error("wrong function name")

}

// funcs for engineer
// registerUser
// recordScore
// joinProject
func (s *SmartContract) readUserInfoById(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 { //함수를 호출한 유저로부터 파라미터를 받는다. 하나만 받는다.
		return shim.Error("function readUserInfoById: wrong args")
	}
	info, _ := stub.GetState(args[0])
	return shim.Success(info)
}

func (s *SmartContract) registerUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// check args
	if len(args) != 1 { //함수를 호출한 유저로부터 파라미터를 받는다. 하나만 받는다.
		return shim.Error("function registerUser: wrong args")
	}
	var dev = Developer{ID: args[0], Avg: 0., CurrentState: 0} //Developer의 인자 세 개를 배열의 형태로 저장하는 변수 dev를 선언한다
	devAsBytes, _ := json.Marshal(dev)                         //dev를 마샬하여 json파일 형식으로 저장한 변수 devAsBytes를 선언한다.
	stub.PutState(args[0], devAsBytes)                         //마샬한 데이터를 월드스테이트에 저장한다.

	return shim.Success([]byte("transaction has been submited"))
}

func (s *SmartContract) recordScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 { //함수를 호출한 유저로부터 파라미터 세 개를 받는다
		return shim.Error("function recordScore: wrong args")
	}
	devAsBytes, err := stub.GetState(args[0]) //월드스테이트에서 첫번째 인자를 호출한 뒤 devAsBytes변수에 저장 및 선언함.
	if err != nil {
		return shim.Error("function recordScore: error from GetState")
	} // if getstate bring error
	if devAsBytes == nil { //if there's no id
		return shim.Error("no ID")
	}
	var dev = Developer{}                //developer배열을 가지는 변수 dev 선언
	_ = json.Unmarshal(devAsBytes, &dev) //월드스테이트에서 가져온 devBytes를 언마샬함

	dev.CurrentProject = args[1]                       //현재 프로젝트를 가져옴
	dev.CurrentProjectScore, _ = strconv.Atoi(args[2]) //현재 프로젝트 점수를 정수형태로

	var newAvg float32                                                                                                  //새로운 평균점수는 소수점 형태
	newAvg = (dev.Avg*float32(dev.NumberOfProject) + float32(dev.CurrentProjectScore)) / float32(dev.NumberOfProject+1) //계산

	dev.Avg = newAvg     //평균 점수를 갱신한다.
	dev.CurrentState = 2 //FINISHED

	devAsBytes, _ = json.Marshal(dev)  //데이터를 마샬해서
	stub.PutState(args[0], devAsBytes) //월드스테이트에 갖다 박는다.

	return shim.Success([]byte("transaction has been submited"))
}

func (s *SmartContract) joinProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// check args
	if len(args) != 2 { //아이디, 프로젝트이름
		return shim.Error("function joinProject: wrong args")
	}
	devAsBytes, err := stub.GetState(args[0]) //월드스테이트에서 첫번째 인자를 호출한 뒤 devAsBytes변수에 저장 및 선언함.
	if err != nil {
		return shim.Error("function joinProject: error from GetState")
	} // if getstate bring error
	if devAsBytes == nil { //if there's no id
		return shim.Error("no user ID")
	}
	projAsBytes, err := stub.GetState(args[1]) //월드스테이트에서 두번째 인자를 호출한 뒤 devAsBytes변수에 저장 및 선언함.
	if err != nil {
		return shim.Error("function joinProject: error from GetState")
	} // if getstate bring error
	if projAsBytes == nil { //if there's no id
		return shim.Error("no project Name")
	}

	var dev = Developer{}                //developer배열을 가지는 변수 dev 선언
	_ = json.Unmarshal(devAsBytes, &dev) //월드스테이트에서 가져온 devBytes를 언마샬함

	dev.CurrentProject = args[1] //현재 프로젝트를 가져옴
	dev.CurrentState = 1

	var proj = Project{}
	_ = json.Unmarshal(projAsBytes, &proj)

	proj.ProjectState = 1
	proj.DeveloperList = append(proj.DeveloperList, args[0])

	devAsBytes, _ = json.Marshal(dev) //데이터를 마샬해서
	projAsBytes, _ = json.Marshal(proj)
	stub.PutState(args[0], devAsBytes) //월드스테이트에 갖다 박는다.
	stub.PutState(args[1], projAsBytes)

	return shim.Success([]byte("transaction has been submited"))
}

// funcs for Project
// registerProject
// finishProject
func (s *SmartContract) readProjectInfoByName(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 { //함수를 호출한 유저로부터 파라미터를 받는다. 하나만 받는다.
		return shim.Error("function readProjectInfoByName: wrong args")
	}
	info, _ := stub.GetState(args[0])
	return shim.Success(info)
}
func (s *SmartContract) registerProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// check args
	if len(args) != 1 { //함수를 호출한 유저로부터 파라미터를 받는다. 하나만 받는다.
		return shim.Error("function registerProject: wrong args")
	}
	var proj = Project{Name: args[0], ProjectState: 0, DeveloperList: make([]string, 0)} //Developer의 인자 세 개를 배열의 형태로 저장하는 변수 dev를 선언한다
	projAsBytes, _ := json.Marshal(proj)                                                 //dev를 마샬하여 json파일 형식으로 저장한 변수 devAsBytes를 선언한다.
	stub.PutState(args[0], projAsBytes)                                                  //마샬한 데이터를 월드스테이트에 저장한다.

	return shim.Success([]byte("transaction has been submited"))
}

func (s *SmartContract) finalizeProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// check args
	if len(args) != 1 { //함수를 호출한 유저로부터 파라미터를 받는다. 하나만 받는다.
		return shim.Error("function finalizeProject: wrong args")
	}
	projAsBytes, err := stub.GetState(args[0]) //월드스테이트에서 첫번째 인자를 호출한 뒤 projAsBytes변수에 저장 및 선언함.
	if err != nil {
		return shim.Error("finalizeProject: error from GetState")
	} // if getstate bring error
	if projAsBytes == nil { //if there's no id
		return shim.Error("no Name")
	}

	var proj = Project{}                   //project배열을 가지는 변수 dev 선언
	_ = json.Unmarshal(projAsBytes, &proj) //월드스테이트에서 가져온 projBytes를 언마샬함

	proj.ProjectState = 2
	for i := 0; i < len(proj.DeveloperList); i++ {

		var developer = proj.DeveloperList[i]
		devAsBytes, err := stub.GetState(developer) //월드스테이트에서 첫번째 인자를 호출한 뒤 projAsBytes변수에 저장 및 선언함.
		if err != nil {
			return shim.Error("finalizeProject: error from GetState")
		} // if getstate bring error
		if devAsBytes == nil { //if there's no id
			return shim.Error("no Id")
		}
		var dev = Developer{}
		_ = json.Unmarshal(devAsBytes, &dev)

		dev.CurrentProject = ""
		dev.CurrentState = 2
		dev.NumberOfProject++ //프로젝트 종료 이후 스코어를 저장하는 함수 이므로 완수한 프로젝트의 개수를 하나 늘린다.

		devAsBytes, _ = json.Marshal(dev)    //데이터를 마샬해서
		stub.PutState(developer, devAsBytes) //월드스테이트에 갖다 박는다.
	}
	projAsBytes, _ = json.Marshal(proj) //데이터를 마샬해서
	stub.PutState(args[0], projAsBytes) //월드스테이트에 갖다 박는다.

	return shim.Success([]byte("transaction has been submited"))
}

// func main
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("Error while creating new SmartContract")
	}
}
