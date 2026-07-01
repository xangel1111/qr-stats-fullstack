import cors from 'cors';
import express, { Express } from 'express';
import helmet from 'helmet';
import pinoHttp from 'pino-http';

import { logger } from './logger';
import { errorHandler } from './middlewares/error';
import { statsRouter } from './routes/stats.routes';

/**
 * createApp builds the Express application. It is exported separately from the
 * server bootstrap so tests can exercise it with supertest without binding a
 * port.
 */
export function createApp(): Express {
  const app = express();

  app.use(helmet());
  app.use(cors());
  app.use(express.json({ limit: '5mb' }));
  app.use(pinoHttp({ logger }));

  app.get('/health', (_req, res) => {
    res.json({ status: 'ok' });
  });

  app.use(statsRouter);

  // Error handler must be registered last.
  app.use(errorHandler);

  return app;
}
