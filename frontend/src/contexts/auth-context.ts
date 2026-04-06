import {createContext} from "react";
import {type AuthStatus} from "../util/authStatus.ts";

export type AuthContextType = {
    authStatus: AuthStatus | null;
    login: () => void;
    logout: () => Promise<void>;
    reloadAuth: () => void;
};

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
