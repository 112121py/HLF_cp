// backend/server.ts
import express from 'express';
import cors from 'cors';

import authRouter from './routes/auth';
import federatedRouter from './routes/federated';

// add by py
import contractRoutes from './routes/contract';

const app = express();
app.use(cors());
app.use(express.json());

app.use('/api', authRouter);
app.use('/api', federatedRouter);

// add by py
app.use('/contract', contractRoutes);

app.listen(3001, () => console.log('Backend running on http://localhost:3001'));
