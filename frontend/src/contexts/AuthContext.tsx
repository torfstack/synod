import React, {useEffect, useState} from "react";
import {config} from "../util/config.ts";
import {deleteAuth, getAuth} from "../util/api.ts";
import {AuthContext} from "./auth-context.ts";
import type {AuthStatus} from "../util/authStatus.ts";
import {
    clearPostLoginHint,
    getInitialAuthStatusFromPostLoginHint,
    setPostLoginHint,
} from "../util/postLoginHint.ts";

export function AuthProvider({children}: { children: React.ReactNode }) {
    const [authStatus, setAuthStatus] = useState<AuthStatus | null>(() => {
        return getInitialAuthStatusFromPostLoginHint(window.sessionStorage);
    });

    const checkAuth = () => {
        getAuth()
            .then(authStatus => {
                clearPostLoginHint(window.sessionStorage);
                setAuthStatus(authStatus)
            })
            .catch(() => setAuthStatus(null))
    }

    useEffect(() => {
        checkAuth();
    }, []);

    useEffect(() => {
        function handleUnauthorized() {
            clearPostLoginHint(window.sessionStorage);
            setAuthStatus(null);
        }

        window.addEventListener("unauthorized", handleUnauthorized);
        return () => window.removeEventListener("unauthorized", handleUnauthorized);
    }, []);

    const login = async () => {
        setPostLoginHint(window.sessionStorage);
        return window.open(config.backendAuthStartUrl, "_self");
    };

    const logout = async () => {
        return deleteAuth().then(() => {
            clearPostLoginHint(window.sessionStorage);
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
