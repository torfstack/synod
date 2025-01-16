import urls from '$lib/config';
import type { Secret } from '$lib/secret';

async function postAuth(token: string) {
	return fetch(urls.backendAuthUrl, {
		method: 'POST',
		mode: 'cors',
		cache: 'no-cache',
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
}

async function getSecrets() {
	return fetch(urls.backendSecretsUrl, {
		method: 'GET',
		mode: 'cors',
		cache: 'no-cache',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});
}

async function postSecrets(secret: Secret) {
	return fetch(urls.backendSecretsUrl, {
		method: 'POST',
		mode: 'cors',
		cache: 'no-cache',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include',
		body: JSON.stringify({
			value: secret.value,
			key: secret.key,
			url: secret.url
		})
	});
}

export default {
	postAuth,
	getSecrets,
	postSecrets
};
