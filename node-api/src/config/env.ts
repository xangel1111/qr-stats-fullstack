import dotenv from 'dotenv';

dotenv.config();

function required(key: string): string {
  const value = process.env[key];
  if (!value) {
    throw new Error(`${key} is required`);
  }
  return value;
}

/**
 * Runtime configuration loaded from the environment. JWT_SECRET is required and
 * must match the Go API's secret (shared HS256 secret for service-to-service).
 */
export const config = {
  port: parseInt(process.env.PORT ?? '3000', 10),
  jwtSecret: required('JWT_SECRET'),
  logLevel: process.env.LOG_LEVEL ?? 'info',
  nodeEnv: process.env.NODE_ENV ?? 'development',
};
