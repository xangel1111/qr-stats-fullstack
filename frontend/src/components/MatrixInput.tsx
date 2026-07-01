interface Props {
  cells: string[][];
  onChange: (cells: string[][]) => void;
}

const MAX_DIM = 8;

function clamp(n: number, lo: number, hi: number): number {
  if (Number.isNaN(n)) return lo;
  return Math.min(hi, Math.max(lo, n));
}

/**
 * Editable matrix grid. Cells are kept as raw strings so intermediate typing
 * (empty, "-", "1.") is preserved; the parent parses to numbers on submit.
 */
export function MatrixInput({ cells, onChange }: Props) {
  const rows = cells.length;
  const cols = cells[0]?.length ?? 0;

  function resize(newRows: number, newCols: number) {
    const next: string[][] = [];
    for (let i = 0; i < newRows; i++) {
      const row: string[] = [];
      for (let j = 0; j < newCols; j++) {
        row.push(cells[i]?.[j] ?? '0');
      }
      next.push(row);
    }
    onChange(next);
  }

  function setCell(i: number, j: number, value: string) {
    const next = cells.map((row) => row.slice());
    next[i][j] = value;
    onChange(next);
  }

  return (
    <div>
      <div className="dims">
        <label>
          Filas
          <input
            type="number"
            min={1}
            max={MAX_DIM}
            value={rows}
            onChange={(e) => resize(clamp(Number(e.target.value), 1, MAX_DIM), cols)}
          />
        </label>
        <label>
          Columnas
          <input
            type="number"
            min={1}
            max={MAX_DIM}
            value={cols}
            onChange={(e) => resize(rows, clamp(Number(e.target.value), 1, MAX_DIM))}
          />
        </label>
      </div>

      <table className="matrix editable">
        <tbody>
          {cells.map((row, i) => (
            <tr key={i}>
              {row.map((value, j) => (
                <td key={j}>
                  <input
                    type="number"
                    value={value}
                    aria-label={`fila ${i + 1} columna ${j + 1}`}
                    onChange={(e) => setCell(i, j, e.target.value)}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
