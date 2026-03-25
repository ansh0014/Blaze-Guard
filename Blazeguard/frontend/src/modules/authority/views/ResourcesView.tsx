import React, { useState, useEffect } from 'react';
import { Truck, Plane, Users, Navigation } from 'lucide-react';
import { fakeBackendAPI } from '../../../services/fakeBackendAPI';
import { useNavigate } from 'react-router-dom';

export const ResourcesView: React.FC = () => {
    const navigate = useNavigate();
    // State for data from backend API
    const [resources, setResources] = useState<any[]>([]);

    // Polling for data updates
    useEffect(() => {
        const stopPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getResources(),
            5000,
            (data) => {
                setResources(data.resources || []);
            }
        );

        return () => stopPolling();
    }, []);

    return (
        <div className="p-6 max-w-6xl mx-auto space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-zinc-900 border border-zinc-800 p-6 rounded-2xl flex items-center justify-between">
                    <div>
                        <p className="text-zinc-500 text-xs font-bold uppercase tracking-wider">Total Units</p>
                        <h2 className="text-3xl font-black text-white mt-2">{resources.length}</h2>
                    </div>
                    <div className="bg-zinc-800 p-3 rounded-xl text-zinc-400">
                        <Truck size={24} />
                    </div>
                </div>
                <div className="bg-zinc-900 border border-zinc-800 p-6 rounded-2xl flex items-center justify-between">
                    <div>
                        <p className="text-zinc-500 text-xs font-bold uppercase tracking-wider">Available</p>
                        <h2 className="text-3xl font-black text-fiery-400 mt-2">
                            {resources.filter(r => r.status === 'idle').length}
                        </h2>
                    </div>
                    <div className="bg-fiery-500/10 p-3 rounded-xl text-fiery-500">
                        <Users size={24} />
                    </div>
                </div>
                <div className="bg-zinc-900 border border-zinc-800 p-6 rounded-2xl flex items-center justify-between">
                    <div>
                        <p className="text-zinc-500 text-xs font-bold uppercase tracking-wider">Active Missions</p>
                        <h2 className="text-3xl font-black text-amber-400 mt-2">
                            {resources.filter(r => r.status !== 'idle').length}
                        </h2>
                    </div>
                    <div className="bg-amber-500/10 p-3 rounded-xl text-amber-500">
                        <Navigation size={24} />
                    </div>
                </div>
            </div>

            <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
                <div className="p-6 border-b border-zinc-800 flex justify-between items-center">
                    <h3 className="text-lg font-bold text-zinc-100">Live Resource Tracking</h3>
                    <button
                        onClick={() => navigate('/dashboard/map')}
                        className="text-sm text-fiery-500 font-medium hover:text-fiery-400"
                    >
                        View All on Map
                    </button>
                </div>
                <div className="divide-y divide-zinc-800">
                    {resources.map(resource => (
                        <div key={resource.id} className="p-4 flex items-center justify-between hover:bg-zinc-800/30 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className="w-10 h-10 rounded-lg bg-zinc-800 flex items-center justify-center text-zinc-400">
                                    {resource.type === 'drone' ? <Plane size={20} /> : <Truck size={20} />}
                                </div>
                                <div>
                                    <h4 className="text-zinc-200 font-medium">{resource.name}</h4>
                                    <p className="text-zinc-500 text-xs capitalize">{resource.type} • ID: {resource.id}</p>
                                </div>
                            </div>

                            <div className="flex items-center gap-6">
                                <div className="text-right">
                                    <div className="text-xs text-zinc-500">Status</div>
                                    <div className={`text-sm font-medium capitalize ${resource.status === 'active' ? 'text-rose-400' :
                                        resource.status === 'en-route' ? 'text-amber-400' : 'text-fiery-400'
                                        }`}>
                                        {resource.status.replace('-', ' ')}
                                    </div>
                                </div>
                                <button
                                    onClick={() => navigate(`/dashboard/map?resourceId=${resource.id}`)}
                                    className="p-2 hover:bg-zinc-800 rounded-lg text-zinc-500 hover:text-zinc-300 transition-colors"
                                    title="View on Map"
                                >
                                    <Navigation size={18} />
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};
