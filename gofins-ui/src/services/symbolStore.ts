/**
 * Global symbol store
 * Single source of truth for all symbol data across the entire UI
 * Components store only ticker lists, actual data is retrieved from here
 */

import type { Symbol } from './api';

type Listener = (tickers: Set<string>) => void;

class SymbolStore {
    private symbols: Map<string, Symbol> = new Map();
    private listeners: Set<Listener> = new Set();

    /**
     * Subscribe to symbol changes
     * Listener receives set of changed tickers
     * Returns unsubscribe function
     */
    subscribe(listener: Listener): () => void {
        this.listeners.add(listener);
        return () => this.listeners.delete(listener);
    }

    /**
     * Get a single symbol by ticker
     */
    get(ticker: string): Symbol | undefined {
        return this.symbols.get(ticker);
    }

    /**
     * Get multiple symbols by tickers
     */
    getMany(tickers: string[]): Symbol[] {
        return tickers.map(t => this.symbols.get(t)).filter((s): s is Symbol => s !== undefined);
    }

    /**
     * Update a single symbol and notify listeners
     */
    update(symbol: Symbol): void {
        this.symbols.set(symbol.ticker, symbol);
        this.notify(new Set([symbol.ticker]));
    }

    /**
     * Bulk update/merge symbols from API response
     * Merges with existing data to preserve fields not in the response
     */
    bulkUpdate(symbols: Symbol[]): void {
        const changedTickers = new Set<string>();
        
        symbols.forEach(symbol => {
            const existing = this.symbols.get(symbol.ticker);
            // Merge with existing to preserve any fields
            const merged = existing ? { ...existing, ...symbol } : symbol;
            this.symbols.set(symbol.ticker, merged);
            changedTickers.add(symbol.ticker);
        });

        if (changedTickers.size > 0) {
            this.notify(changedTickers);
        }
    }

    /**
     * Update specific fields for a symbol (e.g., favorite, rating)
     */
    updateFields(ticker: string, updates: Partial<Symbol>): void {
        const existing = this.symbols.get(ticker);
        if (existing) {
            this.symbols.set(ticker, { ...existing, ...updates });
            this.notify(new Set([ticker]));
        }
    }

    /**
     * Check if symbol exists in store
     */
    has(ticker: string): boolean {
        return this.symbols.has(ticker);
    }

    /**
     * Get all symbols (for debugging)
     */
    getAll(): Map<string, Symbol> {
        return new Map(this.symbols);
    }

    /**
     * Clear all symbols (useful for logout)
     */
    clear(): void {
        this.symbols.clear();
        this.notify(new Set());
    }

    /**
     * Notify all listeners of changes
     */
    private notify(changedTickers: Set<string>): void {
        this.listeners.forEach(listener => listener(changedTickers));
    }
}

// Global singleton instance
export const symbolStore = new SymbolStore();
