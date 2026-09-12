import {config} from "./config.ts";
import type {Secret} from './secret.ts';
import type {AuthStatus} from "./authStatus.ts";

export async function getAuth() {
    return apiFetchJson<AuthStatus>(config.backendAuthUrl, {
        method: 'GET',
    });
}

export async function deleteAuth() {
    return apiFetch(config.backendAuthUrl, {
        method: 'DELETE',
    });
}

export async function getSecrets() {
    return apiFetchJson<Secret[]>(config.backendSecretsUrl, {
        method: 'GET',
    });
}

export async function postSecret(secret: Secret) {
    const tags = secret.tags.toSorted((a, b) => a.localeCompare(b))
    return apiFetch(config.backendSecretsUrl, {
        method: 'POST',
        body: JSON.stringify({
            id: secret.id,
            value: secret.value,
            key: secret.key,
            url: secret.url,
            tags: tags,
        })
    });
}

export async function postThresholdSecret(secret: Secret, threshold: number, sharingIds: string[]) {
    return apiFetchJson<{id: number}>(`${config.backendSecretsUrl}/threshold`, {
        method: "POST",
        body: JSON.stringify({secret, threshold, sharingIds}),
    });
}

export type UnlockRequest = {
    id: number;
    secretId: number;
    requesterName: string;
    secretName: string;
    threshold: number;
    contributions: number;
    contributed: boolean;
    expiresAt: string;
};

export type UnlockResult = {
    ready: boolean;
    contributions: number;
    threshold: number;
    secret?: Secret;
};

export async function startUnlock(secretId: number) {
    return apiFetchJson<{id: number}>(`/api/secrets/${secretId}/unlocks`, {method: "POST"});
}

export async function getUnlockRequests() {
    return apiFetchJson<UnlockRequest[]>("/api/secrets/unlock-requests", {method: "GET"});
}

export async function contributeToUnlock(requestId: number) {
    return apiFetch(`/api/secrets/unlock-requests/${requestId}/contributions`, {method: "POST"});
}

export async function getUnlockResult(requestId: number) {
    return apiFetchJson<UnlockResult>(`/api/secrets/unlock-requests/${requestId}`, {method: "GET"});
}

export type ShareRecipient = {
    sharingId: string;
    fullName: string;
    emailHint: string;
};

export async function searchShareRecipients(search: string, signal: AbortSignal) {
    const url = `/api/users/lookup?find=${encodeURIComponent(search)}`;
    return apiFetchJson<ShareRecipient[]>(url, {method: "GET", signal});
}

export async function shareSecret(secretId: number, sharingId: string) {
    return apiFetch(`/api/secrets/${secretId}/shares`, {
        method: "POST",
        body: JSON.stringify({sharingId}),
    });
}

export async function getSecretRecipients(secretId: number) {
    return apiFetchJson<ShareRecipient[]>(`/api/secrets/${secretId}/shares`, {method: "GET"});
}

export async function revokeSecretAccess(secretId: number, sharingId: string) {
    return apiFetch(`/api/secrets/${secretId}/shares/${sharingId}`, {method: "DELETE"});
}

export async function postSetupPlain() {
    return apiFetch(config.backendSetupPlainUrl, {
        method: 'POST',
    });
}

export async function postSetupPassword(password: string) {
    return apiFetch(config.backendSetupPasswordUrl, {
        method: 'POST',
        body: JSON.stringify({
            password: password
        })
    });
}

export async function postUnsealWithPassword(password: string) {
    return apiFetch(config.backendUnsealUrl, {
        method: "POST",
        body: JSON.stringify({
            password: password
        })
    })
}

async function apiFetch(url: string, options: RequestInit = {}): Promise<Response> {
    const res = await fetch(url, {
        ...options,
        credentials: "include",
        mode: 'cors',
        cache: 'no-cache',
        headers: {
            "Content-Type": "application/json",
            ...(options.headers || {}),
        },
    });

    if (res.status === 401) {
        window.dispatchEvent(new Event("unauthorized"));
        throw new Error("Unauthorized");
    }

    if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText || "API request failed");
    }

    return res;
}

async function apiFetchJson<T>(url: string, options: RequestInit = {}): Promise<T> {
    const res = await apiFetch(url, options)
    return res.json() as Promise<T>;
}
