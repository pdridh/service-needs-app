import { Navigate, useLocation } from "react-router-dom";
import useAuth from "./AuthContext";
import Loading from "../components/Loading";

export default function RequireAuth({ children, roles }) {
    const { user, loading } = useAuth();
    const location = useLocation();

    if (loading) return <Loading></Loading>;

    if (!user) {
        return (
            <Navigate to="/login" state={{ path: location.pathname }} replace />
        );
    }

    if (user.type == "admin") {
        return children;
    }

    if (roles && !roles.includes(user.type)) {
        return <Navigate to="/" replace />;
    }

    return children;
}
