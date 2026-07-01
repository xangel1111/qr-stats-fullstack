/**
 * Pure statistics logic. No Express, no HTTP: trivially unit-testable.
 *
 * The service is intentionally generic — it computes stats over any set of
 * named matrices and knows nothing about QR. The Go API happens to send "q"
 * and "r", but this service would work for any matrices.
 */

export interface NamedMatrix {
  name: string;
  data: number[][];
}

export interface MatrixStats {
  name: string;
  max: number;
  min: number;
  average: number;
  sum: number;
  isDiagonal: boolean;
}

export interface CombinedStats {
  max: number;
  min: number;
  average: number;
  sum: number;
  anyDiagonal: boolean;
}

export interface StatsResult {
  perMatrix: MatrixStats[];
  combined: CombinedStats;
}

// Off-diagonal entries coming from a floating-point QR factorization are never
// exactly zero, so "is diagonal" is checked within a tolerance.
const DIAGONAL_EPSILON = 1e-9;

/**
 * isDiagonal reports whether a matrix is (numerically) diagonal. A diagonal
 * matrix is defined only for square matrices, so non-square inputs are false.
 */
export function isDiagonal(data: number[][]): boolean {
  const rows = data.length;
  const cols = data[0].length;
  if (rows !== cols) {
    return false;
  }
  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && Math.abs(data[i][j]) > DIAGONAL_EPSILON) {
        return false;
      }
    }
  }
  return true;
}

/** aggregate reduces a flat list of values into max/min/sum/average. */
function aggregate(values: number[]): { max: number; min: number; sum: number; average: number } {
  let max = values[0];
  let min = values[0];
  let sum = 0;
  for (const v of values) {
    if (v > max) max = v;
    if (v < min) min = v;
    sum += v;
  }
  return { max, min, sum, average: sum / values.length };
}

function computeMatrixStats(matrix: NamedMatrix): MatrixStats {
  const { max, min, sum, average } = aggregate(matrix.data.flat());
  return {
    name: matrix.name,
    max,
    min,
    average,
    sum,
    isDiagonal: isDiagonal(matrix.data),
  };
}

/**
 * computeStats returns per-matrix statistics and a combined aggregate across
 * every value of every matrix. Inputs are assumed validated (non-empty,
 * rectangular) by the request schema.
 */
export function computeStats(matrices: NamedMatrix[]): StatsResult {
  const perMatrix = matrices.map(computeMatrixStats);

  const allValues = matrices.flatMap((m) => m.data.flat());
  const { max, min, sum, average } = aggregate(allValues);

  return {
    perMatrix,
    combined: {
      max,
      min,
      average,
      sum,
      anyDiagonal: perMatrix.some((m) => m.isDiagonal),
    },
  };
}
