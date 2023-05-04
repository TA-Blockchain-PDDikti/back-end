'use strict';

const FabricCAServices = require('fabric-ca-client');
const fabric = require("../utils/fabric.js")

const registerUser = async(userId, userSecret, organizationName) =>  {
    try {
    
        const ccp = await fabric.getCcp(organizationName)

        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities[`ca.${organizationName}.example.com`].url;
        const ca = new FabricCAServices(caURL);

        const wallet = await fabric.getWallet()

        // Check to see if we've already enrolled the user.
        const userIdentity = await wallet.get(userId);
        if (userIdentity) {
            throw `An identity for the user "${userId}" already exists in the wallet`;
        }

        // Check to see if we've already enrolled the admin user.
        const adminIdentity = await wallet.get('admin');
        if (!adminIdentity) {
            throw "Admin network does not exist"
        }

        // build a user object for authenticating with the CA
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, 'admin');

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register({
            affiliation: 'org1.department1',
            enrollmentID: userId,
            role: 'client'
        }, adminUser);
        const enrollment = await ca.enroll({
            enrollmentID: userId,
            enrollmentSecret: userSecret
        });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: `${organizationName.toUpperCase()}MSP`,
            type: 'X.509',
        };
        await wallet.put(userId, x509Identity);
   
        const response = {
            "success":true,
            "message":"Successfully registered user and imported it into the wallet",
        }
        return response

    } catch (error) {
        const response = {
            "success": false,
            "message": `Failed to register user: ${error}`
        };
        return response
    }
   
}

const enrollAdmin = async(adminId, adminSecret, organizationName) => {
    try {
        const ccp = await fabric.getCcp(organizationName)

        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities[`ca.${organizationName}.example.com`];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        const wallet = await fabric.getWallet()

        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get(adminId);
        if (identity) {
            throw 'An identity for the admin user "admin" already exists in the wallet';
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: adminId, enrollmentSecret: adminSecret });
    
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: `${organizationName.toUpperCase()}MSP`,
            type: 'X.509',
        };
        await wallet.put(adminId, x509Identity);
        console.log(x509Identity)

        const response = {
            "success":true,
            "message":"Successfully registered admin and imported it into the wallet",
        }
        return response

    } catch (error) {
        const response = {
            "success": false,
            "message": `Failed to enroll admin user "admin": ${error}`
        };
        return response
    }
}


module.exports = {enrollAdmin, registerUser}

