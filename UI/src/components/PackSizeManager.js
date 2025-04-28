import React, { useState, useEffect } from 'react';

function PackSizeManager() {
    const [packSizes, setPackSizes] = useState([]);
    const [newSize, setNewSize] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [bulkEdit, setBulkEdit] = useState(false);
    const [bulkSizes, setBulkSizes] = useState('');

    const fetchPackSizes = async () => {
        setLoading(true);
        setError('');
        try {
            const response = await fetch('/api/packs');
            if (!response.ok) throw new Error('Failed to fetch');
            const data = await response.json();
            setPackSizes(data.sizes);
            setBulkSizes(data.sizes.join(', '));
        } catch (err) {
            setError('Failed to load pack sizes');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const addPackSize = async () => {
        const size = parseInt(newSize);
        if (isNaN(size) || size <= 0) {
            setError('Please enter a valid positive number');
            return;
        }

        setLoading(true);
        setError('');
        try {
            const response = await fetch('/api/packs/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ size }),
            });
            if (!response.ok) throw new Error('Failed to add');
            setNewSize('');
            await fetchPackSizes();
        } catch (err) {
            setError('Failed to add pack size');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const removePackSize = async (size) => {
        if (!window.confirm(`Remove pack size ${size}?`)) return;

        setLoading(true);
        setError('');
        try {
            const response = await fetch(`/api/packs/${size}`, {
                method: 'DELETE',
            });
            if (!response.ok) throw new Error('Failed to remove');
            await fetchPackSizes();
        } catch (err) {
            setError('Failed to remove pack size');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const updateAllPackSizes = async () => {
        const sizes = bulkSizes.split(',')
            .map(s => parseInt(s.trim()))
            .filter(s => !isNaN(s) && s > 0);

        if (sizes.length === 0) {
            setError('Please enter at least one valid pack size');
            return;
        }

        setLoading(true);
        setError('');
        try {
            const response = await fetch('/api/packs', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ sizes }),
            });
            if (!response.ok) throw new Error('Failed to update');
            await fetchPackSizes();
            setBulkEdit(false);
        } catch (err) {
            setError('Failed to update pack sizes');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPackSizes();
    }, []);

    return (
        <div>
            <h2 className="text-xl font-semibold mb-4">Pack Sizes</h2>

            {error && (
                <div className="mb-4 p-2 bg-red-100 text-red-700 rounded">
                    {error}
                </div>
            )}

            {loading ? (
                <div className="flex justify-center items-center py-8">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                </div>
            ) : (
                <>
                    {!bulkEdit ? (
                        <>
                            <div className="mb-6">
                                <h3 className="font-medium mb-2">Current Pack Sizes</h3>
                                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2">
                                    {packSizes.map(size => (
                                        <div key={size} className="flex items-center justify-between p-2 border rounded">
                                            <span>{size}</span>
                                            <button
                                                onClick={() => removePackSize(size)}
                                                disabled={loading}
                                                className="text-red-600 hover:text-red-800 disabled:text-red-300"
                                            >
                                                Remove
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="mb-6 p-4 border rounded bg-gray-50">
                                <h3 className="font-medium mb-2">Add New Pack Size</h3>
                                <div className="flex">
                                    <input
                                        type="number"
                                        value={newSize}
                                        onChange={(e) => setNewSize(e.target.value)}
                                        placeholder="Enter pack size"
                                        className="flex-1 p-2 border border-gray-300 rounded-l focus:outline-none focus:ring-2 focus:ring-blue-500"
                                        min="1"
                                    />
                                    <button
                                        onClick={addPackSize}
                                        disabled={loading}
                                        className="bg-blue-600 text-white px-4 py-2 rounded-r hover:bg-blue-700 disabled:bg-blue-400"
                                    >
                                        Add
                                    </button>
                                </div>
                            </div>

                            <button
                                onClick={() => setBulkEdit(true)}
                                className="text-blue-600 hover:text-blue-800"
                            >
                                Edit all sizes at once →
                            </button>
                        </>
                    ) : (
                        <>
                            <div className="mb-4">
                                <label className="block font-medium mb-2">
                                    Edit all pack sizes (comma separated)
                                </label>
                                <textarea
                                    value={bulkSizes}
                                    onChange={(e) => setBulkSizes(e.target.value)}
                                    className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    rows={4}
                                />
                            </div>
                            <div className="flex space-x-2">
                                <button
                                    onClick={updateAllPackSizes}
                                    disabled={loading}
                                    className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:bg-blue-400"
                                >
                                    Save Changes
                                </button>
                                <button
                                    onClick={() => {
                                        setBulkEdit(false);
                                        setBulkSizes(packSizes.join(', '));
                                    }}
                                    className="bg-gray-200 text-gray-800 px-4 py-2 rounded hover:bg-gray-300"
                                >
                                    Cancel
                                </button>
                            </div>
                        </>
                    )}
                </>
            )}
        </div>
    );
}

export default PackSizeManager;