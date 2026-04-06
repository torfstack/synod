import type {AuthStatus} from "./authStatus.ts";

const postLoginHintKey = "synod.post-login-hint";

const pendingUnsealAuthStatus: AuthStatus = {
    isAuthenticated: true,
    isSetup: true,
    needsToUnseal: true,
};

export function getInitialAuthStatusFromPostLoginHint(storage: Storage): AuthStatus | null {
    if (storage.getItem(postLoginHintKey) !== "unseal") {
        return null;
    }

    return pendingUnsealAuthStatus;
}

export function setPostLoginHint(storage: Storage) {
    storage.setItem(postLoginHintKey, "unseal");
}

export function clearPostLoginHint(storage: Storage) {
    storage.removeItem(postLoginHintKey);
}
