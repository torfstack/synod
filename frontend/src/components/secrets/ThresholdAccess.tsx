import {useEffect, useState} from "react";
import type {Secret} from "../../util/secret.ts";
import {addThresholdParticipant, getThresholdParticipants, removeThresholdParticipant, searchShareRecipients, setThresholdParticipantRole, type ShareRecipient, type ThresholdParticipant} from "../../util/api.ts";

export const ThresholdAccess = ({secret, changed}: {secret: Secret; changed: () => void}) => {
    const [participants, setParticipants] = useState<ThresholdParticipant[]>([]);
    const [query, setQuery] = useState("");
    const [matches, setMatches] = useState<ShareRecipient[]>([]);
    const [role, setRole] = useState<"holder" | "maintainer">("holder");
    useEffect(() => { getThresholdParticipants(secret.id!).then(setParticipants); }, [secret.id]);
    useEffect(() => {
        if (query.trim().length < 2) { setMatches([]); return; }
        const controller = new AbortController();
        const timeout = window.setTimeout(() => searchShareRecipients(query.trim(), controller.signal).then(setMatches), 200);
        return () => { window.clearTimeout(timeout); controller.abort(); };
    }, [query]);
    return <fieldset className="fieldset border-t border-base-300 pt-3">
        <legend className="fieldset-legend">Threshold access</legend>
        <p className="text-sm opacity-70">Changing access rotates the encryption key and locks the secret again.</p>
        {participants.map(participant => <div className="flex items-center justify-between gap-2" key={participant.sharingId}>
            <span>{participant.fullName} <span className="badge badge-sm">{participant.role}</span></span>
            {participant.role !== "owner" && <div className="flex gap-2">
                {secret.role === "owner" && <select className="select select-sm" value={participant.role} onChange={async event => {
                    await setThresholdParticipantRole(secret.id!, {sharingId: participant.sharingId, role: event.target.value as "holder" | "maintainer"}); changed();
                }}><option value="holder">Holder</option><option value="maintainer">Maintainer</option></select>}
                {(secret.role === "owner" || participant.role === "holder") && <button type="button" className="btn btn-ghost btn-sm" onClick={async () => { await removeThresholdParticipant(secret.id!, participant.sharingId); changed(); }}>Remove</button>}
            </div>}
        </div>)}
        <div className="flex gap-2">
            <input className="input grow" value={query} onChange={event => setQuery(event.target.value)} placeholder="Find a person"/>
            {secret.role === "owner" && <select className="select" value={role} onChange={event => setRole(event.target.value as "holder" | "maintainer")}><option value="holder">Holder</option><option value="maintainer">Maintainer</option></select>}
        </div>
        {matches.map(match => <button type="button" className="btn btn-ghost justify-start" key={match.sharingId} onClick={async () => {
            await addThresholdParticipant(secret.id!, {sharingId: match.sharingId, role: secret.role === "owner" ? role : "holder"}); changed();
        }}>{match.fullName} <span className="opacity-60">{match.emailHint}</span></button>)}
    </fieldset>;
};
