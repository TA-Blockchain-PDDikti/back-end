'use strict';

const { Wallets, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const { KJUR, KEYUTIL } = require('jsrsasign');
const CryptoJS = require('crypto-js');

const organizationName = 'he1'
const ccpPath = path.resolve(__dirname, '..', '..', '..', 'organizations', 'peerOrganizations', `${organizationName}.example.com`, `connection-${organizationName}.json`);
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet', 'he1');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Collect input parameters
        // user: who initiates this query, can be anyone in the wallet
        // filename: the file to be validated
        const user = process.argv[2];
        const filename = process.argv[3];
        console.log(user, filename, wallet)
        // Check to see if we've already enrolled the user.
        const userExists = await wallet.get(user);
        console.log(userExists)
        if (!userExists) {
            console.log('An identity for the user ' + user + ' does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // calculate Hash from the specified file
        //const fileLoaded = fs.readFileSync(filename, 'utf8');
        var hashToAction = CryptoJS.SHA256(filename).toString();
        console.log("Hash of the file: " + hashToAction);

        // extract certificate info from wallet

        //const walletContents = await wallet.export(user);
        const userPrivateKey = userExists.credentials.privateKey;

        console.log("private key",userPrivateKey)
        var sig = new KJUR.crypto.Signature({"alg": "SHA256withECDSA"});
        sig.init(userPrivateKey, "");
        sig.updateHex(hashToAction);
        var sigValueHex = sig.sign();
        var sigValueBase64 = Buffer.from(sigValueHex, 'hex').toString('base64');
        console.log("Signature: " + sigValueBase64);

        // Create a new gateway for connecting to our peer node.
        // const gateway = new Gateway();
        // await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: false } });

        // // Get the network (channel) our contract is deployed to.
        // const network = await gateway.getNetwork('mychannel');

        // // Get the contract from the network.
        // const contract = network.getContract('docrec');

        // // Submit the specified transaction.
        // await contract.submitTransaction('createDocRecord', hashToAction, sigValueBase64);
        // console.log('Transaction has been submitted');

        // // Disconnect from the gateway.
        // await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
