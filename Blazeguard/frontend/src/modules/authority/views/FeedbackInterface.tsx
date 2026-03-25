import React, { useState } from 'react';
import { MessageSquare, Send, XCircle, CheckCircle, AlertTriangle } from 'lucide-react';
import { fakeBackendAPI } from '../../../services/fakeBackendAPI';

export const FeedbackInterface: React.FC = () => {
    const [submitted, setSubmitted] = useState(false);
    const [type, setType] = useState('false-positive');
    const [incidentId, setIncidentId] = useState('');
    const [description, setDescription] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await fakeBackendAPI.submitFeedback({
                type,
                incidentId,
                description,
                timestamp: new Date().toISOString()
            });
            setSubmitted(true);
            console.log('✅ Feedback submitted successfully');
            setTimeout(() => {
                setSubmitted(false);
                setIncidentId('');
                setDescription('');
            }, 3000);
        } catch (error) {
            console.error('Failed to submit feedback:', error);
        }
    };

    return (
        <div className="p-6 max-w-3xl mx-auto">
            <div className="mb-8">
                <h1 className="text-2xl font-bold text-zinc-100 flex items-center gap-3">
                    <MessageSquare className="text-fiery-500" />
                    Operator Feedback & Reporting
                </h1>
                <p className="text-zinc-400 mt-2">
                    Report false positives, delayed detections, or suggest system improvements.
                    Feedback is analyzed by the technical team to improve model accuracy.
                </p>
            </div>

            <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 md:p-8">
                {submitted ? (
                    <div className="py-12 flex flex-col items-center text-center animate-in fade-in zoom-in duration-300">
                        <div className="w-16 h-16 bg-fiery-500/20 text-fiery-500 rounded-full flex items-center justify-center mb-4">
                            <CheckCircle size={32} />
                        </div>
                        <h3 className="text-xl font-bold text-zinc-100">Feedback Submitted</h3>
                        <p className="text-zinc-500 mt-2 max-w-xs">
                            Thank you for your report. The incident ID has been flagged for review.
                        </p>
                        <button
                            onClick={() => setSubmitted(false)}
                            className="mt-6 text-fiery-500 hover:text-fiery-400 font-medium text-sm"
                        >
                            Submit another report
                        </button>
                    </div>
                ) : (
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div className="grid md:grid-cols-2 gap-6">
                            <div className="space-y-2">
                                <label className="text-sm font-medium text-zinc-300">Report Type</label>
                                <div className="grid grid-cols-1 gap-2">
                                    {[
                                        { id: 'false-positive', label: 'False Positive Detection', icon: XCircle },
                                        { id: 'delayed', label: 'Delayed Detection', icon: ClockIcon },
                                        { id: 'accuracy', label: 'Location/Severity Inaccuracy', icon: AlertTriangle },
                                        { id: 'system', label: 'System Issue', icon: MessageSquare },
                                    ].map((item) => (
                                        <div
                                            key={item.id}
                                            onClick={() => setType(item.id)}
                                            className={`
                        cursor-pointer p-3 rounded-xl border flex items-center gap-3 transition-all
                        ${type === item.id
                                                    ? 'bg-fiery-500/10 border-fiery-500 text-fiery-400'
                                                    : 'bg-zinc-950/50 border-zinc-800 text-zinc-400 hover:bg-zinc-800 hover:border-zinc-700'}
                      `}
                                        >
                                            <item.icon size={18} />
                                            <span className="text-sm font-medium">{item.label}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="space-y-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-zinc-300">Incident ID (Optional)</label>
                                    <input
                                        type="text"
                                        placeholder="e.g., F-102"
                                        value={incidentId}
                                        onChange={(e) => setIncidentId(e.target.value)}
                                        className="w-full bg-zinc-950 border border-zinc-800 rounded-lg p-3 text-zinc-200 placeholder:text-zinc-600 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-colors"
                                    />
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-zinc-300">Description</label>
                                    <textarea
                                        rows={5}
                                        placeholder="Please provide details about the issue..."
                                        value={description}
                                        onChange={(e) => setDescription(e.target.value)}
                                        className="w-full bg-zinc-950 border border-zinc-800 rounded-lg p-3 text-zinc-200 placeholder:text-zinc-600 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-colors resize-none"
                                        required
                                    ></textarea>
                                </div>
                            </div>
                        </div>

                        <div className="pt-4 border-t border-zinc-800 flex justify-end">
                            <button
                                type="submit"
                                className="bg-fiery-600 hover:bg-fiery-500 text-white px-6 py-2.5 rounded-lg font-medium flex items-center gap-2 transition-colors shadow-lg shadow-fiery-900/20"
                            >
                                <Send size={18} />
                                Submit Feedback
                            </button>
                        </div>
                    </form>
                )}
            </div>
        </div>
    );
};

// Helper icon
const ClockIcon = ({ size, className }: { size?: number, className?: string }) => (
    <svg
        xmlns="http://www.w3.org/2000/svg"
        width={size}
        height={size}
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        className={className}
    >
        <circle cx="12" cy="12" r="10" />
        <polyline points="12 6 12 12 16 14" />
    </svg>
);
