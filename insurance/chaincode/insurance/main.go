package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var _mainLogger = shim.NewLogger("insurancesContract")

//SmartContract represents the main entart contract
type SmartContract struct {
	insurance *InsuranceManager
}

// Init initializes chaincode.
func (sc *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_mainLogger.Infof("Inside the init method ")
	sc.insurance = new(InsuranceManager)
	sc.insurance.initLedger(stub)
	return shim.Success(nil)
}

//Invoke is the entry point for any transaction
func (sc *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var response pb.Response
	action, _ := stub.GetFunctionAndParameters()
	switch action {
	case "registeruser":
		response = sc.insurance.registerUser(stub)
	case "queryuser":
		response = sc.insurance.queryUserById(stub)
	case "updatebankdetails":
		response = sc.insurance.updateBankDetails(stub)
	case "createoffer":
		response = sc.insurance.createOffer(stub)
	case "queryoffer":
		response = sc.insurance.queryOfferById(stub)
	case "buyhealthinsurance":
		response = sc.insurance.buyHealthInsurance(stub)
	case "raiseinvoice":
		response = sc.insurance.raiseBulkInvoice(stub)
	// case "hcid":
	// 	response = sc.insurance.getHistoryByinsuranceId(stub)
	case "datapage":
		response = sc.insurance.getDataByPagination(stub)
	default:
		response = shim.Error("Invalid action provided")
	}
	return response
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		_mainLogger.Criticalf("Error starting  chaincode: %v", err)
	}
}
