import useAuth from "../auth/AuthContext";

const Dashboard = () => {
    const { user } = useAuth();

    return (
        <div className="container mx-auto px-4 py-8">
            <div className="bg-white/10 backdrop-blur-lg rounded-3xl p-12 border border-white/20">
                <h1 className="text-4xl font-bold text-white mb-6">
                    Dashboard
                </h1>
                <p className="text-lg text-white/90 leading-relaxed mb-8">
                    Welcome to your dashboard! This is a protected route that
                    requires authentication.
                </p>

                <div className="bg-white/20 rounded-2xl p-6 border border-white/30">
                    <h3 className="text-xl font-semibold text-white mb-4">
                        User Information
                    </h3>
                    <div className="space-y-2 text-white/90">
                        <p>
                            <span className="font-semibold">Name:</span>{" "}
                            {user?.name}
                        </p>
                        <p>
                            <span className="font-semibold">Email:</span>{" "}
                            {user?.email}
                        </p>
                        <p>
                            <span className="font-semibold">User ID:</span>{" "}
                            {user?.id}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
