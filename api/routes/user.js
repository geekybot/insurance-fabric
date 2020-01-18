'use strict'

const express = require('express')
const router = express.Router()
const utils= require("../../utils/validator")
var invoke = require('../app/invoke-transaction.js');
var query = require('../app/query.js');

// configuration pramaters for all the apis coming from globalvars.
var global = require('../globalVars');
let peers = global.peers;
let peers0 = global.peers0;
let channelName = global.channelName;
let chaincodeName = global.insuranceChaincode;

//register user
router.post('/registeruser', async (req, res) => {
    if (Object.keys(req.body).length == 7) {
        res.status(400).send({
            status:  false,
            result: "Invalid no of parameters in the body"
        });
        return;
    }
    var randomNumber= utils.generateRandomNumber(1000000,9999999);
    var userid= "lab"+randomNumber;
    var mobileno= req.body.mobileno;
    var name= req.body.userid;
    var email= req.body.mobileno;
    var dob= req.body.mobileno;
    var gender= req.body.userid;
    var nationality= req.body.mobileno;
    var address= req.body.userid;
    var insurances= [];
    var invoices= [];
   
     if (!mobileno){
    res.status(400).send({ status: false, result: "Invalid input for mobile number" });
      return;
  }else if (!name){
    res.status(400).send({ status: false, result: "Invalid input for name" });
      return;
  }else if (!email){
    res.status(400).send({ status: false, result: "Invalid input for email" });
      return;
  }else if (!dob){
    res.status(400).send({ status: false, result: "Invalid input for dob" });
      return;
  }else if (!gender){
    res.status(400).send({ status: false, result: "Invalid input for gender" });
      return;
  }else if (!nationality){
    res.status(400).send({ status: false, result: "Invalid input for nationality" });
      return;
  }else if (!address){
    res.status(400).send({ status: false, result: "Invalid input for address" });
      return;
  }

    var fcn = "registeruser";
   
   
    var args = [userid, mobileno, name, email, dob, gender, nationality, address, insurances, invoices];
  
    try {
        let message = await invoke.invokeChaincode(peers, channelName, chaincodeName, fcn, args, req.username, req.orgname);
    
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

// query user by user id
router.get('/queryuser', async function (req, res) {
	logger.debug('==================== QUERY BY CHAINCODE ==================');
	

	var userid = req.params.userid;
	

	
	if (!userid) {
		res.status(400).send({ status: false, message: "Invalid input for user id" });
		return;
	}

	

	
	var fcn = "queryuser";
	

	var args = [userid];

	try {
		let message = await query.queryChaincode(peers0, channelName, chaincodeName, args, fcn, req.username, req.orgname)
		       
		res.status(200).send({ status: true, message: message });
		return;
	} catch (err) {
		res.status(503).send({ status: false, message: err.message });
		return;
	}
});

// update bank details
router.post('/updatebankdetails', async (req, res) => {
  if (Object.keys(req.body).length == 4) {
      res.status(400).send({
          status:  false,
          result: "Invalid no of parameters in the body"
      });
      return;
  }

  var userid= req.body.userid;
  var bank1= req.body.bank1;
  var bank2= req.body.bank2;
  var dualAcc= req.body.dualAcc;
 
 
   if (!userid){
  res.status(400).send({ status: false, result: "Invalid input for user id" });
    return;
}else if (!bank1){
  res.status(400).send({ status: false, result: "Invalid input bank1 details" });
    return;
}else if (!bank2){
  res.status(400).send({ status: false, result: "Invalid input for bank2 details" });
    return;
}else if (!dualAcc){
  res.status(400).send({ status: false, result: "Invalid input for dualAcc" });
    return;
}else{

  var fcn = "updatebankdetails";
 
 
  var args = [userid, bank1, bank2, dualAcc];

  try {
      let message = await invoke.invokeChaincode(peers, channelName, chaincodeName, fcn, args, req.username, req.orgname);
  
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
}
});

module.exports = router
