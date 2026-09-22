import {SecretsList} from "../components/secrets/SecretsList.tsx";
import type {Secret} from "../util/secret.ts";
import {useEffect, useState} from "react";
import {ApiError, contributeToUnlock, getSecrets, getUnlockRequests, getUnlockResult, postSecret, postThresholdSecret, startUnlock, type ThresholdParticipantInput, type UnlockRequest} from "../util/api.ts";
import {SecretModal} from "../components/secrets/SecretModal.tsx";
import {UnlockRequests} from "../components/secrets/UnlockRequests.tsx";

export const SecretsScreen = () => {
    const [secrets, setSecrets] = useState<Secret[]>([]);
    const [selectedSecret, setSelectedSecret] = useState<Secret | undefined>(undefined);
    const [filterValue, setFilterValue] = useState("");
    const [isModalOpen, setModalOpen] = useState(false);
    const [unlockRequests, setUnlockRequests] = useState<UnlockRequest[]>([]);
    const [activeUnlock, setActiveUnlock] = useState<number | undefined>();
    const [unlockProgress, setUnlockProgress] = useState("");
    const [requestPollingError, setRequestPollingError] = useState("");
    const [unlockError, setUnlockError] = useState("");

    useEffect(() => {
        retrieveSecrets()
    }, [])

    useEffect(() => {
        let stopped = false;
        let timeout: number | undefined;
        const refresh = async () => {
            try {
                const requests = await getUnlockRequests();
                if (stopped) return;
                setUnlockRequests(current => {
                    if (current.some(previous => !requests.some(request => request.id === previous.id))) {
                        retrieveSecrets();
                    }
                    return requests;
                });
                setRequestPollingError("");
            } catch {
                if (!stopped) setRequestPollingError("Could not refresh unlock requests. Retrying automatically.");
            } finally {
                if (!stopped) timeout = window.setTimeout(refresh, 3000);
            }
        };
        refresh();
        return () => {
            stopped = true;
            if (timeout !== undefined) window.clearTimeout(timeout);
        };
    }, []);

    useEffect(() => {
        if (!activeUnlock) return;
        let stopped = false;
        let timeout: number | undefined;
        const refresh = async () => {
            try {
                const result = await getUnlockResult(activeUnlock);
                if (stopped) return;
                setUnlockProgress(`${result.contributions} of ${result.threshold} shares contributed`);
                if (result.ready && result.secret) {
                    setSelectedSecret({...result.secret, locked: false});
                    setModalOpen(true);
                    setActiveUnlock(undefined);
                    setUnlockProgress("");
                    setUnlockError("");
                    return;
                }
                timeout = window.setTimeout(refresh, 2000);
            } catch {
                if (stopped) return;
                setActiveUnlock(undefined);
                setUnlockProgress("");
                setUnlockError("The unlock request expired or could not be refreshed. Start a new unlock request and try again.");
            }
        };
        refresh();
        return () => {
            stopped = true;
            if (timeout !== undefined) window.clearTimeout(timeout);
        };
    }, [activeUnlock]);

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

    async function uploadThresholdSecret(s: Secret, threshold: number, participants: ThresholdParticipantInput[]) {
        await postThresholdSecret(s, threshold, participants);
        retrieveSecrets();
    }

    async function selectSecret(s: Secret) {
        if (s.locked && s.id) {
            setUnlockError("");
            let request;
            try {
                request = await startUnlock(s.id);
            } catch (error) {
                if (error instanceof ApiError && error.status === 409) {
                    setUnlockRequests(await getUnlockRequests());
                    return;
                }
                throw error;
            }
            setActiveUnlock(request.id);
            setUnlockProgress(`1 of ${s.threshold} shares contributed`);
            return;
        }
        setSelectedSecret(s)
        setModalOpen(true)
    }

    async function contribute(requestId: number) {
        await contributeToUnlock(requestId);
        setUnlockRequests(await getUnlockRequests());
    }

    return <>
        <div className="flex flex-row justify-center bg-base-200 h-full">
            <div className="w-full md:w-3/4 flex flex-col gap-4 p-4">
                {requestPollingError && <div role="alert" className="alert alert-error"><span>{requestPollingError}</span></div>}
                {unlockError && <div role="alert" className="alert alert-error"><span>{unlockError}</span></div>}
                <UnlockRequests requests={unlockRequests} contribute={contribute}/>
                {activeUnlock && <div className="alert alert-info"><span>Waiting for people to participate: {unlockProgress}</span></div>}
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
            handleThresholdSecret={uploadThresholdSecret}
            existingSecret={selectedSecret}
            isOpen={isModalOpen}
            closeModal={() => setModalOpen(false)}
            accessChanged={() => { setModalOpen(false); retrieveSecrets(); }}
        />
    </>
}
