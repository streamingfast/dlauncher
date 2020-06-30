import React, { useState, useEffect, useRef } from 'react';

export type IntervalFunc = () => void
export function useInterval(callback: IntervalFunc, delay :number) {
  const savedCallback = useRef<IntervalFunc>();

  // Remember the latest callback.
  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  // Set up the interval.
  useEffect(() => {

    function tick() {
      if (savedCallback.current) {
        savedCallback.current();
      }
    }

    if (delay !== null) {
      let id = setInterval(tick, delay);
      return () => clearInterval(id);
    }
  }, [delay]);
}
