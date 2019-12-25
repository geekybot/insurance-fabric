#creation of channel insurancecommon
# docker exec cli.bajaj.com bash -c 'peer channel create -c insurancecommon -f ./channels/insurancecommon.tx -o orderer1.insurance:7050 -t 60s --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'

# #joining the peers to that channel
docker exec cli.bajaj.com bash -c 'peer channel join -o orderer1.insurance:7050 -b insurancecommon.block --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer channel join -o orderer1.insurance:7050 -b insurancecommon.block --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'

#updating the anchor.tx file for each of the organizations
docker exec cli.bajaj.com bash -c 'peer channel update -o orderer1.insurance:7050 -c insurancecommon -f ./channels/bajaj-insurancecommon-anchor.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer channel update -o orderer1.insurance:7050 -c insurancecommon -f ./channels/bajajallianz-insurancecommon-anchor.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
