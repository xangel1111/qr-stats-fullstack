import { Router } from 'express';

import { computeStatsHandler } from '../controllers/stats.controller';
import { jwtAuth } from '../middlewares/auth';
import { validateBody } from '../middlewares/validate';
import { statsRequestSchema } from '../schemas/stats.schema';

export const statsRouter = Router();

statsRouter.post('/internal/stats', jwtAuth, validateBody(statsRequestSchema), computeStatsHandler);
