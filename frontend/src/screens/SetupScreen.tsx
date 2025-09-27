import {useState} from "react";
import {postSetupPassword, postSetupPlain} from "../util/api.ts";

export const SetupScreen = () => {
    const [isChecked, setChecked] = useState(true);
    const [password, setPassword] = useState("");

    const handleSetup = () => {
        let promise;
        if (isChecked) {
            if (password.length == 0) return
            promise = postSetupPassword(password)
        } else {
            promise = postSetupPlain()
        }
        promise.then((_) => window.location.reload())
    };

    return <>
        <div className="flex flex-row justify-center bg-base-200 h-full">
            <div className="w-full md:w-3/4 flex flex-col gap-4 p-4">
                {/* Checkbox */}
                <label>
                    <input
                        className="checkbox"
                        type="checkbox"
                        checked={isChecked}
                        onChange={(e) => setChecked(e.target.checked)}
                    />
                    Use Password
                </label>

                {/* Password input */}
                <input
                    className="input"
                    type="text"
                    disabled={!isChecked}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />

                {/* Setup button */}
                <button
                    className="bg-blue-600 text-white px-4 py-2 rounded"
                    onClick={handleSetup}
                >
                    Setup
                </button>
            </div>
        </div>
    </>
}