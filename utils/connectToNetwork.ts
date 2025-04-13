// backend/utils/connectToNetwork.ts
import { Wallets, Gateway } from 'fabric-network';
import fs from 'fs';
import path from 'path';

export async function getContract(username: string, channelName: string, chaincodeName: string) {
  const ccpPath = path.resolve(__dirname, '../connection-org1.json');
  const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

  const walletPath = path.resolve(__dirname, '../wallet');
  const wallet = await Wallets.newFileSystemWallet(walletPath);

  const identity = await wallet.get(username);
  if (!identity) {
    throw new Error(`找不到 ${username} 的身份憑證，請先登入`);
  }

  const gateway = new Gateway();
  await gateway.connect(ccp, {
    wallet,
    identity: username,
    discovery: { enabled: true, asLocalhost: true },
  });

  const network = await gateway.getNetwork(channelName);
  const contract = network.getContract(chaincodeName);
  return contract;
}
