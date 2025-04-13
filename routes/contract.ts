// backend/routes/contract.ts
import express from 'express';
import {
  createTask,
  updateTask,
  getTask,
  listTasks,
  submitModel,
  verifyModel,
  getModel,
  listModels,
  recordContribution,
  queryChannelStats,
  endorseModel,
} from '../controllers/contractController';

// add by py
import verifyToken from '../utils/verifyToken';

const router = express.Router();


// modify by py (each route needs to check verifyToken)
// --- FLTaskContract ---
router.post('/task', verifyToken, createTask);
router.put('/task/:id', verifyToken, updateTask);
router.get('/task/:id', verifyToken, getTask);
router.get('/tasks', verifyToken, listTasks);

// --- ModelContract ---
router.post('/model', verifyToken, submitModel);
router.put('/model/:id/verify', verifyToken, verifyModel);
router.get('/model/:id', verifyToken, getModel);
router.get('/models', verifyToken, listModels);
router.post('/model/contribution', verifyToken, recordContribution);

// --- ChannelStatsContract ---
router.get('/stats', verifyToken, queryChannelStats);

// --- ZeroTrustEndorseContract ---
router.post('/model/:modelId/endorse', verifyToken, endorseModel);

export default router;

