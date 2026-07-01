import jwt from 'jsonwebtoken';
import request from 'supertest';

import { createApp } from '../src/app';

const app = createApp();
const token = jwt.sign({ sub: 'demo' }, process.env.JWT_SECRET as string, { algorithm: 'HS256' });

const auth = (): [string, string] => ['Authorization', `Bearer ${token}`];

describe('POST /internal/stats', () => {
  it('returns 401 without a token', async () => {
    const res = await request(app).post('/internal/stats').send({ matrices: [] });
    expect(res.status).toBe(401);
    expect(res.body.error.code).toBe('UNAUTHORIZED');
  });

  it('returns 422 for a jagged matrix', async () => {
    const res = await request(app)
      .post('/internal/stats')
      .set(...auth())
      .send({ matrices: [{ name: 'q', data: [[1, 2], [3]] }] });
    expect(res.status).toBe(422);
    expect(res.body.error.code).toBe('VALIDATION_ERROR');
  });

  it('returns 200 with correct statistics', async () => {
    const res = await request(app)
      .post('/internal/stats')
      .set(...auth())
      .send({ matrices: [{ name: 'r', data: [[1, 0], [0, 1]] }] });

    expect(res.status).toBe(200);
    expect(res.body.perMatrix[0].isDiagonal).toBe(true);
    expect(res.body.combined.sum).toBe(2);
  });

  it('returns 400 for malformed JSON', async () => {
    const res = await request(app)
      .post('/internal/stats')
      .set(...auth())
      .set('Content-Type', 'application/json')
      .send('{ not json');
    expect(res.status).toBe(400);
    expect(res.body.error.code).toBe('BAD_REQUEST');
  });
});

describe('GET /health', () => {
  it('returns ok', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });
});
