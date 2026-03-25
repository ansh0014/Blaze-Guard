// API Service for BlazeGuard Fake Backend
const API_BASE_URL = 'http://localhost:3001/api';

class FakeBackendAPI {
    // Generic fetch wrapper
    async fetch(endpoint: string): Promise<any> {
        try {
            const response = await fetch(`${API_BASE_URL}${endpoint}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error(`API Error (${endpoint}):`, error);
            throw error;
        }
    }

    // Fire endpoints
    async getFires(): Promise<any> {
        return this.fetch('/fires');
    }

    async getFire(id: string): Promise<any> {
        return this.fetch(`/fires/${id}`);
    }

    // Alert endpoints
    async getAlerts(): Promise<any> {
        return this.fetch('/alerts');
    }

    async getActiveAlerts(): Promise<any> {
        return this.fetch('/alerts/active');
    }

    // Resource endpoints
    async getResources(): Promise<any> {
        return this.fetch('/resources');
    }

    async getResource(id: string): Promise<any> {
        return this.fetch(`/resources/${id}`);
    }

    // Station endpoints
    async getStations(): Promise<any> {
        return this.fetch('/stations');
    }

    // Analytics endpoints
    async getRiskTrend(): Promise<any> {
        return this.fetch('/analytics/risk-trend');
    }

    async getEnvironmentalData(): Promise<any> {
        return this.fetch('/analytics/environmental');
    }

    async getPredictions(): Promise<any> {
        return this.fetch('/analytics/predictions');
    }

    // Evacuation endpoints
    async getEvacuationZones(): Promise<any> {
        return this.fetch('/evacuation/zones');
    }

    async getSafeZones(): Promise<any> {
        return this.fetch('/evacuation/safe-zones');
    }

    // Agent endpoints
    async getAgents(): Promise<any> {
        return this.fetch('/agents');
    }

    // Distance endpoints
    async getDistances(): Promise<any> {
        return this.fetch('/distances');
    }

    // Feedback endpoint
    async submitFeedback(feedbackData: any): Promise<any> {
        try {
            const response = await fetch(`${API_BASE_URL}/feedback`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(feedbackData)
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Feedback submission error:', error);
            throw error;
        }
    }

    // Polling helper - call a function repeatedly
    startPolling(fetchFunction: () => Promise<any>, interval: number = 5000, callback: (data: any) => void): () => void {
        // Initial fetch
        fetchFunction().then(callback).catch(console.error);

        // Set up interval
        const intervalId = setInterval(() => {
            fetchFunction().then(callback).catch(console.error);
        }, interval);

        // Return function to stop polling
        return () => clearInterval(intervalId);
    }
}

// Export singleton instance
export const fakeBackendAPI = new FakeBackendAPI();

// Example usage:
/*
import { fakeBackendAPI } from './services/fakeBackendAPI';

// Simple fetch
const fires = await fakeBackendAPI.getFires();

// Polling example
const stopPolling = fakeBackendAPI.startPolling(
    () => fakeBackendAPI.getFires(),
    5000, // Poll every 5 seconds
    (data) => {
        console.log('New fire data:', data);
        // Update your UI here
    }
);

// Stop polling when component unmounts
stopPolling();
*/
