import {useEffect, useState} from "react";
import {postSetupPassword, postSetupPlain} from "../util/api.ts";
import {Eye, EyeSlash} from "../icons/Eye.tsx";

export const SetupScreen = () => {
    const [isChecked, setChecked] = useState(true);
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("")
    const [passwordVisible, setPasswordVisible] = useState(false)
    const [passwordsWarning, setPasswordsWarning] = useState(false)
    const [needsToConfirm, setNeedsToConfirm] = useState(false);

    useEffect(() => {
        setPasswordsWarning(false)
    }, [password, confirmPassword]);

    const handleSetup = () => {
        if (needsToConfirm) {
            setNeedsToConfirm(false)
            return
        }
        let promise;
        if (isChecked) {
            if (password.length == 0 || confirmPassword != password) {
                setPasswordsWarning(true)
                return
            }
            promise = postSetupPassword(password)
        } else {
            promise = postSetupPlain()
        }
        promise.then((resp) => {
            if (resp.status != 201) {
                // TODO: show warning
                return
            }
            window.location.reload()
        })
    };

    return <div className="flex flex-row justify-center items-center bg-base-200 h-full">
        <div className="w-5/6 md:w-1/4 flex flex-col gap-4 p-4 bg-base-100 rounded-box shadow-md">
            <h2 className="text-xl">Initial Setup</h2>
            <p>You can choose to add a password to protect your secrets.</p>
            <p>With a password: Your secrets will be protected with end-to-end encryption. Only you can decrypt
                them,
                and not even our servers can access your secrets.</p>
            <p>Without a password: Your secrets will still be encrypted, but only with server-side encryption. This
                means we can keep your data safe on our servers, but it won’t be protected if someone gains access
                to
                your account.</p>
            <p className={isChecked ? "text-lg" : "text-warning text-lg"}>
                For maximum security, we recommend setting a password.
            </p>
            {/* Checkbox */}
            <label className="label">
                <input
                    className="checkbox checkbox-primary"
                    type="checkbox"
                    checked={isChecked}
                    onChange={(e) => {
                        setChecked(e.target.checked)
                        setNeedsToConfirm(!e.target.checked)
                        setPassword("")
                        setConfirmPassword("")
                    }}
                />
                Use password
            </label>

            {/* Password input */}
            <label className={passwordsWarning ? "input input-error w-full" : "input w-full"}>
                Password
                <input
                    className="grow"
                    type={passwordVisible ? "text" : "password"}
                    disabled={!isChecked}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <button type="button" onClick={() => setPasswordVisible(prev => !prev)}
                        disabled={!isChecked}
                        className="btn btn-ghost btn-xs">
                    {passwordVisible ? <Eye/> : <EyeSlash/>}
                </button>
            </label>

            {/* Confirm password input */}
            <label className={passwordsWarning ? "input input-error w-full" : "input w-full"}>
                Repeat
                <input
                    className="grow"
                    type={passwordVisible ? "text" : "password"}
                    disabled={!isChecked}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                />
                <button type="button" onClick={() => setPasswordVisible(prev => !prev)}
                        disabled={!isChecked}
                        className="btn btn-ghost btn-xs">
                    {passwordVisible ? <Eye/> : <EyeSlash/>}
                </button>
            </label>

            {/* Setup button */}
            <button
                className={isChecked ? "btn btn-primary" : "btn btn-warning"}
                onClick={handleSetup}
            >
                {isChecked || !needsToConfirm ? "Confirm Setup" : "Continue without password?"}
            </button>
        </div>
    </div>
}