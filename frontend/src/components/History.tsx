import type { ComputationSummary } from '../types';

interface Props {
  items: ComputationSummary[];
}

/** Renders the recent computations (audit trail) from the Go API. */
export function History({ items }: Props) {
  if (items.length === 0) {
    return <p className="muted">Aún no hay cómputos registrados.</p>;
  }

  return (
    <table className="history">
      <thead>
        <tr>
          <th>ID</th>
          <th>Dimensiones</th>
          <th>Usuario</th>
          <th>Fecha</th>
        </tr>
      </thead>
      <tbody>
        {items.map((item) => (
          <tr key={item.id}>
            <td className="mono">{item.id.slice(0, 8)}…</td>
            <td>
              {item.rows}×{item.cols}
            </td>
            <td>{item.username}</td>
            <td>{new Date(item.createdAt).toLocaleString()}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
