import { Request, Response, NextFunction } from 'express';
import { getContract } from '../utils/connectToNetwork';

const CHANNEL_NAME = 'mychannel';

const handleError = (res: Response, error: any, label: string): void => {
  console.error(`${label} error:`, error);
  res.status(500).json({ success: false, message: error.message || 'Internal Server Error' });
};

const parseResult = (result: Buffer) => {
  try {
    return JSON.parse(result.toString());
  } catch {
    return result.toString(); // fallback
  }
};

// === FLTaskContract ===
export const createTask = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const username = (req as any).user?.username; // <-- 從 token 解碼而來，而不是從 req.body
  const { taskId, description } = req.body;
  
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'FLTaskContract');
    await contract.submitTransaction('CreateTask', taskId, description);
    res.json({ success: true, message: '任務建立成功' });
  } catch (error) {
    handleError(res, error, 'CreateTask');
  }
};

export const updateTask = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const { id: taskId } = req.params;
  const username = (req as any).user?.username;
  const { status } = req.body;
  
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'FLTaskContract');
    await contract.submitTransaction('UpdateTaskStatus', taskId, status);
    res.json({ success: true, message: '任務狀態更新成功' });
  } catch (error) {
    handleError(res, error, 'UpdateTask');
  }
};

export const getTask = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const { id: taskId } = req.params;
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'FLTaskContract');
    const result = await contract.evaluateTransaction('GetTask', taskId);
    res.json({ success: true, data: parseResult(result) });
  } catch (error) {
    handleError(res, error, 'GetTask');
  }
};

export const listTasks = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'FLTaskContract');
    const result = await contract.evaluateTransaction('ListTasks');
    res.json({ success: true, data: parseResult(result) });
  } catch (error) {
    handleError(res, error, 'ListTasks');
  }
};

// === ModelContract ===
export const submitModel = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const { modelId, taskId, ipfsHash } = req.body;
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'ModelContract');
    await contract.submitTransaction('SubmitModel', modelId, taskId, ipfsHash);
    res.json({ success: true, message: '模型提交成功' });
  } catch (error) {
    handleError(res, error, 'SubmitModel');
  }
};

export const verifyModel = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const { id: modelId } = req.params;
  const { result } = req.body;
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'ModelContract');
    await contract.submitTransaction('VerifyModel', modelId, result);
    res.json({ success: true, message: '模型驗證結果提交成功' });
  } catch (error) {
    handleError(res, error, 'VerifyModel');
  }
};

export const getModel = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const { id: modelId } = req.params;
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'ModelContract');
    const result = await contract.evaluateTransaction('GetModel', modelId);
    res.json({ success: true, data: parseResult(result) });
  } catch (error) {
    handleError(res, error, 'GetModel');
  }
};

export const listModels = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'ModelContract');
    const result = await contract.evaluateTransaction('ListModels');
    res.json({ success: true, data: parseResult(result) });
  } catch (error) {
    handleError(res, error, 'ListModels');
  }
};

export const recordContribution = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const { modelId, score } = req.body;
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'ModelContract');
    await contract.submitTransaction('RecordContribution', modelId, score.toString());
    res.json({ success: true, message: '貢獻分數已紀錄' });
  } catch (error) {
    handleError(res, error, 'RecordContribution');
  }
};

// === ChannelStatsContract ===
export const queryChannelStats = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'ChannelStatsContract');
    const result = await contract.evaluateTransaction('GetStats');
    res.json({ success: true, data: parseResult(result) });
  } catch (error) {
    handleError(res, error, 'ChannelStats');
  }
};

// === ZeroTrustEndorseContract ===
export const endorseModel = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  const { modelId } = req.params;
  const { reason } = req.body;
  
  const username = (req as any).user?.username;
  if (!username) {
    res.status(400).json({ success: false, message: '找不到使用者資訊，請確認已登入' });
    return;
  }
  
  try {
    const contract = await getContract(username, CHANNEL_NAME, 'ZeroTrustEndorseContract');
    await contract.submitTransaction('EndorseModel', modelId, reason);
    res.json({ success: true, message: '模型背書成功' });
  } catch (error) {
    handleError(res, error, 'EndorseModel');
  }
};

