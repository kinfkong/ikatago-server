const fs = require('fs')
const jwt = require('jsonwebtoken')
const crypto = require('crypto')

const issuer = 'kinfkong'

const rawWorldFile = '../../../ikatago-credentials/world/world-raw.json'
const privateKeyFile = '../../../ikatago-credentials/world/private.pem'

const targetTokenFile = '../../../ikatago-credentials/world/aistudio-8v-tokens.txt'

const privateKey = fs.readFileSync(privateKeyFile)

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

const tokens = []
for (const platform of world.platforms) {
    if (platform.name !== 'aistudio-8v') {
        continue
    }
    const startDate = new Date('2020-10-12T08:30+08:00')
    const endDate = new Date('2020-10-16T00:00+08:00')
    const duration = Math.round(2.5 * 3600 * 1000)
    let currentDate = startDate
    while (currentDate < endDate) {
        const expiresAt = new Date(currentDate.getTime() + duration)
        const signStuff = { 
            algorithm: 'RS256', 
            expiresIn: Math.round((expiresAt.getTime() - new Date().getTime()) / 1000), 
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
            platform.oss.encryptedData = encrypt(platform.token.dataEncryptKeyPrefix + '#ikatago', JSON.stringify(platform.oss.data))
            delete platform.oss.data
        }
        tokens.push({
            token,
            expiresAt: expiresAt
        })
        currentDate = expiresAt
    }
    
}
let content = ''
for (const token of tokens) {
    content += '有效期到: ' + token.expiresAt + '\n\n\n'
    content += `PLATFORM_TOKEN='${token.token}'` + '\n'
    content += '\n\n\n\n\n'
}
// write tokens
fs.writeFileSync(targetTokenFile, content)

