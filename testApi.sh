#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
cd api
# echo "POST request Enroll on Org1  ..."
# echo
# ORG1_TOKEN=$(curl -s -X POST \
#   http://localhost:4000/users \
#   -H "content-type: application/x-www-form-urlencoded" \
#   -d 'username=Jim&orgName=Org1')
# echo $ORG1_TOKEN
# ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
# echo
# echo "ORG1 token is $ORG1_TOKEN"
# echo
# echo "POST request Enroll on Org2 ..."
# echo
# ORG2_TOKEN=$(curl -s -X POST \
#   http://localhost:4000/users \
#   -H "content-type: application/x-www-form-urlencoded" \
#   -d 'username=Barry&orgName=Org2')
# echo $ORG2_TOKEN
# ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".token" | sed "s/\"//g")
# echo
# echo "ORG2 token is $ORG2_TOKEN"
# echo
# echo

cd .././insurance
# #joining the peers to that channel
docker exec cli.bajaj.com bash -c 'peer channel join -o orderer1.insurance:7050 -b insurancecommon.block --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer channel join -o orderer1.insurance:7050 -b insurancecommon.block --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'

#updating the anchor.tx file for each of the organizations
docker exec cli.bajaj.com bash -c 'peer channel update -o orderer1.insurance:7050 -c insurancecommon -f ./channels/bajaj-insurancecommon-anchor.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer channel update -o orderer1.insurance:7050 -c insurancecommon -f ./channels/bajajallianz-insurancecommon-anchor.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'

#installing the chaincode on all the orgs
docker exec cli.bajaj.com bash -c 'peer chaincode install -o orderer1.insurance:7050 -p chainc -n insurance -v 1.0 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer chaincode install -o orderer1.insurance:7050 -p chainc -n insurance -v 1.0 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'


#instantiating the chaincode on org1
docker exec cli.bajaj.com bash -c "peer chaincode instantiate -o orderer1.insurance:7050 -C insurancecommon -n insurance -v 1.0 -c '{\"Args\":[]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem"


echo "Total execution time : $(($(date +%s)-starttime)) secs ..."