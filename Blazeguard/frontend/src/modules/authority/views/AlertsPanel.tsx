import React, { useState, useEffect } from 'react';
import { AlertTriangle, MoreHorizontal } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';
import { fakeBackendAPI } from '../../../services/fakeBackendAPI';

export const AlertsPanel: React.FC = () => {
    // State for data from backend API
    const [alerts, setAlerts] = useState<any[]>([]);
    const [fires, setFires] = useState<any[]>([]);

    // Polling for data updates
    useEffect(() => {
        // Start polling for alerts
        const stopAlertsPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getAlerts(),
            5000,
            (data) => {
                setAlerts(data.alerts || []);
                console.log('🚨 Alerts updated:', data.alerts?.length || 0);
            }
        );

        // Start polling for fires
        const stopFiresPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getFires(),
            5000,
            (data) => {
                setFires(data.fires || []);
            }
        );

        // Cleanup on unmount
        return () => {
            stopAlertsPolling();
            stopFiresPolling();
        };
    }, []);

    return (
        <div className="p-6 max-w-5xl mx-auto space-y-8">
            <div>
                <h1 className="text-2xl font-bold text-zinc-100 mb-2">Alerts & Incidents</h1>
                <div className="flex gap-4 border-b border-zinc-800 pb-4">
                    <button className="text-fiery-400 border-b-2 border-fiery-400 pb-4 -mb-4.5 font-medium text-sm">Active Alerts ({alerts.length})</button>
                    <button className="text-zinc-500 hover:text-zinc-300 transition-colors pb-4 -mb-4.5 font-medium text-sm">Resolved History</button>
                </div>
            </div>

            {/* System Alerts Section */}
            <section className="space-y-4">
                <h2 className="text-xs font-bold text-zinc-500 uppercase tracking-widest pl-1">System Broadcasts</h2>
                <div className="grid gap-4">
                    {alerts.map(alert => (
                        <div key={alert.id} className="bg-zinc-900 border border-zinc-800 rounded-xl p-5 flex items-start gap-4 hover:border-zinc-700 transition-colors group">
                            <div className={`
                        shrink-0 w-10 h-10 rounded-full flex items-center justify-center
                        ${alert.level === 'danger' ? 'bg-rose-500/10 text-rose-500' : 'bg-amber-500/10 text-amber-500'}
                    `}>
                                <AlertTriangle size={20} />
                            </div>
                            <div className="flex-1">
                                <div className="flex justify-between items-start">
                                    <h3 className="text-zinc-200 font-semibold">{alert.location}</h3>
                                    <span className="text-zinc-500 text-xs font-mono">{formatDistanceToNow(alert.timestamp, { addSuffix: true })}</span>
                                </div>
                                <p className="text-zinc-400 text-sm mt-1">{alert.message}</p>
                            </div>
                            <button className="text-zinc-600 hover:text-zinc-300 opacity-0 group-hover:opacity-100 transition-opacity">
                                <MoreHorizontal size={18} />
                            </button>
                        </div>
                    ))}
                </div>
            </section>

            {/* Active Incidents List derived from Fires */}
            <section className="space-y-4">
                <h2 className="text-xs font-bold text-zinc-500 uppercase tracking-widest pl-1">Active Detected Incidents</h2>
                <div className="bg-zinc-900 border border-zinc-800 rounded-xl overflow-hidden">
                    <table className="w-full text-left text-sm">
                        <thead className="bg-zinc-950/50 text-zinc-500 font-medium border-b border-zinc-800">
                            <tr>
                                <th className="px-6 py-3">ID</th>
                                <th className="px-6 py-3">Severity</th>
                                <th className="px-6 py-3">Zone</th>
                                <th className="px-6 py-3">Detected</th>
                                <th className="px-6 py-3">Source</th>
                                <th className="px-6 py-3 text-right">Action</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-zinc-800">
                            {fires.map(fire => (
                                <tr key={fire.id} className="hover:bg-zinc-800/30 transition-colors">
                                    <td className="px-6 py-4 font-mono text-zinc-300">#{fire.id}</td>
                                    <td className="px-6 py-4">
                                        <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border ${fire.severity === 'critical' ? 'bg-rose-950/30 text-rose-400 border-rose-500/20' :
                                            fire.severity === 'medium' ? 'bg-amber-950/30 text-amber-400 border-amber-500/20' :
                                                'bg-emerald-950/30 text-fiery-400 border-fiery-500/20'
                                            }`}>
                                            <div className={`w-1.5 h-1.5 rounded-full ${fire.severity === 'critical' ? 'bg-rose-500' :
                                                fire.severity === 'medium' ? 'bg-amber-500' : 'bg-fiery-500'
                                                }`} />
                                            {fire.severity.toUpperCase()}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 text-zinc-400">{fire.zoneId}</td>
                                    <td className="px-6 py-4 text-zinc-500">{formatDistanceToNow(fire.detectedAt)} ago</td>
                                    <td className="px-6 py-4 text-zinc-400">{fire.source}</td>
                                    <td className="px-6 py-4 text-right">
                                        <button className="text-fiery-500 hover:text-fiery-400 font-medium text-xs border border-emerald-500/20 bg-fiery-500/10 px-3 py-1.5 rounded-md transition-colors">
                                            Deploy Resources
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </section>
        </div>
    );
};
