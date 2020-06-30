import * as React from "react";
import styled from "@emotion/styled";
import { formatNumberWithCommas } from "../utils/format";

const BlockNumWrapper = styled.span`
  font-weight: bold;
  font-style: italic;
`;

export const BlockNum: React.FC<{ blockNum: number }> = props => (
  <BlockNumWrapper>{formatNumberWithCommas(props.blockNum)}</BlockNumWrapper>
);
