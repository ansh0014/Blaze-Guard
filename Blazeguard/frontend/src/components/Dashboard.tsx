import React from 'react';
import { useAuth } from '../context/AuthContext';

export const Dashboard: React.FC = () => {
    const { user, logout } = useAuth();

    return (
        <div style={{
            minHeight: '100vh',
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            padding: '2rem',
        }}>
            <div style={{
                maxWidth: '1200px',
                margin: '0 auto',
            }}>
                <div style={{
                    background: 'rgba(255, 255, 255, 0.1)',
                    backdropFilter: 'blur(10px)',
                    borderRadius: '20px',
                    padding: '2rem',
                    boxShadow: '0 8px 32px 0 rgba(31, 38, 135, 0.37)',
                    border: '1px solid rgba(255, 255, 255, 0.18)',
                }}>
                    <div style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        marginBottom: '2rem',
                    }}>
                        <h1 style={{
                            color: 'white',
                            fontSize: '2.5rem',
                            margin: 0,
                        }}>
                            Fire Authority Dashboard
                        </h1>
                        <button
                            onClick={logout}
                            style={{
                                background: 'rgba(255, 255, 255, 0.2)',
                                color: 'white',
                                border: 'none',
                                padding: '0.75rem 1.5rem',
                                borderRadius: '10px',
                                cursor: 'pointer',
                                fontSize: '1rem',
                                fontWeight: 'bold',
                                transition: 'all 0.3s ease',
                            }}
                            onMouseEnter={(e) => {
                                e.currentTarget.style.background = 'rgba(255, 255, 255, 0.3)';
                            }}
                            onMouseLeave={(e) => {
                                e.currentTarget.style.background = 'rgba(255, 255, 255, 0.2)';
                            }}
                        >
                            Logout
                        </button>
                    </div>

                    {user && (
                        <div style={{
                            background: 'rgba(255, 255, 255, 0.15)',
                            borderRadius: '15px',
                            padding: '2rem',
                            marginBottom: '2rem',
                        }}>
                            <div style={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: '2rem',
                            }}>
                                {user.picture && (
                                    <img
                                        src={user.picture}
                                        alt={user.name}
                                        style={{
                                            width: '100px',
                                            height: '100px',
                                            borderRadius: '50%',
                                            border: '3px solid white',
                                        }}
                                    />
                                )}
                                <div>
                                    <h2 style={{
                                        color: 'white',
                                        fontSize: '2rem',
                                        margin: '0 0 0.5rem 0',
                                    }}>
                                        Welcome, {user.name}!
                                    </h2>
                                    <p style={{
                                        color: 'rgba(255, 255, 255, 0.8)',
                                        fontSize: '1.1rem',
                                        margin: '0.25rem 0',
                                    }}>
                                        <strong>Email:</strong> {user.email}
                                    </p>
                                    <p style={{
                                        color: 'rgba(255, 255, 255, 0.8)',
                                        fontSize: '1.1rem',
                                        margin: '0.25rem 0',
                                    }}>
                                        <strong>Role:</strong> {user.role}
                                    </p>
                                    <p style={{
                                        color: 'rgba(255, 255, 255, 0.6)',
                                        fontSize: '0.9rem',
                                        margin: '0.5rem 0 0 0',
                                    }}>
                                        Member since: {new Date(user.created_at).toLocaleDateString()}
                                    </p>
                                </div>
                            </div>
                        </div>
                    )}

                    <div style={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
                        gap: '1.5rem',
                    }}>
                        {[
                            { title: 'Active Incidents', value: '12', icon: '🔥' },
                            { title: 'Units Available', value: '8', icon: '🚒' },
                            { title: 'Personnel On Duty', value: '45', icon: '👨‍🚒' },
                            { title: 'Response Time Avg', value: '4.2 min', icon: '⏱️' },
                        ].map((stat, index) => (
                            <div
                                key={index}
                                style={{
                                    background: 'rgba(255, 255, 255, 0.1)',
                                    borderRadius: '15px',
                                    padding: '1.5rem',
                                    textAlign: 'center',
                                    transition: 'transform 0.3s ease',
                                }}
                                onMouseEnter={(e) => {
                                    e.currentTarget.style.transform = 'translateY(-5px)';
                                }}
                                onMouseLeave={(e) => {
                                    e.currentTarget.style.transform = 'translateY(0)';
                                }}
                            >
                                <div style={{ fontSize: '3rem', marginBottom: '0.5rem' }}>
                                    {stat.icon}
                                </div>
                                <h3 style={{
                                    color: 'white',
                                    fontSize: '1.2rem',
                                    margin: '0.5rem 0',
                                }}>
                                    {stat.title}
                                </h3>
                                <p style={{
                                    color: 'white',
                                    fontSize: '2rem',
                                    fontWeight: 'bold',
                                    margin: '0.5rem 0',
                                }}>
                                    {stat.value}
                                </p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
