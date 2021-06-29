// 모듈 포함
const express = require('express');
const app = express();
var bodyParser = require('body-parser');

// Constants
const PORT = 3000;
const HOST = "0.0.0.0";

 // 패브릭 연결설정
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

app.use(express.static(path.join(__dirname, 'views')));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// Index page
app.get('/', (req, res) => {
    res.sendFile(__dirname + '/views/index.html');
})

// rest 라우팅
// url mate POST(type1 생성 type2 프로젝트 점수 추가)
app.post('/mate', async(req, res) => {
    const id = req.body.id;
    console.log("your id is: " + id);

    // Wallet 가져오기
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('No User1 in your wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // Gateway에 연결하기
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false }});
    // Channel에 연결하기
    const network = await gateway.getNetwork('mychannel');
    // ChainCode 클래스 가져오기
    const contract = network.getContract('teamate');
    // submit transaction
    const result = await contract.submitTransaction("registerUser", id);

    // 체인코드 수행결과를 클라이언트에 알려주기
    console.log("good");
    res.status(200).send(`{result: "success", msg: "TX has been submitted"}`);
})

app.post('/join', async(req, res) => {
    const id = req.body.id;
    const name = req.body.name;
    console.log("id is " + id + " and project name is " + name);

    // Wallet 가져오기
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // Gateway에 연결하기
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false }});
    // Channel에 연결하기
    const network = await gateway.getNetwork('mychannel');
    // ChainCode 클래스 가져오기
    const contract = network.getContract('teamate');
    // submit transaction
    const result = await contract.submitTransaction("joinProject", id, name);
    // 체인코스 수행결과를 클라이언트에 알려주기
    console.log("well joined");
    res.status(200).send(`{result: "success", msg: "TX has been submitted"}`);
})

// url mate GET(type1 조회 type2 이력조회)
app.get('/mate', async(req, res) => {
    const id = req.query.id;

    // Wallet 가져오기
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // Gateway에 연결하기
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false }});
    // Channel에 연결하기
    const network = await gateway.getNetwork('mychannel');
    // ChainCode 클래스 가져오기
    const contract = network.getContract('teamate');
    // submit transactiond
    const result = await contract.evaluateTransaction("readUserInfoById", id);
    // 체인코스 수행결과를 클라이언트에 알려주기
    console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
    res.status(200).send(result);
})

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
