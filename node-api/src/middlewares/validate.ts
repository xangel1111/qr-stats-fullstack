import { RequestHandler } from 'express';
import { ZodTypeAny } from 'zod';

import { ApiError } from '../errors';

/**
 * validateBody parses and replaces req.body with the schema's typed output.
 * On failure it forwards a 422 with field-level details.
 */
export const validateBody =
  (schema: ZodTypeAny): RequestHandler =>
  (req, _res, next) => {
    const result = schema.safeParse(req.body);
    if (!result.success) {
      const details = result.error.issues.map(
        (issue) => `${issue.path.join('.') || '(root)'}: ${issue.message}`,
      );
      return next(ApiError.validation('invalid request body', details));
    }
    req.body = result.data;
    return next();
  };
