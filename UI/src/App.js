// src/App.js
import React, { useState, useEffect } from 'react';
import PackSizeManager from './components/PackSizeManager';
import PackCalculator from './components/PackCalculator';

function App() {
    const [activeTab, setActiveTab] = useState('calculator');

    return (
        <div className="min-h-screen bg-gray-50 p-6">
            <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-md overflow-hidden">
                <h1 className="text-2xl font-bold bg-blue-600 text-white p-4">
                    Order Packs Calculator
                </h1>

                <div className="flex border-b">
                    <button
                        className={`px-4 py-2 font-medium ${activeTab === 'calculator' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500'}`}
                        onClick={() => setActiveTab('calculator')}
                    >
                        Calculate Packs
                    </button>
                    <button
                        className={`px-4 py-2 font-medium ${activeTab === 'manage' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500'}`}
                        onClick={() => setActiveTab('manage')}
                    >
                        Manage Pack Sizes
                    </button>
                </div>

                <div className="p-6">
                    {activeTab === 'calculator' ? <PackCalculator /> : <PackSizeManager />}
                </div>
            </div>
        </div>
    );
}

export default App;