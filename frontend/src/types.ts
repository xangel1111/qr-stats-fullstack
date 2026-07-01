// Types mirror the Go API response contract.

export interface MatrixStats {
  name?: string;
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

export interface Statistics {
  q: MatrixStats;
  r: MatrixStats;
  combined: CombinedStats;
}

export interface ComputationResponse {
  id: string;
  input: { matrix: number[][]; rows: number; cols: number };
  qr: { q: number[][]; r: number[][] };
  statistics: Statistics;
  createdAt: string;
}

export interface ComputationSummary {
  id: string;
  rows: number;
  cols: number;
  username: string;
  createdAt: string;
}
