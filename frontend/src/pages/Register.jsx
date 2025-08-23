import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    AUTH_REGISTER_BUSINESS_URL,
    AUTH_REGISTER_CONSUMER_URL,
} from "../config";

const Register = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [userType, setUserType] = useState("consumer");
    const [consumerFirstName, setConsumerFirstName] = useState("");
    const [consumerLastName, setConsumerLastName] = useState("");

    const [businessName, setBusinessName] = useState("");
    const [businessCategory, setBusinessCategory] = useState("");

    const [checkPassword, setCheckPassword] = useState("");
    const [error, setError] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setIsLoading(true);

        try {
            if (password != checkPassword) {
                setError("Confirm password must match the password.");
                return;
            }

            if (userType === "consumer") {
                const res = await fetch(AUTH_REGISTER_CONSUMER_URL, {
                    method: "POST",
                    credentials: "include",
                    body: JSON.stringify({
                        email,
                        password,
                        firstName: consumerFirstName,
                        lastName: consumerLastName,
                    }),
                });

                if (res.ok) {
                    navigate("/login");
                }
            } else {
                const res = await fetch(AUTH_REGISTER_BUSINESS_URL, {
                    method: "POST",
                    credentials: "include",
                    body: JSON.stringify({
                        email,
                        password,
                        name: businessName,
                        category: businessCategory,
                        longitude: Math.random() * 24.0, // TODO get this from a location that the user will pin on map
                        latitude: 100.0 * Math.random(),
                    }),
                });

                if (res.ok) {
                    navigate("/login");
                }
            }
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
                        Sign Up
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

                        <div>
                            <label className="block text-sm font-semibold text-gray-700 mb-2">
                                Confirm Password
                            </label>
                            <input
                                type="password"
                                className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                                value={checkPassword}
                                onChange={(e) =>
                                    setCheckPassword(e.target.value)
                                }
                                placeholder="Confirm password"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-semibold text-gray-700 mb-2">
                                Register as
                            </label>
                            <select
                                name="usertypes"
                                onChange={(e) => setUserType(e.target.value)}
                                className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                            >
                                <option value="consumer">Consumer</option>
                                <option value="business">Business</option>
                            </select>
                        </div>

                        {userType === "consumer" && (
                            <>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                                        First name
                                    </label>
                                    <input
                                        value={consumerFirstName}
                                        onChange={(e) =>
                                            setConsumerFirstName(e.target.value)
                                        }
                                        type="text"
                                        required
                                        placeholder="John"
                                        className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                                        Last name
                                    </label>
                                    <input
                                        value={consumerLastName}
                                        onChange={(e) =>
                                            setConsumerLastName(e.target.value)
                                        }
                                        type="text"
                                        required
                                        placeholder="Doe"
                                        className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                                    />
                                </div>
                            </>
                        )}

                        {userType === "business" && (
                            <>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                                        Name
                                    </label>
                                    <input
                                        value={businessName}
                                        onChange={(e) =>
                                            setBusinessName(e.target.value)
                                        }
                                        type="text"
                                        required
                                        placeholder="Something inc."
                                        className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                                        Category
                                    </label>
                                    <input
                                        value={businessCategory}
                                        onChange={(e) =>
                                            setBusinessCategory(e.target.value)
                                        }
                                        type="text"
                                        required
                                        placeholder="Plumbing"
                                        className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:outline-none transition-colors duration-300"
                                    />
                                </div>
                            </>
                        )}

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
                                    Signing Up
                                </div>
                            ) : (
                                "Sign Up"
                            )}
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default Register;
