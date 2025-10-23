import { type SymbolProfile } from '../services/api';

interface SymbolProfileProps {
    profile: SymbolProfile | null;
    loading: boolean;
    error: string | null;
}

const formatMarketCap = (marketCap: number | undefined): string => {
    if (!marketCap) return 'N/A';
    if (marketCap >= 1e12) return `$${(marketCap / 1e12).toFixed(2)}T`;
    if (marketCap >= 1e9) return `$${(marketCap / 1e9).toFixed(2)}B`;
    if (marketCap >= 1e6) return `$${(marketCap / 1e6).toFixed(2)}M`;
    return `$${marketCap.toLocaleString()}`;
};

export default function SymbolProfileView({ profile, loading, error }: SymbolProfileProps) {
    if (loading) {
        return (
            <div className="text-center py-8">
                <p className="text-gray-500">Loading profile...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-center py-8">
                <p className="text-red-600">Error: {error}</p>
            </div>
        );
    }

    if (!profile) {
        return (
            <div className="text-center py-8">
                <p className="text-gray-500">No profile data available</p>
            </div>
        );
    }

    return (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">Profile Information</h2>
            <div className="grid grid-cols-2 gap-4">
                {profile.exchange && (
                    <div>
                        <span className="text-gray-600">Exchange:</span>
                        <span className="ml-2 font-medium">{profile.exchange}</span>
                    </div>
                )}
                {profile.type && (
                    <div>
                        <span className="text-gray-600">Type:</span>
                        <span className="ml-2 font-medium">{profile.type}</span>
                    </div>
                )}
                {profile.sector && (
                    <div>
                        <span className="text-gray-600">Sector:</span>
                        <span className="ml-2 font-medium">{profile.sector}</span>
                    </div>
                )}
                {profile.industry && (
                    <div>
                        <span className="text-gray-600">Industry:</span>
                        <span className="ml-2 font-medium">{profile.industry}</span>
                    </div>
                )}
                {profile.country && (
                    <div>
                        <span className="text-gray-600">Country:</span>
                        <span className="ml-2 font-medium">{profile.country}</span>
                    </div>
                )}
                {profile.currency && (
                    <div>
                        <span className="text-gray-600">Currency:</span>
                        <span className="ml-2 font-medium">{profile.currency}</span>
                    </div>
                )}
                {profile.marketCap && (
                    <div>
                        <span className="text-gray-600">Market Cap:</span>
                        <span className="ml-2 font-medium">{formatMarketCap(profile.marketCap)}</span>
                    </div>
                )}
                {profile.inception && (
                    <div>
                        <span className="text-gray-600">Inception:</span>
                        <span className="ml-2 font-medium">{profile.inception}</span>
                    </div>
                )}
                {profile.oldestPrice && (
                    <div>
                        <span className="text-gray-600">Oldest Price:</span>
                        <span className="ml-2 font-medium">{profile.oldestPrice}</span>
                    </div>
                )}
                {profile.currentPriceUsd !== undefined && (
                    <div>
                        <span className="text-gray-600">Current Price:</span>
                        <span className="ml-2 font-medium">${profile.currentPriceUsd.toFixed(2)}</span>
                    </div>
                )}
                {profile.ath12m !== undefined && (
                    <div>
                        <span className="text-gray-600">ATH (12m):</span>
                        <span className="ml-2 font-medium">${profile.ath12m.toFixed(2)}</span>
                    </div>
                )}
            </div>
            {profile.description && (
                <div className="mt-4">
                    <span className="text-gray-600">Description:</span>
                    <p className="mt-2 text-gray-700">{profile.description}</p>
                </div>
            )}
            {profile.website && (
                <div className="mt-4">
                    <span className="text-gray-600">Website:</span>
                    <a 
                        href={profile.website} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        className="ml-2 text-blue-600 hover:underline"
                    >
                        {profile.website}
                    </a>
                </div>
            )}
        </div>
    );
}
