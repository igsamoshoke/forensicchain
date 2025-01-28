#!/bin/bash

echo "Querying installed chaincode..."
peer lifecycle chaincode queryinstalled \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem > query_result.txt

PACKAGE_ID=$(grep -oP 'Package ID: \K[^,]+' query_result.txt)
echo "Package ID: $PACKAGE_ID"

# Save the PACKAGE_ID for reuse
echo "$PACKAGE_ID" > package_id.txt
