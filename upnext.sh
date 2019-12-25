#installing the chaincode on all the orgs
docker exec cli.bajaj.com bash -c 'peer chaincode install -o orderer1.insurance:7050 -p chainc -n simplyfi -v 1.0 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajajallianz.com bash -c 'peer chaincode install -o orderer1.insurance:7050 -p chainc -n simplyfi -v 1.0 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'


#instantiating the chaincode on org1
docker exec cli.bajaj.com bash -c "peer chaincode instantiate -o orderer1.insurance:7050 -C insurancecommon -n simplyfi -v 1.0 -c '{\"Args\":[]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem"


