import { computeStats, isDiagonal } from '../src/services/stats.service';

describe('isDiagonal', () => {
  it('is true for a diagonal square matrix', () => {
    expect(isDiagonal([[5, 0], [0, 3]])).toBe(true);
  });

  it('is false when an off-diagonal entry is non-zero', () => {
    expect(isDiagonal([[5, 1], [0, 3]])).toBe(false);
  });

  it('is false for a non-square matrix', () => {
    expect(isDiagonal([[1, 0], [0, 1], [0, 0]])).toBe(false);
  });

  it('tolerates tiny floating-point noise off the diagonal', () => {
    expect(isDiagonal([[5, 1e-12], [1e-12, 3]])).toBe(true);
  });
});

describe('computeStats', () => {
  it('computes per-matrix and combined statistics', () => {
    const result = computeStats([
      { name: 'q', data: [[1, 2], [3, 4]] },
      { name: 'r', data: [[10, 0], [0, 20]] },
    ]);

    const q = result.perMatrix[0];
    expect(q).toMatchObject({ name: 'q', max: 4, min: 1, sum: 10, average: 2.5, isDiagonal: false });

    const r = result.perMatrix[1];
    expect(r.isDiagonal).toBe(true);

    expect(result.combined).toEqual({
      max: 20,
      min: 0,
      sum: 40,
      average: 5,
      anyDiagonal: true,
    });
  });

  it('handles negative values (e.g. an orthonormal Q factor)', () => {
    const result = computeStats([{ name: 'q', data: [[-0.5, 0.5], [0.5, -0.5]] }]);
    expect(result.combined.min).toBe(-0.5);
    expect(result.combined.max).toBe(0.5);
    expect(result.combined.sum).toBe(0);
  });
});
