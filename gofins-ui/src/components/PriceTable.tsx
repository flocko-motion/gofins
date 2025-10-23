import { useState, useEffect } from 'react';
import { api, type PriceData } from '../services/api';
import { formatPrice, formatDate } from '../utils/format';

interface PriceTableProps {
    title: string;
    symbol: string;
    interval: 'monthly' | 'weekly';
    expanded: boolean;
    onToggle: () => void;
    dateFormat?: 'monthly' | 'weekly';
    sectionRef?: React.RefObject<HTMLDivElement | null>;
}

export default function PriceTable({
    title,
    symbol,
    interval,
    expanded,
    onToggle,
    dateFormat = 'monthly',
    sectionRef
}: PriceTableProps) {
    const [prices, setPrices] = useState<PriceData[]>([]);
    const [loading, setLoading] = useState(false);
    const [fetched, setFetched] = useState(false);

    useEffect(() => {
        const fetchPrices = async () => {
            if (fetched || !expanded) return;
            setLoading(true);
            try {
                const data = await api.get<{ prices: PriceData[] }>(`prices/${interval}/${symbol}`);
                setPrices(data.prices || []);
                setFetched(true);
            } catch (err) {
                console.error(`Failed to fetch ${interval} prices:`, err);
            } finally {
                setLoading(false);
            }
        };

        if (expanded) {
            fetchPrices();
        }
    }, [expanded, fetched, interval, symbol]);

    // Reset when symbol changes
    useEffect(() => {
        setPrices([]);
        setFetched(false);
    }, [symbol]);

    return (
        <div ref={sectionRef} className="mb-8">
            <button
                onClick={onToggle}
                className="flex items-center gap-2 text-lg font-semibold mb-4 hover:text-gray-700"
            >
                <span>{expanded ? '▼' : '▶'}</span>
                <span>{title}</span>
            </button>
            {expanded && (
                <div>
                    {loading ? (
                        <div className="text-center py-4">
                            <p className="text-gray-500">Loading prices...</p>
                        </div>
                    ) : prices.length > 0 ? (
                        <div className="overflow-x-auto">
                            <table className="min-w-full border border-gray-200 text-sm">
                                <thead className="bg-gray-50">
                                    <tr>
                                        <th className="px-4 py-2 text-left font-medium text-gray-600 border-b">Date</th>
                                        <th className="px-4 py-2 text-right font-medium text-gray-600 border-b">Open</th>
                                        <th className="px-4 py-2 text-right font-medium text-gray-600 border-b">High</th>
                                        <th className="px-4 py-2 text-right font-medium text-gray-600 border-b">Low</th>
                                        <th className="px-4 py-2 text-right font-medium text-gray-600 border-b">Close</th>
                                        <th className="px-4 py-2 text-right font-medium text-gray-600 border-b">Avg</th>
                                        <th className="px-4 py-2 text-right font-medium text-gray-600 border-b">YoY %</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {prices.map((price, idx) => (
                                        <tr key={idx} className="hover:bg-gray-50">
                                            <td className="px-4 py-2 border-b">
                                                {formatDate(price.Date, dateFormat === 'monthly' ? 'short' : 'long')}
                                            </td>
                                            <td className="px-4 py-2 text-right border-b">{formatPrice(price.Open)}</td>
                                            <td className="px-4 py-2 text-right border-b">{formatPrice(price.High)}</td>
                                            <td className="px-4 py-2 text-right border-b">{formatPrice(price.Low)}</td>
                                            <td className="px-4 py-2 text-right border-b font-medium">{formatPrice(price.Close)}</td>
                                            <td className="px-4 py-2 text-right border-b">{formatPrice(price.Avg)}</td>
                                            <td className={`px-4 py-2 text-right border-b ${
                                                price.YoY === null ? 'text-gray-400' :
                                                price.YoY > 0 ? 'text-green-600' :
                                                price.YoY < 0 ? 'text-red-600' : 'text-gray-600'
                                            }`}>
                                                {price.YoY === null ? 'N/A' : `${price.YoY > 0 ? '+' : ''}${price.YoY.toFixed(1)}%`}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    ) : (
                        <p className="text-gray-500">No {dateFormat} price data available</p>
                    )}
                </div>
            )}
        </div>
    );
}
