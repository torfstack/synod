import urls from '$lib/config';
import type { Secret } from '$lib/secret';

async function getAuth() {
	return fetch(urls.backendAuthUrl, {
		method: 'GET',
		mode: 'cors',
		cache: 'no-cache',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});
}

async function deleteAuth() {
	return fetch(urls.backendAuthUrl, {
		method: 'DELETE',
		mode: 'cors',
		cache: 'no-cache',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
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
	deleteAuth,
	getAuth,
	getSecrets,
	postSecrets
};
