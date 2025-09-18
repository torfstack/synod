import {SecretsList} from "../components/secrets/SecretsList.tsx";
import type {Secret} from "../util/secret.ts";
import {useEffect, useState} from "react";
import {getSecrets, postSecret} from "../util/api.ts";
import {SecretModal} from "../components/secrets/SecretModal.tsx";
import {showModal} from "../util/modal.ts";

export const SecretsScreen = () => {
    const [secrets, setSecrets] = useState<Secret[]>([]);
    const [selectedSecret, setSelectedSecret] = useState<Secret | undefined>(undefined);

    useEffect(() => {
        retrieveSecrets()
    }, [])

    function retrieveSecrets() {
        getSecrets().then(resp => resp.json()).then(
            (json) => {
                setSecrets(json)
            }
        )
    }

    async function uploadSecret(s: Secret) {
        await postSecret(s)
        retrieveSecrets()
        setSelectedSecret(undefined)
    }

    function selectSecret(s: Secret) {
        setSelectedSecret(s)
        showModal()
    }

    return <>
        <div className="flex flex-row justify-center bg-base-200 h-full">
            <div className="w-full md:w-3/4 flex flex-col gap-4 p-4">
                <div className="flex flex-row gap-4">
                    <input
                        type="text"
                        placeholder="Search"
                        className="input input-bordered w-3/4 md:w-1/2"
                    />
                    <button className="btn btn-neutral" onClick={() => {
                        setSelectedSecret(undefined);
                        showModal()
                    }}>
                        Add Secret
                    </button>
                </div>
                <SecretsList secrets={secrets} clickedSecret={selectSecret}/>
            </div>
        </div>

        <SecretModal handleSecret={uploadSecret} existingSecret={selectedSecret}/>
    </>
}

