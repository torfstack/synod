import {config} from "./config.ts";
import type {Secret} from './secret.ts';

export async function getAuth() {
  return fetch(config.backendAuthUrl, {
    method: 'GET',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include'
  });
}

export async function deleteAuth() {
  return fetch(config.backendAuthUrl, {
    method: 'DELETE',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include'
  });
}

export async function getSecrets() {
  return fetch(config.backendSecretsUrl, {
    method: 'GET',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include'
  });
}

export async function postSecret(secret: Secret) {
  const tags = secret.tags.toSorted((a, b) => a.localeCompare(b))
  return fetch(config.backendSecretsUrl, {
    method: 'POST',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify({
      id: secret.id,
      value: secret.value,
      key: secret.key,
      url: secret.url,
      tags: tags,
    })
  });
}
