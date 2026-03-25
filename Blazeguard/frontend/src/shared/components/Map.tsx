import React, { useEffect, useRef, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

// User provided token
const MAPBOX_TOKEN = import.meta.env.VITE_MAPBOX_TOKEN;

mapboxgl.accessToken = MAPBOX_TOKEN;

interface MapProps {
    initialCenter?: [number, number];
    initialZoom?: number;
    className?: string;
    children?: React.ReactNode;
    onLoad?: (map: mapboxgl.Map) => void;
}

export const MapComponent: React.FC<MapProps> = ({
    initialCenter = [-40, 25], // Atlantic Ocean center for global view
    initialZoom = 1.5,
    className = "h-full w-full",
    children,
    onLoad
}) => {
    const mapContainerRef = useRef<HTMLDivElement>(null);
    const mapRef = useRef<mapboxgl.Map | null>(null);
    const [isLoaded, setIsLoaded] = useState(false);

    useEffect(() => {
        if (mapRef.current || !mapContainerRef.current) return;

        console.log("Initializing Mapbox map...");
        console.log("Container dimensions:", mapContainerRef.current.clientWidth, mapContainerRef.current.clientHeight);

        try {
            mapRef.current = new mapboxgl.Map({
                container: mapContainerRef.current,
                style: 'mapbox://styles/mapbox/satellite-streets-v12', // Try satellite to verify tile loading
                center: initialCenter,
                zoom: initialZoom,
                attributionControl: false,
                failIfMajorPerformanceCaveat: false
            });

            mapRef.current.addControl(new mapboxgl.NavigationControl(), 'bottom-right');
            mapRef.current.addControl(new mapboxgl.ScaleControl(), 'bottom-left');

            // Debug listeners
            mapRef.current.on('style.load', () => console.log('Mapbox: Style loaded'));
            mapRef.current.on('styledata', () => console.log('Mapbox: Style data loading...'));
            mapRef.current.on('tiledata', () => { /* too noisy to log */ });
            mapRef.current.on('error', (e) => console.error('Mapbox Error:', e));

            mapRef.current.on('load', () => {
                console.log("Map loaded successfully");
                setIsLoaded(true);
                mapRef.current?.resize();
                if (onLoad && mapRef.current) {
                    onLoad(mapRef.current);
                }
            });

        } catch (err) {
            console.error("Error initializing map:", err);
        }

        return () => {
            mapRef.current?.remove();
            mapRef.current = null;
        };
    }, []);

    // Handle resize observer to fix blank canvas issues
    useEffect(() => {
        if (!mapRef.current || !mapContainerRef.current) return;

        const resizeObserver = new ResizeObserver(() => {
            if (mapRef.current) {
                console.log("Resizing map...");
                mapRef.current.resize();
            }
        });

        resizeObserver.observe(mapContainerRef.current);

        return () => resizeObserver.disconnect();
    }, [isLoaded]);

    return (
        <div className={`relative ${className}`}>
            <div
                ref={mapContainerRef}
                className="absolute inset-0 rounded-lg overflow-hidden border border-zinc-800 shadow-inner"
                style={{ width: '100%', height: '100%' }}
            />

            {/* Loading State */}
            {!isLoaded && (
                <div className="absolute inset-0 flex items-center justify-center bg-zinc-900 z-10">
                    <div className="flex flex-col items-center gap-2">
                        <div className="w-8 h-8 border-4 border-fiery-500 border-t-transparent rounded-full animate-spin"></div>
                        <span className="text-zinc-500 text-sm">Loading Map...</span>
                    </div>
                </div>
            )}

            {/* Render children only when map is loaded */}
            {isLoaded && children}
        </div>
    );
};
