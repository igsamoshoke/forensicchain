#!/bin/bash

echo "Checking commit readiness..."

peer lifecycle chaincode checkcommitreadiness \
  --channelID mychannel \
  --name coc_chain \
  --version 1.0 \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --output json > readiness_check.json

echo "Commit readiness status saved to 'readiness_check.json'."
cat readiness_check.json

read -p "Proceed with committing the chaincode? (yes/no): " CONFIRM
if [[ "$CONFIRM" != "yes" ]]; then
  echo "Chaincode commit aborted."
  exit 0
fi

echo "Committing the chaincode..."

peer lifecycle chaincode commit \
  --orderer localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --channelID mychannel \
  --name coc_chain \
  --version 1.0 \
  --sequence 1 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --peerAddresses localhost:7051 --tlsRootCertFiles $PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 --tlsRootCertFiles $PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

if [ $? -eq 0 ]; then
  echo "Chaincode committed successfully!"
else
  echo "Failed to commit chaincode."
  exit 1
fi
