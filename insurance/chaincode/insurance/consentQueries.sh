###########################################
#consents
###########################################
peer chaincode invoke -C telcocommon -n consent -c '{"args" :["recordConsent","[{\"urn\":\"1ab32c275ae\",\"msisdn\":\"9876543210\",\"cstid\":\"temp1001\",\"eid\":\"e1001\",\"cli\":\"testheader\",\"sts\":\"1\",\"uorg\":\"AI\",\"cmode\":\"1\",\"pur\":\"1\",\"cts\":\"1511468532\",\"uts\":\"123456789\"},{\"urn\":\"1ab32c275af\",\"msisdn\":\"9876543219\",\"cstid\":\"temp1001\",\"eid\":\"e1001\",\"cli\":\"testheader1\",\"sts\":\"1\",\"uorg\":\"AI\",\"cmode\":\"1\",\"pur\":\"1\",\"cts\":\"1511468532\",\"uts\":\"123456789\"}]"]}'
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
