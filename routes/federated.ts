import express, { Request, Response } from 'express';
import Docker from 'dockerode';
import multer from 'multer';
import path from 'path';

const router = express.Router();
const docker = new Docker({ socketPath: '/var/run/docker.sock' });
const upload = multer({ dest: 'uploads/' });

interface StartContainerRequestBody {
  username: string;
}

interface MulterRequest extends Request {
  file: Express.Multer.File;
}

router.post(
  '/upload-training-data',
  upload.single('file'),
  (req: Request, res: Response): void => {
    (async () => {
      try {
        const file = (req as any).file; // 或自行定義 interface 才有 .file

        if (!file) {
          return res.status(400).json({ success: false, message: '沒有接收到檔案' });
        }

        console.log(`接收到檔案: ${file.originalname}`);
        res.json({ success: true });
      } catch (error) {
        console.error('處理檔案時發生錯誤:', error);
        res.status(500).json({ success: false, message: '處理檔案時發生錯誤' });
      }
    })();
  }
);

router.post('/start-container', (req: Request, res: Response) => {
  (async () => {
    const { username, cid } = req.body;

    if (!cid) {
      return res.status(400).json({ success: false, message: '缺少 CID' });
    }

    const containerName = `trainer-${username}-${Date.now()}`;
    const imageName = 'ztfedbc_trainer_image';

    try {
      const container = await docker.createContainer({
        Image: imageName,
        name: containerName,
        Tty: true,
        Env: [
          `USERNAME=${username}`,
          `CID=${cid}`
        ],
        HostConfig: {
          NetworkMode: 'fabric_test',
          Binds: [
            `${path.resolve('uploads')}:/app/uploads`
          ]
        }
      });

      await container.start();
      console.log(`容器 ${containerName} 已啟動`);
      res.json({ success: true, container: containerName });
    } catch (err: any) {
      console.error('啟動容器失敗:', err);
      res.status(500).json({ success: false, message: '啟動容器失敗', error: err.message });
    }
  })();
});

export default router;
