docker exec cli.org1 bash -c "peer chaincode invoke -C testcommon -n training -c '{\"Args\":[\"create\",\"I0036\",\"Sai Lakshman\", \"Innovation Engineer\", \"9948035024\", \"Playing VolleyBall\"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/test/orderers/orderer1.test/msp/tlscacerts/tlsca.test-cert.pem"
sleep 3
docker exec cli.org1 bash -c "peer chaincode invoke -C testcommon -n training -c '{\"Args\":[\"retrieve\",\"I0036\"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/test/orderers/orderer1.test/msp/tlscacerts/tlsca.test-cert.pem"
sleep 3
docker exec cli.org1 bash -c "peer chaincode invoke -C testcommon -n training -c '{\"Args\":[\"update\",\"I0036\",\"Playing Cricket\"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/test/orderers/orderer1.test/msp/tlscacerts/tlsca.test-cert.pem"
sleep 3
docker exec cli.org1 bash -c "peer chaincode invoke -C testcommon -n training -c '{\"Args\":[\"queryByName\",\"Sai Lakshman\"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/test/orderers/orderer1.test/msp/tlscacerts/tlsca.test-cert.pem"
sleep 3
docker exec cli.org1 bash -c "peer chaincode invoke -C testcommon -n training -c '{\"Args\":[\"history\",\"I0036\"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/test/orderers/orderer1.test/msp/tlscacerts/tlsca.test-cert.pem"


docker exec cli.bajaj.com bash -c 'peer channel create -c insurancecommon -f ./channels/insurancecommon.tx -o orderer1.insurance:7050 -t 60s --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem'
docker exec cli.bajaj.com bash -c "peer chaincode invoke -C insurancecommon -n simplyfi -c '{\"Args\":[\"setPO\",\"I0036\",\"Sai Lakshman\", \"Innovation Engineer\", \"9948035024\", \"Playing VolleyBall\"]}' --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/insurance/orderers/orderer1.insurance/msp/tlscacerts/tlsca.insurance-cert.pem"