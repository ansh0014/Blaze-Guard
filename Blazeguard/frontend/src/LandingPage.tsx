import { Link } from 'react-router-dom';
import { Shield, MapPin, Activity, ArrowRight } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAuth } from './context/AuthContext';

export const LandingPage = () => {
    const { isAuthenticated, user, logout } = useAuth();

    const handleLogout = async () => {
        console.log('🟢 Logout button clicked');
        await logout();
    };

    return (
        <div className="min-h-screen bg-background text-zinc-100 flex flex-col relative overflow-hidden font-sans selection:bg-fiery-500/30 selection:text-fiery-500">

            {/* Ambient Background Effects */}
            <div className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] bg-fiery-600/20 rounded-full blur-[120px] pointer-events-none mix-blend-screen" />
            <div className="absolute bottom-[-10%] right-[-5%] w-[600px] h-[600px] bg-fiery-500/10 rounded-full blur-[150px] pointer-events-none mix-blend-screen" />

            {/* Grid Pattern Overlay */}
            <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:40px_40px] pointer-events-none" />

            {/* Navbar */}
            <nav className="relative z-50 px-6 py-6 flex items-center justify-between max-w-7xl mx-auto w-full">
                <div className="flex items-center gap-2 font-black text-xl tracking-tighter">
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-fiery-500 to-fiery-400 flex items-center justify-center shadow-lg shadow-fiery-500/20">
                        <Activity size={18} className="text-white" />
                    </div>
                    <span>BLAZE<span className="text-fiery-500">GUARD</span></span>
                </div>
                <div className="flex items-center gap-6 text-sm font-medium text-zinc-400">
                    <span className="hidden md:block hover:text-white transition-colors cursor-pointer">System Status</span>
                    <span className="hidden md:block hover:text-white transition-colors cursor-pointer">Global Ops</span>
                    <Link to="/citizen" className="px-5 py-2 rounded-full border border-zinc-700 hover:border-fiery-500 text-zinc-300 hover:text-fiery-500 transition-all hover:bg-fiery-500/5">
                        Client Login
                    </Link>
                    {isAuthenticated ? (
                        <div className="flex items-center gap-3">
                            <span className="text-xs text-zinc-500">Welcome, {user?.name}</span>
                            <button
                                onClick={handleLogout}
                                className="px-5 py-2 rounded-full bg-fiery-500 hover:bg-fiery-600 text-white transition-all"
                            >
                                Logout
                            </button>
                        </div>
                    ) : (
                        <Link to="/login" className="px-5 py-2 rounded-full bg-fiery-500 hover:bg-fiery-600 text-white transition-all">
                            Authority Login
                        </Link>
                    )}
                </div>
            </nav>

            {/* Hero Section */}
            <main className="flex-1 flex flex-col items-center justify-center relative z-10 px-4">
                <div className="text-center max-w-4xl mx-auto space-y-8">

                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6 }}
                        className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-zinc-900/50 border border-zinc-800 backdrop-blur-md text-xs font-semibold text-fiery-400 uppercase tracking-widest mb-4"
                    >
                        <span className="w-1.5 h-1.5 rounded-full bg-fiery-500 animate-pulse"></span>
                        System Operational
                    </motion.div>

                    <motion.h1
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6, delay: 0.1 }}
                        className="text-6xl md:text-8xl font-black tracking-tighter leading-[0.9] text-transparent bg-clip-text bg-gradient-to-b from-white to-zinc-500"
                    >
                        INTELLIGENT <br />
                        <span className="text-stroke-zinc">FIRE DEFENSE</span>
                    </motion.h1>

                    <motion.p
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6, delay: 0.2 }}
                        className="text-lg md:text-xl text-zinc-400 max-w-2xl mx-auto font-light leading-relaxed"
                    >
                        Deploy advanced satellite analytics and real-time response heuristics to mitigate wildfire risks before they escalate.
                    </motion.p>

                    <motion.div
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.6, delay: 0.3 }}
                        className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full max-w-3xl mx-auto mt-12"
                    >
                        {/* Citizen Card */}
                        <Link to="/citizen" className="group relative overflow-hidden rounded-3xl bg-zinc-900/50 border border-zinc-800 hover:border-fiery-500/50 transition-all duration-500 p-8 text-left hover:shadow-[0_0_40px_-10px_rgba(255,77,0,0.3)]">
                            <div className="absolute inset-0 bg-gradient-to-br from-fiery-500/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                            <div className="relative z-10 flex flex-col h-full justify-between gap-8">
                                <div className="w-12 h-12 rounded-2xl bg-zinc-800 flex items-center justify-center group-hover:bg-fiery-500 text-fiery-500 group-hover:text-white transition-colors duration-500">
                                    <MapPin size={24} />
                                </div>
                                <div>
                                    <h3 className="text-2xl font-bold text-white mb-2 group-hover:translate-x-1 transition-transform">Citizen Portal</h3>
                                    <p className="text-zinc-500 text-sm">Real-time safety status, evacuation routing, and community alerts.</p>
                                </div>
                                <div className="flex items-center gap-2 text-sm font-semibold text-fiery-500 group-hover:text-fiery-400">
                                    Access Portal <ArrowRight size={16} className="group-hover:translate-x-1 transition-transform" />
                                </div>
                            </div>
                        </Link>

                        {/* Authority Card */}
                        <Link to="/dashboard" className="group relative overflow-hidden rounded-3xl bg-zinc-900/50 border border-zinc-800 hover:border-fiery-500/50 transition-all duration-500 p-8 text-left hover:shadow-[0_0_40px_-10px_rgba(255,77,0,0.3)]">
                            <div className="absolute inset-0 bg-gradient-to-br from-fiery-500/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                            <div className="relative z-10 flex flex-col h-full justify-between gap-8">
                                <div className="w-12 h-12 rounded-2xl bg-zinc-800 flex items-center justify-center group-hover:bg-fiery-500 text-fiery-500 group-hover:text-white transition-colors duration-500">
                                    <Shield size={24} />
                                </div>
                                <div>
                                    <h3 className="text-2xl font-bold text-white mb-2 group-hover:translate-x-1 transition-transform">Command Center</h3>
                                    <p className="text-zinc-500 text-sm">Strategic resource deployment, live satellite feed, and unit telemetry.</p>
                                </div>
                                <div className="flex items-center gap-2 text-sm font-semibold text-fiery-500 group-hover:text-fiery-400">
                                    Initialize Dashboard <ArrowRight size={16} className="group-hover:translate-x-1 transition-transform" />
                                </div>
                            </div>
                        </Link>
                    </motion.div>
                </div>
            </main>

            {/* Footer Status */}
            <footer className="relative z-10 py-8 text-center text-xs font-mono text-zinc-600 border-t border-zinc-900/50 mt-12 bg-background/50 backdrop-blur-sm">
                <p>SECURE CONNECTION • V4.2.0 • SERVER: ASIA-SOUTH-1</p>
            </footer>
        </div>
    );
};
