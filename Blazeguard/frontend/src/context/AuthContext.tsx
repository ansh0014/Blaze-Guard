import React, { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { authService, type User } from '../services/AuthService';
import { firebaseAuth } from '../services/firebase';
import { type User as FirebaseUser } from 'firebase/auth';

interface AuthContextType {
    user: User | null;
    firebaseUser: FirebaseUser | null;
    loading: boolean;
    isAuthenticated: boolean;
    loginWithGoogle: () => Promise<void>;
    loginWithEmail: (email: string, password: string) => Promise<void>;
    signUpWithEmail: (email: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [firebaseUser, setFirebaseUser] = useState<FirebaseUser | null>(null);
    const [loading, setLoading] = useState(true);

    const fetchUser = async () => {
        try {
            const userData = await authService.getProfile();
            setUser(userData);
        } catch (error) {
            setUser(null);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const unsubscribe = firebaseAuth.onAuthStateChanged(async (fbUser) => {
            setFirebaseUser(fbUser);

            if (fbUser) {
                const idToken = await fbUser.getIdToken();
                try {
                    const backendUser = await authService.loginWithToken(idToken);
                    setUser(backendUser);
                } catch (error) {
                    setUser(null);
                }
            } else {
                setUser(null);
            }

            setLoading(false);
        });

        return () => unsubscribe();
    }, []);

    const loginWithGoogle = async () => {
        try {
            setLoading(true);
            const fbUser = await firebaseAuth.signInWithGoogle();
            const idToken = await fbUser.getIdToken();
            const backendUser = await authService.loginWithToken(idToken);
            setUser(backendUser);
            setFirebaseUser(fbUser);
        } catch (error) {
            throw error;
        } finally {
            setLoading(false);
        }
    };

    const loginWithEmail = async (email: string, password: string) => {
        try {
            setLoading(true);
            const fbUser = await firebaseAuth.signInWithEmail(email, password);
            const idToken = await fbUser.getIdToken();
            const backendUser = await authService.loginWithToken(idToken);
            setUser(backendUser);
            setFirebaseUser(fbUser);
        } catch (error) {
            throw error;
        } finally {
            setLoading(false);
        }
    };

    const signUpWithEmail = async (email: string, password: string) => {
        try {
            setLoading(true);
            const fbUser = await firebaseAuth.signUpWithEmail(email, password);
            const idToken = await fbUser.getIdToken();
            const backendUser = await authService.loginWithToken(idToken);
            setUser(backendUser);
            setFirebaseUser(fbUser);
        } catch (error) {
            throw error;
        } finally {
            setLoading(false);
        }
    };

    const logout = async () => {
        console.log('🔴 Logout started');
        try {
            console.log('🔴 Signing out from Firebase...');
            await firebaseAuth.signOut();
            console.log('🔴 Firebase sign out complete');

            console.log('🔴 Calling backend logout...');
            await authService.logout();
            console.log('🔴 Backend logout complete');

            setUser(null);
            setFirebaseUser(null);
            console.log('🔴 State cleared, redirecting...');
            window.location.href = '/';
        } catch (error) {
            console.error('🔴 Logout error:', error);
        }
    };

    const refreshUser = async () => {
        setLoading(true);
        await fetchUser();
    };

    const value: AuthContextType = {
        user,
        firebaseUser,
        loading,
        isAuthenticated: user !== null,
        loginWithGoogle,
        loginWithEmail,
        signUpWithEmail,
        logout,
        refreshUser,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

export default AuthContext;
