import {useEffect, useState} from "react";
import {searchShareRecipients, type ShareRecipient} from "../../util/api.ts";

type Props = {
    enabled: boolean;
    setEnabled: (enabled: boolean) => void;
    threshold: number;
    setThreshold: (threshold: number) => void;
    recipients: ShareRecipient[];
    setRecipients: (recipients: ShareRecipient[]) => void;
};

export const ThresholdOptions = ({enabled, setEnabled, threshold, setThreshold, recipients, setRecipients}: Props) => {
    const [query, setQuery] = useState("");
    const [matches, setMatches] = useState<ShareRecipient[]>([]);

    useEffect(() => {
        if (query.trim().length < 2) {
            setMatches([]);
            return;
        }
        const controller = new AbortController();
        const timeout = window.setTimeout(() => searchShareRecipients(query.trim(), controller.signal)
            .then(results => setMatches(results.filter(result => !recipients.some(item => item.sharingId === result.sharingId))))
            .catch(error => { if (!controller.signal.aborted) throw error; }), 200);
        return () => { window.clearTimeout(timeout); controller.abort(); };
    }, [query, recipients]);

    const participantCount = recipients.length + 1;
    return <fieldset className="fieldset rounded-box border border-base-300 p-3">
        <label className="label cursor-pointer justify-start gap-3">
            <input type="checkbox" className="toggle toggle-primary" checked={enabled} onChange={event => setEnabled(event.target.checked)}/>
            Require several people to unlock this secret
        </label>
        {enabled && <div className="flex flex-col gap-3">
            <p className="text-sm opacity-70">You receive one share. Add people who should hold the other shares.</p>
            <input className="input input-bordered w-full" value={query} onChange={event => setQuery(event.target.value)} placeholder="Find by name or exact email"/>
            {matches.length > 0 && <div className="rounded-box border border-base-300 p-1">
                {matches.map(match => <button type="button" className="btn btn-ghost w-full justify-start" key={match.sharingId} onClick={() => {
                    setRecipients([...recipients, match]); setQuery(""); setMatches([]);
                }}>{match.fullName} <span className="opacity-60">{match.emailHint}</span></button>)}
            </div>}
            {recipients.map(recipient => <div className="flex items-center justify-between" key={recipient.sharingId}>
                <span>{recipient.fullName} <span className="text-sm opacity-60">{recipient.emailHint}</span></span>
                <button type="button" className="btn btn-ghost btn-sm" onClick={() => setRecipients(recipients.filter(item => item.sharingId !== recipient.sharingId))}>Remove</button>
            </div>)}
            <label className="label flex-col items-start">
                Shares required ({threshold} of {participantCount})
                <input type="range" className="range range-primary w-full" min={2} max={Math.max(2, participantCount)} value={Math.min(threshold, participantCount)} onChange={event => setThreshold(Number(event.target.value))}/>
            </label>
        </div>}
    </fieldset>;
};
