import { Link, useNavigate } from "react-router-dom";
import useAuth from "../auth/AuthContext";

const Navbar = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate("/login");
    };

    return (
        <nav className="bg-white/10 backdrop-blur-lg border border-white/20 rounded-2xl mx-4 mt-4 p-4">
            <div className="container mx-auto flex justify-between items-center">
                <Link to="/" className="text-2xl font-bold text-white">
                    Servicer
                </Link>

                <div className="flex items-center space-x-6">
                    {user != null ? (
                        <>
                            <Link
                                to="/dashboard"
                                className="text-white hover:bg-white/20 px-4 py-2 rounded-lg transition-all duration-300"
                            >
                                Dashboard
                            </Link>
                            <span className="text-white/80">
                                Welcome, {user.email}
                            </span>
                            <button
                                onClick={handleLogout}
                                className="bg-white/20 hover:bg-white/30 text-white px-6 py-2 rounded-lg border border-white/30 transition-all duration-300"
                            >
                                Logout
                            </button>
                        </>
                    ) : (
                        <>
                            <Link
                                to="/register"
                                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg transition-all duration-300 transform hover:-translate-y-0.5"
                            >
                                Sign Up
                            </Link>

                            <Link
                                to="/login"
                                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg transition-all duration-300 transform hover:-translate-y-0.5"
                            >
                                Login
                            </Link>
                        </>
                    )}
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
