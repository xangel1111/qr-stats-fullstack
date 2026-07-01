import { RequestHandler } from 'express';

import { StatsRequest } from '../schemas/stats.schema';
import { computeStats } from '../services/stats.service';

/**
 * computeStatsHandler handles POST /internal/stats. The body is already
 * validated by the validate middleware, so it delegates straight to the
 * pure service.
 */
export const computeStatsHandler: RequestHandler = (req, res) => {
  const body = req.body as StatsRequest;
  const result = computeStats(body.matrices);
  res.status(200).json(result);
};
