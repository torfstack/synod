export function showModal() {
    getModal()?.showModal()
}

export function closeModal() {
    getModal()?.close()
}

function getModal(): HTMLDialogElement | null {
    const dialog = document.getElementById('add_secret_modal')
    if (dialog instanceof HTMLDialogElement) {
        return dialog
    }
    return null;
}

