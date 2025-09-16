import {useState} from "react";
import type {Secret} from "../../util/secret.ts";

export function showModal() {
  let dialog: any = document.getElementById('add_secret_modal');
  dialog.showModal();
}

export function closeModal() {
  let dialog: any = document.getElementById('add_secret_modal');
  dialog.open = false
}

type SecretModalProps = {
  handleSecret: (s: Secret) => Promise<void>
  selectedSecret?: Secret
}

export const SecretModal = (props: SecretModalProps) => {
  const [name, setName] = useState(props.selectedSecret?.key ?? "")
  const [secret, setSecret] = useState(props.selectedSecret?.value ?? "")
  const [url, setUrl] = useState(props.selectedSecret?.url ?? "")

  async function onSubmit() {
    if (!checkInput()) {
      return
    }
    const s: Secret = {
      id: props.selectedSecret?.id,
      key: name,
      value: secret,
      url: url,
      tags: []
    }
    await props.handleSecret(s)
    closeModal()
  }

  function checkInput(): boolean {
    return name.length > 0 && secret.length > 0 && url.length > 0
  }

  return <dialog id="add_secret_modal" className="modal">
    <div className="modal-box">
      <h3 className="font-bold text-lg">New Secret</h3>
      <div className="flex flex-col gap-2">
        <input type="text" placeholder="Name" value={name}
               onChange={(e) => setName(e.target.value)}
               className="input input-bordered w-full max-w-xs"/>
        <input type="text" placeholder="Secret" value={secret}
               onChange={(e) => setSecret(e.target.value)}
               className="input input-bordered w-full max-w-xs"/>
        <input type="text" placeholder="URL" value={url}
               onChange={(e) => setUrl(e.target.value)}
               className="input input-bordered w-full max-w-xs"/>
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