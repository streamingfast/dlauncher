// max drift value we display on the gaph
export const MAX_GRAPH_DRIFT = 15 * 60 * 60;
export const INFINITE_DRIFT_THRESHOLD = 9000000000;

// retain half an hour of data
export const DRIFT_RETENTION = 60;

export const MINDREADER_APP_ID = 'mindreader';

export type MetricConfig = {
  headBlockNumber: boolean;
  headBlockDrift: boolean;
};
