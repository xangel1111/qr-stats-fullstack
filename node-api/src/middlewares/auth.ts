import { RequestHandler } from 'express';
import jwt from 'jsonwebtoken';

import { config } from '../config/env';
import { ApiError } from '../errors';

/**
 * jwtAuth verifies the Bearer token using the shared HS256 secret. The Go API
 * propagates its token here for service-to-service calls.
 */
export const jwtAuth: RequestHandler = (req, _res, next) => {
  const header = req.headers.authorization;
  if (!header || !header.startsWith('Bearer ')) {
    return next(ApiError.unauthorized('missing or malformed Authorization header'));
  }

  const token = header.slice('Bearer '.length);
  try {
    jwt.verify(token, config.jwtSecret, { algorithms: ['HS256'] });
    return next();
  } catch {
    return next(ApiError.unauthorized('invalid or expired token'));
  }
};
