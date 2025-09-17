import type {Secret} from "../../util/secret.ts";
import {useEffect, useState} from "react";

export function showModal() {
    let dialog: any = document.getElementById('add_secret_modal')
    dialog.showModal()
}

export function closeModal() {
    let dialog: any = document.getElementById('add_secret_modal')
    dialog.close()
}

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
        return name.length > 0 && secret.length > 0 && url.length > 0
    }

    function removeTag(tag: string): () => void {
        return () => setTags(prevState => prevState.filter(t => t !== tag));
    }

    return <dialog id="add_secret_modal" className="modal">
        <div className="modal-box">
            <h3 className="font-bold text-lg pb-4">Add Secret</h3>
            <div className="flex flex-col gap-2">
                <label className="floating-label">
                    <span>Name</span>
                    <input type="text" placeholder="Name" value={name}
                           onChange={(e) => setName(e.target.value)}
                           className="input input-bordered w-full max-w-xs"/>
                </label>

                <label className="floating-label">
                    <span>Secret</span>
                    <input type="text" placeholder="Secret" value={secret}
                           onChange={(e) => setSecret(e.target.value)}
                           className="input input-bordered w-full max-w-xs"/>
                </label>
                <label className="floating-label">
                    <span>URL</span>
                    <input type="text" placeholder="URL" value={url}
                           onChange={(e) => setUrl(e.target.value)}
                           className="input input-bordered w-full max-w-xs"/>
                </label>
                <label className="floating-label">
                    <span>Add Tag</span>
                    <input type="text" placeholder="Tags" value={tag}
                           disabled={tags.length >= 3}
                           onChange={(e) => setTag(e.target.value)}
                           onKeyDown={(e) => {
                               if (e.key == "Enter") {
                                   setTags([...tags, tag])
                                   setTag("")
                               }
                           }}
                           className="input input-bordered w-full max-w-xs"/>
                </label>
                <div className="flex flex-row gap-2">
                    {tags.map((tag) => (
                        <span className="badge btn" onClick={removeTag(tag)}>{tag}</span>
                    ))}
                </div>
            </div>
            <div className="modal-action">
                <form method="dialog">
                    <button className="btn btn-primary" onClick={onSubmit}>Submit</button>
                </form>
            </div>
        </div>
        <form method="dialog" className="modal-backdrop">
            <button>close</button>
        </form>
    </dialog>
}