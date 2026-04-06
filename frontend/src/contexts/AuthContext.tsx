import React, {useEffect, useState} from "react";
import {config} from "../util/config.ts";
import {deleteAuth, getAuth} from "../util/api.ts";
import {AuthContext} from "./auth-context.ts";

export function AuthProvider({children}: { children: React.ReactNode }) {
    const [authStatus, setAuthStatus] = useState<AuthStatus | null>(null);

    const checkAuth = () => {
        getAuth()
            .then(authStatus => setAuthStatus(authStatus))
            .catch(() => setAuthStatus(null))
    }

    useEffect(() => {
        checkAuth();
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
        checkAuth();
    }

    return (
        <AuthContext.Provider value={{authStatus, login, logout, reloadAuth}}>
            {children}
        </AuthContext.Provider>
    );
}
