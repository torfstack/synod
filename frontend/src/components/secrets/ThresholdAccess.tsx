import {useEffect, useState} from "react";
import type {Secret} from "../../util/secret.ts";
import {ApiError, getThresholdParticipants, searchShareRecipients, setThresholdParticipants, type ShareRecipient, type ThresholdParticipant, type ThresholdParticipantInput} from "../../util/api.ts";

const editableParticipants = (participants: ThresholdParticipant[]): ThresholdParticipantInput[] =>
    participants.filter(participant => participant.role !== "owner")
        .map(({sharingId, role}) => ({sharingId, role: role as "holder" | "maintainer"}));

export const ThresholdAccess = ({secret, changed}: {secret: Secret; changed: () => void}) => {
    const [participants, setParticipants] = useState<ThresholdParticipant[]>([]);
    const [expectedParticipants, setExpectedParticipants] = useState<ThresholdParticipant[]>([]);
    const [query, setQuery] = useState("");
    const [matches, setMatches] = useState<ShareRecipient[]>([]);
    const [role, setRole] = useState<"holder" | "maintainer">("holder");
    const [dirty, setDirty] = useState(false);
    const [error, setError] = useState("");
    useEffect(() => {
        getThresholdParticipants(secret.id!).then(current => {
            setParticipants(current);
            setExpectedParticipants(current);
        });
    }, [secret.id]);
    useEffect(() => {
        if (query.trim().length < 2) { setMatches([]); return; }
        const controller = new AbortController();
        const timeout = window.setTimeout(() => searchShareRecipients(query.trim(), controller.signal).then(setMatches), 200);
        return () => { window.clearTimeout(timeout); controller.abort(); };
    }, [query]);
    return <fieldset className="fieldset border-t border-base-300 pt-3">
        <legend className="fieldset-legend">Threshold access</legend>
        <div role="status" className="alert alert-warning py-2">
            <span>Applying access changes rotates the encryption key and locks the secret again for everyone.</span>
        </div>
        {error && <div role="alert" className="alert alert-error py-2"><span>{error}</span></div>}
        {participants.map(participant => <div className="flex items-center justify-between gap-2" key={participant.sharingId}>
            <span>{participant.fullName} <span className="badge badge-sm">{participant.role}</span></span>
            {participant.role !== "owner" && <div className="flex gap-2">
                {secret.role === "owner" && <select className="select select-sm" value={participant.role} onChange={async event => {
                    setParticipants(current => current.map(item => item.sharingId === participant.sharingId ? {...item, role: event.target.value as "holder" | "maintainer"} : item)); setDirty(true);
                }}><option value="holder">Holder</option><option value="maintainer">Maintainer</option></select>}
                {(secret.role === "owner" || participant.role === "holder") && <button type="button" className="btn btn-ghost btn-sm" onClick={() => { setParticipants(current => current.filter(item => item.sharingId !== participant.sharingId)); setDirty(true); }}>Remove</button>}
            </div>}
        </div>)}
        <div className="flex gap-2">
            <input className="input grow" value={query} onChange={event => setQuery(event.target.value)} placeholder="Find a person"/>
            {secret.role === "owner" && <select className="select" value={role} onChange={event => setRole(event.target.value as "holder" | "maintainer")}><option value="holder">Holder</option><option value="maintainer">Maintainer</option></select>}
        </div>
        {matches.map(match => <button type="button" className="btn btn-ghost justify-start" key={match.sharingId} onClick={async () => {
            setParticipants(current => [...current, {...match, role: secret.role === "owner" ? role : "holder"}]); setQuery(""); setMatches([]); setDirty(true);
        }}>{match.fullName} <span className="opacity-60">{match.emailHint}</span></button>)}
        <button type="button" className="btn btn-primary self-end" disabled={!dirty || participants.length < (secret.threshold ?? 0)} onClick={async () => {
            try {
                await setThresholdParticipants(secret.id!, editableParticipants(expectedParticipants), editableParticipants(participants));
                changed();
            } catch (caught) {
                if (!(caught instanceof ApiError) || caught.status !== 409) throw caught;
                setError("The participant list changed while you were editing it. Nothing was changed; review the refreshed list and try again.");
                const current = await getThresholdParticipants(secret.id!);
                setParticipants(current);
                setExpectedParticipants(current);
                setDirty(false);
            }
        }}>Apply access changes</button>
    </fieldset>;
};
