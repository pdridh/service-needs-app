import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import AuthProvider from "./auth/AuthProvider";
import Navbar from "./components/Navbar";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import Home from "./pages/Home";
import RequireAuth from "./auth/RequireAuth";
import WSProvider from "./ws/WSProvider";
import Register from "./pages/Register";
import Browse from "./pages/Browse";

function App() {
    return (
        <div className="min-h-screen bg-black">
            <BrowserRouter>
                <AuthProvider>
                    <WSProvider>
                        <Navbar />
                        <Routes>
                            <Route path="/" element={<Home />} />
                            <Route
                                path="/browse"
                                element={
                                    <RequireAuth roles={["consumer"]}>
                                        <Browse />
                                    </RequireAuth>
                                }
                            />
                            <Route path="/login" element={<Login />} />
                            <Route path="/register" element={<Register />} />
                            <Route
                                path="/dashboard"
                                element={
                                    <RequireAuth>
                                        <Dashboard />
                                    </RequireAuth>
                                }
                            />
                            <Route
                                path="*"
                                element={<Navigate to="/" replace />}
                            />
                        </Routes>
                    </WSProvider>
                </AuthProvider>
            </BrowserRouter>
        </div>
    );
}

export default App;
