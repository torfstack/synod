import './App.css'
import {Navbar} from "./components/Navbar.tsx";
import {SecretsScreen} from "./screens/SecretsScreen.tsx";
import {StartScreen} from "./screens/StartScreen.tsx";
import {SetupScreen} from "./screens/SetupScreen.tsx";
import {UnsealScreen} from "./screens/UnsealScreen.tsx";
import {useAuth} from "./contexts/AuthContext.tsx";

export const App = () => {
    const {authStatus} = useAuth()

    const isAuthenticated = authStatus && authStatus.isAuthenticated
    const isSetup = authStatus && authStatus.isSetup
    const needsToUnseal = authStatus && authStatus.needsToUnseal

    return (
        <div className="flex flex-col app">
            <Navbar/>
            {isAuthenticated && isSetup && !needsToUnseal && <SecretsScreen/>}
            {isAuthenticated && isSetup && needsToUnseal && <UnsealScreen/>}
            {isAuthenticated && !isSetup && <SetupScreen/>}
            {!isAuthenticated && <StartScreen/>}
        </div>
    )
}
