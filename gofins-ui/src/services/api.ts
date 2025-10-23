// Relative to Vite's base config (/gofins/ in production, / in dev)
const API_BASE_URL = 'api';

// Helper to build API URLs consistently (internal use only)
const apiUrl = (endpoint: string): string => {
    // Remove leading slash if present
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint.slice(1) : endpoint;
    return `${API_BASE_URL}/${cleanEndpoint}`;
};

// Generic API call with error handling
export async function apiCall<T>(
    endpoint: string,
    options?: RequestInit
): Promise<T> {
    const response = await fetch(apiUrl(endpoint), options);
    
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`);
    }
    
    return response.json();
}

// Convenience methods
export const api = {
    get: <T>(endpoint: string) => apiCall<T>(endpoint),
    
    post: <T>(endpoint: string, data?: any) => apiCall<T>(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: data ? JSON.stringify(data) : undefined,
    }),
    
    put: <T>(endpoint: string, data?: any) => apiCall<T>(endpoint, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: data ? JSON.stringify(data) : undefined,
    }),
    
    delete: <T>(endpoint: string) => apiCall<T>(endpoint, {
        method: 'DELETE',
    }),
};

// User types
export interface User {
    id: string;
    name: string;
    createdAt: string;
    isAdmin: boolean;
}

export interface UserRating {
    id: number;
    ticker: string;
    rating: number;
    notes?: string;
    createdAt: string;
}

export interface Note {
    id: number;
    ticker: string;
    rating: number;
    notes: string;
    createdAt: string;
}

// Symbol types
export interface Symbol {
    ticker: string;
    exchange?: string;
    name?: string;
    type?: string;
    sector?: string;
    industry?: string;
    country?: string;
    inception?: string;
    oldestPrice?: string;
    isActivelyTrading?: boolean;
    marketCap?: number;
    ath12m?: number;
    currentPriceUsd?: number;
    isFavorite?: boolean;
    latestRating?: number;
    userRating?: number;
}

export interface PriceData {
    Date: string;
    Open: number;
    High: number;
    Low: number;
    Close: number;
    Avg: number;
    YoY: number | null;
    SymbolTicker: string;
}

// Error types
export interface ErrorEntry {
    id: number;
    timestamp: string;
    source: string;
    errorType: string;
    message: string;
    details?: string;
}

// Analysis types
export interface AnalysisPackage {
    ID: string;
    Name: string;
    CreatedAt: string;
    Interval: string;
    TimeFrom: string;
    TimeTo: string;
    HistBins: number;
    HistMin: number;
    HistMax: number;
    McapMin?: number;
    InceptionMax?: string;
    SymbolCount: number;
    Status: string;
}

export interface AnalysisResult {
    symbol: string;
    mean: number;
    stddev: number;
    min: number;
    max: number;
    inception?: string;
}

export interface SymbolProfile {
    ticker: string;
    exchange?: string;
    name?: string;
    type?: string;
    currency?: string;
    sector?: string;
    industry?: string;
    country?: string;
    description?: string;
    website?: string;
    isin?: string;
    inception?: string;
    oldestPrice?: string;
    isActivelyTrading?: boolean;
    marketCap?: number;
    ath12m?: number;
    currentPriceUsd?: number;
    isFavorite?: boolean;
}

export interface CreateAnalysisRequest {
    name: string;
    interval?: string;
    time_from?: string;
    time_to?: string;
    hist_bins?: number;
    hist_min?: number;
    hist_max?: number;
    mcap_min?: string;
    inception_max?: string;
}

// No specialized API objects - just use api.get/post/put/delete directly

