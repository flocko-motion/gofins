import { useEffect, useState } from 'react';
import { symbolStore } from '../services/symbolStore';
import type { Symbol } from '../services/api';

/**
 * Hook to get a single symbol and subscribe to its changes
 */
export function useSymbol(ticker: string): Symbol | undefined {
    const [symbol, setSymbol] = useState(() => symbolStore.get(ticker));

    useEffect(() => {
        // Subscribe to changes
        const unsubscribe = symbolStore.subscribe((changedTickers) => {
            if (changedTickers.has(ticker)) {
                setSymbol(symbolStore.get(ticker));
            }
        });

        // Update with current value in case it changed before subscription
        setSymbol(symbolStore.get(ticker));

        return unsubscribe;
    }, [ticker]);

    return symbol;
}

/**
 * Hook to get multiple symbols and subscribe to their changes
 * Pass in array of tickers, returns array of symbols
 */
export function useSymbols(tickers: string[]): Symbol[] {
    const [symbols, setSymbols] = useState(() => symbolStore.getMany(tickers));

    useEffect(() => {
        const tickerSet = new Set(tickers);
        
        // Subscribe to changes
        const unsubscribe = symbolStore.subscribe((changedTickers) => {
            // Check if any of our tickers changed
            const hasChanges = Array.from(changedTickers).some(t => tickerSet.has(t));
            if (hasChanges) {
                setSymbols(symbolStore.getMany(tickers));
            }
        });

        // Update with current values
        setSymbols(symbolStore.getMany(tickers));

        return unsubscribe;
    }, [tickers.join(',')]); // Re-subscribe if ticker list changes

    return symbols;
}
