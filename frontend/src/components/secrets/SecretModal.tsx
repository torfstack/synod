import type {Secret} from "../../util/secret.ts";
import React, {useEffect, useRef, useState} from "react";
import {Eye, EyeSlash} from "../../icons/Eye.tsx";
import {ShareSecret} from "./ShareSecret.tsx";

interface SecretModalProps {
    handleSecret: (s: Secret) => Promise<void>;
    existingSecret?: Secret;
    isOpen: boolean;
    closeModal: () => void;
}

export const SecretModal: React.FC<SecretModalProps> = ({handleSecret, existingSecret, isOpen, closeModal}) => {
	const readOnly = existingSecret?.owned === false;
    const [name, setName] = useState(existingSecret?.key ?? "")
    const [secret, setSecret] = useState(existingSecret?.value ?? "")
    const [url, setUrl] = useState(existingSecret?.url ?? "")
    const [tags, setTags] = useState<string[]>(existingSecret?.tags ?? [])
    const [tag, setTag] = useState("")
    const [passwordVisible, setPasswordVisible] = useState(false)
    const dialogRef = useRef<HTMLDialogElement>(null)
    const [sharePortal, setSharePortal] = useState<HTMLDivElement | null>(null)

    useEffect(() => {
        const dialog = dialogRef.current;
        if (!dialog) return;
        const handleClose = () => {
            closeModal();
        };
        dialog.addEventListener("close", handleClose);
        if (isOpen && !dialog.open) {
            dialog.showModal();
        } else if (!isOpen && dialog.open) {
            dialog.close();
        }
        return () => {
            dialog.removeEventListener("close", handleClose);
        };
    }, [isOpen, closeModal]);

    useEffect(() => {
        if (!isOpen) return
        setName(existingSecret?.key ?? "")
        setSecret(existingSecret?.value ?? "")
        setUrl(existingSecret?.url ?? "")
        setTags(existingSecret?.tags ?? [])
        setTag("")
        setPasswordVisible(false)
    }, [existingSecret, isOpen]);

    async function onSubmit() {
        if (!checkInput()) {
            return
        }
        const s: Secret = {
            id: existingSecret?.id,
            key: name,
            value: secret,
            url: url,
            tags: tags,
        }
        await handleSecret(s)
        closeModal()
    }

    function checkInput(): boolean {
        return name.length > 0 && secret.length > 0
    }

    function removeTag(tag: string): () => void {
        const tagFilter = (t: string) => t !== tag;
        return () => setTags(prevState => prevState.filter(tagFilter));
    }

    function togglePassword() {
        const isPassword = document.getElementById("input-password")?.getAttribute("type") == "password"
        setPasswordVisible(isPassword)
    }

    const title = readOnly ? "Shared secret" : existingSecret ? "Edit secret" : "Add secret";

    return (
        <dialog ref={dialogRef} className="modal">
            <div className="modal-box w-11/12 max-w-2xl overflow-visible">
                <form>
                    <fieldset className="fieldset">
                        <legend className="fieldset-legend">{title}</legend>
                        <div className="flex flex-col gap-4">
                            {readOnly && <div role="status" className="alert alert-info py-2">
                                <span>This secret was shared with you. Only its owner can modify it or manage access.</span>
                            </div>}
                            <label className="input w-full">
                                Name
                                <input type="text" placeholder="MyNewSecret" value={name}
                                       onChange={(e) => setName(e.target.value)}
                                       className="grow validator" minLength={1} required
                                       title="Can not be empty" readOnly={readOnly}/>
                            </label>
                            <label className="input w-full">
                                Secret
                                <input id="input-password" type={passwordVisible ? "text" : "password"}
                                       placeholder="*****"
                                       value={secret}
                                       onChange={(e) => setSecret(e.target.value)}
                                       className="grow validator" minLength={1} required
                                       title="Can not be empty" readOnly={readOnly}/>
                                <button type="button" onClick={togglePassword} className="btn btn-ghost btn-xs">
                                    {passwordVisible ? <Eye/> : <EyeSlash/>}
                                </button>
                            </label>
                            <label className="input w-full">
                                URL
                                <input type="text" placeholder="https://example.com" value={url}
                                       onChange={(e) => setUrl(e.target.value)}
                                       className="grow" readOnly={readOnly}/>
                                <span className="badge badge-neutral badge-xs">Optional</span>
                            </label>
                            <label className="input w-full">
                                Add Tag
                                <input type="text" placeholder="example" value={tag}
                                       onChange={(e) => {
                                           const lastChar = e.target.value.charAt(e.target.value.length - 1)
                                           if (lastChar == " " || lastChar == ",") {
                                               return
                                           }
                                           setTag(e.target.value)
                                       }}
                                       onKeyDown={(e) => {
                                           if (e.code == "Enter") {
                                               if (tag.length == 0) return
                                               setTags([...tags, tag])
                                               setTag("")
                                           }
                                       }}
                                       disabled={readOnly || tags.length >= 3} className="grow"/>
                                <span className="badge badge-neutral badge-xs">&lt;4</span>
                                <kbd className="kbd kbd-sm">↵</kbd>
                            </label>
                            <div className="flex flex-col gap-2">
                                {tags.map((tag) => (
                                    <div key={tag} className={`badge badge-neutral ${readOnly ? "" : "btn"}`}
                                         onClick={readOnly ? undefined : removeTag(tag)}>
                                        <p className="truncate max-w-56">{tag}</p>
                                    </div>
                                ))}
                            </div>
                            <div className="modal-action">
                                {!readOnly && <button type="button" className="btn btn-primary" onClick={onSubmit}>Submit</button>}
                                {readOnly && <button type="button" className="btn" onClick={closeModal}>Close</button>}
                            </div>
                            <div ref={setSharePortal}>
                                {existingSecret?.id && existingSecret.owned !== false && sharePortal &&
                                    <ShareSecret secretId={existingSecret.id} portalContainer={sharePortal}/>}
                            </div>
                        </div>
                    </fieldset>
                </form>
            </div>
            <form method="dialog" className="modal-backdrop">
                <button>close</button>
            </form>
        </dialog>
    )
}
