'use strict';
var log4js = require('log4js');
var logger = log4js.getLogger('SampleWebApp');
var express = require('express');


require('./config.js');
var hfc = require('fabric-client');


var invoke = require('./app/invoke-transaction.js');
var query = require('./app/query.js');


   var router = express.Router();


const utils= require("../utils/validator")

// configuration pramaters for all the apis coming from globalvars.
var global = require('../globalVars');
let peer = global.peer;
let peer0 = global.peer0;

let channelName = global.channelName;
let chaincodeName = global.hulChaincode;

//Register Vendor
router.post('/registervendor', async (req, res) => {
    if (Object.keys(req.body).length < 3) {
        res.status(400).send({
            status:  false,
            result: "Invalid no of parameters in the body"
        });
        return;
    }
    var name= req.body.name;
    var location= req.body.location;
    var contactNo= req.body.contactNo;
    var randomNumber= utils.generateRandomNumber(1000,9999);
    var vendorid= "v"+randomNumber;
  
    if (!name) {
        res.status(400).send({ status: false, result: "Invalid input for VendorName" });
        return;
    }
    if (!location) {
        res.status(400).send({ status: false, result: "Invalid input for Location" });
        return;
    }
    if (!contactNo) {
        res.status(400).send({ status: false, result: "Invalid input for Contact Number" });
        return;
    }
    var fcn = "vendorRegistration";
   
   
    var args = [name, location, contactNo, vendorid];
  
    try {
        let message = await invoke.invokeChaincode(peer, channelName, chaincodeName, fcn, args, req.username, req.orgname);
        let msg = await query.getTransactionByID(peer0, channelName, message.TxtID, req.username, req.orgname);
      
        
        message["block_num"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['data']['actions'][0]['payload']['action']['proposal_response_payload']['extension']['results']['ns_rwset'][0]['rwset']['reads'][0]['version']['block_num']));
        message["tx_time"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['header']['channel_header']['timestamp']));
        res.status(200).send({
            status: true,
			message:message
        });
        return;
    }
    catch (err) {
        let erm = JSON.parse(err.toString().substring(7));

        res.status(503).send({ status: false, result: erm });
        return;
    }
});



router.post('/vendorlogin', async (req, res) => {
    if (Object.keys(req.body).length < 2) {
        res.status(400).send({
            status:  false,
            result: "Invalid no of parameters in the body"
        });
        return;
    }
    var vendorid= req.body.vendorid;
    var password= req.body.password;
   

    if (!vendorid) {
        res.status(400).send({ status: false, result: "Invalid input for VendorID" });
        return;
    }
   
    if (!password) {
        res.status(400).send({ status: false, result: "Invalid input for PASSWORD" });
        return;
    }
    var fcn = "vendorLogIn";
   
   
    var args = [vendorid, password];
  
    try {
        let message = await invoke.invokeChaincode(peer, channelName, chaincodeName, fcn, args, req.username, req.orgname);
        let msg = await query.getTransactionByID(peer0, channelName, message.TxtID, req.username, req.orgname);
        message["block_num"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['data']['actions'][0]['payload']['action']['proposal_response_payload']['extension']['results']['ns_rwset'][0]['rwset']['reads'][0]['version']['block_num']));
        message["tx_time"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['header']['channel_header']['timestamp']));
        res.status(200).send({
            status: true,
			message:message
        });
        return;
    }
    catch (err) {
        let erm = JSON.parse(err.toString().substring(7));

        res.status(503).send({ status: false, result: erm });
        return;
    }
});




// Invoke transaction on chaincode on target peer for assign truck
router.post('/assigntruck', async (req, res) => {

	

	var shipment_id = req.body.shipment_id;
	var truck_num = req.body.truck_num;
	
	


	if (Object.keys(req.body).length != 2) {
		res.status(400).send({ status: false, message: "Invalid number of parameters " });
		return;
	}

	if (!shipment_id) {
		res.status(400).send({ status: false, message: "Shipment ID is mandatory" });
		return;
	}
	if (!truck_num) {
		res.status(400).send({ status: false, message: "Truck number is mandatory" });
		return;
	}
	
	

	var fcn = "assignTruck";
	

	var args = [shipment_id, truck_num];
	try {
        let message = await invoke.invokeChaincode(peer, channelName, chaincodeName, fcn, args, req.username, req.orgname);
        let msg = await query.getTransactionByID(peer0, channelName, message.TxtID, req.username, req.orgname);
        message["block_num"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['data']['actions'][0]['payload']['action']['proposal_response_payload']['extension']['results']['ns_rwset'][0]['rwset']['reads'][0]['version']['block_num']));
        message["tx_time"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['header']['channel_header']['timestamp']));
        res.status(200).send({
            status: true,
			message:message
        });
        return;
    }
    catch (err) {
        let erm = JSON.parse(err.toString().substring(7));

        res.status(503).send({ status: false, result: erm });
        return;
    }
});

//history for a key
router.post('/gethistory', async (req, res) => {
    if (Object.keys(req.body).length != 1) {
        res.status(400).send({
            status:  false,
            result: "Invalid no of parameters in the body"
        });
        return;
    }
    if (!req.body.id) {
        res.status(400).send({ status: false, result: "Invalid input for id" });
        return;
    }
    var fcn = "getHistory";
    var args = [JSON.stringify(req.body)];
    try {
        let message = await query.queryChaincode(peer0, channelName, chaincodeName, args, fcn, req.username, req.orgname);
        if (typeof (message) != "object") {
            res.status(503).send({ status: false, result: message });
            return;
        }
        res.status(200).send({ status: true,
			message:message });
        return;
    }
    catch (err) {
        console.log("found error: ", err.toString());
        res.status(503).send({ status: false, result: "Could not get data!" });
        return;
    }
});


// Query on chaincode on target peers for vendor
router.get('/getvendordetails', async function (req, res) {
	logger.debug('==================== QUERY BY CHAINCODE ==================');
	

	var vendorid = req.body.vendorid;
	
	
	


	if (Object.keys(req.body).length != 1) {
		res.status(400).send({ status: false, message: "Invalid number of parameters " });
		return;
	}

	if (!vendorid) {
		res.status(400).send({ status: false, message: "Vendor ID is mandatory" });
		return;
	}
	

	var fcn = "getVendorDetails";

	var args = [vendorid];

	

	try {
		let message = await query.queryChaincode(peer0, channelName, chaincodeName, args, fcn, req.username, req.orgname)
		
		console.log("Transaction Successful: " + message);
       
		res.status(200).send({ status: true, message: message });
		return;
	} catch (err) {
		res.status(503).send({ status: false, message: err.message });
		return;
	}
});

// Query on chaincode on target peers for ALL vendor
router.get('/getallvendordetails', async function (req, res) {
	logger.debug('==================== QUERY BY CHAINCODE ==================');
	



	
	var fcn = "getAllVendors";
	
  
	var args = [];

	try {
		let message = await query.queryChaincode(peer0, channelName, chaincodeName, args, fcn, req.username, req.orgname)
		
		console.log("Transaction Successful: " + message);
       
		res.status(200).send({ status: true, message: message });
		return;
	} catch (err) {
		res.status(503).send({ status: false, message: err.message });
		return;
	}
});

// Invoke transaction on chaincode on target peer for invoice settlement
router.post("/updatevendorstatus", async (req, res) => {
	
	
	var shipment_id = req.body.shipment_id;
	var settlement_date = req.body.settlement_date;
	var vendor_status = req.body.vendor_status;
	
	


	if (Object.keys(req.body).length != 3) {
		res.status(400).send({ status: false, message: "Invalid number of parameters " });
		return;
	}

	if (!shipment_id) {
		res.status(400).send({ status: false, message: "Shipment ID is mandatory" });
		return;
	}
	

	if (!settlement_date) {
		res.status(400).send({ status: false, message: "Settlement Date is mandatory" });
		return;
	}

	if (!vendor_status) {
		res.status(400).send({ status: false, message: "Vendor Status is mandatory" });
		return;
	}

	
	var fcn = "updateVendorStatus";

  
	var args = [shipment_id, settlement_date, vendor_status];

	try {
		let message = await invoke.invokeChaincode(peer, channelName, chaincodeName, fcn, args, req.username, req.orgname);
        let msg = await query.getTransactionByID(peer0, channelName, message.TxtID, req.username, req.orgname);
        message["block_num"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['data']['actions'][0]['payload']['action']['proposal_response_payload']['extension']['results']['ns_rwset'][0]['rwset']['reads'][0]['version']['block_num']));
        message["tx_time"] = JSON.parse(JSON.stringify(msg['transactionEnvelope']['payload']['header']['channel_header']['timestamp']));
		res.status(200).send({ status: true, message: message });
		return;
	} catch (err) {
		res.status(503).send({ status: false, message: err.message });
		return;
	}
});

module.exports = router;
