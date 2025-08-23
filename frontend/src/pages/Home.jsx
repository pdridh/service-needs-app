import { Link } from "react-router-dom";

const Home = () => {
    return (
        <div className="container mx-auto px-4 py-8">
            <div className="bg-white/10 backdrop-blur-lg rounded-3xl p-12 border border-white/20">
                <h1 className="text-5xl font-bold text-white mb-6">
                    Connect with businesses or sum shit
                </h1>
                <p className="text-xl text-white/90 leading-relaxed mb-8">
                    Good businesses here like a yellow page
                </p>

                <div className="flex gap-4">
                    <Link
                        to="/browse"
                        className="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-3 rounded-xl font-semibold transition-all duration-300 transform hover:-translate-y-0.5"
                    >
                        Browse
                    </Link>
                    <Link
                        to="/dashboard"
                        className="bg-white/20 hover:bg-white/30 text-white px-8 py-3 rounded-xl font-semibold border border-white/30 transition-all duration-300"
                    >
                        View Dashboard
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default Home;
