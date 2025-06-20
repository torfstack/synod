import { loadConfig } from '$lib/config';

export const ssr = false;
export const prerender = true;

export async function load() {
	await loadConfig();
	return {};
}
