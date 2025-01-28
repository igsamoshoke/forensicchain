#!/bin/bash

# Set environment variables for the Fabric CA client
export FABRIC_CA_CLIENT_HOME=$HOME/fabric-samples/fabric-ca-client
CA_NAME="ca-org1"
CA_URL="https://localhost:7054"
ADMIN_NAME="admin"
ADMIN_PASS="adminpw"

# Users and their roles
declare -A USERS
USERS=(
    ["firstresponder1"]="first responder"
    ["investigator2"]="second investigator"
    ["prosecutor1"]="prosecutor"
    ["defense1"]="defense"
    ["court1"]="court"
)

# Enroll the default CA admin
if [ ! -d "$FABRIC_CA_CLIENT_HOME/$ADMIN_NAME/msp" ]; then
    echo "Enrolling the CA admin..."
    fabric-ca-client enroll -u https://$ADMIN_NAME:$ADMIN_PASS@$CA_URL --caname $CA_NAME -M $FABRIC_CA_CLIENT_HOME/$ADMIN_NAME/msp
    if [ $? -ne 0 ]; then
        echo "Failed to enroll CA admin"
        exit 1
    fi
else
    echo "CA admin is already enrolled."
fi

# Register and enroll users
for USER in "${!USERS[@]}"; do
    ROLE=${USERS[$USER]}
    PASSWORD="${USER}pw"

    echo "Registering user $USER with role $ROLE..."
    fabric-ca-client register --caname $CA_NAME --id.name $USER --id.secret $PASSWORD --id.type client --id.attrs "role=$ROLE:ecert"
    if [ $? -ne 0 ]; then
        echo "Failed to register user $USER"
        continue
    fi

    echo "Enrolling user $USER..."
    fabric-ca-client enroll -u https://$USER:$PASSWORD@$CA_URL --caname $CA_NAME -M $FABRIC_CA_CLIENT_HOME/$USER/msp
    if [ $? -ne 0 ]; then
        echo "Failed to enroll user $USER"
        continue
    fi

    echo "User $USER registered and enrolled successfully."
done

echo "All users have been processed successfully."
