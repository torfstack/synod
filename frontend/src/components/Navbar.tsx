import {ThemeSwitcher} from "./ThemeSwitcher.tsx";

type NavbarProps = {
    isAuthenticated: boolean
    loginButtonPressed: () => void
    logoutButtonPressed: () => void
}

export const Navbar = (props: NavbarProps) => {
    return (
        <div className="flex flex-col">
            <div className="navbar bg-base-100 shadow-sm">
                <div className="flex-1">
                    <a className="btn btn-ghost text-l" href={"/"}>
                        <span className="badge">
                            KayVault
                        </span>
                    </a>
                </div>

                <div className="flex flex-row gap-4 items-center">
                    {props.isAuthenticated
                        ? <LogoutButton logoutButtonPressed={props.logoutButtonPressed}/>
                        : <LoginButton loginButtonPressed={props.loginButtonPressed}/>}

                    <ThemeSwitcher/>
                </div>
            </div>
        </div>
    );
};

type LoginButtonProps = {
    loginButtonPressed: () => void
}

const LoginButton = (props: LoginButtonProps) => {
    return <button className="btn btn-primary" onClick={props.loginButtonPressed}>
        <svg aria-label="Google logo" width="16" height="16" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
            <g>
                <path fill="#34a853" d="M153 292c30 82 118 95 171 60h62v48A192 192 0 0190 341"></path>
                <path fill="#4285f4" d="m386 400a140 175 0 0053-179H260v74h102q-7 37-38 57"></path>
                <path fill="#fbbc02" d="m90 341a208 200 0 010-171l63 49q-12 37 0 73"></path>
                <path fill="#ea4335" d="m153 219c22-69 116-109 179-50l55-54c-78-75-230-72-297 55"></path>
            </g>
        </svg>
        Login with Google
    </button>
}

type LogoutButtonProps = {
    logoutButtonPressed: () => void
}

const LogoutButton = (props: LogoutButtonProps) => {
    return <button className="btn btn-accent" onClick={props.logoutButtonPressed}>
        Logout
    </button>
}
