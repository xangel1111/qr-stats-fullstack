import { formatNumber } from '../format';
import type { Statistics } from '../types';

interface Props {
  stats: Statistics;
}

/** Renders per-matrix (Q, R) and combined statistics as a table. */
export function StatsPanel({ stats }: Props) {
  const perMatrix = [
    { label: 'Q', s: stats.q },
    { label: 'R', s: stats.r },
  ];

  return (
    <table className="stats">
      <thead>
        <tr>
          <th></th>
          <th>Máx</th>
          <th>Mín</th>
          <th>Promedio</th>
          <th>Suma</th>
          <th>Diagonal</th>
        </tr>
      </thead>
      <tbody>
        {perMatrix.map(({ label, s }) => (
          <tr key={label}>
            <th>{label}</th>
            <td>{formatNumber(s.max)}</td>
            <td>{formatNumber(s.min)}</td>
            <td>{formatNumber(s.average)}</td>
            <td>{formatNumber(s.sum)}</td>
            <td>{s.isDiagonal ? 'Sí' : 'No'}</td>
          </tr>
        ))}
        <tr className="combined">
          <th>Combinado</th>
          <td>{formatNumber(stats.combined.max)}</td>
          <td>{formatNumber(stats.combined.min)}</td>
          <td>{formatNumber(stats.combined.average)}</td>
          <td>{formatNumber(stats.combined.sum)}</td>
          <td>{stats.combined.anyDiagonal ? 'Sí' : 'No'}</td>
        </tr>
      </tbody>
    </table>
  );
}
