package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	id "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/hyperledger/fabric/protos/peer"
)

var _consentLogger = shim.NewLogger("Consent")

const _CreateEvent = "CREATE_CONSENT"
const _UpdateEvent = "UPDATE_CONSENT"
const _BulkCreateEvent = "BULK_CREATE"

// admin data for two operators
type Admin struct {
	ObjType      string `json:"obj"`         // DocType
	AdminId      string `json:"adminid"`     // Admin Id unique key
	Name         string `json:"name"`        //Name of the party
	Operator     string `json:"operator"`    //Operator name
	Collectibles int    `json:"collectible"` // amount to collect
	Collected    int    `json:"collected"`   // Amount collected after discount
}

//invoice to roll out due payments
type Invoice struct {
	ObjType        string `json:"obj"`            // DocType
	UserId         string `json:"userid"`         // user id to mamtch
	InvoiceId      string `json:"invoiceid"`      // Inovoice id unique key
	ForMonth       string `json:"formonth"`       //For which month invoice is generated
	InsuranceId    string `json:"insuranceid"`    //Invoice raised against which invoice id
	PremiumPayable int    `json:"premiumpayable"` // Payable premium due for this month
	PremiumPaid    int    `json:"premiumpaid"`    // Premium paid after discount applied
	BajajAllianz   int    `json:"bajajallianz"`   // Premium paid to bajaj allianz
	Allianz        int    `json:"allianz"`        // premium paid to allianz
}

//user structure to register user on the platform
type User struct {
	ObjType     string   `json:"obj"`         // DocType
	UserId      string   `json:"userid"`      // user id unique key
	Name        string   `json:"name"`        //Name of the party
	Email       string   `json:"email"`       //email of the user
	MobileNo    string   `json:"mobileno"`    // Mobile no of the user
	DOB         string   `json:"dob"`         // Date of Birth of the user
	Gender      string   `json:"gender"`      // Gender of the user
	Nationality string   `json:"nationality"` //Nationality of the user
	Address     string   `json:"address"`     // address of the user
	Insurances  []string `json:"insurances"`  //Insurances taken by the user
	Invoices    []string `json:"invoices"`    //Invoice generated for the user
	Payable     int      `json:"payable"`     // Amount payable
	Paid        int      `json:"paid"`        // amount paid after disocunt
	Bank1       string   `json:"bank1"`       // bank 1 details
	Bank2       string   `json:"bank2"`       //bank 2 details
	DualAcc     bool     `json:"dualAcc"`     // dual account enabled or not
}

// offer structure to rollout offers
type Offer struct {
	ObjType        string `json:"obj"`            // DocType
	OfferId        string `json:"offerid"`        // offer id unique key
	Name           string `json:"name"`           //Name of the offer
	Type           string `json:"type"`           // type of insurance offer
	TermsCondition string `json:"termscondition"` // Terms and condition for this offer
	BasePrice      int    `json:"baseprice"`      // base price of the offer
	Active         bool   `json:"active"`         // Is offer active true, false
}

// string `json:""`
type HealthInsurance struct {
	ObjType       string `json:"obj"`           // DocType
	InsuranceId   string `json:"insuranceid"`   // insurance id unique key
	UserId        string `json:"userid"`        //User id to save the refernece of the user
	Name          string `json:"name"`          //Name of the offer
	Gender        string `json:"gender"`        // Gender of the user
	Email         string `json:"email"`         //email of the user
	MobileNo      string `json:"mobileno"`      // Mobile no of the user
	City          string `json:"city"`          // City of the user/buyer
	Age           int    `json:"age"`           // Age of the applicant
	Smoking       bool   `json:"smoking"`       // Smoker/ non smoker
	DOB           string `json:"dob"`           // Date of Birth of the user
	Occupation    string `json:"occupation"`    // Occupation of the  buyer
	Height        string `json:"height"`        // height of the user in x"y'
	Weight        string `json:"weight"`        // Weight in kilograms (floor, no decimal point)
	HealthRemarks string `json:"healthremarks"` // Pre existing healt condition if any
	Active        bool   `json:"active"`        // Is offer active true, false
	DualAcc       bool   `json:"dualAcc"`       // dual account enabled or not
}

//Dashboard

type Dashboard struct {
	ObjType           string `json:"obj"`              // DocType
	DashboardId       string `json:"dashboardid"`      // dashboard id unique key
	InsuranceSold     int    `json:"insurancesold"`    // Total Insurance sold
	DualUserCount     int    `json:"dualusercount"`    // Dual user count
	InsuranceForDual  int    `json:"isurancefordual"`  //Insurance sold to dual users
	TotalCollectibles int    `json:"totalcollectibles` //total amount collectibles
	TotalCollected    int    `json:"totalcollected`    //total amount collected
}

//InsuranceManager manages Consent related transactions
type InsuranceManager struct {
}

var errorDetails, errKey, jsonResp string

func (s *InsuranceManager) initLedger(APIstub shim.ChaincodeStubInterface) peer.Response {
	admins := []Admin{
		Admin{ObjType: "Admin", AdminId: "bajaj", Name: "Bajaj", Operator: "Operator 1", Collectibles: 0, Collected: 0},
		Admin{ObjType: "Admin", AdminId: "bajajallianz", Name: "Bajaj Allianz", Operator: "Operator 2", Collectibles: 0, Collected: 0},
	}

	i := 0
	for i < len(admins) {
		fmt.Println("i is ", i)
		adminAsBytes, _ := json.Marshal(admins[i])
		APIstub.PutState(admins[i].AdminId, adminAsBytes)
		// fmt.Println("Added", admins[i])
		i = i + 1
	}
	var dashboard = Dashboard{ObjType: "Dashboard", DashboardId: "dashboard001", InsuranceSold: 0, DualUserCount: 0, InsuranceForDual: 0, TotalCollectibles: 0, TotalCollected: 0}
	dashboardAsBytes, _ := json.Marshal(dashboard)
	APIstub.PutState("dashboard001", dashboardAsBytes)
	return shim.Success(nil)
}

//creating User record in the ledger
func (s *InsuranceManager) registerUser(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("createConsent: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var userToSave User
	err := json.Unmarshal([]byte(args[0]), &userToSave)
	if err != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerUser: " + jsonResp)
		return shim.Error(jsonResp)
	}
	//Second Check if the user id is already existing or not
	if recordBytes, _ := stub.GetState(userToSave.UserId); len(recordBytes) > 0 {
		errKey = string(userToSave.UserId)
		errorDetails = "User with this user id already exists, please enter correct userid"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerUser: " + jsonResp)
		return shim.Error(jsonResp)
	}
	userToSave.ObjType = "User"
	// userToSave.Insurances = []
	// userToSave.Inovoices = []
	userToSave.Payable = 0
	userToSave.Paid = 0
	userToSave.DualAcc = false

	//Save the entry
	userJSON, _ := json.Marshal(userToSave)
	_consentLogger.Info("Saving User to the ledger with id----------", userToSave.UserId)
	err = stub.PutState(userToSave.UserId, userJSON)
	if err != nil {
		errKey = string(userJSON)
		errorDetails = "Unable to save user with UserId - " + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerUser: " + jsonResp)
		return shim.Error(jsonResp)
	}
	retErr := stub.SetEvent(_CreateEvent, userJSON)
	if retErr != nil {
		errKey = string(userJSON)
		errorDetails = "Event not generated for event : CREATE_CONSENT- " + string(retErr.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerUser: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"urn":     userToSave.UserId,
		"message": "User registered Successfully",
		"status":  "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

//updating the bank records for first one
func (s *InsuranceManager) updateBankDetails(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var updatedUser User
	errUser := json.Unmarshal([]byte(args[0]), &updatedUser)
	if errUser != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	userRecord, err := stub.GetState(updatedUser.UserId)
	if err != nil {
		errKey = string(updatedUser.UserId)
		errorDetails = "Could not fetch the details for the User- " + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	} else if userRecord == nil {
		errKey = string(updatedUser.UserId)
		errorDetails = "User does not exist with UserId"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var existingUser User
	err = json.Unmarshal([]byte(userRecord), &existingUser)
	if err != nil {
		errKey = string(userRecord)
		errorDetails = "Invalid JSON for storing" + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	existingUser.Bank1 = updatedUser.Bank1
	existingUser.Bank2 = updatedUser.Bank2
	existingUser.DualAcc = updatedUser.DualAcc
	userJSON, _ := json.Marshal(existingUser)

	err = stub.PutState(existingUser.UserId, userJSON)
	if err != nil {
		errKey = string(userJSON)
		errorDetails = "Unable to save User with UserId -" + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}

	if updatedUser.DualAcc == true {
		dashBoardRecord, err := stub.GetState("dashboard001")
		if err != nil {
			errKey = string("dashboard001")
			errorDetails = "Could not fetch the dashboard details for the User- " + string(err.Error())
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		} else if dashBoardRecord == nil {
			errKey = string("dashboard001")
			errorDetails = "Dashboard does not exist with this ID"
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		}
		var existingDashboard Dashboard
		err = json.Unmarshal([]byte(dashBoardRecord), &existingDashboard)
		if err != nil {
			errKey = string(dashBoardRecord)
			errorDetails = "Invalid JSON for storing" + string(err.Error())
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		}
		existingDashboard.DualUserCount = existingDashboard.DualUserCount + 1
		dashbaordJSON, _ := json.Marshal(existingDashboard)
		err = stub.PutState(existingDashboard.DashboardId, dashbaordJSON)
		if err != nil {
			errKey = string(dashbaordJSON)
			errorDetails = "Unable to save update bankdetails with UserId -" + string(err.Error())
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		}
	}

	retErr := stub.SetEvent(_UpdateEvent, userJSON)
	if retErr != nil {
		errKey = string(userJSON)
		errorDetails = "Event not generated for event : UPDATE_CONSENT-  " + string(retErr.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"urn":     updatedUser.UserId,
		"message": "Bank details updated successfully",
		"status":  "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

//querying the consent record from the ledger given consentId
func (s *InsuranceManager) queryUserById(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryUserById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var qUser User
	err := json.Unmarshal([]byte(args[0]), &qUser)
	if err != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryUserById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	userRecord, retErr := stub.GetState(qUser.UserId)
	if retErr != nil {
		errKey = string(qUser.UserId)
		errorDetails = "Unable to fetch the consent - " + string(retErr.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryUserById: " + jsonResp)
		return shim.Error(jsonResp)
	} else if userRecord == nil {
		errKey = string(qUser.UserId)
		errorDetails = "User does not exist with UserId"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryUserById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var record User
	err1 := json.Unmarshal(userRecord, &record)
	if err1 != nil {
		errKey = string(userRecord)
		errorDetails = "Invalid JSON - " + string(err1.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryUserById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"data":   record,
		"status": "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

//creating offers in the ledger
func (s *InsuranceManager) createOffer(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("createConsent: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var offerToSave Offer
	err := json.Unmarshal([]byte(args[0]), &offerToSave)
	if err != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("createOffer: " + jsonResp)
		return shim.Error(jsonResp)
	}
	//Second Check if the offer id is already existing or not
	if recordBytes, _ := stub.GetState(offerToSave.OfferId); len(recordBytes) > 0 {
		errKey = string(offerToSave.OfferId)
		errorDetails = "Offer with this offer id already exists, please enter correct offerid"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("createOffer: " + jsonResp)
		return shim.Error(jsonResp)
	}
	offerToSave.ObjType = "Offer"
	offerToSave.Active = true
	//Save the entry
	offerJSON, _ := json.Marshal(offerToSave)
	_consentLogger.Info("Saving Offer to the ledger with id----------", offerToSave.OfferId)
	err = stub.PutState(offerToSave.OfferId, offerJSON)
	if err != nil {
		errKey = string(offerJSON)
		errorDetails = "Unable to save offer with OfferId - " + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("createOffer: " + jsonResp)
		return shim.Error(jsonResp)
	}
	retErr := stub.SetEvent(_CreateEvent, offerJSON)
	if retErr != nil {
		errKey = string(offerJSON)
		errorDetails = "Event not generated for event : CREATE_CONSENT- " + string(retErr.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("createOffer: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"urn":     offerToSave.OfferId,
		"message": "Offer created Successfully",
		"status":  "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

//querying the consent record from the ledger given consentId
func (s *InsuranceManager) queryOfferById(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryOfferById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var qOffer Offer
	err := json.Unmarshal([]byte(args[0]), &qOffer)
	if err != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryOfferById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	offerRecord, retErr := stub.GetState(qOffer.OfferId)
	if retErr != nil {
		errKey = string(qOffer.OfferId)
		errorDetails = "Unable to fetch the consent - " + string(retErr.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryOfferById: " + jsonResp)
		return shim.Error(jsonResp)
	} else if offerRecord == nil {
		errKey = string(qOffer.OfferId)
		errorDetails = "Offer does not exist with OfferId"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryOfferById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var record Offer
	err1 := json.Unmarshal(offerRecord, &record)
	if err1 != nil {
		errKey = string(offerRecord)
		errorDetails = "Invalid JSON - " + string(err1.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryOfferById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"data":   record,
		"status": "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

//Buy healthinsurance
func (s *InsuranceManager) buyHealthInsurance(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("createConsent: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var healthToSave HealthInsurance
	err := json.Unmarshal([]byte(args[0]), &healthToSave)
	if err != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerHealthInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	//Second Check if the health id is already existing or not
	if recordBytes, _ := stub.GetState(healthToSave.InsuranceId); len(recordBytes) > 0 {
		errKey = string(healthToSave.InsuranceId)
		errorDetails = "HealthInsurance with this health id already exists, please enter correct healthid"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerHealthInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	healthToSave.ObjType = "Insurance"
	healthToSave.Active = true
	//Save the entry
	healthJSON, _ := json.Marshal(healthToSave)
	_consentLogger.Info("Saving HealthInsurance to the ledger with id----------", healthToSave.InsuranceId)
	err = stub.PutState(healthToSave.InsuranceId, healthJSON)
	if err != nil {
		errKey = string(healthJSON)
		errorDetails = "Unable to save health with HealthInsuranceId - " + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerHealthInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	userRecord, err := stub.GetState(healthToSave.UserId)
	if err != nil {
		errKey = string(healthToSave.UserId)
		errorDetails = "Could not fetch the user details for the User- " + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	} else if userRecord == nil {
		errKey = string(healthToSave.UserId)
		errorDetails = "User does not exist with this ID"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var existingUser User
	err = json.Unmarshal([]byte(userRecord), &existingUser)
	if err != nil {
		errKey = string(userRecord)
		errorDetails = "Invalid JSON for storing" + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	existingUser.Insurances = append(existingUser.Insurances, healthToSave.InsuranceId)
	userJSON, _ := json.Marshal(existingUser)
	err = stub.PutState(existingUser.UserId, userJSON)
	if err != nil {
		errKey = string(userJSON)
		errorDetails = "Unable to save update bankdetails with UserId -" + string(err.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("Update bank details: " + jsonResp)
		return shim.Error(jsonResp)
	}
	if healthToSave.DualAcc == true {
		dashBoardRecord, err := stub.GetState("dashboard001")
		if err != nil {
			errKey = string("dashboard001")
			errorDetails = "Could not fetch the dashboard details for the User- " + string(err.Error())
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		} else if dashBoardRecord == nil {
			errKey = string("dashboard001")
			errorDetails = "Dashboard does not exist with this ID"
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		}
		var existingDashboard Dashboard
		err = json.Unmarshal([]byte(dashBoardRecord), &existingDashboard)
		if err != nil {
			errKey = string(dashBoardRecord)
			errorDetails = "Invalid JSON for storing" + string(err.Error())
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		}
		existingDashboard.InsuranceForDual = existingDashboard.InsuranceForDual + 1
		dashbaordJSON, _ := json.Marshal(existingDashboard)
		err = stub.PutState(existingDashboard.DashboardId, dashbaordJSON)
		if err != nil {
			errKey = string(dashbaordJSON)
			errorDetails = "Unable to save update bankdetails with UserId -" + string(err.Error())
			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
			_consentLogger.Errorf("Update bank details: " + jsonResp)
			return shim.Error(jsonResp)
		}
	}
	retErr := stub.SetEvent(_CreateEvent, healthJSON)
	if retErr != nil {
		errKey = string(healthJSON)
		errorDetails = "Event not generated for event : CREATE_CONSENT- " + string(retErr.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("registerHealthInsurance: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"urn":     healthToSave.InsuranceId,
		"message": "HealthInsurance registered Successfully",
		"status":  "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

//querying the healtinsurances record from the ledger given insuranceid
func (s *InsuranceManager) queryHealthInsuranceById(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryInsuranceById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var qInsurance HealthInsurance
	err := json.Unmarshal([]byte(args[0]), &qInsurance)
	if err != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryInsuranceById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	insuranceRecord, retErr := stub.GetState(qInsurance.InsuranceId)
	if retErr != nil {
		errKey = string(qInsurance.InsuranceId)
		errorDetails = "Unable to fetch the consent - " + string(retErr.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryInsuranceById: " + jsonResp)
		return shim.Error(jsonResp)
	} else if insuranceRecord == nil {
		errKey = string(qInsurance.InsuranceId)
		errorDetails = "Insurance does not exist with InsuranceId"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryInsuranceById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var record HealthInsurance
	err1 := json.Unmarshal(insuranceRecord, &record)
	if err1 != nil {
		errKey = string(insuranceRecord)
		errorDetails = "Invalid JSON - " + string(err1.Error())
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("queryInsuranceById: " + jsonResp)
		return shim.Error(jsonResp)
	}
	resultData := map[string]interface{}{
		"data":   record,
		"status": "true",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

//InitiateBulkConsents creates a consent record in the ledger
// func (s *InsuranceManager) createBulkConsent(stub shim.ChaincodeStubInterface) peer.Response {
// 	_, args := stub.GetFunctionAndParameters()
// 	if len(args) < 1 {
// 		errKey = string(len(args))
// 		errorDetails = "Invalid Number of Arguments"
// 		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
// 		_consentLogger.Errorf("createBulkConsent: " + jsonResp)
// 		return shim.Error(jsonResp)
// 	}
// 	var listConsent []Consent
// 	err := json.Unmarshal([]byte(args[0]), &listConsent)
// 	if err != nil {
// 		errKey = args[0]
// 		errorDetails = "Invalid JSON provided"
// 		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
// 		_consentLogger.Errorf("createBulkConsent: " + jsonResp)
// 		return shim.Error(jsonResp)
// 	}
// 	var rejectedConsents []string
// 	_, creator := s.getInvokerIdentity(stub)
// 	for i := 0; i < len(listConsent); i++ {
// 		var consentToSave Consent
// 		consentToSave = listConsent[i]
// 		if recordBytes, _ := stub.GetState(consentToSave.ConsentId); len(recordBytes) > 0 {
// 			rejectedConsents = append(rejectedConsents, consentToSave.ConsentId)
// 			continue
// 		}
// 		consentToSave.ObjType = "Consent"
// 		consentToSave.Creator = creator
// 		consentToSave.UpdatedBy = creator
// 		consentToSave.UpdateTs = consentToSave.CreateTS
// 		consentJSON, _ := json.Marshal(consentToSave)
// 		if isValid, err := IsValidConsentPresent(consentToSave); !isValid {
// 			errKey = string(consentJSON)
// 			errorDetails = string(err)
// 			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
// 			_consentLogger.Errorf("createBulkConsent: " + jsonResp)
// 			rejectedConsents = append(rejectedConsents, consentToSave.ConsentId)
// 			continue
// 		}
// 		_consentLogger.Info("Saving Consent to the ledger with id----------", consentToSave.ConsentId)
// 		err = stub.PutState(consentToSave.ConsentId, consentJSON)
// 		if err != nil {
// 			errKey = string(consentJSON)
// 			errorDetails = "Unable to save with ConsentId -" + string(err.Error())
// 			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
// 			_consentLogger.Errorf("createBulkConsent: " + jsonResp)
// 			rejectedConsents = append(rejectedConsents, consentToSave.ConsentId)
// 			continue
// 		}
// 		retErr := stub.SetEvent(_BulkCreateEvent, consentJSON)
// 		if retErr != nil {
// 			errKey = string(consentJSON)
// 			errorDetails = "Event not generated for event : BULK_CREATE- " + string(retErr.Error())
// 			jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
// 			_consentLogger.Errorf("createBulkConsent: " + jsonResp)
// 			rejectedConsents = append(rejectedConsents, consentToSave.ConsentId)
// 			continue
// 		}
// 	}
// 	resultData := map[string]interface{}{
// 		"trxnID":     stub.GetTxID(),
// 		"consents_f": rejectedConsents,
// 		"message":    "Consents Registered Successfully",
// 		"status":     "true",
// 	}
// 	respJSON, _ := json.Marshal(resultData)
// 	return shim.Success(respJSON)
// }

// getDataByPagination will query the ledger on the selector input, and display using the pagination
func (s *InsuranceManager) getDataByPagination(stub shim.ChaincodeStubInterface) peer.Response {
	type Query struct {
		SQuery   string `json:"sq"`
		PageSize string `json:"ps"`
		Bookmark string `json:"bm"`
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		errKey = string(len(args))
		errorDetails = "Invalid Number of Arguments"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("getDataByPagination: " + jsonResp)
		return shim.Error(jsonResp)
	}
	var tempQuery Query
	err := json.Unmarshal([]byte(args[0]), &tempQuery)
	if err != nil {
		errKey = args[0]
		errorDetails = "Invalid JSON provided"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("getDataByPagination: " + jsonResp)
		return shim.Error(jsonResp)
	}
	queryString := tempQuery.SQuery
	pageSize, err1 := strconv.ParseInt(tempQuery.PageSize, 10, 32)
	if err1 != nil {
		errKey = string(pageSize)
		errorDetails = "PageSize should be a Number"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("getDataByPagination: " + jsonResp)
		return shim.Error(jsonResp)
	}
	bookMark := tempQuery.Bookmark
	paginationResults, err2 := getQueryResultForQueryStringWithPagination(stub, queryString, int32(pageSize), bookMark)
	if err2 != nil {
		errKey = queryString + "," + string(pageSize) + "," + bookMark
		errorDetails = "Could not fetch the data"
		jsonResp = "{\"Data\":\"" + errKey + "\",\"ErrorDetails\":\"" + errorDetails + "\"}"
		_consentLogger.Errorf("getDataByPagination: " + jsonResp)
		return shim.Error(jsonResp)
	}
	fmt.Println(paginationResults)
	return shim.Success([]byte(paginationResults))
}

//function used for fetching the data from ledger using pagination
func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) (string, error) {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return "", err
	}
	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())
	return bufferWithPaginationInfo.String(), nil
}

//adding pagination metadata to results
func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *peer.QueryResponseMetadata) *bytes.Buffer {
	buffer.WriteString(",\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}")
	return buffer
}

//constructing result to the pagination query
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{\"Records\":[")
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
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return &buffer, nil
}

//Returns the complete identity in the format
//Certitificate issuer orgs's domain name
//Returns string Unkown if not able parse the invoker certificate
func (s *InsuranceManager) getInvokerIdentity(stub shim.ChaincodeStubInterface) (bool, string) {
	//Following id comes in the format X509::<Subject>::<Issuer>>
	enCert, err := id.GetX509Certificate(stub)
	if err != nil {
		return false, "Unknown."
	}
	issuersOrgs := enCert.Issuer.Organization
	if len(issuersOrgs) == 0 {
		return false, "Unknown.."
	}
	return true, fmt.Sprintf("%s", issuersOrgs[0])
}
