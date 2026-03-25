import React, { useState, useRef, useEffect } from 'react';
import { MapComponent } from '../../../shared/components/Map';
import { ShieldAlert, Bell, Phone, ArrowRight, Navigation, TriangleAlert, MapPin } from 'lucide-react';
import { clsx } from 'clsx';
import mapboxgl from 'mapbox-gl';
import { fakeBackendAPI } from '../../../services/fakeBackendAPI';

// Components for tabs
const StatusView = () => (
    <div className="space-y-6 pb-20 md:pb-0 p-4 md:p-6">
        {/* HIGH ALERT CARD */}
        <div className="bg-gradient-to-br from-red-900/80 to-red-950/90 border border-red-500 rounded-3xl p-8 text-center shadow-lg shadow-red-900/50 relative overflow-hidden shrink-0 animate-pulse-slow">
            <div className="absolute top-0 left-0 w-full h-1.5 bg-gradient-to-r from-red-500 via-orange-500 to-red-500 animate-shimmer"></div>

            <div className="w-24 h-24 bg-red-500/20 rounded-full flex items-center justify-center mx-auto mb-6 animate-bounce-subtle border border-red-500/30">
                <ShieldAlert size={48} className="text-red-500" />
            </div>

            <h2 className="text-3xl font-bold text-white mb-2 uppercase tracking-tight">High Alert</h2>
            <p className="text-red-200 text-lg font-medium leading-relaxed">
                Immediate Evacuation Ordered for <span className="text-white font-bold underline">Sector Amazon-1</span>.
            </p>

            <div className="mt-8 flex items-center justify-center gap-4 text-xs font-mono text-red-400/80 uppercase tracking-widest bg-red-950/50 py-2 rounded-full border border-red-900/50">
                <span className="animate-pulse text-red-400">● LIVE UPDATES</span>
                <span>•</span>
                <span>AMAZONAS, BR</span>
            </div>
        </div>

        <div className="space-y-4">
            <h3 className="text-zinc-400 text-xs font-bold uppercase tracking-wider px-2">Critical Alerts</h3>
            <div className="bg-red-950/30 border border-red-900/50 rounded-2xl p-4 flex gap-4 items-start hover:bg-red-900/10 transition-colors cursor-pointer group">
                <div className="bg-red-500/20 p-2.5 rounded-xl text-red-500 shrink-0 group-hover:scale-110 transition-transform">
                    <TriangleAlert size={20} />
                </div>
                <div>
                    <h4 className="text-red-100 font-bold text-sm flex items-center gap-2">
                        Wildfire Apporaching
                        <span className="text-[10px] bg-red-500 text-white px-1.5 py-0.5 rounded font-mono">NOW</span>
                    </h4>
                    <p className="text-red-200/60 text-xs mt-1 leading-relaxed">Active fire front moving West at 40km/h. Smoke visibility &lt; 50m.</p>
                </div>
            </div>
        </div>

        <div className="space-y-4">
            <h3 className="text-zinc-400 text-xs font-bold uppercase tracking-wider px-2">Evacuation Plan</h3>
            <div className="grid grid-cols-1 gap-3">
                <button className="w-full bg-emerald-500 hover:bg-emerald-600 border border-emerald-400/50 p-4 rounded-2xl flex items-center justify-between group transition-all shadow-lg shadow-emerald-900/20">
                    <div className="flex items-center gap-4">
                        <div className="bg-emerald-900/30 text-emerald-100 p-3 rounded-xl border border-emerald-400/30">
                            <Navigation size={22} />
                        </div>
                        <div className="text-left">
                            <div className="text-white font-bold text-lg">Navigate to Safe Zone</div>
                            <div className="text-emerald-100/80 text-xs mt-0.5">Route Verified • 12 mins away</div>
                        </div>
                    </div>
                    <ArrowRight size={20} className="text-white group-hover:translate-x-1 transition-transform" />
                </button>

                <button className="w-full bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 p-4 rounded-2xl flex items-center justify-between group transition-all hover:scale-[1.01]">
                    <div className="flex items-center gap-4">
                        <div className="bg-zinc-500/10 text-zinc-400 p-3 rounded-xl">
                            <Phone size={22} />
                        </div>
                        <div className="text-left">
                            <div className="text-zinc-200 font-medium">Emergency Contacts</div>
                            <div className="text-zinc-500 text-xs mt-0.5">SOS Signal Active</div>
                        </div>
                    </div>
                    <ArrowRight size={18} className="text-zinc-600 group-hover:text-zinc-300" />
                </button>
            </div>
        </div>
    </div>
);

export const SafetyPortal: React.FC = () => {
    const [activeTab, setActiveTab] = useState<'status' | 'map'>('status');
    const mapRef = useRef<mapboxgl.Map | null>(null);

    // State for data from backend API
    const [fires, setFires] = useState<any[]>([]);
    const [evacuationZones, setEvacuationZones] = useState<any[]>([]);

    // Polling for data updates
    useEffect(() => {
        const stopFiresPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getFires(),
            5000,
            (data) => {
                setFires(data.fires || []);
                console.log('🔥 Citizen fires updated:', data.fires?.length || 0);
            }
        );

        const stopZonesPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getEvacuationZones(),
            5000,
            (data) => {
                setEvacuationZones(data.zones || []);
            }
        );

        return () => {
            stopFiresPolling();
            stopZonesPolling();
        };
    }, []);

    // React to data updates and update map
    useEffect(() => {
        const map = mapRef.current;
        if (!map) return;

        // Update Fires
        const fireSource = map.getSource('fires-source') as mapboxgl.GeoJSONSource;
        if (fireSource) {
            const fireFeatures = fires.map(fire => ({
                type: 'Feature',
                geometry: { type: 'Point', coordinates: fire.location },
                properties: { id: fire.id }
            }));
            fireSource.setData({ type: 'FeatureCollection', features: fireFeatures } as any);
        }

        // Update Evacuation Zones (Markers)
        // Clear existing markers (naive approach for now, or manage ref array)
        const markers = document.getElementsByClassName('evac-marker');
        while (markers[0]) {
            markers[0].parentNode?.removeChild(markers[0]);
        }

        evacuationZones.forEach(zone => {
            const el = document.createElement('div');
            el.className = `flex flex-col items-center justify-center group cursor-pointer animate-bounce-slow evac-marker`;
            const colorClass = zone.status === 'safe' ? 'bg-emerald-500 shadow-emerald-500/40 border-white' : 'bg-rose-500 shadow-rose-500/40 border-white';
            const icon = zone.status === 'safe'
                ? `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg>`
                : `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`;

            el.innerHTML = `
                <div class="${colorClass} text-white p-2 rounded-full shadow-lg border-2 transform group-hover:scale-110 transition-transform">
                    ${icon}
                </div>
                <div class="bg-zinc-900/90 text-white text-xs font-bold px-2 py-1 rounded mt-1 shadow-sm border border-zinc-700 whitespace-nowrap uppercase">${zone.name} • ${zone.status}</div>
            `;

            new mapboxgl.Marker({ element: el })
                .setLngLat(zone.location)
                .addTo(map);
        });

    }, [fires, evacuationZones]);

    const handleMapLoad = (map: mapboxgl.Map) => {
        mapRef.current = map;

        // Initialize Fires Source
        map.addSource('fires-source', {
            type: 'geojson',
            data: { type: 'FeatureCollection', features: [] }
        });

        // Fire Outer Glow
        map.addLayer({
            id: 'fires-glow',
            type: 'circle',
            source: 'fires-source',
            paint: {
                'circle-radius': 20,
                'circle-color': '#ef4444',
                'circle-opacity': 0.4,
                'circle-blur': 0.6
            }
        });

        // Fire Core
        map.addLayer({
            id: 'fires-core',
            type: 'circle',
            source: 'fires-source',
            paint: {
                'circle-radius': 6,
                'circle-color': '#b91c1c',
                'circle-stroke-width': 2,
                'circle-stroke-color': '#ffffff'
            }
        });

        // Animation Loop
        const animate = () => {
            const time = Date.now() / 1000;
            const opacity = (Math.sin(time * 4) + 1) / 2 * 0.5 + 0.3;
            if (map.getLayer('fires-glow')) {
                map.setPaintProperty('fires-glow', 'circle-opacity', opacity);
                map.setPaintProperty('fires-glow', 'circle-radius', 20 + Math.sin(time * 4) * 8);
            }
            requestAnimationFrame(animate);
        };
        animate();

        // Fly to initial view
        map.flyTo({
            center: [-63.40, -5.20],
            zoom: 8,
            speed: 1.5
        });
    };

    return (
        <div className="h-screen bg-zinc-950 text-zinc-100 flex flex-col md:flex-row overflow-hidden">

            {/* Sidebar (Desktop) / Header & Status (Mobile) */}
            <aside className={clsx(
                "flex flex-col shrink-0 bg-zinc-950/50 backdrop-blur z-20 transition-all duration-300",
                "md:w-[450px] md:border-r md:border-zinc-800 md:h-screen",
                activeTab === 'status' ? "h-screen" : "h-auto"
            )}>
                {/* Universal Header */}
                <header className="h-16 shrink-0 flex items-center justify-between px-6 border-b border-zinc-800 bg-zinc-900/50">
                    <div className="font-bold text-xl tracking-tight flex items-center gap-2">
                        <ShieldAlert className="text-red-500 animate-pulse" size={24} />
                        <span>BLAZE<span className="text-fiery-500">GUARD</span></span>
                    </div>
                    <div className="flex items-center gap-3">
                        <div className="hidden md:flex items-center gap-2 px-3 py-1 bg-zinc-900 rounded-full border border-zinc-800">
                            <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                            <span className="text-xs font-mono text-zinc-400">SYSTEM LIVE</span>
                        </div>
                        <div className="w-9 h-9 bg-zinc-800 rounded-full flex items-center justify-center text-xs text-zinc-400 hover:bg-zinc-700 transition-colors cursor-pointer border border-zinc-700">
                            Me
                        </div>
                    </div>
                </header>

                {/* Status Ticker */}
                <div className="bg-red-900/20 border-b border-red-900/30 py-2 px-4 overflow-hidden relative">
                    <div className="flex items-center gap-4 animate-marquee whitespace-nowrap">
                        <span className="text-xs font-bold text-red-400 flex items-center gap-2">
                            <TriangleAlert size={12} /> CRITICAL EVACUATION ORDER IN EFFECT FOR ZONE BETA
                        </span>
                        <span className="text-zinc-600">•</span>
                        <span className="text-xs font-bold text-amber-400 flex items-center gap-2">
                            HIGH WIND WARNING (40KM/H)
                        </span>
                        <span className="text-zinc-600">•</span>
                        <span className="text-xs font-bold text-emerald-400 flex items-center gap-2">
                            ZONE ALPHA CONFIRMED SAFE
                        </span>
                    </div>
                </div>

                {/* Mobile Tab Toggle */}
                <div className="md:hidden flex border-b border-zinc-800 shrink-0">
                    <button
                        onClick={() => setActiveTab('status')}
                        className={clsx("flex-1 py-3 text-sm font-medium border-b-2 transition-colors", activeTab === 'status' ? "border-red-500 text-red-400 bg-red-500/5" : "border-transparent text-zinc-500")}
                    >
                        My Status
                    </button>
                    <button
                        onClick={() => setActiveTab('map')}
                        className={clsx("flex-1 py-3 text-sm font-medium border-b-2 transition-colors", activeTab === 'map' ? "border-red-500 text-red-400 bg-red-500/5" : "border-transparent text-zinc-500")}
                    >
                        Live Map
                    </button>
                </div>

                {/* Status Content - Visible on Desktop OR when 'status' tab is active on Mobile */}
                <div className={clsx(
                    "overflow-y-auto custom-scrollbar flex-1",
                    "md:block", // Always visible on desktop side-by-side
                    activeTab === 'status' ? "block" : "hidden"
                )}>
                    <StatusView />
                </div>
            </aside>

            {/* Main Map Area */}
            <main className={clsx(
                "relative bg-zinc-900 transition-all duration-300",
                "md:flex-1 md:h-screen md:block", // Desktop: fills remaining width, full height
                activeTab === 'map' ? "flex-1 block" : "hidden" // Mobile: fills remaining height if active
            )}>
                <MapComponent className="h-full w-full" onLoad={handleMapLoad}>
                    {/* Map Layers */}
                </MapComponent>

                {/* Map Overlay Controls */}
                <div className="absolute top-6 left-6 pointer-events-none z-10 w-[240px]">
                    <div className="bg-zinc-950/90 backdrop-blur border border-zinc-800 p-4 rounded-2xl shadow-2xl flex flex-col gap-1 pointer-events-auto">
                        <div className="flex items-center gap-3 mb-2">
                            <div className="w-3 h-3 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_10px_rgba(16,185,129,0.6)]" />
                            <span className="text-sm font-bold text-zinc-200">Evacuation Zones</span>
                        </div>
                        {evacuationZones.slice(0, 3).map(zone => (
                            <div key={zone.id} className={`flex justify-between items-center bg-${zone.status === 'safe' ? 'emerald' : 'red'}-500/10 p-2 rounded-lg border border-${zone.status === 'safe' ? 'emerald' : 'red'}-500/20 mb-1`}>
                                <span className="text-xs text-zinc-400">{zone.name}</span>
                                <span className={`text-xs font-bold text-${zone.status === 'safe' ? 'emerald' : 'red'}-400 flex items-center gap-1 uppercase`}>
                                    {zone.status === 'safe' ? <MapPin size={10} /> : <TriangleAlert size={10} />} {zone.status}
                                </span>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="md:hidden absolute bottom-8 left-0 right-0 text-center pointer-events-none">
                    <div className="inline-block bg-red-600/90 backdrop-blur text-white font-bold text-[10px] px-4 py-1.5 rounded-full border border-red-400 shadow-lg shadow-red-900/40 animate-pulse">
                        CRITICAL ALERT ACTIVE
                    </div>
                </div>

                {/* Legend */}
                <div className="absolute bottom-6 right-6 bg-zinc-900/80 backdrop-blur p-3 rounded-xl border border-zinc-800 text-xs hidden md:block">
                    <div className="flex items-center gap-2 mb-1">
                        <span className="w-2 h-2 rounded-full bg-red-500"></span> <span className="text-zinc-400">Fire Incident</span>
                    </div>
                    <div className="flex items-center gap-2 mb-1">
                        <span className="w-2 h-2 rounded-full bg-rose-500 border border-white"></span> <span className="text-zinc-400">Evacuation Zone</span>
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="w-2 h-2 rounded-full bg-emerald-500 border border-white"></span> <span className="text-zinc-400">Safe Zone</span>
                    </div>
                </div>
            </main>

        </div>
    );
};
