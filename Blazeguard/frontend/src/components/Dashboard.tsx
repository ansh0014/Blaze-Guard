import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { motion } from 'framer-motion';
import { Flame, Truck, Users, Clock, ShieldAlert, LogOut, Activity, MapPin, Wind, Layers, Bell } from 'lucide-react';

export const Dashboard: React.FC = () => {
    const { user, logout } = useAuth();
    const [activeTab, setActiveTab] = useState<'incidents' | 'logistics' | 'ai-status'>('incidents');

    // Mock logs representing NASA FIRMS stream & Go agent event communications
    const [incidents] = useState([
        {
            id: 'INC-8292',
            source: 'NASA FIRMS VIIRS',
            location: '28.6304° N, 77.2177° E (Delhi)',
            severity: 'CRITICAL',
            confidence: '94%',
            time: 'Just now',
            status: 'Active',
            action: 'Logistics Route Dispatched',
        },
        {
            id: 'INC-8291',
            source: 'Citizen Upload',
            location: '28.5912° N, 77.1950° E (Chanakyapuri)',
            severity: 'HIGH',
            confidence: 'YOLOv8 89%',
            time: '12 mins ago',
            status: 'Deploying',
            action: 'Chanakyapuri Fire Station dispatched (ETA: 4.5 min)',
        },
        {
            id: 'INC-8290',
            source: 'NASA FIRMS MODIS',
            location: '28.5494° N, 77.2516° E (Nehru Place)',
            severity: 'MEDIUM',
            confidence: '78%',
            time: '45 mins ago',
            status: 'Contained',
            action: 'Completed',
        },
    ]);

    const [logisticsRoutes] = useState([
        {
            from: 'Connaught Place Fire Station',
            to: 'Zone NCR-Delhi (INC-8292)',
            distance: '2.4 km',
            duration: '5.2 mins',
            status: 'En Route',
            trucks: 3,
        },
        {
            from: 'Chanakyapuri Fire Station',
            to: 'Zone NCR-Chanakya (INC-8291)',
            distance: '3.1 km',
            duration: '7.8 mins',
            status: 'Dispatched',
            trucks: 2,
        },
    ]);

    const [aiEvolutionLogs] = useState([
        { model: 'YOLOv8 Visual Detector', version: 'v1.4.2', driftScore: '0.02', status: 'Stable' },
        { model: 'MODIS Satellite Classifier', version: 'v2.1.0', driftScore: '0.04', status: 'Stable' },
        { model: 'Spread Evacuation Heuristics', version: 'v1.0.5', driftScore: '0.01', status: 'Optimal' },
    ]);

    return (
        <div className="min-h-screen bg-zinc-950 text-zinc-100 flex flex-col relative overflow-hidden font-sans">
            {/* Ambient Lighting Gradients */}
            <div className="absolute top-[-10%] right-[-10%] w-[600px] h-[600px] bg-fiery-500/10 rounded-full blur-[150px] pointer-events-none mix-blend-screen" />
            <div className="absolute bottom-[-10%] left-[-5%] w-[500px] h-[500px] bg-fiery-600/10 rounded-full blur-[120px] pointer-events-none mix-blend-screen" />

            {/* Grid Overlay */}
            <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.015)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.015)_1px,transparent_1px)] bg-[size:30px_30px] pointer-events-none" />

            {/* Navbar */}
            <header className="relative z-50 border-b border-zinc-900 bg-zinc-950/80 backdrop-blur-md">
                <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
                    <div className="flex items-center gap-2 font-black text-xl tracking-tighter">
                        <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-fiery-500 to-fiery-400 flex items-center justify-center shadow-lg shadow-fiery-500/20 animate-pulse">
                            <Activity size={18} className="text-white" />
                        </div>
                        <span>BLAZE<span className="text-fiery-500">GUARD</span></span>
                    </div>

                    {user && (
                        <div className="flex items-center gap-4">
                            <div className="flex items-center gap-3 bg-zinc-900/60 border border-zinc-800 rounded-full py-1.5 pl-2.5 pr-4">
                                {user.picture ? (
                                    <img
                                        src={user.picture}
                                        alt={user.name}
                                        className="w-8 h-8 rounded-full border border-fiery-500"
                                    />
                                ) : (
                                    <div className="w-8 h-8 rounded-full bg-fiery-600 flex items-center justify-center text-white font-bold text-sm">
                                        {user.name.charAt(0)}
                                    </div>
                                )}
                                <div className="text-left">
                                    <div className="text-xs font-bold text-white tracking-tight leading-none mb-0.5">{user.name}</div>
                                    <div className="text-[10px] text-zinc-400 uppercase tracking-wider leading-none">{user.role}</div>
                                </div>
                            </div>
                            <button
                                onClick={logout}
                                className="p-2 rounded-full bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-400 hover:text-white transition-all cursor-pointer"
                                title="Sign Out"
                            >
                                <LogOut size={18} />
                            </button>
                        </div>
                    )}
                </div>
            </header>

            {/* Content Container */}
            <main className="flex-1 max-w-7xl mx-auto px-6 py-8 relative z-10 w-full">
                {/* System State Banner */}
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8 bg-zinc-900/40 border border-zinc-800/80 rounded-2xl p-6 backdrop-blur-md">
                    <div>
                        <div className="flex items-center gap-2 mb-1">
                            <span className="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse" />
                            <h2 className="text-lg font-black text-white">Central Operations Center</h2>
                        </div>
                        <p className="text-zinc-400 text-xs font-mono">
                            Supabase database active &bull; Redis State caching enabled &bull; Kafka event listeners attached
                        </p>
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="text-xs text-zinc-500 font-mono">Incident Detection Threshold:</span>
                        <span className="px-2 py-1 rounded bg-fiery-950/60 border border-fiery-800/40 text-xs font-bold text-fiery-400 font-mono">
                            YOLOv8 &gt; 50%
                        </span>
                    </div>
                </div>

                {/* Dashboard Metrics Grid */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
                    {[
                        { title: 'Active Hotspots', value: '3', sub: 'NASA Satellite confirmed', icon: Flame, color: 'text-amber-500 bg-amber-500/10 border-amber-500/20' },
                        { title: 'Seeded Stations', value: '5', sub: 'Delhi region operational', icon: Truck, color: 'text-blue-500 bg-blue-500/10 border-blue-500/20' },
                        { title: 'Dispatch Vehicles', value: '22', sub: '8 Trucks &bull; 6 Ambulances', icon: Users, color: 'text-emerald-500 bg-emerald-500/10 border-emerald-500/20' },
                        { title: 'Mean Dispatch ETA', value: '6.5 min', sub: 'Mapbox dynamic traffic', icon: Clock, color: 'text-rose-500 bg-rose-500/10 border-rose-500/20' },
                    ].map((stat, index) => {
                        const Icon = stat.icon;
                        return (
                            <motion.div
                                key={index}
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: index * 0.1 }}
                                className={`bg-zinc-900/30 border border-zinc-800 rounded-2xl p-6 backdrop-blur-sm relative overflow-hidden`}
                            >
                                <div className="flex items-start justify-between">
                                    <div>
                                        <p className="text-zinc-500 text-xs font-bold uppercase tracking-wider mb-1">{stat.title}</p>
                                        <h3 className="text-3xl font-black text-white tracking-tight mb-2">{stat.value}</h3>
                                        <p className="text-[10px] text-zinc-400 font-mono" dangerouslySetInnerHTML={{ __html: stat.sub }} />
                                    </div>
                                    <div className={`w-12 h-12 rounded-xl flex items-center justify-center border ${stat.color}`}>
                                        <Icon size={24} />
                                    </div>
                                </div>
                            </motion.div>
                        );
                    })}
                </div>

                {/* Main Dashboard Section Split */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Left & Middle Column (Panels) */}
                    <div className="lg:col-span-2 flex flex-col gap-6">
                        {/* Tab Headers */}
                        <div className="flex border-b border-zinc-800 gap-6">
                            {[
                                { id: 'incidents', label: 'NASA & Citizen Incidents', icon: Flame },
                                { id: 'logistics', label: 'Dispatches & Mapbox routes', icon: Truck },
                                { id: 'ai-status', label: 'Oracle 23ai Drift Logs', icon: Layers },
                            ].map((tab) => {
                                const TabIcon = tab.icon;
                                const isActive = activeTab === tab.id;
                                return (
                                    <button
                                        key={tab.id}
                                        onClick={() => setActiveTab(tab.id as any)}
                                        className={`flex items-center gap-2 pb-4 text-sm font-bold tracking-tight transition-all relative cursor-pointer ${
                                            isActive ? 'text-fiery-500 font-black' : 'text-zinc-400 hover:text-white'
                                        }`}
                                    >
                                        <TabIcon size={16} />
                                        {tab.label}
                                        {isActive && (
                                            <motion.div
                                                layoutId="activeTabIndicator"
                                                className="absolute bottom-0 left-0 right-0 h-0.5 bg-fiery-500"
                                            />
                                        )}
                                    </button>
                                );
                            })}
                        </div>

                        {/* Tab Content Panels */}
                        <div className="flex-1 bg-zinc-900/20 border border-zinc-800/80 rounded-2xl p-6 backdrop-blur-md">
                            {activeTab === 'incidents' && (
                                <div className="space-y-4">
                                    {incidents.map((incident) => (
                                        <div key={incident.id} className="p-4 rounded-xl bg-zinc-900/60 border border-zinc-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
                                            <div>
                                                <div className="flex items-center gap-2 mb-1.5">
                                                    <span className="text-xs font-mono font-bold text-zinc-500">{incident.id}</span>
                                                    <span className="px-2 py-0.5 rounded bg-zinc-800 text-[10px] text-zinc-400 font-bold uppercase tracking-wider">{incident.source}</span>
                                                    <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                                                        incident.severity === 'CRITICAL' ? 'bg-red-950 border border-red-800 text-red-400' : 'bg-amber-950 border border-amber-800 text-amber-400'
                                                    }`}>{incident.severity}</span>
                                                </div>
                                                <div className="flex items-center gap-1 text-sm font-bold text-white mb-1">
                                                    <MapPin size={14} className="text-fiery-500" />
                                                    {incident.location}
                                                </div>
                                                <p className="text-xs text-zinc-400 font-mono">
                                                    Event Confidence: <strong className="text-white">{incident.confidence}</strong> &bull; Status: <strong className="text-fiery-400">{incident.action}</strong>
                                                </p>
                                            </div>
                                            <div className="flex items-center gap-3 text-right">
                                                <span className="text-[10px] text-zinc-500 font-mono">{incident.time}</span>
                                                <span className="px-3 py-1 rounded-full bg-emerald-950/60 border border-emerald-800 text-emerald-400 text-xs font-black">
                                                    {incident.status}
                                                </span>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}

                            {activeTab === 'logistics' && (
                                <div className="space-y-4">
                                    {logisticsRoutes.map((route, i) => (
                                        <div key={i} className="p-4 rounded-xl bg-zinc-900/60 border border-zinc-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
                                            <div>
                                                <h4 className="text-sm font-black text-white mb-1">Route calculation successfully retrieved</h4>
                                                <p className="text-xs text-zinc-400 font-mono mb-2">
                                                    Origin: <strong>{route.from}</strong> &rarr; Target: <strong>{route.to}</strong>
                                                </p>
                                                <div className="flex items-center gap-4">
                                                    <span className="text-xs text-zinc-400 font-mono">Distance: <strong className="text-white">{route.distance}</strong></span>
                                                    <span className="text-xs text-zinc-400 font-mono">Duration: <strong className="text-fiery-400">{route.duration}</strong></span>
                                                    <span className="text-xs text-zinc-400 font-mono">Units: <strong className="text-white">{route.trucks} Engines</strong></span>
                                                </div>
                                            </div>
                                            <div className="px-3 py-1 rounded bg-blue-950 border border-blue-800 text-blue-400 text-xs font-bold">
                                                {route.status}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}

                            {activeTab === 'ai-status' && (
                                <div className="space-y-4">
                                    <div className="p-4 rounded-xl bg-zinc-900/40 border border-zinc-800 text-xs font-mono text-zinc-400 mb-2">
                                        Oracle Database 23ai tracking model drift counters via self-evolving triggers:
                                    </div>
                                    {aiEvolutionLogs.map((logItem, i) => (
                                        <div key={i} className="p-4 rounded-xl bg-zinc-900/60 border border-zinc-800 flex items-center justify-between">
                                            <div>
                                                <h4 className="text-sm font-bold text-white mb-0.5">{logItem.model}</h4>
                                                <p className="text-[10px] text-zinc-500 font-mono">Active Version: {logItem.version}</p>
                                            </div>
                                            <div className="flex items-center gap-6">
                                                <div className="text-right">
                                                    <span className="text-[10px] text-zinc-500 block">Drift Index</span>
                                                    <span className="text-xs font-bold text-emerald-400 font-mono">{logItem.driftScore}</span>
                                                </div>
                                                <span className="px-2 py-0.5 rounded bg-emerald-950 border border-emerald-800 text-emerald-400 text-[10px] font-bold">
                                                    {logItem.status}
                                                </span>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Right Column: Live Map Alert Indicator */}
                    <div className="flex flex-col gap-6">
                        <div className="bg-zinc-900/20 border border-zinc-800 rounded-2xl p-6 backdrop-blur-md flex-1">
                            <h3 className="text-md font-bold text-white mb-4 flex items-center gap-2">
                                <Bell size={18} className="text-fiery-500" />
                                Real-Time Operations Map
                            </h3>

                            {/* Mock Visual Map Placeholder */}
                            <div className="h-64 rounded-xl border border-zinc-800 relative bg-zinc-950 overflow-hidden mb-6 flex items-center justify-center">
                                {/* Grid map grid lines */}
                                <div className="absolute inset-0 bg-[radial-gradient(rgba(244,63,94,0.06)_1.5px,transparent_1.5px)] bg-[size:16px_16px]" />
                                
                                {/* Map Marker 1 */}
                                <div className="absolute top-[35%] left-[45%] flex flex-col items-center">
                                    <div className="w-3 h-3 rounded-full bg-rose-600 animate-ping absolute" />
                                    <div className="w-3 h-3 rounded-full bg-rose-500 border border-white relative" />
                                    <span className="px-1.5 py-0.5 rounded bg-zinc-900/90 border border-zinc-800 text-[8px] font-mono text-white mt-1 shadow-lg">Delhi Core</span>
                                </div>

                                {/* Map Marker 2 */}
                                <div className="absolute top-[60%] left-[65%] flex flex-col items-center">
                                    <div className="w-2.5 h-2.5 rounded-full bg-orange-600 animate-ping absolute" />
                                    <div className="w-2.5 h-2.5 rounded-full bg-orange-500 border border-white relative" />
                                    <span className="px-1.5 py-0.5 rounded bg-zinc-900/90 border border-zinc-800 text-[8px] font-mono text-white mt-1 shadow-lg">Dwarka Sec-6</span>
                                </div>

                                {/* Map center crosshairs */}
                                <div className="w-6 h-6 border border-zinc-800/40 rounded-full flex items-center justify-center text-zinc-800/40 text-[10px] font-mono">+</div>
                                
                                <span className="absolute bottom-2 right-2 px-2 py-1 rounded bg-zinc-900/80 border border-zinc-800 text-[9px] font-mono text-zinc-500">
                                    NCR Grid System
                                </span>
                            </div>

                            {/* Map Information */}
                            <div className="space-y-4">
                                <div className="flex items-start gap-3">
                                    <div className="w-8 h-8 rounded-lg bg-fiery-500/10 border border-fiery-500/20 flex items-center justify-center text-fiery-500 flex-shrink-0">
                                        <Wind size={16} />
                                    </div>
                                    <div>
                                        <h4 className="text-xs font-bold text-white mb-0.5">Atmospheric Conditions</h4>
                                        <p className="text-[10px] text-zinc-400 leading-normal">
                                            Delhi reporting 34°C, 42% Humidity, Wind 12 km/h from South-West. Spread vectors predicted North-East.
                                        </p>
                                    </div>
                                </div>

                                <div className="flex items-start gap-3">
                                    <div className="w-8 h-8 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center text-blue-500 flex-shrink-0">
                                        <Layers size={16} />
                                    </div>
                                    <div>
                                        <h4 className="text-xs font-bold text-white mb-0.5">PostGIS Grid Lookup</h4>
                                        <p className="text-[10px] text-zinc-400 leading-normal">
                                            Scanning 5 regional stations. Proximity calculations validated using native PostgreSQL geography functions.
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default Dashboard;
