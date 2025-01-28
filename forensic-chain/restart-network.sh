#!/bin/bash

# Navigate to the Fabric samples test-network
cd $HOME/fabric-samples/test-network || exit

# Stop the network and remove old containers and data
./network.sh down
docker system prune -f
docker volume prune -f
echo "Old network and data cleared."

# Start the network
./network.sh up createChannel -ca
echo "Fabric network restarted and channel created."

# Deploy the chaincode
CHAINCODE_NAME="forensic-chain"
CHAINCODE_PATH="../forensic-chain"
CHAINCODE_LANG="golang"
CHAINCODE_VERSION="1.0"
CHAINCODE_LABEL="${CHAINCODE_NAME}_${CHAINCODE_VERSION}"

./network.sh deployCC -ccn $CHAINCODE_NAME -ccp $CHAINCODE_PATH -ccl $CHAINCODE_LANG -ccv $CHAINCODE_VERSION
echo "Chaincode deployed successfully."
