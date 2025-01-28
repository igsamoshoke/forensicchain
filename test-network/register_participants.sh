#!/bin/bash

echo "Registering participants on the blockchain..."

# Define the participant data
participants=(
  "Org1MSP.participant1 firstresponder"
  "Org2MSP.participant2 secondinvestigator"
  "Org1MSP.participant3 prosecutor"
  "Org1MSP.participant4 defense"
  "Org2MSP.participant5 court"
)

# Common arguments for invoking the chaincode
ORDERER="localhost:7050"
ORDERER_TLS_HOSTNAME="orderer.example.com"
CAFILE="$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
CHANNEL="mychannel"
CHAINCODE_NAME="coc_chain"
PEER1_ADDRESS="localhost:7051"
PEER1_TLS_CERT="$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
PEER2_ADDRESS="localhost:9051"
PEER2_TLS_CERT="$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# Loop through the participants and invoke the chaincode for each
for participant in "${participants[@]}"; do
  IFS=' ' read -r MSP_ID ROLE <<< "$participant"

  echo "Registering participant: $MSP_ID with role: $ROLE"
  
  peer chaincode invoke \
    -o "$ORDERER" \
    --ordererTLSHostnameOverride "$ORDERER_TLS_HOSTNAME" \
    --tls --cafile "$CAFILE" \
    -C "$CHANNEL" -n "$CHAINCODE_NAME" \
    --peerAddresses "$PEER1_ADDRESS" --tlsRootCertFiles "$PEER1_TLS_CERT" \
    --peerAddresses "$PEER2_ADDRESS" --tlsRootCertFiles "$PEER2_TLS_CERT" \
    -c "{\"function\":\"RegisterParticipant\",\"Args\":[\"$MSP_ID\", \"$ROLE\"]}"

  if [ $? -eq 0 ]; then
    echo "Successfully registered $MSP_ID as $ROLE."
  else
    echo "Failed to register $MSP_ID as $ROLE."
    exit 1
  fi

  echo "--------------------------------------"
done

echo "All participants registered successfully!"
