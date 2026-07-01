/**
 * ApiError carries an HTTP status and a machine-readable code. The central
 * error handler turns it into the shared JSON error envelope, matching the Go
 * API's contract.
 */
export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string,
    public readonly details?: string[],
  ) {
    super(message);
    this.name = 'ApiError';
  }

  static badRequest(message: string, details?: string[]): ApiError {
    return new ApiError(400, 'BAD_REQUEST', message, details);
  }

  static unauthorized(message: string): ApiError {
    return new ApiError(401, 'UNAUTHORIZED', message);
  }

  static validation(message: string, details?: string[]): ApiError {
    return new ApiError(422, 'VALIDATION_ERROR', message, details);
  }

  static internal(message: string): ApiError {
    return new ApiError(500, 'INTERNAL_ERROR', message);
  }
}
