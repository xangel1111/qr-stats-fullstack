import { formatNumber } from '../format';

interface Props {
  data: number[][];
}

/** Renders a numeric matrix as a read-only table. */
export function MatrixTable({ data }: Props) {
  return (
    <table className="matrix">
      <tbody>
        {data.map((row, i) => (
          <tr key={i}>
            {row.map((value, j) => (
              <td key={j}>{formatNumber(value)}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
