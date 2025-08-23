import { useEffect, useState } from "react";
import { API_GET_BUSINESSES_URL } from "../config";

const Browse = () => {
    const [page, setPage] = useState(1);
    const [pageSize] = useState(10);
    const [filters, setFilters] = useState({
        sortOrder: "",
        sortBy: "",
    });
    const [debouncedFilters, setDebouncedFilters] = useState({
        sortOrder: "",
        sortBy: "",
    });
    const [search, setSearch] = useState("");
    const [debouncedSearch, setDebouncedSearch] = useState("");
    const [businesses, setBusinesses] = useState([]);
    const [total, setTotal] = useState(0);

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(search);
            setPage(1);
        }, 400);

        return () => clearTimeout(timer);
    }, [search]);

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedFilters(filters);
            setPage(1);
        }, 400);

        return () => clearTimeout(timer);
    }, [filters]);

    useEffect(() => {
        const urlParams = new URLSearchParams({
            page,
            pageSize,
            ...debouncedFilters,
            search: debouncedSearch,
        });
        fetch(`${API_GET_BUSINESSES_URL}?${urlParams}`, {
            credentials: "include",
        })
            .then((res) => res.json())
            .then((jsonData) => {
                setBusinesses(jsonData.data.businesses);
                setTotal(jsonData.data.total);
            });
    }, [page, pageSize, debouncedSearch, debouncedFilters]);

    return (
        <div className="container mx-auto px-4 py-8">
            <div className="bg-white/10 backdrop-blur-lg rounded-3xl p-12 border border-white/20">
                <h1 className="text-4xl font-bold text-white mb-6">
                    Businesses
                </h1>
                <p className="text-lg text-white/90 leading-relaxed mb-8">
                    Browse services and businesses. Use filters to narrow your
                    search.
                </p>

                {/* Filters */}
                <div className="flex flex-col gap-5 bg-white/20 rounded-2xl p-6 border border-white/30 mb-8">
                    <input
                        value={search}
                        className="w-full px-4 py-2 rounded-xl bg-white/20 border border-white/30 text-white placeholder-white/60 focus:outline-none focus:ring-2 focus:ring-white/50"
                        placeholder="Search businesses..."
                        onChange={(e) => {
                            setSearch(e.target.value);
                            setPage(1);
                        }}
                    />
                    {/* TODO change how this select looks  */}
                    <select
                        className="w-full px-4 py-2 rounded-xl bg-white/20 border border-white/30 text-white"
                        value={filters.sortBy}
                        onChange={(e) =>
                            setFilters((f) => ({
                                ...f,
                                sortBy: e.target.value,
                            }))
                        }
                    >
                        <option value="">Sort by</option>
                        <option value="name">Name</option>
                    </select>

                    <select
                        className="w-full px-4 py-2 rounded-xl bg-white/20 border border-white/30 text-white"
                        value={filters.sortOrder}
                        onChange={(e) =>
                            setFilters((f) => ({
                                ...f,
                                sortOrder: e.target.value,
                            }))
                        }
                    >
                        <option value="asc">Ascending</option>
                        <option value="desc">Descending</option>
                    </select>
                </div>

                {/* Business list */}
                <div className="flex flex-col gap-5">
                    {businesses &&
                        businesses.map((b) => (
                            <div
                                key={b.id}
                                className="bg-white/20 rounded-2xl p-6 border border-white/30 shadow-lg"
                            >
                                <h2 className="text-xl font-semibold text-white mb-2">
                                    {b.name}
                                </h2>
                                <p className="text-white/80 font-bold">
                                    {b.category}
                                </p>
                                <p className="text-white/80">
                                    {b.description
                                        ? b.description
                                        : "No description available..."}
                                </p>
                            </div>
                        ))}
                </div>

                {/* Pagination */}
                <div className="flex justify-center items-center gap-4 mt-8">
                    <button
                        disabled={page === 1}
                        onClick={() => setPage(page - 1)}
                        className="px-4 py-2 rounded-xl bg-white/20 border border-white/30 text-white disabled:opacity-40 hover:bg-white/30"
                    >
                        Prev
                    </button>
                    <span className="text-white">{page}</span>
                    <button
                        disabled={page * pageSize >= total}
                        onClick={() => setPage(page + 1)}
                        className="px-4 py-2 rounded-xl bg-white/20 border border-white/30 text-white disabled:opacity-40 hover:bg-white/30"
                    >
                        Next
                    </button>
                </div>
            </div>
        </div>
    );
};

export default Browse;
