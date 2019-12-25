package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var _insuranceLogger = shim.NewLogger("INSURANCE")

const _CreateEvent = "CREATE_EVENT"
const _UpdateEvent = "UPDATE_EVENT"

//CarInsurance structure
type CarInsurance struct {

	ObjType             string `json:"carInsurance"`
	InsuranceID         string `json:"incuranceID"`
	TermsAndConditoion  string `json:"termsAndConditoion"`
	Count               string `json:"count"` 
	DualUser            string `json:"dualUser"` 
	
}

type Admin struct {

	ObjType             string `json:"admin"`
	OperatorName        string `json:"operatorName"`
	Collictible         string `json:"collictible"`
	Collected           string `json:"collected"` 
	
}

type User struct {
	ObjType             string   `json:"user"`
	UserID              string   `json:"userID"`
	UserName            string   `json:"userName"`
	UserAddress         string   `json:"userAddress"`
	UserContactNo       string   `json:"userContactNo"`
	UserEmail           string   `json:"userEmail"`
	UserGender          string   `json:"userGender"`
	UserDOB             string   `json:"userDOB"`
	VehicleType         string   `json:"vehicleType"`
	VehicleModel        string   `json:"vehicleModle"`
	ChassisNo           string   `json:"chassisNo"`
	DateOfRegistration  string   `json:"dateOfRegistration"`
	InsuranceID         string[] `json:"insuranceID"`
	InvoiceID           string[] `json:"invoiceID"`
	Payble              string   `json:"payble"`
	Paid                string   `json:"paid"`
	RegistraionDate     string   `json:"registrationDate"`
	UpdateDate          string   `json:"updateDate"`
}

type Invoices struct {
	ObjType             string   `json:"invoice"`
	InvoiceID           string   `json:"invoiceID"`
	Month               string   `json:"month"`          // Format (Jan,2019)
	InsuranceID         string   `json:"insuranceID"`
	Payble              string   `json:"payble"`
	Paid                string   `json:"paid"`
	BajajAllianzAmount  string   `json:"bajajAllianzAmount"`
	AllianzAmount       string   `json:"allianzAmount"`
	GenerateDate        string   `json:"generateDate"`
	UpdateDate          string   `json:"updateDate"`
}

//Insurance structure
type Insurance struct {
}

//global vars
var errorDetails, errKey, jsonResp, repError string

//Init function of chaincode
func (s *Insurance) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	var jsonResp string

	ObjType          := "admin"     
	OperatorName     := "BajajAllianz"  
	OperatorName1    := "Allianz"
	Collictible      := "0"   
	Collected        := "0"  
   



	// ==== Bajaj object and marshal to JSON ====
	BajajAllianzRegistartion        := &Admin{ObjType, OperatorName, Collictible, Collected}
	BajajAllianzJSONasBytes, err    := json.Marshal(BajajAllianzRegistartion)
	if err != nil {
		jsonResp = "{\"ErrorDetails\":\"Error in marshaling while registring Bajaj Allianz\"}"
			return shim.Error(jsonResp)
		}

    // ==== BajajAllianz object and marshal to JSON ====
	AllianzRegistartion            := &Admin{ObjType, OperatorName, Collictible, Collected}
	AllianzJSONasBytes, err        := json.Marshal(AllianzRegistartion)
	if err != nil {
		jsonResp = "{\"ErrorDetails\":\"Error in marshaling while registring Allianz\"}"
			return shim.Error(jsonResp)
		}
	

	// === Save Bajaj Allianz to state ===
	err = stub.PutState(OperatorName, BajajAllianzJSONasBytes)
	if err != nil {
		jsonResp = "{\"ErrorDetails\":\"Error while registring Bajaj Allianz\"}"
		return shim.Error(jsonResp)
	}

	// === Save Allianz to state ===
	err = stub.PutState(OperatorName, AllianzJSONasBytes)
	if err != nil {
		jsonResp = "{\"ErrorDetails\":\"Error while registring Allianz\"}"
		return shim.Error(jsonResp)
	}


}

//Invoke function of chaincode
func (s *Insurance) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	var response pb.Response
	// Handle different functions
	switch function {
	case "insuranceRegistration":
		response = s.insuranceRegistration(stub)
	case "getInsuranceByUserID":
		response = s.getInsuranceByUserID(stub)

	default:
		var jsonResp string
		jsonResp = "{\"ErrorDetails\":\" Received unknown function invocation \"}"
		return shim.Error(jsonResp)
	}
	return response
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//insurance registration
func (s *Insurance) userRegistration(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var userToSave User
	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	err := json.Unmarshal([]byte(args[0]), &userToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid JSON provided- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	getBytes, err := stub.GetState(userToSave.UserID)
	if err == nil || len(getBytes) > 0 {
		errKey = userToSave.UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Insurance with this ID already exist, provide unique UserID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	userToSave.ObjType = "user"
	userJSON, marshalErr := json.Marshal(userToSave)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("userRegister: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Saving Insurance to the ledger with InsuranceId = ", userToSave.UserID)
	err = stub.PutState(userToSave.UserID, userJSON)
	if err != nil {
		//errKey = string(userJSON)
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to save user with UserID- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("userRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	//setting event after storing into the ledger
	retErr := stub.SetEvent(_CreateEvent, insuranceJSON)
	if retErr != nil {
		repError = strings.Replace(retErr.Error(), "\"", " ", -1)
		errorDetails = "Event not generated for event : Insurance Registration- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("userRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"txID"   :      stub.GetTxID(),
		"userId" :      userToSave.UserID,
		"message":     "Insurance Registered Successfully",
		"status" :      "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

func (s *Insurance) getInsuranceByUserID(stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	var jsonResp string

	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("getInsuranceByUserID: " + jsonResp)
		return shim.Error(jsonResp)
	}

	UserID := strings.ToLower(args[0])


}


