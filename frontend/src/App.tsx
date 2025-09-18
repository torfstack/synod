import './App.css'
import {useEffect, useState} from "react";
import {deleteAuth, getAuth} from "./util/api.ts";
import {Navbar} from "./components/Navbar.tsx";
import {config} from "./util/config.ts";
import {StartScreen} from "./screens/StartScreen.tsx";
import {SecretsScreen} from "./screens/SecretsScreen.tsx";

export const App = () => {
    const [isAuthenticated, setAuthenticated] = useState(false);

    useEffect(() => {
        getAuth().then(res => {
            if (res.status == 200) {
                setAuthenticated(true);
            }
        })
    }, []);

    function signInWithProvider() {
        window.open(config.backendAuthStartUrl, "_self");
    }

    function signOut() {
        deleteAuth().then(res => {
            setAuthenticated(res.status == 204)
        })
    }

    return (
        <div className="flex flex-col app">
            <Navbar isAuthenticated={isAuthenticated}
                    loginButtonPressed={signInWithProvider}
                    logoutButtonPressed={signOut}/>
            {isAuthenticated ? <SecretsScreen/> : <StartScreen/>}
        </div>
    )
}
