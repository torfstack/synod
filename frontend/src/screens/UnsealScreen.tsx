import {useState} from "react";
import {postUnsealWithPassword} from "../util/api.ts";
import {Eye, EyeSlash} from "../icons/Eye.tsx";

export const UnsealScreen = () => {
    const [password, setPassword] = useState("");
    const [passwordVisible, setPasswordVisible] = useState(false)
    const [passwordsWarning, setPasswordsWarning] = useState(false)

    const handleUnseal = () => {
        postUnsealWithPassword(password).then((resp) => {
            if (resp.status != 204) {
                setPasswordsWarning(true)
                return
            }
            window.location.reload()
        })
    };

    return <div className="flex flex-row justify-center items-center bg-base-200 h-full">
        <div className="w-5/6 md:w-1/4 flex flex-col gap-4 p-4 bg-base-100 rounded-box shadow-md">
            <h2 className="text-xl">Unseal with password</h2>
            <p>Provide the password to unseal your secrets.</p>

            {/* Password input */}
            <label className={passwordsWarning ? "input input-error w-full" : "input w-full"}>
                Password
                <input
                    className="grow"
                    type={passwordVisible ? "text" : "password"}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <button type="button" onClick={() => setPasswordVisible(prev => !prev)}
                        className="btn btn-ghost btn-xs">
                    {passwordVisible ? <Eye/> : <EyeSlash/>}
                </button>
            </label>

            {/* Unseal button */}
            <button
                className="btn btn-primary"
                onClick={handleUnseal}
            >
                Unseal secrets
            </button>
        </div>
    </div>
}