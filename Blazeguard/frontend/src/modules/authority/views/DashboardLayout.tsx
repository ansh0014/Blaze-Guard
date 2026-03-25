import React from 'react';
import { SidebarIcon, Bell, Settings, LogOut, Flame, ShieldAlert, Truck, BarChart3, Map as MapIcon, Menu, FileText } from 'lucide-react';
import { NavLink, Outlet } from 'react-router-dom';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { useAuth } from '../../../context/AuthContext';

function cn(...inputs: (string | undefined | null | false)[]) {
    return twMerge(clsx(inputs));
}

const NavItem = ({ to, icon: Icon, label }: { to: string, icon: any, label: string }) => (
    <NavLink
        to={to}
        className={({ isActive }) => cn(
            "flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium",
            isActive
                ? "bg-fiery-500/10 text-fiery-500 border border-fiery-500/20"
                : "text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800"
        )}
    >
        <Icon size={18} />
        <span>{label}</span>
    </NavLink>
);

export const DashboardLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = React.useState(true);
    const { logout } = useAuth();

    return (
        <div className="flex h-screen w-full bg-zinc-950 text-zinc-100 overflow-hidden">
            {/* Sidebar */}
            <aside
                className={cn(
                    "flex flex-col border-r border-zinc-800 bg-zinc-900/50 transition-all duration-300 ease-in-out",
                    sidebarOpen ? "w-64" : "w-16"
                )}
            >
                <div className="h-16 flex items-center px-4 border-b border-zinc-800 justify-between">
                    {sidebarOpen ? (
                        <div className="flex items-center gap-2 text-fiery-500 font-bold tracking-wider">
                            <Flame className="fill-fiery-500" size={20} />
                            <span>BLAZE<span className="text-zinc-100">GUARD</span></span>
                        </div>
                    ) : (
                        <Flame className="mx-auto text-fiery-500 fill-fiery-500" size={24} />
                    )}
                    <button onClick={() => setSidebarOpen(!sidebarOpen)} className="text-zinc-500 hover:text-zinc-100 p-1">
                        {sidebarOpen ? <SidebarIcon size={16} /> : <Menu size={16} />}
                    </button>
                </div>

                <div className="flex-1 py-6 px-2 space-y-1 overflow-y-auto">
                    <NavItem to="/dashboard/map" icon={MapIcon} label={sidebarOpen ? "Live Map" : ""} />
                    <NavItem to="/dashboard/alerts" icon={ShieldAlert} label={sidebarOpen ? "Alerts & Incidents" : ""} />
                    <NavItem to="/dashboard/resources" icon={Truck} label={sidebarOpen ? "Resources" : ""} />
                    <NavItem to="/dashboard/analytics" icon={BarChart3} label={sidebarOpen ? "Risk Analytics" : ""} />
                    <NavItem to="/dashboard/report" icon={FileText} label={sidebarOpen ? "AI Report" : ""} />
                </div>

                <div className="p-2 border-t border-zinc-800">
                    <NavItem to="/settings" icon={Settings} label={sidebarOpen ? "Settings" : ""} />
                    <div className="mt-2 pt-2 border-t border-zinc-800/50">
                        <button
                            onClick={logout}
                            className={cn(
                                "flex items-center gap-3 px-3 py-2 w-full text-left rounded-md transition-colors text-sm font-medium text-rose-400 hover:bg-rose-500/10",
                                !sidebarOpen && "justify-center"
                            )}
                        >
                            <LogOut size={18} />
                            {sidebarOpen && <span>Logout</span>}
                        </button>
                    </div>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 flex flex-col min-w-0">
                <header className="h-16 border-b border-zinc-800 flex items-center justify-between px-6 bg-zinc-900/30">
                    <h1 className="font-semibold text-lg text-zinc-100">Live Operations</h1>
                    <div className="flex items-center gap-4">
                        <button className="relative p-2 text-zinc-400 hover:text-fiery-400 transition-colors">
                            <Bell size={20} />
                            <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-rose-500 rounded-full animate-pulse" />
                        </button>
                        <div className="w-8 h-8 rounded-full bg-fiery-500/20 border border-fiery-500/50 flex items-center justify-center text-xs font-bold text-fiery-500">
                            AD
                        </div>
                    </div>
                </header>
                <div className="flex-1 relative overflow-hidden p-4">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};
