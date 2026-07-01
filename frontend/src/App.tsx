import { useCallback, useState } from 'react';

import { ApiError, computeQR, listComputations, login } from './api';
import { History } from './components/History';
import { MatrixInput } from './components/MatrixInput';
import { MatrixTable } from './components/MatrixTable';
import { StatsPanel } from './components/StatsPanel';
import type { ComputationResponse, ComputationSummary } from './types';

const EXAMPLE: string[][] = [
  ['1', '2'],
  ['3', '4'],
  ['5', '6'],
];

/** Parses a string grid into numbers, or null if any cell is not a number. */
function parseMatrix(cells: string[][]): number[][] | null {
  const matrix: number[][] = [];
  for (const row of cells) {
    const parsed: number[] = [];
    for (const cell of row) {
      const n = Number(cell);
      if (cell.trim() === '' || Number.isNaN(n)) {
        return null;
      }
      parsed.push(n);
    }
    matrix.push(parsed);
  }
  return matrix;
}

export default function App() {
  const [username, setUsername] = useState('demo');
  const [password, setPassword] = useState('demo123');
  const [token, setToken] = useState<string | null>(null);
  const [authError, setAuthError] = useState<string | null>(null);

  const [cells, setCells] = useState<string[][]>(EXAMPLE);
  const [result, setResult] = useState<ComputationResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const [history, setHistory] = useState<ComputationSummary[]>([]);

  const refreshHistory = useCallback(async (tk: string) => {
    try {
      setHistory(await listComputations(tk));
    } catch {
      // History is non-critical; ignore failures.
    }
  }, []);

  async function connect() {
    setAuthError(null);
    try {
      const tk = await login(username, password);
      setToken(tk);
      void refreshHistory(tk);
    } catch (err) {
      setToken(null);
      setAuthError(err instanceof ApiError ? err.message : 'No se pudo conectar con la API');
    }
  }

  async function compute() {
    setError(null);
    setResult(null);
    if (!token) return;

    const matrix = parseMatrix(cells);
    if (!matrix) {
      setError('Todos los valores de la matriz deben ser números.');
      return;
    }

    setLoading(true);
    try {
      const res = await computeQR(token, matrix);
      setResult(res);
      void refreshHistory(token);
    } catch (err) {
      setError(err instanceof ApiError ? `${err.code}: ${err.message}` : 'Error de red');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="app">
      <header>
        <h1>Factorización QR + Estadísticas</h1>
        <p>Frontend que consume la API de Go (que orquesta la API de Node).</p>
      </header>

      <section className="card">
        <h2>Conexión</h2>
        <form
          className="row"
          onSubmit={(e) => {
            e.preventDefault();
            void connect();
          }}
        >
          <label>
            Usuario
            <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} />
          </label>
          <label>
            Contraseña
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>
          <button type="submit">{token ? 'Reconectar' : 'Conectar'}</button>
          <span className="status">
            <span className={`dot ${token ? 'on' : ''}`} />
            {token ? 'Autenticado' : 'Sin conexión'}
          </span>
        </form>
        {authError && <p className="error">{authError}</p>}
      </section>

      <div className="grid">
        <section className="card">
          <h2>Matriz de entrada</h2>
          <MatrixInput cells={cells} onChange={setCells} />
          <div className="row" style={{ marginTop: 12 }}>
            <button type="button" onClick={() => void compute()} disabled={!token || loading}>
              {loading ? 'Calculando…' : 'Calcular QR'}
            </button>
            <button type="button" className="secondary" onClick={() => setCells(EXAMPLE)}>
              Cargar ejemplo
            </button>
          </div>
          {!token && <p className="muted" style={{ marginTop: 8 }}>Conéctate para calcular.</p>}
          {error && <p className="error">{error}</p>}
        </section>

        <section className="card">
          <h2>Resultado</h2>
          {result ? (
            <>
              <p>
                <span className="badge">id {result.id.slice(0, 8)}…</span>{' '}
                <span className="muted">
                  {result.input.rows}×{result.input.cols}
                </span>
              </p>
              <div className="matrices">
                <div className="block">
                  <h3>Q</h3>
                  <MatrixTable data={result.qr.q} />
                </div>
                <div className="block">
                  <h3>R</h3>
                  <MatrixTable data={result.qr.r} />
                </div>
              </div>
              <h3 className="muted">Estadísticas</h3>
              <StatsPanel stats={result.statistics} />
            </>
          ) : (
            <p className="muted">Ejecuta un cálculo para ver Q, R y las estadísticas.</p>
          )}
        </section>
      </div>

      <section className="card">
        <h2>Historial (auditoría)</h2>
        <History items={history} />
      </section>
    </div>
  );
}
