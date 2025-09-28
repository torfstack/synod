import {SecretsList} from "../components/secrets/SecretsList.tsx";
import type {Secret} from "../util/secret.ts";
import {useEffect, useState} from "react";
import {getSecrets, postSecret} from "../util/api.ts";
import {SecretModal} from "../components/secrets/SecretModal.tsx";

export const SecretsScreen = () => {
    const [secrets, setSecrets] = useState<Secret[]>([]);
    const [selectedSecret, setSelectedSecret] = useState<Secret | undefined>(undefined);
    const [filterValue, setFilterValue] = useState("");
    const [isModalOpen, setModalOpen] = useState(false);

    useEffect(() => {
        retrieveSecrets()
    }, [])

    const filteredSecrets = filterSecrets(secrets, filterValue)

    function filterSecrets(secrets: Secret[], filterValue: string): Secret[] {
        return secrets.filter(secret => {
            return secret.key?.toLowerCase().includes(filterValue.toLowerCase()) ||
                secret.url?.toLowerCase().includes(filterValue.toLowerCase()) ||
                secret.tags?.some(tag => tag.toLowerCase().includes(filterValue.toLowerCase()))
        }).toSorted((s1, s2) => s1.id! - s2.id!);
    }

    function retrieveSecrets() {
        getSecrets().then(secrets => setSecrets(secrets))
    }

    async function uploadSecret(s: Secret) {
        return postSecret(s).then(() => retrieveSecrets())
    }

    function selectSecret(s: Secret) {
        setSelectedSecret(s)
        setModalOpen(true)
    }

    return <>
        <div className="flex flex-row justify-center bg-base-200 h-full">
            <div className="w-full md:w-3/4 flex flex-col gap-4 p-4">
                <div className="flex flex-row gap-4">
                    <input type="text" placeholder="Search" value={filterValue}
                           className="input input-bordered w-3/4 md:w-1/2"
                           onChange={(e) => setFilterValue(e.target.value)}/>
                    <button className="btn btn-neutral" onClick={() => {
                        setSelectedSecret(undefined);
                        setModalOpen(true)
                    }}>
                        Add Secret
                    </button>
                </div>
                <SecretsList secrets={filteredSecrets} clickedSecret={selectSecret}/>
            </div>
        </div>

        <SecretModal
            handleSecret={uploadSecret}
            existingSecret={selectedSecret}
            isOpen={isModalOpen}
            closeModal={() => setModalOpen(false)}
        />
    </>
}

