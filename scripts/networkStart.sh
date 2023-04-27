#!/bin/bash

# Import utils
. scripts/utils.sh


function runCaContainer() {
  infoln "Starting docker container for Certificate Authority"
  println ""

  CA_COMPOSE_FILES="-f compose/docker-compose-ca.yaml"
  docker-compose ${CA_COMPOSE_FILES} up -d 2>&1

  println ""
}

function runOrgContainer() {
  infoln "Starting docker container for All Organizations"
  println ""

  COMPOSE_FILES="-f compose/docker-compose-orderer.yaml -f compose/docker-compose-kemdikbud.yaml -f compose/docker-compose-he1.yaml"
  docker-compose ${COMPOSE_FILES} up -d 2>&1

  println ""
}

function createAllCertificatesForKemdikbud() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/kemdikbud.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/kemdikbud.example.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca.kemdikbud.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  println 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-kemdikbud-example-com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-kemdikbud-example-com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-kemdikbud-example-com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-kemdikbud-example-com.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy kemdikbud's CA cert to kemdikbud's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem" "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/msp/tlscacerts/ca.crt"

  # Copy kemdikbud's CA cert to kemdikbud's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem" "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/tlsca/tlsca.kemdikbud.example.com-cert.pem"

  # Copy kemdikbud's CA cert to kemdikbud's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/ca"
  cp "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem" "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/ca/ca.kemdikbud.example.com-cert.pem"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca.kemdikbud.example.com --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca.kemdikbud.example.com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca.kemdikbud.example.com --id.name kemdikbudadmin --id.secret kemdikbudadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca.kemdikbud.example.com -M "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/msp" --csr.hosts peer0.kemdikbud.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca.kemdikbud.example.com -M "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/tls" --enrollment.profile tls --csr.hosts peer0.kemdikbud.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/peers/peer0.kemdikbud.example.com/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca.kemdikbud.example.com -M "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/users/User1@kemdikbud.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/users/User1@kemdikbud.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://kemdikbudadmin:kemdikbudadminpw@localhost:7054 --caname ca.kemdikbud.example.com -M "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/users/Admin@kemdikbud.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/kemdikbud/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/kemdikbud.example.com/users/Admin@kemdikbud.example.com/msp/config.yaml"
}

function createAllCertificatesForHE1() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/he1.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/he1.example.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca.he1.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  println 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-he1-example-com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-he1-example-com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-he1-example-com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-he1-example-com.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/he1.example.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy he1's CA cert to he1's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/he1.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/he1/ca-cert.pem" "${PWD}/organizations/peerOrganizations/he1.example.com/msp/tlscacerts/ca.crt"

  # Copy he1's CA cert to he1's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/he1.example.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/he1/ca-cert.pem" "${PWD}/organizations/peerOrganizations/he1.example.com/tlsca/tlsca.he1.example.com-cert.pem"

  # Copy he1's CA cert to he1's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/he1.example.com/ca"
  cp "${PWD}/organizations/fabric-ca/he1/ca-cert.pem" "${PWD}/organizations/peerOrganizations/he1.example.com/ca/ca.he1.example.com-cert.pem"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca.he1.example.com --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca.he1.example.com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca.he1.example.com --id.name he1admin --id.secret he1adminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca.he1.example.com -M "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/msp" --csr.hosts peer0.he1.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/he1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca.he1.example.com -M "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/tls" --enrollment.profile tls --csr.hosts peer0.he1.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/he1.example.com/peers/peer0.he1.example.com/tls/server.key"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca.he1.example.com -M "${PWD}/organizations/peerOrganizations/he1.example.com/users/User1@he1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/he1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/he1.example.com/users/User1@he1.example.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://he1admin:he1adminpw@localhost:8054 --caname ca.he1.example.com -M "${PWD}/organizations/peerOrganizations/he1.example.com/users/Admin@he1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/he1/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/he1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/he1.example.com/users/Admin@he1.example.com/msp/config.yaml"
}

function createAllCertificatesForOrderer() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/example.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/example.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca.example.com --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  println 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-example-com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-example-com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-example-com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-example-com.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy orderer org's CA cert to orderer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

  # Copy orderer org's CA cert to orderer org's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/ordererOrganizations/example.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca.example.com --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca.example.com --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca.example.com -M "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp" --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca.example.com -M "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls" --enrollment.profile tls --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key"

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca.example.com -M "${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/config.yaml"
}

function createAllCertificates() {
  # Creating certificates for all organizations

  infoln "Creating Certificates for All Organizations"
  println ""

  infoln "Creating Kemdikbud Certificates"
  createAllCertificatesForKemdikbud
  println ""

  infoln "Creating HE1 Certificates"
  createAllCertificatesForHE1
  println ""

  infoln "Creating Orderer Certificates"
  createAllCertificatesForOrderer
  println ""
}

function generateGenesisBlock() {
  # Generating channel genesis block

  infoln "Generating orderer genesis block"
  set -x
  configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID $CHANNEL_NAME
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    errorln "Failed to generate orderer genesis block..."
    println ""
    exit 1
  fi
  BLOCKFILE="./channel-artifacts/genesis.block"

  println ""
}

function generateChannelConfigTx() {
  # Generating channel configuration transaction: channel.tx

  infoln "Generating channel configuration transaction: channel.tx"
  set -x
  configtxgen -profile AcademicChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    errorln "Failed to generate channel configuration transaction..."
    println ""
    exit 1
  fi

  println ""
}

function generateAnchorPeerKemdikbud() {
  # Generating anchor peer update for KemdikbudMSP

  infoln "Generating anchor peer update for KemdikbudMSP"
  set -x
  configtxgen -profile AcademicChannel -outputAnchorPeersUpdate ./channel-artifacts/KemdikbudMSPanchors.tx -channelID $CHANNEL_NAME -asOrg KemdikbudMSP
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    errorln "Failed to generate anchor peer update for KemdikbudMSP..."
    println ""
    exit 1
  fi

  println ""
}

function generateAnchorPeerHE1() {
  # Generating anchor peer update for HE1MSP

  infoln "Generating anchor peer update for HE1MSP"
  set -x
  configtxgen -profile AcademicChannel -outputAnchorPeersUpdate ./channel-artifacts/HE1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg HE1MSP
  res=$?
  set +x
  if [ $res -ne 0 ]; then
      errorln "Failed to generate anchor peer update for HE1MSP..."
      println ""
      exit 1
  fi

  println ""
}

function createChannelArtifact() {
  # Check configtxgen tool

  which configtxgen
	if [ "$?" -ne 0 ]; then
		fatalln "configtxgen tool not found."
	fi

  if [ ! -d "channel-artifacts" ]; then
    mkdir channel-artifacts
  fi

  infoln "Generating channel artifacts"
  println ""

  # Generate block and channel transaction

  generateGenesisBlock
  generateChannelConfigTx

  # Generate anchor peer

  generateAnchorPeerKemdikbud
  generateAnchorPeerHE1

  println ""
}

function startNetwork() {
  # Start the network

  CHANNEL_NAME=$1

  infoln "Start the network"

  println "###########################################################################"
  runCaContainer

  while :
    do
      if [ ! -f "organizations/fabric-ca/kemdikbud/tls-cert.pem" ]; then
        sleep 1
      else
        break
      fi
    done

  println "###########################################################################"
  createAllCertificates

  println "###########################################################################"
  createChannelArtifact

  println "###########################################################################"
  runOrgContainer

  println ""
}
