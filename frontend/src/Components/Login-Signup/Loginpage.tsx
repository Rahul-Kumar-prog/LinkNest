import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const API_BASE = 'http://localhost:8080';

export default function Loginpage() {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const navigate = useNavigate();

    const handleLogin = async () => {
        try {
            const response = await fetch(`${API_BASE}/api/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ email, password }),
            });
            const data = await response.json();
            if (!response.ok) {
                alert(data.error || 'Login failed');
                return;
            }
            localStorage.setItem('token', data.token);
            alert(data.message);
            navigate('/');
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to connect to backend from the frontend on login');
        }
    };

    const handleGoogleLogin = () => {
        window.location.href = `${API_BASE}/api/auth/google`;
    };

    return (
        <div className="flex min-h-screen items-center justify-center bg-stone-950 px-4">
            <div className="w-full max-w-md space-y-6 rounded-[2rem] border border-white/10 bg-slate-900/80 p-10 shadow-2xl shadow-black/30">
                <div className="space-y-2 text-center">
                    <h2 className="text-3xl font-black text-white">Sign in to LinkNest</h2>
                    <p className="text-sm text-stone-400">Google is the primary login path for the MVP. Email login still works for existing users.</p>
                </div>

                <button
                    type="button"
                    onClick={handleGoogleLogin}
                    className="w-full rounded-full bg-sky-400 py-3 font-semibold text-slate-950 transition hover:bg-sky-300"
                >
                    Continue with Google
                </button>

                <div className="flex items-center gap-3 text-xs uppercase tracking-[0.25em] text-stone-500">
                    <div className="h-px flex-1 bg-white/10" />
                    <span>or use email</span>
                    <div className="h-px flex-1 bg-white/10" />
                </div>

                <div className="space-y-2">
                    <label htmlFor="email" className="text-sm font-medium text-stone-300">Email Address</label>
                    <input
                        id="email"
                        type="email"
                        name="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        className="w-full rounded-2xl border border-white/10 bg-white px-4 py-3 text-stone-900 outline-none transition focus:border-sky-400"
                    />
                </div>

                <div className="space-y-2">
                    <label htmlFor="password" className="text-sm font-medium text-stone-300">Password</label>
                    <input
                        id="password"
                        type="password"
                        name="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        className="w-full rounded-2xl border border-white/10 bg-white px-4 py-3 text-stone-900 outline-none transition focus:border-sky-400"
                    />
                </div>

                <button
                    type='button'
                    onClick={handleLogin}
                    className="w-full rounded-full border border-white/10 py-3 font-semibold text-white transition hover:bg-white/5"
                >
                    Sign In
                </button>

                <div className="flex justify-between text-sm text-white">
                    <a href="/Signuppage" className="text-sky-300 transition-colors hover:text-sky-200">Create an account</a>
                    <a href="/" className="text-stone-400 transition-colors hover:text-stone-200">Back home</a>
                </div>
            </div>
        </div>
    );
}
