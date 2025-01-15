import urls from '$lib/config';

async function doAuth(token: string) {
	await fetch(urls.backendAuthUrl, {
		method: 'POST',
		mode: 'cors',
		cache: 'no-cache',
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	return token;
}

export default doAuth;
