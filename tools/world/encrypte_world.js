const fs = require('fs')
const jwt = require('jsonwebtoken')
const crypto = require('crypto')

const issuer = 'kinfkong'

const rawWorldFile = '../../../ikatago-credentials/world/world-raw.json'
const publicKeyFile = '../../../ikatago-credentials/world/public.pem'
const privateKeyFile = '../../../ikatago-credentials/world/private.pem'

const targetWorldFile = '../../../ikatago-credentials/world/world.json'
const targetTokenFile = '../../../ikatago-credentials/world/platform-tokens.json'

const privateKey = fs.readFileSync(privateKeyFile)
const publicKey = fs.readFileSync(publicKeyFile)
const world = JSON.parse(fs.readFileSync(rawWorldFile))

const encrypteAlgorithm = 'aes-256-cbc'

// Encrypts plain text into cipher text
function encrypt(dataEncryptKey, plainText) {
    const iv = crypto.randomBytes(16);
    const padKey = dataEncryptKey.padEnd(32)
    const cipher = crypto.createCipheriv(encrypteAlgorithm, padKey, iv);
    let cipherText;
    try {
      cipherText = cipher.update(plainText, 'utf8', 'hex');
      cipherText += cipher.final('hex');
      cipherText = iv.toString('hex') + cipherText
    } catch (e) {
      cipherText = null;
    }
    return cipherText;
  }



const newWorld = {
    publicKey: publicKey.toString('base64'),
    ...world,
}

newWorld.platforms = []
const tokens = []
for (const platform of world.platforms) {
    const signStuff = { 
        algorithm: 'RS256', 
        expiresIn: Math.round((new Date(platform.token.expiresAt).getTime() - new Date().getTime()) / 1000), 
        issuer: issuer, 
        audience: platform.name
    }
    if (platform.token.username) {
        signStuff.subject = platform.token.username
    }
    const token = jwt.sign({ 
        dataEncryptKeyPrefix: platform.token.dataEncryptKeyPrefix 
    }, privateKey, signStuff)
    if (platform.oss && platform.oss.data) {
        platform.oss.encryptedData = encrypt(platform.token.dataEncryptKeyPrefix + '#new-ikatago', JSON.stringify(platform.oss.data))
        delete platform.oss.data
    }
    tokens.push({
        platform: platform.name,
        token,
        createdAt: new Date()
    })
    delete platform.token
    newWorld.platforms.push(platform)
}
// write tokens
// write world
fs.writeFileSync(targetWorldFile, JSON.stringify(newWorld, null, 2))
fs.writeFileSync(targetTokenFile, JSON.stringify(tokens, null, 2))

