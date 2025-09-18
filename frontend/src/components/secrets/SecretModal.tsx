import type {Secret} from "../../util/secret.ts";
import {useEffect, useState} from "react";
import {closeModal} from "../../util/modal.ts";

type SecretModalProps = {
    handleSecret: (s: Secret) => Promise<void>
    existingSecret?: Secret
}

export const SecretModal = (props: SecretModalProps) => {
    const [name, setName] = useState(props.existingSecret?.key ?? "")
    const [secret, setSecret] = useState(props.existingSecret?.value ?? "")
    const [url, setUrl] = useState(props.existingSecret?.url ?? "")
    const [tags, setTags] = useState<string[]>(props.existingSecret?.tags ?? [])
    const [tag, setTag] = useState("")

    useEffect(() => {
        setName(props.existingSecret?.key ?? "")
        setSecret(props.existingSecret?.value ?? "")
        setUrl(props.existingSecret?.url ?? "")
        setTags(props.existingSecret?.tags ?? [])
    }, [props.existingSecret]);

    async function onSubmit() {
        if (!checkInput()) {
            return
        }
        const s: Secret = {
            id: props.existingSecret?.id,
            key: name,
            value: secret,
            url: url,
            tags: tags,
        }
        await props.handleSecret(s)
        closeModal()
    }

    function checkInput(): boolean {
        let valid = true
        if (name.length == 0) {
            valid = false
        }
        if (secret.length == 0) {
            valid = false
        }
        return valid
    }

    function removeTag(tag: string): () => void {
        return () => setTags(prevState => prevState.filter(t => t !== tag));
    }

    return <dialog id="add_secret_modal" className="modal">
        <div className="modal-box">
            <form>

                <fieldset className="fieldset">
                    <legend className="fieldset-legend">Add Secret</legend>
                    <div className="flex flex-col gap-4">
                        <label className="input">
                            Name
                            <input type="text" placeholder="MyNewSecret" value={name}
                                   onChange={(e) => setName(e.target.value)}
                                   className="grow validator" minLength={1} required title="Can not be empty"/>
                            <p id="name-input-error" className="validator-hint">Can not be empty</p>
                        </label>
                        <label className="input">
                            Secret
                            <input type="password" placeholder="*****" value={secret}
                                   onChange={(e) => setSecret(e.target.value)}
                                   className="grow validator" minLength={1} required title="Can not be empty"/>
                            <p id="secret-input-error" className="validator-hint">Can not be empty</p>
                        </label>
                        <label className="input">
                            URL
                            <input type="text" placeholder="https://example.com" value={url}
                                   onChange={(e) => setUrl(e.target.value)}
                                   className="grow"/>
                            <span className="badge badge-neutral badge-xs">Optional</span>
                        </label>
                        <label className="input">
                            Add Tag
                            <input type="text" placeholder="example" value={tag}
                                   onChange={(e) => setTag(e.target.value)}
                                   onKeyDown={(e) => {
                                       if (e.code == "Space") {
                                           setTags([...tags, tag])
                                           setTag("")
                                       }
                                   }}
                                   disabled={tags.length >= 3} className="grow"/>
                            <span className="badge badge-neutral badge-xs">&lt;4</span>
                            <kbd className="kbd kbd-sm">␣</kbd>
                        </label>
                        <div className="flex flex-row gap-2">
                            {tags.map((tag) => (
                                <span className="badge badge-neutral btn" onClick={removeTag(tag)}>{tag}</span>
                            ))}
                        </div>
                        <div className="modal-action">
                            <button className="btn btn-primary" onClick={onSubmit}>Submit</button>
                        </div>
                    </div>
                </fieldset>
            </form>
        </div>
        <form method="dialog" className="modal-backdrop">
            <button>close</button>
        </form>
    </dialog>
}