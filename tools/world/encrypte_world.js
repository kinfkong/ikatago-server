const fs = require('fs')
const jwt = require('jsonwebtoken')
const crypto = require('crypto')

const issuer = 'kinfkong'

const rawWorldFile = '../../ikatago-credentials/world/world-raw.json'
const publicKeyFile = '../../ikatago-credentials/world/public.pem'
const privateKeyFile = '../../ikatago-credentials/world/private.pem'

const targetWorldFile = '../../ikatago-credentials/world/world.json'
const targetTokenFile = '../../ikatago-credentials/world/platform-tokens.json'

const privateKey = fs.readFileSync(privateKeyFile)
const publicKey = fs.readFileSync(publicKeyFile)
const world = JSON.parse(fs.readFileSync(rawWorldFile))

const encrypteAlgorithm = 'aes256'

function encrypt(dataEncryptKey, text) {
    const cipher = crypto.createCipher(encrypteAlgorithm, dataEncryptKey)
    const encrypted = cipher.update(text, 'utf8', 'hex') + cipher.final('hex')
    return encrypted
}

function decrypt(dataEncryptKey, encrypted) {
    var decipher = crypto.createDecipher(encrypteAlgorithm, dataEncryptKey);
    var decrypted = decipher.update(encrypted, 'hex', 'utf8') + decipher.final('utf8');
    return decrypted
}

const newWorld = {
    publicKey,
    ...world,
}

newWorld.platforms = []
const tokens = []
for (const platform of world.platforms) {
    const token = jwt.sign({ 
        dataEncryptKeyPrefix: platform.token.dataEncryptKeyPrefix 
    }, 
    privateKey, { 
        algorithm: 'RS256', 
        expiresIn: Math.round((new Date(platform.token.expiresAt).getTime() - new Date().getTime()) / 1000), 
        issuer: issuer, 
        audience: platform.name 
    })
    if (platform.oss && platform.oss.data) {
        platform.oss.data = encrypt(platform.token.dataEncryptKeyPrefix + '#ikatago', JSON.stringify(platform.oss.data))
    }
    tokens.push({
        platform: platform.name,
        token
    })
    delete platform.token
    newWorld.platforms.push(platform)
}
// write tokens
// write world
fs.writeFileSync(targetWorldFile, JSON.stringify(newWorld, null, 2))
fs.writeFileSync(targetTokenFile, JSON.stringify(tokens, null, 2))

