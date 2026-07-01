import { ErrorRequestHandler } from 'express';

import { ApiError } from '../errors';
import { logger } from '../logger';

/**
 * errorHandler renders every error as the shared envelope
 * { error: { code, message, details } }.
 */
export const errorHandler: ErrorRequestHandler = (err, _req, res, _next) => {
  if (err instanceof ApiError) {
    res.status(err.status).json({
      error: { code: err.code, message: err.message, details: err.details },
    });
    return;
  }

  // Malformed JSON body from express.json().
  if (err && typeof err === 'object' && (err as { type?: string }).type === 'entity.parse.failed') {
    res.status(400).json({
      error: { code: 'BAD_REQUEST', message: 'invalid JSON body' },
    });
    return;
  }

  logger.error({ err }, 'unhandled error');
  res.status(500).json({
    error: { code: 'INTERNAL_ERROR', message: 'internal server error' },
  });
};
