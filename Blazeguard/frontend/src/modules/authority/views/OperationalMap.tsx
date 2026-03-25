import React, { useRef, useState, useEffect } from 'react';
import { MapComponent } from '../../../shared/components/Map';
import { useSearchParams } from 'react-router-dom';
import mapboxgl from 'mapbox-gl';
import { fakeBackendAPI } from '../../../services/fakeBackendAPI';

export const OperationalMap: React.FC = () => {
    const mapRef = useRef<mapboxgl.Map | null>(null);
    const resourceMarkersRef = useRef<mapboxgl.Marker[]>([]);
    const stationMarkersRef = useRef<mapboxgl.Marker[]>([]);

    // State for data from backend API
    const [fires, setFires] = useState<any[]>([]);
    const [resources, setResources] = useState<any[]>([]);
    const [stations, setStations] = useState<any[]>([]);

    // Polling for data updates
    useEffect(() => {
        // Start polling for fires
        const stopFiresPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getFires(),
            5000,
            (data) => {
                setFires(data.fires || []);
                console.log('🔥 Fires updated:', data.fires?.length || 0);
            }
        );

        // Start polling for resources
        const stopResourcesPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getResources(),
            5000,
            (data) => {
                setResources(data.resources || []);
                console.log('🚒 Resources updated:', data.resources?.length || 0);
            }
        );

        // Start polling for stations
        const stopStationsPolling = fakeBackendAPI.startPolling(
            () => fakeBackendAPI.getStations(),
            5000,
            (data) => {
                setStations(data.stations || []);
                console.log('🏢 Stations updated:', data.stations?.length || 0);
            }
        );

        // Cleanup on unmount
        return () => {
            stopFiresPolling();
            stopResourcesPolling();
            stopStationsPolling();
        };
    }, []);

    // Layer State
    const [layers, setLayers] = React.useState({
        fires: true,
        heatmap: true,
        resources: true,
        stations: true
    });

    // Toggle Layer Visibility
    const toggleLayer = (layerKey: keyof typeof layers) => {
        setLayers(prev => {
            const newState = { ...prev, [layerKey]: !prev[layerKey] };
            const map = mapRef.current;
            if (!map) return newState;

            const isVisible = newState[layerKey];
            const visibility = isVisible ? 'visible' : 'none';

            if (layerKey === 'fires') {
                if (map.getLayer('fires-core')) map.setLayoutProperty('fires-core', 'visibility', visibility);
                if (map.getLayer('fires-glow')) map.setLayoutProperty('fires-glow', 'visibility', visibility);
            }

            if (layerKey === 'heatmap') {
                if (map.getLayer('risk-heat-layer')) map.setLayoutProperty('risk-heat-layer', 'visibility', visibility);
            }

            if (layerKey === 'resources') {
                resourceMarkersRef.current.forEach(marker => {
                    const el = marker.getElement();
                    el.style.display = isVisible ? 'flex' : 'none';
                });
            }

            if (layerKey === 'stations') {
                stationMarkersRef.current.forEach(marker => {
                    const el = marker.getElement();
                    // Stations are display:flex by default in our CSS class creation
                    el.style.display = isVisible ? 'flex' : 'none';
                });
            }

            return newState;
        });
    };

    const handleMapLoad = (map: mapboxgl.Map) => {
        mapRef.current = map;

        // Add Fire Station Markers (Global)
        stations.forEach(station => {
            const stationEl = document.createElement('div');
            stationEl.className = 'flex flex-col items-center justify-center group cursor-pointer';
            stationEl.innerHTML = `
                <div class="w-10 h-10 bg-blue-600 border-2 border-white rounded-lg flex items-center justify-center shadow-lg transform group-hover:scale-110 transition-transform">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 21h18"/><path d="M5 21V7l8-4 8 4v14"/><path d="M13 11V7"/><path d="M7 21V12h10v9"/></svg>
                </div>
                <div class="bg-zinc-900/90 text-zinc-100 text-[10px] px-2 py-0.5 rounded mt-1 font-bold border border-zinc-700 whitespace-nowrap">${station.name}</div>
            `;

            const marker = new mapboxgl.Marker({ element: stationEl })
                .setLngLat(station.location)
                .addTo(map);

            stationMarkersRef.current.push(marker);
        });

        // Initialize Sources with Empty Data (Data will be populated by useEffect)

        // Heatmap Source
        map.addSource('risk-heat', {
            type: 'geojson',
            data: { type: 'FeatureCollection', features: [] }
        });

        // Add Heatmap Layer FIRST (so it draws underneath)
        map.addLayer({
            id: 'risk-heat-layer',
            type: 'heatmap',
            source: 'risk-heat',
            maxzoom: 15,
            paint: {
                'heatmap-weight': [
                    'interpolate',
                    ['linear'],
                    ['get', 'intensity'],
                    0, 0,
                    1, 1
                ],
                'heatmap-intensity': [
                    'interpolate',
                    ['linear'],
                    ['zoom'],
                    0, 1, // Reduced intensity back to 1 for softer look
                    9, 3
                ],
                'heatmap-color': [
                    'interpolate',
                    ['linear'],
                    ['heatmap-density'],
                    0, 'rgba(0,0,0,0)', // Fully transparent start
                    0.2, 'rgba(234, 88, 12, 0.2)', // Very faint orange (fiery-600)
                    0.4, 'rgba(234, 88, 12, 0.5)', // Medium orange
                    0.6, 'rgba(220, 38, 38, 0.7)', // Red (red-600)
                    0.8, 'rgba(185, 28, 28, 0.8)', // Dark Red
                    1, 'rgba(127, 29, 29, 0.9)'  // Deepest Red
                ],
                'heatmap-radius': [
                    'interpolate',
                    ['linear'],
                    ['zoom'],
                    0, 25, // Reduced from 50 to 25 for better separation
                    9, 40  // Smoother local view
                ],
                'heatmap-opacity': [
                    'interpolate',
                    ['linear'],
                    ['zoom'],
                    7, 0.8, // Slightly more transparent overall
                    15, 0.4
                ],
            },
            layout: {
                visibility: 'visible' // Explicit initial visibility
            }
        });

        // Fires Source
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
                'circle-radius': 15,
                'circle-color': '#f43f5e',
                'circle-opacity': 0.3,
                'circle-blur': 0.5
            },
            layout: {
                visibility: 'visible'
            }
        });

        // Fire Core
        map.addLayer({
            id: 'fires-core',
            type: 'circle',
            source: 'fires-source',
            paint: {
                'circle-radius': 5, // Slightly larger again (4->5)
                'circle-color': '#ff4d4d', // Brighter Red for visibility against heatmap
                'circle-stroke-width': 1, // 1px Stroke to define edges
                'circle-stroke-color': '#fff'
            },
            layout: {
                visibility: 'visible'
            }
        });

        // Blinking Animation Loop
        let animationFrameId: number;
        const animate = () => {
            if (!map || !map.getStyle()) return; // Safety check

            const time = Date.now() / 1000;
            const opacity = (Math.sin(time * 3) + 1) / 2 * 0.6 + 0.2;
            const coreOpacity = (Math.sin(time * 3) + 1) / 2 * 0.2 + 0.8;

            if (map.getLayer('fires-glow')) {
                map.setPaintProperty('fires-glow', 'circle-opacity', opacity);
                map.setPaintProperty('fires-glow', 'circle-radius', 15 + Math.sin(time * 3) * 5);
            }
            if (map.getLayer('fires-core')) {
                map.setPaintProperty('fires-core', 'circle-opacity', coreOpacity);
            }

            animationFrameId = requestAnimationFrame(animate);
        };

        animate();

        // Cleanup function for this specific map load closure
        // Note: The main component cleanup handles map removal, but we should stop animation
        return () => cancelAnimationFrame(animationFrameId); // Start animation

        // Add Click Interaction for Fires
        map.on('click', 'fires-core', (e) => {
            if (!e.features || e.features.length === 0) return;

            const feature = e.features[0];
            const coordinates = (feature.geometry as any).coordinates.slice();
            const props = feature.properties;

            while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
                coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
            }

            new mapboxgl.Popup({ offset: 15, className: 'text-zinc-900' })
                .setLngLat(coordinates)
                .setHTML(`
                    <div class="p-2 min-w-[200px]">
                        <div class="flex items-center gap-2 mb-2">
                             <span class="text-rose-600 font-bold bg-rose-100 px-2 py-0.5 rounded text-xs uppercase">${props?.severity}</span>
                             <span class="text-zinc-500 text-xs">${props?.id}</span>
                        </div>
                        <div class="text-sm font-medium text-zinc-900">Confidence: ${(props?.confidence * 100).toFixed(0)}%</div>
                         <div class="text-xs text-zinc-500 mt-1">Source: ${props?.source}</div>
                    </div>
                `)
                .addTo(map);
        });

        // Cursor pointer events
        map.on('mouseenter', 'fires-core', () => {
            map.getCanvas().style.cursor = 'pointer';
        });
        map.on('mouseleave', 'fires-core', () => {
            map.getCanvas().style.cursor = '';
        });

        // Add Markers for Resources
        resources.forEach(unit => {
            const el = document.createElement('div');
            el.className = 'flex flex-col items-center justify-center group cursor-pointer';
            // Fire Truck Icon
            el.innerHTML = `
                <div class="p-1.5 bg-zinc-900/90 border border-fiery-500 rounded-md shadow-lg transform group-hover:scale-110 transition-transform">
                     <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#f97316" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect width="16" height="12" x="2" y="6" rx="2" />
                        <circle cx="7" cy="18" r="2" />
                        <circle cx="17" cy="18" r="2" />
                        <path d="M18 6V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v2" />
                    </svg>
                </div>
                <div class="opacity-0 group-hover:opacity-100 transition-opacity absolute top-full mt-1 bg-zinc-900 text-fiery-500 text-[10px] font-bold px-1.5 py-0.5 rounded border border-fiery-500 whitespace-nowrap z-10 pointer-events-none">
                    ${unit.name}
                </div>
            `;

            const marker = new mapboxgl.Marker({ element: el })
                .setLngLat(unit.location)
                .addTo(map);

            resourceMarkersRef.current.push(marker);
        });

        // Ensure Heatmap is below everything else
        if (map.getLayer('risk-heat-layer')) map.moveLayer('risk-heat-layer');
        // Ensure Fire Glow is above Heatmap
        if (map.getLayer('fires-glow')) map.moveLayer('fires-glow');
        // Ensure Fire Core is ON TOP of everything
        if (map.getLayer('fires-core')) map.moveLayer('fires-core');
    };

    // React to data updates and update map sources
    useEffect(() => {
        const map = mapRef.current;
        if (!map) return;

        // Update Fires Source
        const fireSource = map.getSource('fires-source') as mapboxgl.GeoJSONSource;
        if (fireSource) {
            const fireFeatures = fires.map(fire => ({
                type: 'Feature',
                geometry: { type: 'Point', coordinates: fire.location },
                properties: {
                    id: fire.id,
                    severity: fire.severity,
                    confidence: fire.confidence,
                    source: fire.source
                }
            }));
            fireSource.setData({
                type: 'FeatureCollection',
                features: fireFeatures
            } as any);

            // Auto-zoom to first fire if available
            if (fires.length > 0) {
                const bounds = new mapboxgl.LngLatBounds();
                fires.forEach(fire => bounds.extend(fire.location));

                map.fitBounds(bounds, {
                    padding: 100,
                    maxZoom: 12,
                    duration: 2000
                });
            }
        }

        // Update Heatmap Source
        const heatmapSource = map.getSource('risk-heat') as mapboxgl.GeoJSONSource;
        if (heatmapSource) {
            const heatmapData: GeoJSON.FeatureCollection = {
                type: 'FeatureCollection',
                features: fires.flatMap(fire => {
                    return Array.from({ length: 25 }).map(() => ({
                        type: 'Feature',
                        properties: { intensity: Math.random() },
                        geometry: {
                            type: 'Point',
                            coordinates: [
                                fire.location[0] + (Math.random() - 0.5) * 0.5,
                                fire.location[1] + (Math.random() - 0.5) * 0.5
                            ]
                        }
                    }));
                }) as any
            };
            heatmapSource.setData(heatmapData);
        }

    }, [fires]); // Re-run when fires data changes

    const [searchParams] = useSearchParams();
    const activeResourceId = searchParams.get('resourceId');

    // Effect for Logistics Visualization
    useEffect(() => {
        const map = mapRef.current;
        if (!map || !activeResourceId || resources.length === 0 || fires.length === 0) return;

        const resource = resources.find(r => r.id === activeResourceId);
        if (!resource) return;

        // Find nearest fire
        let nearestFire = fires[0];
        let minDist = Infinity;

        fires.forEach(fire => {
            const dist = Math.sqrt(
                Math.pow(fire.location[0] - resource.location[0], 2) +
                Math.pow(fire.location[1] - resource.location[1], 2)
            );
            if (dist < minDist) {
                minDist = dist;
                nearestFire = fire;
            }
        });

        // Draw Route Line
        const routeGeoJSON = {
            type: 'Feature',
            geometry: {
                type: 'LineString',
                coordinates: [resource.location, nearestFire.location]
            }
        };

        if (map.getSource('logistics-route')) {
            (map.getSource('logistics-route') as mapboxgl.GeoJSONSource).setData(routeGeoJSON as any);
        } else {
            map.addSource('logistics-route', {
                type: 'geojson',
                data: routeGeoJSON as any
            });

            map.addLayer({
                id: 'logistics-route-line',
                type: 'line',
                source: 'logistics-route',
                layout: {
                    'line-join': 'round',
                    'line-cap': 'round'
                },
                paint: {
                    'line-color': '#f59e0b', // Amber-500
                    'line-width': 4,
                    'line-dasharray': [2, 1]
                }
            });
        }

        // Ensure Route is visible (move to top)
        if (map.getLayer('logistics-route-line')) {
            map.moveLayer('logistics-route-line');
        }
        // Keep markers on top of line
        if (map.getLayer('fires-core')) map.moveLayer('fires-core');

        // Fly to fit bounds
        const bounds = new mapboxgl.LngLatBounds();
        bounds.extend(resource.location);
        bounds.extend(nearestFire.location);
        map.fitBounds(bounds, { padding: 150, maxZoom: 13, duration: 2000 });

    }, [activeResourceId, resources, fires]);

    // Check if we are in logistics mode but waiting for data
    const isLoadingLogistics = activeResourceId && (resources.length === 0 || fires.length === 0);

    return (
        <div className="h-full w-full relative group">
            <MapComponent className="h-full w-full" onLoad={handleMapLoad}>
            </MapComponent>

            {/* Loading Indicator for Logistics Mode */}
            {isLoadingLogistics && (
                <div className="absolute inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
                    <div className="bg-zinc-900 border border-zinc-700 p-4 rounded-xl shadow-2xl flex items-center gap-3">
                        <div className="w-5 h-5 border-2 border-amber-500 border-t-transparent rounded-full animate-spin"></div>
                        <span className="text-zinc-200 font-medium">Locating Resource & Calculating Route...</span>
                    </div>
                </div>
            )}

            {/* Map Floating Controls */}
            <div className="absolute top-4 left-4 flex flex-col gap-2 pointer-events-none">
                <div className="bg-zinc-900/90 backdrop-blur-md border border-zinc-700 p-3 rounded-lg shadow-xl pointer-events-auto">
                    <h2 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider mb-2">Active Layers</h2>
                    <div className="space-y-1">
                        <label className="flex items-center gap-2 text-sm text-zinc-400 cursor-pointer hover:text-fiery-400">
                            <input
                                type="checkbox"
                                checked={layers.fires}
                                onChange={() => toggleLayer('fires')}
                                className="accent-fiery-500 rounded bg-zinc-800 border-zinc-600"
                            />
                            <span>Fire Incidents</span>
                        </label>
                        <label className="flex items-center gap-2 text-sm text-zinc-400 cursor-pointer hover:text-fiery-400">
                            <input
                                type="checkbox"
                                checked={layers.heatmap}
                                onChange={() => toggleLayer('heatmap')}
                                className="accent-fiery-500 rounded bg-zinc-800 border-zinc-600"
                            />
                            <span>Risk Heatmap</span>
                        </label>
                        <label className="flex items-center gap-2 text-sm text-zinc-400 cursor-pointer hover:text-fiery-400">
                            <input
                                type="checkbox"
                                checked={layers.resources}
                                onChange={() => toggleLayer('resources')}
                                className="accent-fiery-500 rounded bg-zinc-800 border-zinc-600"
                            />
                            <span>Fire Trucks</span>
                        </label>
                        <label className="flex items-center gap-2 text-sm text-zinc-400 cursor-pointer hover:text-fiery-400">
                            <input
                                type="checkbox"
                                checked={layers.stations}
                                onChange={() => toggleLayer('stations')}
                                className="accent-fiery-500 rounded bg-zinc-800 border-zinc-600"
                            />
                            <span>Stations</span>
                        </label>
                    </div>
                </div>
            </div>

            {/* Logistics Active Overlay */}
            {activeResourceId && (
                <div className="absolute top-4 right-4 pointer-events-none">
                    <div className="bg-amber-500/90 backdrop-blur-md text-zinc-900 px-4 py-2 rounded-lg shadow-xl animate-pulse flex items-center gap-3">
                        <span className="font-bold font-mono">LOGISTICS AGENT ACTIVE</span>
                        <div className="w-2 h-2 bg-zinc-900 rounded-full animate-ping"></div>
                    </div>
                    <div className="bg-zinc-900/90 backdrop-blur-md border border-amber-500/50 p-3 rounded-lg shadow-xl mt-2 text-zinc-300 text-xs">
                        <div>Optimizing Route...</div>
                        <div className="font-mono text-amber-400 text-lg font-bold">ETA: 12 MINS</div>
                        <div className="text-zinc-500">Distance: 15.4 km</div>
                    </div>
                </div>
            )}

            <div className="absolute bottom-6 left-1/2 -translate-x-1/2 flex gap-4 pointer-events-none">
                <div className="bg-zinc-900/90 backdrop-blur-md border border-zinc-700 px-4 py-2 rounded-full shadow-xl flex items-center gap-3 pointer-events-auto">
                    <div className="flex items-center gap-2">
                        {fires.length > 0 ? (
                            <>
                                <span className="w-2 h-2 rounded-full bg-rose-500 animate-pulse"></span>
                                <span className="text-xs font-bold text-zinc-200">{fires.length} ACTIVE FIRE{fires.length !== 1 && 'S'} DETECTED</span>
                            </>
                        ) : (
                            <>
                                <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
                                <span className="text-xs font-bold text-zinc-200">SYSTEM NOMINAL - MONITORING</span>
                            </>
                        )}
                    </div>
                    <div className="h-4 w-px bg-zinc-700"></div>
                    <div className="flex items-center gap-2">
                        <span className="w-1.5 h-1.5 rounded-full bg-fiery-500"></span>
                        <span className="text-xs text-zinc-400 font-mono">AGENTS ACTIVE</span>
                    </div>
                </div>
            </div>
        </div>
    );
};
