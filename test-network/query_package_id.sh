#!/bin/bash

# Set FABRIC_CFG_PATH
export FABRIC_CFG_PATH=/home/fabric/fabric-samples/config

# Org1 Environment Variables (querying chaincode from Org1 peer)
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# Debugging: Print Environment Variables
echo "FABRIC_CFG_PATH=$FABRIC_CFG_PATH"
echo "CORE_PEER_LOCALMSPID=$CORE_PEER_LOCALMSPID"
echo "CORE_PEER_MSPCONFIGPATH=$CORE_PEER_MSPCONFIGPATH"

# Query Installed Chaincode
echo "Querying installed chaincode..."
peer lifecycle chaincode queryinstalled \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem > query_result.txt

# Extract Package ID
PACKAGE_ID=$(grep -oP 'Package ID: \K[^,]+' query_result.txt)
echo "Package ID: $PACKAGE_ID"

# Save the PACKAGE_ID for reuse
echo "$PACKAGE_ID" > package_id.txt
