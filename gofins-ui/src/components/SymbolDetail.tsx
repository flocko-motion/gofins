import { useState, useEffect, useRef } from 'react';
import { api, imageUrl, favorites, ratings, type SymbolProfile, type UserRating } from '../services/api';
import PriceTable from './PriceTable';
import SymbolProfileView from './SymbolProfile';
import RatingSection from './RatingSection';
import ChartSection, { type ChartSectionHandle } from './ChartSection';

interface SymbolDetailProps {
    symbol: string;
    analysisId?: string; // If provided, use analysis-specific chart endpoints
    onClose?: () => void; // If provided, ESC will trigger close
}

export default function SymbolDetail({ symbol, analysisId, onClose }: SymbolDetailProps) {
    const [profile, setProfile] = useState<SymbolProfile | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [rating, setRating] = useState<number | null>(null);
    const [notes, setNotes] = useState<string>('');
    const [ratingHistory, setRatingHistory] = useState<UserRating[]>([]);
    const [submitting, setSubmitting] = useState(false);
    const [pricesExpanded, setPricesExpanded] = useState(false);
    const [weeklyExpanded, setWeeklyExpanded] = useState(false);
    const [isFavorite, setIsFavorite] = useState(false);
    const ratingSectionRef = useRef<HTMLDivElement>(null);
    const chartSectionRef = useRef<ChartSectionHandle>(null);
    const chartScrollRef = useRef<HTMLDivElement>(null);
    const pricesSectionRef = useRef<HTMLDivElement>(null);
    const weeklySectionRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const fetchProfile = async () => {
            setLoading(true);
            setError(null);
            try {
                const endpoint = analysisId
                    ? `analysis/${analysisId}/profile/${symbol}`
                    : `symbol/${symbol}`;
                const profileData = await api.get<SymbolProfile>(endpoint);
                setProfile(profileData);
                setIsFavorite(profileData.isFavorite || false);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load profile');
            } finally {
                setLoading(false);
            }
        };

        fetchProfile();
    }, [symbol, analysisId]);

    useEffect(() => {
        const fetchRatingHistory = async () => {
            try {
                const data = await ratings.getHistory(symbol);
                setRatingHistory(data || []);
            } catch (err) {
                console.error('Failed to fetch rating history:', err);
            }
        };
        fetchRatingHistory();
    }, [symbol]);

    const toggleFavorite = async () => {
        try {
            const data = await favorites.toggle(symbol);
            setIsFavorite(data.isFavorite);
        } catch (err) {
            console.error('Failed to toggle favorite:', err);
        }
    };

    const togglePrices = () => setPricesExpanded(!pricesExpanded);
    const toggleWeekly = () => setWeeklyExpanded(!weeklyExpanded);

    // Reset expanded state when symbol changes
    useEffect(() => {
        setPricesExpanded(false);
        setWeeklyExpanded(false);
    }, [symbol]);

    const handleDeleteRating = async (ratingId: number) => {
        if (!confirm('Delete this rating?')) return;

        try {
            await ratings.delete(ratingId);
            // Refresh rating history
            const data = await ratings.getHistory(symbol);
            setRatingHistory(data || []);
            // Clear symbol list cache to force refresh when user returns to stocks tab
            sessionStorage.setItem('symbolCacheInvalidated', Date.now().toString());
        } catch (err) {
            alert('Error deleting rating');
        }
    };

    const handleSubmitRating = async () => {
        if (rating === null) {
            alert('Please select a rating');
            return;
        }
        setSubmitting(true);
        try {
            await ratings.add(symbol, rating, notes || undefined);
            // Refresh rating history
            const data = await ratings.getHistory(symbol);
            setRatingHistory(data || []);
            // Reset form and blur textarea
            setRating(null);
            setNotes('');
            const textarea = document.querySelector('textarea[placeholder*="Optional notes"]') as HTMLTextAreaElement;
            if (textarea) textarea.blur();
            // Clear symbol list cache to force refresh when user returns to stocks tab
            sessionStorage.setItem('symbolCacheInvalidated', Date.now().toString());
        } catch (err) {
            alert('Error submitting rating');
        } finally {
            setSubmitting(false);
        }
    };



    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            const target = event.target as HTMLElement;
            const isTyping = target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.tagName === 'SELECT';

            // If typing, only allow Enter in textarea to submit, ignore everything else
            if (isTyping) {
                if (event.key === 'Enter' && target.tagName === 'TEXTAREA') {
                    event.preventDefault();
                    handleSubmitRating();
                }
                return; // Ignore all other hotkeys when typing
            }

            // Hotkeys only work when NOT typing
            if (event.key === 'Escape' && onClose) {
                onClose();
                return;
            }

            if (event.key.toLowerCase() === 'r') {
                event.preventDefault();
                ratingSectionRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
                return;
            }

            if (event.key.toLowerCase() === 'c') {
                event.preventDefault();
                chartSectionRef.current?.toggleChart();
                return;
            }

            if (event.key.toLowerCase() === 'h') {
                event.preventDefault();
                chartSectionRef.current?.toggleHistogram();
                return;
            }

            if (event.key.toLowerCase() === 'm') {
                event.preventDefault();
                pricesSectionRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
                if (!pricesExpanded) {
                    togglePrices();
                }
                return;
            }

            if (event.key.toLowerCase() === 'w') {
                event.preventDefault();
                weeklySectionRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
                if (!weeklyExpanded) {
                    toggleWeekly();
                }
                return;
            }

            if (event.key.toLowerCase() === 'f') {
                event.preventDefault();
                toggleFavorite();
                return;
            }

            if (event.key === 'Enter') {
                event.preventDefault();
                handleSubmitRating();
                return;
            }

            // Number keys for rating (1-5, Shift+1-5 for negative, 0 for neutral)
            // Also support special characters for negative ratings: ! " § $ %
            const negativeRatingMap: Record<string, number> = {
                '!': -1, '"': -2, '§': -3, '$': -4, '%': -5
            };
            
            if (event.key in negativeRatingMap) {
                setRating(negativeRatingMap[event.key]);
                setTimeout(() => {
                    const textarea = document.querySelector('textarea[placeholder*="Optional notes"]') as HTMLTextAreaElement;
                    if (textarea) textarea.focus();
                }, 0);
            } else {
                const num = parseInt(event.key);
                if (!isNaN(num) && num >= 0 && num <= 5) {
                    const newRating = event.shiftKey && num > 0 ? -num : num;
                    setRating(newRating);
                    // Auto-focus the notes textarea
                    setTimeout(() => {
                        const textarea = document.querySelector('textarea[placeholder*="Optional notes"]') as HTMLTextAreaElement;
                        if (textarea) textarea.focus();
                    }, 0);
                }
            }
            
            if (event.key.toLowerCase() === 't') {
                const tradingViewUrl = `https://www.tradingview.com/chart/?symbol=${symbol}`;
                window.open(tradingViewUrl, '_blank', 'noopener,noreferrer');
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [symbol, onClose, handleSubmitRating]);

    // Build image URLs using imageUrl helper
    const chartUrl = analysisId
        ? imageUrl(`analysis/${analysisId}/chart/${symbol}`)
        : imageUrl(`symbol/${symbol}/chart`);

    const histogramUrl = analysisId
        ? imageUrl(`analysis/${analysisId}/histogram/${symbol}`)
        : imageUrl(`symbol/${symbol}/histogram`);

    return (
        <>
            {/* Header */}
            <div className="flex justify-between items-center mb-6">
                <div>
                    <div className="flex items-baseline gap-3">
                        <span
                            onClick={toggleFavorite}
                            className={`cursor-pointer text-3xl hover:scale-125 inline-block transition-transform ${isFavorite ? 'text-yellow-500' : 'text-gray-300'}`}
                        >
                            {isFavorite ? '★' : '☆'}
                        </span>
                        <h2 className="text-2xl font-semibold">{symbol}</h2>
                        {ratingHistory.length > 0 && (
                            <span className={`font-mono text-2xl font-bold ${
                                ratingHistory[0].rating === 0 ? 'text-gray-500' :
                                ratingHistory[0].rating > 0 ? 'text-green-600' : 'text-red-600'
                            }`}>
                                {ratingHistory[0].rating > 0 ? '+' : ''}{ratingHistory[0].rating}
                            </span>
                        )}
                    </div>
                    {profile?.name && (
                        <p className="text-lg text-gray-600 mt-1">{profile.name}</p>
                    )}
                </div>
                <div className="flex items-center gap-2">
                    <button
                        onClick={() => chartScrollRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' })}
                        className="px-2 py-1 text-xs text-gray-400 hover:text-gray-600 border border-gray-300 rounded"
                    >
                        [C]hart / [H]istogram
                    </button>
                    <button
                        onClick={togglePrices}
                        className="px-2 py-1 text-xs text-gray-400 hover:text-gray-600 border border-gray-300 rounded"
                    >
                        [M]onthly
                    </button>
                    <button
                        onClick={toggleWeekly}
                        className="px-2 py-1 text-xs text-gray-400 hover:text-gray-600 border border-gray-300 rounded"
                    >
                        [W]eekly
                    </button>
                    <button
                        onClick={() => ratingSectionRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' })}
                        className="px-2 py-1 text-xs text-gray-400 hover:text-gray-600 border border-gray-300 rounded"
                    >
                        [R]atings
                    </button>
                    <a
                        href={`https://www.tradingview.com/chart/?symbol=${symbol}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="px-2 py-1 text-xs text-gray-400 hover:text-gray-600 border border-gray-300 rounded"
                    >
                        [T]radingView
                    </a>
                    {onClose && (
                        <button
                            onClick={onClose}
                            className="px-2 py-1 text-xs text-gray-400 hover:text-gray-600 border border-gray-300 rounded"
                        >
                            [ESC] close
                        </button>
                    )}
                </div>
            </div>

            {/* Charts */}
            <ChartSection
                ref={chartSectionRef}
                chartUrl={chartUrl}
                histogramUrl={histogramUrl}
                symbol={symbol}
                sectionRef={chartScrollRef}
            />

            {/* Monthly Prices Table */}
            <PriceTable
                title="Monthly Prices"
                symbol={symbol}
                interval="monthly"
                expanded={pricesExpanded}
                onToggle={togglePrices}
                dateFormat="monthly"
                sectionRef={pricesSectionRef}
            />

            {/* Weekly Prices Table */}
            <PriceTable
                title="Weekly Prices"
                symbol={symbol}
                interval="weekly"
                expanded={weeklyExpanded}
                onToggle={toggleWeekly}
                dateFormat="weekly"
                sectionRef={weeklySectionRef}
            />

            {/* Profile Information */}
            <SymbolProfileView profile={profile} loading={loading} error={error} />

            {/* Rating Section */}
            {profile && (
                <RatingSection
                    symbol={symbol}
                    rating={rating}
                    setRating={setRating}
                    notes={notes}
                    setNotes={setNotes}
                    submitting={submitting}
                    ratingHistory={ratingHistory}
                    onSubmit={handleSubmitRating}
                    onDelete={handleDeleteRating}
                    sectionRef={ratingSectionRef}
                />
            )}
        </>
    );
}
