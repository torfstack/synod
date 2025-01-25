import api from '$lib/api';

export const load: () => Promise<any> = async () => {
	const res = await api.getAuth();
	if (res.status == 200) {
		return {
			isAuthenticated: true
		};
	}
	return {
		isAuthenticated: false
	};
};
