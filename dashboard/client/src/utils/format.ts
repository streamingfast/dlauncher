import numeral from "numeral";

export function formatNumberWithCommas(input: number | string): string {
  return numeral(input).format("0,0");
}
