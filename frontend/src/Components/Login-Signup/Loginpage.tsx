import { useState } from 'react';

export default function Loginpage() {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');

    const handleLogin = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "include",
                body: JSON.stringify({ email, password }),
            });
            const data = await response.json();
            if (!response.ok) {
                alert(data.message || 'Signup failed');
                return;
            }
            // localStorage.setItem("token", data.token)
            alert(data.message);
        } catch (error) {
            console.error('Error:', error);
            alert('Failed to connect to backend from the frontend on login');
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-900">
            <div className="w-full max-w-md p-10 space-y-6 bg-gray-800 rounded-2xl shadow-2xl">
                <h2 className="text-3xl font-bold text-center text-white">Sign In</h2>

                <div className="space-y-2">
                    <label htmlFor="email" className="text-sm font-medium text-gray-300">Email Address</label>
                    <input
                        id="email"
                        type="email"
                        name="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        className="w-full px-4 py-3 text-gray-900 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
                    />
                </div>

                <div className="space-y-2">
                    <label htmlFor="password" className="text-sm font-medium text-gray-300">Password</label>
                    <input
                        id="password"
                        type="password"
                        name="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        className="w-full px-4 py-3 text-gray-900 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
                    />
                </div>

                <button
                    type='button'
                    onClick={handleLogin}
                    className="w-full py-3 text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-800 transition-all duration-200 font-medium"

                >
                    Sign In
                </button>

                <div className="text-white justify-between flex">
                    <a href="/Signuppage" className="text-blue-400 hover:text-blue-300 transition-colors">Create an account</a>
                    <a href="#" className="text-blue-400 hover:text-gray-300 transition-colors">Forgot Your password?</a>
                </div>
            </div>
        </div>
    );
}
