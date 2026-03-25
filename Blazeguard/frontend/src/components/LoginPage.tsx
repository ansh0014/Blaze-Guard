import React from 'react';
import { useAuth } from '../context/AuthContext';
import { motion } from 'framer-motion';
import { Shield, Activity } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export const LoginPage: React.FC = () => {
    const { loginWithGoogle, loading } = useAuth();
    const navigate = useNavigate();

    const handleGoogleSignIn = async () => {
        try {
            await loginWithGoogle();
            navigate('/dashboard');
        } catch (error) {
            // Silent error handling
        }
    };

    return (
        <div className="min-h-screen bg-background text-zinc-100 flex flex-col relative overflow-hidden font-sans">
            {/* Ambient Background Effects */}
            <div className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] bg-fiery-600/20 rounded-full blur-[120px] pointer-events-none mix-blend-screen" />
            <div className="absolute bottom-[-10%] right-[-5%] w-[600px] h-[600px] bg-fiery-500/10 rounded-full blur-[150px] pointer-events-none mix-blend-screen" />

            {/* Grid Pattern Overlay */}
            <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:40px_40px] pointer-events-none" />

            {/* Navbar */}
            <nav className="relative z-50 px-6 py-6 flex items-center justify-between max-w-7xl mx-auto w-full">
                <div className="flex items-center gap-2 font-black text-xl tracking-tighter cursor-pointer" onClick={() => navigate('/')}>
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-fiery-500 to-fiery-400 flex items-center justify-center shadow-lg shadow-fiery-500/20">
                        <Activity size={18} className="text-white" />
                    </div>
                    <span>BLAZE<span className="text-fiery-500">GUARD</span></span>
                </div>
            </nav>

            {/* Login Content */}
            <main className="flex-1 flex flex-col items-center justify-center relative z-10 px-4">
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6 }}
                    className="w-full max-w-md"
                >
                    <div className="bg-zinc-900/50 border border-zinc-800 rounded-3xl p-8 backdrop-blur-md shadow-[0_0_40px_-10px_rgba(255,77,0,0.2)]">
                        {/* Icon */}
                        <div className="flex justify-center mb-6">
                            <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-fiery-500 to-fiery-600 flex items-center justify-center shadow-lg shadow-fiery-500/30">
                                <Shield size={40} className="text-white" />
                            </div>
                        </div>

                        {/* Title */}
                        <h1 className="text-3xl font-black text-center text-white mb-2 tracking-tight">
                            Authority Access
                        </h1>
                        <p className="text-center text-zinc-400 mb-8 text-sm">
                            Secure authentication for fire authority personnel
                        </p>

                        {/* Status Badge */}
                        <div className="flex items-center justify-center gap-2 px-3 py-1.5 rounded-full bg-zinc-800/50 border border-zinc-700 text-xs font-semibold text-fiery-400 uppercase tracking-widest mb-8 w-fit mx-auto">
                            <span className="w-1.5 h-1.5 rounded-full bg-fiery-500 animate-pulse"></span>
                            System Operational
                        </div>

                        {/* Google Sign-In Button */}
                        <button
                            onClick={handleGoogleSignIn}
                            disabled={loading}
                            className="w-full group relative overflow-hidden rounded-xl bg-white hover:bg-zinc-100 text-zinc-900 font-bold py-4 px-6 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg"
                        >
                            <div className="relative flex items-center justify-center gap-3">
                                <svg width="24" height="24" viewBox="0 0 24 24" className="flex-shrink-0">
                                    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                                    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                                    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
                                    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
                                </svg>
                                <span className="text-lg">
                                    {loading ? 'Authenticating...' : 'Sign in with Google'}
                                </span>
                            </div>
                        </button>

                        {/* Footer Note */}
                        <p className="text-center text-zinc-500 text-xs mt-6 font-mono">
                            AUTHORIZED PERSONNEL ONLY
                        </p>
                    </div>

                    {/* Additional Info */}
                    <p className="text-center text-zinc-600 text-xs mt-6">
                        Secure connection • End-to-end encrypted
                    </p>
                </motion.div>
            </main>
        </div>
    );
};

export default LoginPage;
