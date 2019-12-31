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
	Count               int    `json:"count"`                    
	DualUserCount       int    `json:"dualUserCount"`
	TotalAmount         int    `json:"totalAmount"` 
	PaidAmount          int    `json:"paidAmount"`
	GenerateDate        string `json:"generateDate"`
	
}

type Admin struct {

	ObjType             string `json:"admin"`
	OperatorName        string `json:"operatorName"`
	Collictible         int    `json:"collictible"`
	Collected           int    `json:"collected"` 
	
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
	InsuranceID      [100]string `json:"insuranceID"`
	InvoiceID       [1000]string `json:"invoiceID"`
	Payble              int      `json:"payble"`
	Paid                int      `json:"paid"`
	Bank1               string   `json:"bank1"`
	Bank2               string   `json:"bank2"`
	DualAcc             bool     `json:"dualAcc"`
	RegistraionDate     string   `json:"registrationDate"`
}

type Invoice struct {
	ObjType                   string   `json:"invoice"`
	InvoiceID                 string   `json:"invoiceID"`
	UserID                    string   `json:"userID"`
	InsuranceID               string   `json:"insuranceID"`
	Month                     string   `json:"month"`          // Format (Jan,2019)
	Payble                    int      `json:"payble"`
	Paid                      int      `json:"paid"`
	BajajAllianzPaybleAmount  int      `json:"bajajAllianzPaybleAmount"`
	AllianzPaybleAmount       int      `json:"allianzPaybleAmount"`
	BajajAllianzPaidAmount    int      `json:"bajajAllianzPaidAmount"`
	AllianzPaidAmount         int      `json:"allianzPaidAmount"`
	GenerateDate              string   `json:"generateDate"`
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
	Collictible      := 0   
	Collected        := 0  
   



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
	case "userRegistration":
		response = s.userRegistration(stub)
    case "buyInsurance":
		response = s.buyInsurance(stub)
	case "getInsuranceByUserID":
		response = s.getInsuranceByUserID(stub)
	case "generateInvoice":
		response = s.generateInvoice(stub)
	case "getInvoiceByInvoiceID":
		response = s.getInvoiceByInvoiceID(stub)
	case "getInvoiceByUserID":
		response = s.getInvoiceByUserID(stub)
	case "getInsuranceByUserID":
		response = s.getInsuranceByUserID(stub)

	default:
		var jsonResp string
		jsonResp = "{\"ErrorDetails\":\" Received unknown function invocation \"}"
		return shim.Error(jsonResp)
	}
	return response
}

// insurance regestration (Input Care Insurance schema)

func (s *Insurance) insuranceRegistration(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var insuranceToSave CarInsurance
	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	err := json.Unmarshal([]byte(args[0]), &insuranceToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid JSON provided- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	getBytes, err := stub.GetState(insuranceToSave.InsuranceID)
	if err == nil || len(getBytes) > 0 {
		errKey = insuranceToSave.InsuranceID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Insurance with this ID already exist"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	insuranceToSave.ObjType = "carInsurance"
	insuranceJSON, marshalErr := json.Marshal(insuranceToSave)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegestration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Saving Insurance to the ledger with InsuranceId = ", insuranceToSave.InsuranceID)
	err = stub.PutState(insuranceToSave.InsuranceID, insuranceJSON)
	if err != nil {
		//errKey = string(userJSON)
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to save user with UserID- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("insuranceRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	
	resultData := map[string]interface{}{
		"txID"   :      stub.GetTxID(),
		"userId" :      insuranceToSave.InsuranceID,
		"message":     "Insurance Registered Successfully",
		"status" :      "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}


//user registration (User schema)

func (s *Insurance) userRegistration(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var userToSave User
	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("userRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	err := json.Unmarshal([]byte(args[0]), &userToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid JSON provided- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("userRegistration: " + jsonResp)
		return shim.Error(jsonResp)
	}
	getBytes, err := stub.GetState(userToSave.UserID)
	if err == nil || len(getBytes) > 0 {
		errKey = userToSave.UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "User with this ID already exist, provide unique UserID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("userRegistration: " + jsonResp)
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
	_insuranceLogger.Info("Saving Insurance to the ledger with User ID = ", userToSave.UserID)
	err = stub.PutState(userToSave.UserID, userJSON)
	if err != nil {
		//errKey = string(userJSON)
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to save user with UserID- " + repError
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



//buy insurance (User id and insurance id)

func (s *Insurance) buyInsurance(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 2 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	UserID      := args[0];
	InsuranceID := args[1];
	var userToSave User
	var insuranceToSave CarInsurance

	getUser, err := stub.GetState(UserID)
	if err != nil  {
		errKey = UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid UserID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	getInsurance, err := stub.GetState(InsuranceID)
	if err != nil  {
		errKey = InsuranceID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid InsuranceID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	
	err := json.Unmarshal([]byte(getUser, &userToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Error in unmarshling user data" + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	err := json.Unmarshal([]byte(getInsurance, &insuranceToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Error in unmarshling insurance data" + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	userToSave.InsuranceID = append(InsuranceID)
	if (userToSave.DualAcc){
		insuranceToSave.DualUserCount = insuranceToSave.DualUserCount++;
	}
	insuranceToSave.Count = insuranceToSave.Count++;
	

	userJSON, marshalErr := json.Marshal(userToSave)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Updating insurance id of User ID = ", userToSave.UserID)
	err = stub.PutState(userToSave.UserID, userJSON)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to user with User- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}


	insuranceJSON, marshalErr := json.Marshal(insuranceToSave)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Updating insurance with insurance id = ", insuranceToSave.InsuranceID)
	err = stub.PutState(InsuranceID.InsuranceID, insuranceJSON)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to update with Insurance " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	resultData := map[string]interface{}{
		"txID"   :      stub.GetTxID(),
		"message":     "User had buyed the insurance",
		"status" :      "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}


// get insurance by user id (input user id)

func (s *Insurance) getInsuranceByUserID(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var err error
	var jsonResp string

	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("getInsuranceByUserID: " + jsonResp)
		return shim.Error(jsonResp)
	}

	UserID := args[0]
	var userToSave User
	var insuranceToSave CarInsurance

	getUser, err := stub.GetState(UserID)
	if err != nil  {
		errKey = UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid UserID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	err := json.Unmarshal([]byte(getUser, &userToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Error in unmarshling user data" + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	InsuranceID := userToSave.InsuranceID

	if (!InsuranceID[0]){

		resultData := map[string]interface{}{
			"txID"   :      stub.GetTxID(),
			"message":     "User have no insurance",
			"status" :      "true",
			"data"   :       ""
		}
		respJSON, _ := json.Marshal(resultData)
		return shim.Success(respJSON)

	}else{
		var data[len(InsuranceID.length]byte;
		for (i=0;i<len(InsuranceID.length);i++){
			getInsurance, err := stub.GetState(InsuranceID[i])
			if err != nil  {
				errKey = InsuranceID
				repError = strings.Replace(err.Error(), "\"", " ", -1)
				errorDetails = "Invalid InsuranceID"
				jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
				_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
				return shim.Error(jsonResp)
			}else{
				data = append(getInsurance)
				
				}

				resultData := map[string]interface{}{
					"txID"   :      stub.GetTxID(),
					"message":     "Success",
					"status" :      "true",
					"data"   :      data
				}
				respJSON, _ := json.Marshal(resultData)
				return shim.Success(respJSON)	
			}


		}
	}


// generateInvoice for user (input user id and insaurance id)

func (s *Insurance) generateInvoice(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var invoiceToSave Invoice
	var err error
	var jsonResp string

	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("generateInvoice: " + jsonResp)
		return shim.Error(jsonResp)
	}
	err := json.Unmarshal([]byte(args[0]), &invoiceToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid JSON provided- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("generateInvoice: " + jsonResp)
		return shim.Error(jsonResp)
	}
	getBytes, err := stub.GetState(invoiceToSave.InvoiceID)
	if err == nil || len(getBytes) > 0 {
		errKey = invoiceToSave.InvoiceID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Insurance with this ID already exist"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("generateInvoice: " + jsonResp)
		return shim.Error(jsonResp)
	}
	invoiceToSave.ObjType = "invoice"
	invoiceJSON, marshalErr := json.Marshal(invoiceToSave)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("generateInvoice: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Saving Invoice to the ledger with Invoice Id  = ", invoiceToSave.InvoiceID)
	err = stub.PutState(insuranceToSave.InsuranceID, insuranceJSON)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to save user with UserID- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("generateInvoice: " + jsonResp)
		return shim.Error(jsonResp)
	}

	InsuranceID := invoiceToSave.InsuranceID
	UserID := insuranceToSave.UserID
	var userToSave User
	var insuranceToSave CarInsurance

	getUser, err := stub.GetState(UserID)
	if err != nil  {
		errKey = UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid UserID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("generateInvoice: " + jsonResp)
		return shim.Error(jsonResp)
	}

	err := json.Unmarshal([]byte(getUser, &userToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Error in unmarshling user data" + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("generateInvoice: " + jsonResp)
		return shim.Error(jsonResp)
	}

	userToSave.InvoiceID = append(invoiceToSave.InvoiceID)
	userToSAVE.Payble = userToSAVE.Payble + invoiceToSave.Payble
	userToSave.Paid = userToSave.Paid + invoiceToSave.Paid

	userJSON, marshalErr := json.Marshal(userToSave)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Updating invoice detaile of User ID = ", userToSave.UserID)
	err = stub.PutState(userToSave.UserID, userJSON)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to user with User- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}


	
	getInsurance, err := stub.GetState(InsuranceID)
	if err != nil  {
		errKey = UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid InsuranceID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	err := json.Unmarshal([]byte(getInsurance, &insuranceToSave)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Error in unmarshling user data" + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	insuranceToSave.PaidAmount = insuranceToSave.PaidAmount + invoiceToSave.Paid

    insuranceJSON, marshalErr := json.Marshal(insuranceToSave)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Updating invoice details with insurance id = ", insuranceToSave.InsuranceID)
	err = stub.PutState(InsuranceID.InsuranceID, insuranceJSON)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to update with Insurance " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var bajajAllianz Admin
	var allianz Admin
	getAdminBajaj, err := stub.GetState("BajajAllianz")
	if err != nil  {
		errKey = UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid InsuranceID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	err := json.Unmarshal([]byte(getAdminBajaj, &bajajAllianz)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Error in unmarshling user data" + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	bajajAllianz.Collictible = bajajAllianz.Collictible + invoiceToSave.BajajAllianzPaybleAmount
	bajajAllianz.Collected   = bajajAllianz.Collected   + invoiceToSave.BajajAllianzPaidAmount

    bajajAllianzJSON, marshalErr := json.Marshal(bajajAllianz)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Updating bajaj allianz details")
	err = stub.PutState("BajajAllianz", bajajAllianzJSON)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to update bajaj allianz " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}	

	getAdminBajaj, err := stub.GetState("BajajAllianz")
	if err != nil  {
		errKey = UserID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid InsuranceID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	err := json.Unmarshal([]byte(getAdminBajaj, &bajajAllianz)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Error in unmarshling user data" + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}

	allianz.Collictible = allianz.Collictible + invoiceToSave.AllianzPaybleAmount
	allianz.Collected   = allianz.Collected   + invoiceToSave.AllianzPaidAmount

    bajajAllianzJSON, marshalErr := json.Marshal(allianz)
	if marshalErr != nil {
		repError = strings.Replace(marshalErr.Error(), "\"", " ", -1)
		errorDetails = "Cannot Marshal the JSON- " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	_insuranceLogger.Info("Updating allianz details")
	err = stub.PutState("Allianz", allianzJSON)
	if err != nil {
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Unable to update allianz " + repError
		jsonResp = "{\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("buyInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}	
	resultData := map[string]interface{}{
		"txID"   :      stub.GetTxID(),
		"message":     "Success gnerated invoice",
		"status" :      "true"
	
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)	
}
}
	
// get invoice by invoice id (input invoice id)

func (s *Insurance) getInvoiceByInvoiceID(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var err error
	var jsonResp string

	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("getInvoiceByInvoiceID: " + jsonResp)
		return shim.Error(jsonResp)
	}

	InvoiceID := args[0]


	getInvoice, err := stub.GetState(InvoiceID)
	if err != nil  {
		errKey = InvoiceID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid InvoiceID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("getInvoiceByInvoiceID: " + jsonResp)
		return shim.Error(jsonResp)
	}



	
		resultData := map[string]interface{}{
			"txID"   :      stub.GetTxID(),
			"message":     "Success",
			"status" :      "true",
			"data"   :       getInvoice
		}
		respJSON, _ := json.Marshal(resultData)
		return shim.Success(respJSON)

	}


// get invoice by invoice id (input invoice id)

func (s *Insurance) getInvoiceByInvoiceID(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var err error
	var jsonResp string

	if len(args) != 1 {
		errKey = strconv.Itoa(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("getInvoiceByInvoiceID: " + jsonResp)
		return shim.Error(jsonResp)
	}

	InvoiceID := args[0]


	getInvoice, err := stub.GetState(InvoiceID)
	if err != nil  {
		errKey = InvoiceID
		repError = strings.Replace(err.Error(), "\"", " ", -1)
		errorDetails = "Invalid InvoiceID"
		jsonResp = "{\"Data\":" + errKey + ",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_insuranceLogger.Errorf("getInvoiceByInvoiceID: " + jsonResp)
		return shim.Error(jsonResp)
	}



		resultData := map[string]interface{}{
			"txID"   :      stub.GetTxID(),
			"message":     "Success",
			"status" :      "true",
			"data"   :       getInvoice
		}
		respJSON, _ := json.Marshal(resultData)
		return shim.Success(respJSON)

	}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {


	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
	

// get invoice by user id (input user id)
func (t *SimpleChaincode) getInvoiceByUserID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}


// get invoice by insurance id (input insurance id)
func (t *SimpleChaincode) getInvoiceByInsuranceID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
