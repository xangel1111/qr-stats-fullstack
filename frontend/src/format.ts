/** formatNumber rounds for display and normalizes -0 to 0. */
export function formatNumber(value: number, precision = 4): string {
  if (Object.is(value, -0)) {
    return '0';
  }
  return String(Number(value.toFixed(precision)));
}
