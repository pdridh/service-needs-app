import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import useAuth from "../auth/AuthContext";

const Login = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    const { login } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

    const from = location?.state?.path || "/dashboard";

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setIsLoading(true);

        try {
            await login(email, password);
            navigate(from, { replace: true });
        } catch (error) {
            setError(
                error?.message ||
                    "An unexpected error occurred. Please try again."
            );
        } finally {
            // Always reset loading state
            setIsLoading(false);
        }

        setIsLoading(false);
    };

    return (
        <div className="container mx-auto px-4 py-8">
            <div className="flex justify-center items-center min-h-[80vh]">
                <div className="bg-white/95 backdrop-blur-xl rounded-3xl p-10 shadow-2xl border border-white/30 w-full max-w-md">
                    <h2 className="text-3xl font-bold text-center mb-8 text-gray-900">
                        Welcome Back
                    </h2>

                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div>
                            <label className="block text-sm font-semibold text-gray-700 mb-2">
                                Email
                            </label>
                            <input
                                type="email"
                                className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="Enter your email"
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-semibold text-gray-700 mb-2">
                                Password
                            </label>
                            <input
                                type="password"
                                className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Enter your password"
                                required
                            />
                        </div>

                        {error && (
                            <div className="text-red-500 text-sm text-center">
                                {error}
                            </div>
                        )}

                        <button
                            type="submit"
                            className="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-3 px-6 rounded-xl font-semibold transition-all duration-300 transform hover:-translate-y-0.5 disabled:opacity-50 disabled:cursor-not-allowed"
                            disabled={isLoading}
                        >
                            {isLoading ? (
                                <div className="flex items-center justify-center">
                                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                                    Logging in...
                                </div>
                            ) : (
                                "Login"
                            )}
                        </button>
                    </form>

                    <div className="mt-8 text-center text-gray-600 space-y-1">
                        <p className="font-medium">Demo credentials:</p>
                        <p>
                            <span className="font-semibold">Email:</span>{" "}
                            demo@example.com
                        </p>
                        <p>
                            <span className="font-semibold">Password:</span>{" "}
                            password
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Login;
