import { createApp } from './app';
import { config } from './config/env';
import { logger } from './logger';

const app = createApp();

app.listen(config.port, () => {
  logger.info(`stats API listening on port ${config.port}`);
});
