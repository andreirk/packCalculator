import React, { useState } from 'react';

function PackCalculator() {
    const [items, setItems] = useState('');
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleCalculate = async () => {
        if (!items || isNaN(items) || items <= 0) {
            setError('Please enter a valid positive number');
            return;
        }

        setLoading(true);
        setError('');
        try {
            const response = await fetch(`/api/calculate?items=${items}`);
            if (!response.ok) throw new Error('Calculation failed');
            const data = await response.json();
            setResult(data.packs);
        } catch (err) {
            setError('Failed to calculate packs');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2 className="text-xl font-semibold mb-4">Calculate packs for order</h2>

            <div className="flex items-center mb-6">
                <input
                    type="number"
                    value={items}
                    onChange={(e) => setItems(e.target.value)}
                    placeholder="Items ordered"
                    className="flex-1 p-2 border border-gray-300 rounded-l focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="1"
                />
                <button
                    onClick={handleCalculate}
                    disabled={loading}
                    className="bg-blue-600 text-white px-4 py-2 rounded-r hover:bg-blue-700 disabled:bg-blue-400"
                >
                    {loading ? 'Calculating...' : 'Calculate'}
                </button>
            </div>

            {error && (
                <div className="mb-4 p-2 bg-red-100 text-red-700 rounded">
                    {error}
                </div>
            )}

            {result && (
                <div>
                    <h3 className="font-medium mb-2">Packages to send:</h3>
                    <table className="min-w-full border">
                        <thead>
                        <tr className="bg-gray-100">
                            <th className="py-2 px-4 border text-left">Pack Size</th>
                            <th className="py-2 px-4 border text-left">Quantity</th>
                        </tr>
                        </thead>
                        <tbody>
                        {Object.entries(result).map(([size, quantity]) => (
                            <tr key={size} className="border-b">
                                <td className="py-2 px-4 border">{size}</td>
                                <td className="py-2 px-4 border">{quantity}</td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}

export default PackCalculator;