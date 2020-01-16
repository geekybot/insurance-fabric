###########################################
#consents
###########################################
# user records
peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["registeruser","{\"userid\":\"1ab32c275ae\",\"mobileno\":\"9876543210\",\"name\":\"utpal pal\",\"email\":\"utpal@test.com\",\"dob\":\"05\/09\/1993\",\"gender\":\"Male\",\"nationality\":\"Indian\",\"address\":\"Bangalore\",\"insurances\":[],\"invoices\":[]}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem
peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["queryuser","{\"userid\":\"1ab32c275ae\"}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem
peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["updatebankdetails","{\"userid\":\"1ab32c275ae\",\"bank1\":\"HDFC, Bangalore, Account: 12545886\", \"bank2\":\"Deutche, Munich, Account: 65468651456\", \"dualAcc\": true}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem


# offers
peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["createoffer","{\"offerid\":\"5ab345875ae\",\"name\":\"Health Insurance Plus\",\"type\":\"Health Insurance\",\"termscondition\":\"lorem ipsum thingy\",\"baseprice\":2000}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem
peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["queryoffer","{\"offerid\":\"5ab345875ae\"}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem
peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["datapage","{\"sq\": {\"selector\": {\"obj\": \"Offer\",\"active\": true}},\"ps\": 10,\"bm\": \"\"}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem


# buyhealthinsurance

peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["buyhealthinsurance","{\"insuranceid\":\"ins12345\",\"userid\":\"1ab32c275ae\",\"name\":\"Utpal Pal\",\"gender\":\"Male\",\"email\":\"utpal@test.com\",\"mobileno\":\"9876543210\",\"city\":\"Bangalore\",\"age\":25,\"smoking\":true,\"dob\":\"05\/09\/1993\",\"occupation\":\"IT engineer\",\"height\":\"5ft 2 inch\",\"weight\":\"60kg\",\"premium\":2000,\"healthremarks\":\"Nothing specific\",\"dualAcc\":true}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem

# Raise invoice to existing userid, insurance ids
peer chaincode invoke -C insurancecommon -n insurance6 -c '{"args" :["raiseinvoice","{\"formonth\":\"January\",\"userid\":[\"1ab32c275ae\",\"1easfhsgj1a\"]}"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem


# {\"insuranceid\":\"ins12345\",\"userid\":\"1ab32c275ae\",\"name\":\"Utpal Pal\",\"gender\":\"Male\",\"email\":\"utpal@test.com\",\"mobileno\":\"9876543210\",\"city\":\"Bangalore\",\"age\":25,\"smoking\":true,\"dob\":\"05\/09\/1993\",\"occupation\":\"IT engineer\",\"height\":\"5ft 2 inch\",\"weight\":\"60kg\",\"healthremarks\":\"Nothing specific\",\"dualAcc\":true}
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["getConsent","{\"type\":\"msisdn\",\"msisdn\":\"9876543210\"}"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["getConsent","{\"type\":\"urn\",\"urn\":\"1ab32c275ae\"}"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["getConsent","{\"type\":\"entity\",\"entity\":\"e1001\"}"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["getConsent","{\"type\":\"header\",\"header\":\"testheader\"}"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["updateConsentStatus","1ab32c275ae","2","123456"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["updateConsentStatusByHeaderAndMsisdn","[{\"msisdn\":\"9876543210\",\"cli\":\"testheader\"}]","2","123456"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["updateConsentStatusByIDs","[\"1ab32c275ae\"]","4","123456"]}' 
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["updateConsentExpiryByIDs","[\"1ab32c275ae\"]","155123456","123456"]}' 
#command to work the status of the consent should be active
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["updateConsentExpiryByHeaderAndMsisdn","[{\"msisdn\":\"9876543210\",\"cli\":\"testheader\"},{\"msisdn\":\"9876543210\",\"cli\":\"testheader1\"}]","155123456","123456"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["getActiveConsentsByMSISDN","9876543210","2"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["revokeActiveConsentsByMsisdn","9876543210","123456789","3"]}'
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["getHistory","1ab32c275ae"]}'
