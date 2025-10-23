// API base URL: configurable via VITE_API_URL env var, localStorage override, or defaults
// You can override in browser console: 
//   localStorage.setItem('apiUrl', 'https://your-domain.com/api')
//   localStorage.setItem('apiAuth', btoa('username:password'))  // for HTTP Basic Auth
const getApiBaseUrl = () => {
    const localStorageUrl = localStorage.getItem('apiUrl');
    if (localStorageUrl) return localStorageUrl;
    
    if (import.meta.env.VITE_API_URL) return import.meta.env.VITE_API_URL;
    
    return import.meta.env.DEV ? 'http://localhost:8080/api' : 'api';
};

const API_BASE_URL = getApiBaseUrl();

// Get HTTP Basic Auth credentials from localStorage if set
const getAuthHeader = (): HeadersInit => {
    const auth = localStorage.getItem('apiAuth');
    if (auth) {
        return { 'Authorization': `Basic ${auth}` };
    }
    return {};
};

// Helper to build API URLs consistently (internal use only)
const apiUrl = (endpoint: string): string => {
    // Remove leading slash if present
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint.slice(1) : endpoint;
    return `${API_BASE_URL}/${cleanEndpoint}`;
};

// Helper to build image URLs (charts, histograms, etc.)
export const imageUrl = (endpoint: string): string => {
    return apiUrl(endpoint);
};

// Generic API call with error handling
export async function apiCall<T>(
    endpoint: string,
    options?: RequestInit
): Promise<T> {
    // Merge auth headers with any provided headers
    const authHeaders = getAuthHeader();
    const headers = { ...authHeaders, ...options?.headers };
    
    const response = await fetch(apiUrl(endpoint), {
        ...options,
        headers,
    });

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

// Favorites API helpers
export const favorites = {
    list: () => api.get<string[]>('favorites'),
    toggle: (ticker: string) => api.post<{ isFavorite: boolean }>(`favorites/${ticker}`),
};

// Ratings API helpers
export const ratings = {
    getHistory: (ticker: string) => api.get<UserRating[]>(`ratings/${ticker}/history`),
    add: (ticker: string, rating: number, notes?: string) =>
        api.post(`ratings/${ticker}`, { rating, notes: notes || undefined }),
    delete: (ratingId: number) => api.delete(`ratings/${ratingId}`),
};

