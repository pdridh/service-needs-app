import React, { createContext, useEffect, useRef, useState } from "react";
import useAuth from "../auth/AuthContext";

const WSContext = createContext(null);

const WSProvider = ({ children }) => {
    const { user } = useAuth();
    const wsRef = useRef(null);
    const [connected, setConnected] = useState(false);

    useEffect(() => {
        if (!user) return; // only connect if logged in

        const ws = new WebSocket(`ws://localhost:8080/ws`);
        wsRef.current = ws;

        ws.onopen = () => {
            setConnected(true);
            ws.send(
                JSON.stringify({
                    code: "hello",
                })
            );
        };
        ws.onclose = () => setConnected(false);
        ws.onerror = () => {
            console.log("ERROR CONNECTING TO WS");
        };

        ws.onmessage = (ev) => {
            console.log(ev.data);
        };

        return () => {
            ws.close(); // close on logout or unmount
        };
    }, [user]);

    return (
        <WSContext.Provider value={{ ws: wsRef.current, connected }}>
            {children}
        </WSContext.Provider>
    );
};

export default WSProvider;
