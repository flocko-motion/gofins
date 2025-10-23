import { type UserRating } from '../services/api';
import { getRatingColor } from '../utils/format';

interface RatingSectionProps {
    symbol: string;
    rating: number | null;
    setRating: (rating: number | null) => void;
    notes: string;
    setNotes: (notes: string) => void;
    submitting: boolean;
    ratingHistory: UserRating[];
    onSubmit: () => void;
    onDelete: (ratingId: number) => void;
    sectionRef?: React.RefObject<HTMLDivElement | null>;
}

export default function RatingSection({
    symbol,
    rating,
    setRating,
    notes,
    setNotes,
    submitting,
    ratingHistory,
    onSubmit,
    onDelete,
    sectionRef
}: RatingSectionProps) {
    return (
        <div ref={sectionRef} className="bg-white rounded-lg shadow p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">Rate {symbol}</h2>

            {/* Rating Buttons */}
            <div className="flex gap-2 mb-4">
                {[-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5].map((r) => (
                    <button
                        key={r}
                        onClick={() => setRating(r)}
                        className={`px-4 py-2 rounded border-2 transition-all ${rating === r
                                ? `${getRatingColor(r)} bg-white border-current font-bold`
                                : `${getRatingColor(r)} bg-gray-50 border-gray-300 hover:bg-white hover:border-current`
                            }`}
                    >
                        {r > 0 ? `+${r}` : r}
                    </button>
                ))}
            </div>

            {/* Notes Textarea */}
            <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Optional notes about this stock..."
                className="w-full border border-gray-300 rounded p-2 mb-4 h-24"
            />

            {/* Submit Button */}
            <button
                onClick={onSubmit}
                disabled={submitting || rating === null}
                className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:bg-gray-400"
            >
                {submitting ? 'Submitting...' : 'Submit Rating'}
            </button>

            {/* Rating History */}
            {ratingHistory.length > 0 && (
                <div className="mt-6">
                    <h3 className="text-lg font-semibold mb-3">Rating History</h3>
                    <div className="space-y-2">
                        {ratingHistory.map((r) => (
                            <div key={r.id} className="border border-gray-200 rounded p-3 flex justify-between items-start">
                                <div className="flex-1">
                                    <div className="flex items-center gap-2">
                                        <span className={`font-bold ${getRatingColor(r.rating)}`}>
                                            {r.rating > 0 ? `+${r.rating}` : r.rating}
                                        </span>
                                        <span className="text-gray-500 text-sm">
                                            {new Date(r.createdAt).toLocaleString()}
                                        </span>
                                    </div>
                                    {r.notes && (
                                        <p className="text-gray-700 mt-1">{r.notes}</p>
                                    )}
                                </div>
                                <button
                                    onClick={() => onDelete(r.id)}
                                    className="text-red-600 hover:text-red-800 text-sm ml-4"
                                >
                                    Delete
                                </button>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
