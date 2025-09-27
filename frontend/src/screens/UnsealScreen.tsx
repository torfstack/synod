import {useState} from "react";
import {postUnsealWithPassword} from "../util/api.ts";

export const UnsealScreen = () => {
    const [password, setPassword] = useState("");

    const handleUnseal = () => {
        postUnsealWithPassword(password).then((_) => window.location.reload())
    };

    return <div className="flex flex-row justify-center bg-base-200 h-full">
        <div className="w-full md:w-3/4 flex flex-col gap-4 p-4">
            {/* Password input */}
            <input
                className="input"
                type="text"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />

            {/* Unseal button */}
            <button
                className="bg-blue-600 text-white px-4 py-2 rounded"
                onClick={handleUnseal}
            >
                Unseal
            </button>
        </div>
    </div>
}