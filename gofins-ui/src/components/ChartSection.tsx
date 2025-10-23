import { useState, useImperativeHandle, forwardRef } from 'react';

interface ChartSectionProps {
    chartUrl: string;
    histogramUrl: string;
    symbol: string;
    sectionRef?: React.RefObject<HTMLDivElement | null>;
}

export interface ChartSectionHandle {
    showChart: () => void;
    showHistogram: () => void;
    toggleChart: () => void;
    toggleHistogram: () => void;
    closeFullscreen: () => void;
}

const ChartSection = forwardRef<ChartSectionHandle, ChartSectionProps>(
    ({ chartUrl, histogramUrl, symbol, sectionRef }, ref) => {
        const [fullscreenImage, setFullscreenImage] = useState<string | null>(null);

        useImperativeHandle(ref, () => ({
            showChart: () => setFullscreenImage(chartUrl),
            showHistogram: () => setFullscreenImage(histogramUrl),
            toggleChart: () => setFullscreenImage(prev => prev === chartUrl ? null : chartUrl),
            toggleHistogram: () => setFullscreenImage(prev => prev === histogramUrl ? null : histogramUrl),
            closeFullscreen: () => setFullscreenImage(null),
        }));

    return (
        <>
            {/* Charts */}
            <div ref={sectionRef} className="mb-8">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    <div>
                        <img
                            src={chartUrl}
                            alt={`Chart for ${symbol}`}
                            className="w-full h-[50vh] object-contain border border-gray-200 rounded cursor-pointer hover:opacity-90"
                            onClick={() => setFullscreenImage(chartUrl)}
                        />
                    </div>
                    <div>
                        <img
                            src={histogramUrl}
                            alt={`Histogram for ${symbol}`}
                            className="w-full h-[50vh] object-contain border border-gray-200 rounded cursor-pointer hover:opacity-90"
                            onClick={() => setFullscreenImage(histogramUrl)}
                        />
                    </div>
                </div>
            </div>

            {/* Fullscreen Image Modal */}
            {fullscreenImage && (
                <div
                    className="fixed inset-0 bg-black z-50 flex items-center justify-center"
                    onClick={() => setFullscreenImage(null)}
                >
                    <img
                        src={fullscreenImage}
                        alt="Fullscreen view"
                        className="w-full h-full object-contain"
                        onClick={(e) => e.stopPropagation()}
                    />
                    <button
                        className="absolute top-2 right-2 text-white text-3xl font-bold hover:text-gray-400 px-3 py-1"
                        onClick={() => setFullscreenImage(null)}
                    >
                        ×
                    </button>
                </div>
            )}
        </>
    );
});

ChartSection.displayName = 'ChartSection';

export default ChartSection;
