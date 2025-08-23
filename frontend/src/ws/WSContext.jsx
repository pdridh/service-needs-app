import { createContext, useContext } from "react";

export const WSContext = createContext();

export default function useWS() {
    return useContext(WSContext);
}
