import {useContext, useEffect, useState} from "react";
import {
    Button,
    ComboBox,
    ComboBoxStateContext,
    Input,
    Label,
    ListBox,
    ListBoxItem,
    Popover
} from "react-aria-components";
import {getSecretRecipients, revokeSecretAccess, searchShareRecipients, shareSecret, type ShareRecipient} from "../../util/api.ts";

type ShareSecretProps = {
    secretId: number;
    portalContainer: HTMLElement;
};

const OpenResults = ({open}: {open: boolean}) => {
    const state = useContext(ComboBoxStateContext);

    useEffect(() => {
        if (open) {
            state?.open(null, "input");
        } else {
            state?.close();
        }
    }, [open, state]);

    return null;
};

export const ShareSecret = ({secretId, portalContainer}: ShareSecretProps) => {
    const [query, setQuery] = useState("");
    const [recipients, setRecipients] = useState<ShareRecipient[]>([]);
    const [selected, setSelected] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [sharedWith, setSharedWith] = useState<ShareRecipient[]>([]);
    const [open, setOpen] = useState(false);

    useEffect(() => {
        getSecretRecipients(secretId).then(setSharedWith);
    }, [secretId]);

    useEffect(() => {
        if (query.trim().length < 2) {
            setRecipients([]);
            setOpen(false);
            return;
        }
        setOpen(true);
        const controller = new AbortController();
        const timeout = window.setTimeout(async () => {
            setLoading(true);
            try {
                setRecipients(await searchShareRecipients(query.trim(), controller.signal));
                setOpen(true);
            } catch (error) {
                if (!controller.signal.aborted) throw error;
            } finally {
                if (!controller.signal.aborted) setLoading(false);
            }
        }, 200);
        return () => {
            window.clearTimeout(timeout);
            controller.abort();
        };
    }, [query]);

    const submit = async () => {
        if (!selected) return;
        await shareSecret(secretId, selected);
        setSelected(null);
        setQuery("");
        setRecipients([]);
        setOpen(false);
        setSharedWith(await getSecretRecipients(secretId));
    };

    const revoke = async (recipient: ShareRecipient) => {
        await revokeSecretAccess(secretId, recipient.sharingId);
        setSharedWith(current => current.filter(item => item.sharingId !== recipient.sharingId));
    };

    return <fieldset className="fieldset border-t border-base-300 pt-4">
        <legend className="fieldset-legend">Share this secret</legend>
        <div className="flex gap-2 items-end">
            <ComboBox<ShareRecipient>
                className="grow"
                inputValue={query}
                onInputChange={value => {
                    setQuery(value);
                    setSelected(null);
                }}
                items={recipients}
                selectedKey={selected}
                onSelectionChange={key => {
                    setSelected(key?.toString() ?? null);
                    if (key) setOpen(false);
                }}
                onOpenChange={setOpen}
                menuTrigger="focus"
                allowsEmptyCollection
            >
                <OpenResults open={open}/>
                <Label className="label">Find a person</Label>
                <div className="join w-full">
                    <Input className="input input-bordered join-item w-full" placeholder="Name or exact email"/>
                    <Button className="btn join-item" aria-label="Show suggestions">⌄</Button>
                </div>
                <Popover
                    UNSTABLE_portalContainer={portalContainer}
                    className="z-50 w-[--trigger-width] rounded-box border border-base-300 bg-base-100 text-base-content shadow-xl"
                >
                    {loading && <div className="p-3 text-sm text-base-content/70">Searching…</div>}
                    {!loading && recipients.length === 0 &&
                        <div className="p-3 text-sm text-base-content/70">No matching users</div>}
                    <ListBox<ShareRecipient> className="max-h-64 overflow-y-auto p-1 outline-none">
                        {recipient => <ListBoxItem
                            id={recipient.sharingId}
                            textValue={recipient.fullName}
                            className="cursor-pointer rounded-lg p-3 text-base-content outline-none data-[focused]:bg-primary data-[focused]:text-primary-content data-[selected]:bg-primary data-[selected]:text-primary-content"
                        >
                            <div className="font-medium">{recipient.fullName}</div>
                            <div className="text-sm opacity-75">{recipient.emailHint}</div>
                        </ListBoxItem>}
                    </ListBox>
                </Popover>
            </ComboBox>
            <button type="button" className="btn btn-secondary w-24 shrink-0" disabled={!selected} onClick={submit}>
                Share
            </button>
        </div>
        {sharedWith.length > 0 && <div className="mt-3 flex flex-col gap-2">
            <span className="text-sm opacity-70">People with access</span>
            {sharedWith.map(recipient => <div className="flex items-center justify-between" key={recipient.sharingId}>
                <div>
                    <div>{recipient.fullName}</div>
                    <div className="text-sm opacity-60">{recipient.emailHint}</div>
                </div>
                <button type="button" className="btn btn-ghost btn-sm" onClick={() => revoke(recipient)}>Remove</button>
            </div>)}
        </div>}
    </fieldset>;
};
