const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:2000';

export interface User {
    id: number;
    firebase_uid: string;
    email: string;
    name: string;
    picture: string;
    role: string;
    created_at: string;
    updated_at: string;
}

class AuthService {
    async loginWithToken(idToken: string): Promise<User> {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ idToken }),
        });

        if (!response.ok) {
            throw new Error('Login failed');
        }

        return await response.json();
    }

    async logout(): Promise<void> {
        console.log('🔵 AuthService.logout called');
        try {
            const response = await fetch(`${API_BASE_URL}/auth/logout`, {
                method: 'POST',
                credentials: 'include',
            });
            console.log('🔵 Logout response:', response.status);
        } catch (error) {
            console.error('🔵 Logout fetch error:', error);
            throw error;
        }
    }

    async getProfile(): Promise<User | null> {
        try {
            const response = await fetch(`${API_BASE_URL}/api/profile`, {
                credentials: 'include',
            });

            if (response.status === 401) {
                return null;
            }

            if (!response.ok) {
                throw new Error('Failed to fetch profile');
            }

            return await response.json();
        } catch (error) {
            return null;
        }
    }

    async isAuthenticated(): Promise<boolean> {
        const user = await this.getProfile();
        return user !== null;
    }

    async healthCheck(): Promise<boolean> {
        try {
            const response = await fetch(`${API_BASE_URL}/health`);
            return response.ok;
        } catch (error) {
            return false;
        }
    }
}

export const authService = new AuthService();
export default authService;
