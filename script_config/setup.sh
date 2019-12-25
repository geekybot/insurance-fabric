#generating the cryptos
sudo rm -r channels orderer crypto-config
cryptogen generate --config crypto-config.yaml

mkdir channels orderer
export FABRIC_CFG_PATH=./


# #generating genesis block
configtxgen -profile InsuranceOrdererGenesis -outputBlock ./orderer/genesis.block
#generating channel artifacts
configtxgen -profile InsuranceCommon -outputCreateChannelTx ./channels/insurancecommon.tx -channelID insurancecommon
configtxgen -profile InsuranceCommon -outputAnchorPeersUpdate ./channels/bajaj-insurancecommon-anchor.tx -channelID insurancecommon -asOrg BajajMSP
configtxgen -profile InsuranceCommon -outputAnchorPeersUpdate ./channels/bajajallianz-insurancecommon-anchor.tx -channelID insurancecommon -asOrg BajajAllianzMSP
