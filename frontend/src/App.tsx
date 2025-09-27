import './App.css'
import {useEffect, useState} from "react";
import {deleteAuth, getAuth} from "./util/api.ts";
import {Navbar} from "./components/Navbar.tsx";
import {config} from "./util/config.ts";
import {SecretsScreen} from "./screens/SecretsScreen.tsx";
import {StartScreen} from "./screens/StartScreen.tsx";
import {SetupScreen} from "./screens/SetupScreen.tsx";
import {UnsealScreen} from "./screens/UnsealScreen.tsx";
import {type AuthStatus, EmptyAuthStatus} from "./util/authStatus.ts";

export const App = () => {
    const [authStatus, setAuthStatus] = useState<AuthStatus>(EmptyAuthStatus);

    useEffect(() => {
        getAuth().then(res => {
            if (res.status == 200) {
                res.json().then((json) => setAuthStatus(json))
            }
        })
    }, []);

    function signInWithProvider() {
        window.open(config.backendAuthStartUrl, "_self");
    }

    function signOut() {
        deleteAuth().then(_ => {
            setAuthStatus(EmptyAuthStatus)
        })
    }

    const isAuthenticated = authStatus.isAuthenticated
    const isSetup = authStatus.isSetup
    const needsToUnseal = authStatus.needsToUnseal

    return (
        <div className="flex flex-col app">
            <Navbar isAuthenticated={authStatus.isAuthenticated}
                    loginButtonPressed={signInWithProvider}
                    logoutButtonPressed={signOut}/>
            {isAuthenticated && isSetup && !needsToUnseal && <SecretsScreen/>}
            {isAuthenticated && isSetup && needsToUnseal && <UnsealScreen/>}
            {isAuthenticated && !isSetup && <SetupScreen/>}
            {!isAuthenticated && <StartScreen/>}
        </div>
    )
}
