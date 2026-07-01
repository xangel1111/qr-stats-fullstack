// Runs before the test modules are loaded, so config/env (which requires
// JWT_SECRET at import time) can be imported without throwing.
process.env.JWT_SECRET = process.env.JWT_SECRET ?? 'test-secret';
process.env.LOG_LEVEL = 'silent';
