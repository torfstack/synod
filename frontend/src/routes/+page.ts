import api from '$lib/api';

export const load: () => Promise<any> = async () => {
	const res = await api.getAuth();
	if (res.status == 200) {
		console.log('Authenticated');
		return {
			isAuthenticated: true
		};
	}
	console.log('Not authenticated');
	return {
		isAuthenticated: false
	};
};
