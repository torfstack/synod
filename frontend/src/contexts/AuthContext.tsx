import React, {createContext, useContext, useEffect, useState} from "react";
import {config} from "../util/config.ts";
import {type AuthStatus} from "../util/authStatus.ts";
import {deleteAuth, getAuth} from "../util/api.ts";

type AuthContextType = {
    authStatus: AuthStatus | null;
    login: () => void;
    logout: () => Promise<void>;
    reloadAuth: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({children}: { children: React.ReactNode }) {
    const [authStatus, setAuthStatus] = useState<AuthStatus | null>(null);

    const checkAuth = async () => {
        try {
            const authStatus = await getAuth();
            setAuthStatus(authStatus)
        } catch {
            setAuthStatus(null);
        }
    }

    useEffect(() => {
        (async () => {
            await checkAuth();
        })();
    }, []);

    useEffect(() => {
        function handleUnauthorized() {
            setAuthStatus(null);
        }

        window.addEventListener("unauthorized", handleUnauthorized);
        return () => window.removeEventListener("unauthorized", handleUnauthorized);
    }, []);

    const login = async () => {
        return window.open(config.backendAuthStartUrl, "_self");
    };

    const logout = async () => {
        return deleteAuth().then(() => {
            setAuthStatus(null)
        })
    };

    const reloadAuth = () => {
        (async () => {
            await checkAuth();
        })();
    }

    return (
        <AuthContext.Provider value={{authStatus, login, logout, reloadAuth}}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth(): AuthContextType {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
