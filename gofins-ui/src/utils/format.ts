/**
 * Formatting utilities for prices, numbers, dates, etc.
 */

/**
 * Format price with adaptive decimal places based on magnitude
 * - Large prices (>=100): 2 decimals
 * - Medium prices (>=10): 2 decimals  
 * - Small prices (>=1): 2 decimals
 * - Very small prices (<1): more decimals to show significant digits
 */
export const formatPrice = (price: number | undefined | null): string => {
    if (price === undefined || price === null) return 'N/A';
    
    const absValue = Math.abs(price);
    if (absValue === 0) return '0.00';

    // Use log10 to determine decimal places
    const log = Math.log10(absValue);
    const decimals = Math.min(10, Math.max(0, 3 - Math.ceil(Math.abs(log))));
    
    return price.toFixed(decimals);
};

/**
 * Format market cap in human-readable format (T, B, M, K)
 */
export const formatMarketCap = (marketCap: number | undefined): string => {
    if (!marketCap) return 'N/A';
    
    const absValue = Math.abs(marketCap);
    if (absValue >= 1e12) return `$${(marketCap / 1e12).toFixed(2)}T`;
    if (absValue >= 1e9) return `$${(marketCap / 1e9).toFixed(2)}B`;
    if (absValue >= 1e6) return `$${(marketCap / 1e6).toFixed(2)}M`;
    if (absValue >= 1e3) return `$${(marketCap / 1e3).toFixed(2)}K`;
    
    return `$${marketCap.toLocaleString()}`;
};

/**
 * Format number with fallback
 */
export const formatNumber = (
    value: number | null | undefined, 
    decimals: number = 2, 
    fallback: string = 'N/A'
): string => {
    if (value == null || value === undefined) return fallback;
    return value.toFixed(decimals);
};

/**
 * Format percentage with sign
 */
export const formatPercentage = (value: number | null | undefined, decimals: number = 1): string => {
    if (value === null || value === undefined) return 'N/A';
    const sign = value > 0 ? '+' : '';
    return `${sign}${value.toFixed(decimals)}%`;
};

/**
 * Get color class for rating (-5 to +5)
 */
export const getRatingColor = (rating: number): string => {
    if (rating === 0) return 'text-gray-500';
    
    if (rating > 0) {
        const intensity = Math.min(rating / 5, 1);
        if (intensity > 0.6) return 'text-green-700 font-bold';
        if (intensity > 0.3) return 'text-green-600';
        return 'text-green-500';
    } else {
        const intensity = Math.min(Math.abs(rating) / 5, 1);
        if (intensity > 0.6) return 'text-red-700 font-bold';
        if (intensity > 0.3) return 'text-red-600';
        return 'text-red-500';
    }
};

/**
 * Format date for display
 */
export const formatDate = (dateStr: string, format: 'short' | 'long' = 'short'): string => {
    const date = new Date(dateStr);
    
    if (format === 'short') {
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
    }
    
    return date.toLocaleDateString('en-US', { 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric' 
    });
};
