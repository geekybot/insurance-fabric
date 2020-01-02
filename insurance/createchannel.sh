# cd insurance
# Create channel
docker-compose up -d
sleep 10

docker exec cli.bajaj.com bash -c 'peer channel create -o orderer1.insurance:7050 -c insurancecommon -f ./channels/insurancecommon.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
# #joining the peers to that channel
docker exec cli.bajaj.com bash -c 'peer channel join -o orderer1.insurance:7050 -b insurancecommon.block --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer channel join -o orderer1.insurance:7050 -b insurancecommon.block --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'

#updating the anchor.tx file for each of the organizations
docker exec cli.bajaj.com bash -c 'peer channel update -o orderer1.insurance:7050 -c insurancecommon -f ./channels/bajaj-insurancecommon-anchor.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer channel update -o orderer1.insurance:7050 -c insurancecommon -f ./channels/bajajallianz-insurancecommon-anchor.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'

# #installing the chaincode on all the orgs
# docker exec cli.bajaj.com bash -c 'peer chaincode install -o orderer1.insurance:7050 -p consent -n consent -v 1.0 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
# docker exec cli.bajajallianz.com bash -c 'peer chaincode install -o orderer1.insurance:7050 -p consent -n consent -v 1.0 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'


# #instantiating the chaincode on org1
# docker exec cli.bajaj.com bash -c "peer chaincode instantiate -o orderer1.insurance:7050 -C insurancecommon -n consent -v 1.0 -c '{\"Args\":[]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem"


# echo "Total execution time : $(($(date +%s)-starttime)) secs ..."
