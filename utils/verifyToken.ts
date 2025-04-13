import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';

const verifyToken = (req: Request, res: Response, next: NextFunction) => {
  const token = req.headers['authorization']?.split(' ')[1];
  if (!token) {
      res.status(403).json({ success: false, message: '未提供 token' });
      return;
  }

  try {
    const decoded = jwt.verify(token, 'my-secret-key');
    (req as any).user = decoded; // 存入 req.user 供後續使用
    next();
  } catch (err) {
    res.status(401).json({ success: false, message: 'token 無效或過期' });
  }
};

export default verifyToken;

