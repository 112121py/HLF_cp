# Backend
```
backend/
├── server.ts          # 啟動 Express 伺服器的進入點
├── routes/
│   └── auth.ts        # 處理登入的 API 路由（連接 Fabric CA 做驗證）
│   └── federated.ts 
├── wallet/            # (可選) 儲存登入成功者的身份憑證檔案（目前選擇不寫入，部儲存 credential 在 local）
├── package.json       # Node.js 套件與指令定義
└── tsconfig.json      # TypeScript 設定
```

## 準備後端
### 生成 `package.json`
初始檔案，產生 `package.json`，如果已經有這個檔案的話就可以略過此步驟
```
npm init -y
```
### 下載需要的東西
#### 如果是執行 `npm init -y` 才有 `package.json`
安裝執行時會用到的套件（Runtime Dependency）
```
npm install express cors fabric-network fabric-ca-client dockerode multer
```
安裝型別定義（TypeScript 用），這是 TypeScript 所需要的「型別資訊檔 (.d.ts)」
```
npm install -D typescript ts-node @types/node @types/express @types/cors @types/dockerode @types/multer
```
建立 `tsconfig.json`：

```bash
npx tsc --init
```
#### 如果原先已經有 `package.json`
```
npm install
```


### 完整總表對照

| 功能來源 | 所需套件 |
|----------|----------|
| `express`, `express.Router` | `express`, `@types/express` |
| `cors` | `cors` |
| `ts-node`, `typescript` | `ts-node`, `typescript` |
| `fs`, `path` | 內建（Node.js），無需安裝 |
| `fabric-network`（Chaincode交互） | `fabric-network` |
| `fabric-ca-client`（CA 登入） | `fabric-ca-client` |
| `Wallets`, `Gateway` | 從 `fabric-network` 匯入 |
| 型別支援與執行工具 | `@types/node`, `ts-node`, `typescript` |



## 啟動後端伺服器
成功啟動會看到 `Backend running on http://localhost:3001`
```bash
npx ts-node server.ts
```

## 登入 API 測試
請確保已經啟動後端 server，理論上你會看到 `{"success":true,"message":"登入成功"}` 這樣的結果。
```
curl -X POST http://localhost:3001/api/login \
   -H "Content-Type: application/json" \
   -d '{"username":"user1", "password":"user1pw"}'
```
> 登入失敗的話，請先確定容器都有正確 running

## 登入流程邏輯（auth.ts）
1. 接收前端登入請求
    ```
    const { username, password } = req.body;
    ```
2. 連接對應的 Fabric CA server
    - Hyperledger Fabric 官方提供的 CA (Certificate Authority) 用戶端 SDK，用來與 Fabric CA server 通訊。
    ```
    Hyperledger Fabric 官方提供的 CA (Certificate Authority) 用戶端 SDK，用來與 Fabric CA server 通訊。
    ```
3. 嘗試 enroll 驗證帳密
    ```
    await ca.enroll({
        enrollmentID: username,
        enrollmentSecret: password,
    });
    ```
4. 回傳登入結果給前端
    ```
    return res.status(200).json({ success: true, message: '登入成功' });
    ```