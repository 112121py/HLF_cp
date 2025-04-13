import express, { Request, Response } from 'express';
import FabricCAServices from 'fabric-ca-client';
import path from 'path';
import fs from 'fs';

// add by py
import jwt from 'jsonwebtoken';
import { Wallets } from 'fabric-network';

const router = express.Router();

router.post('/login', async (req: Request, res: Response): Promise<void> => {
  const { username, password } = req.body;

  // modify by py
  try {
    const ca = new FabricCAServices(
      'https://localhost:7054',
      {
        trustedRoots: fs.readFileSync(path.resolve(__dirname, '../../test-network/organizations/fabric-ca/org1/tls-cert.pem')),
        verify: false,
      },
      'ca-org1'
    );

    const enrollment = await ca.enroll({ enrollmentID: username, enrollmentSecret: password });

    const walletPath = path.join(__dirname, '../../wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    const identity = {
      credentials: {
        certificate: enrollment.certificate,
        privateKey: enrollment.key.toBytes(),
      },
      mspId: 'Org1MSP',
      type: 'X.509',
    };
    await wallet.put(username, identity);


    const token = jwt.sign({ username }, 'my-secret-key', { expiresIn: '1h' });

    res.status(200).json({ success: true, message: '登入成功', token });
  } catch (err) {
    res.status(401).json({ success: false, message: '登入失敗，請確認帳號或密碼' });
  }
  // modify end by py
});
export default router;
