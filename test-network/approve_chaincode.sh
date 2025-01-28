#!/bin/bash

# Set FABRIC_CFG_PATH
export FABRIC_CFG_PATH=/home/fabric/fabric-samples/config

# Org1 Environment Variables for Querying Package ID
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# Query Installed Chaincode
echo "Querying installed chaincode on Org1..."
peer lifecycle chaincode queryinstalled \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem > query_result.txt

# Extract Package ID
PACKAGE_ID=$(grep -oP 'Package ID: \K[^,]+' query_result.txt)
if [ -z "$PACKAGE_ID" ]; then
  echo "Error: Failed to extract Package ID. Please ensure the chaincode is installed."
  exit 1
fi
echo "Package ID: $PACKAGE_ID"

# Approve Chaincode for Org1
echo "Approving chaincode for Org1..."
peer lifecycle chaincode approveformyorg \
  --orderer localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --channelID mychannel \
  --name coc_chain \
  --version 1.0 \
  --package-id "$PACKAGE_ID" \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# Org2 Environment Variables for Approval
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# Approve Chaincode for Org2
echo "Approving chaincode for Org2..."
peer lifecycle chaincode approveformyorg \
  --orderer localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --channelID mychannel \
  --name coc_chain \
  --version 1.0 \
  --package-id "$PACKAGE_ID" \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# Completion Message
echo "Chaincode approved for both Org1 and Org2."
