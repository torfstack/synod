export interface AppConfig {
	backendSecretsUrl: string;
	backendAuthUrl: string;
	backendAuthStartUrl: string;
}

export let config: AppConfig;

const secretPath = '/api/secrets';
const authPath = '/api/auth';

export async function loadConfig(): Promise<void> {
	config = {
		backendSecretsUrl: secretPath,
		backendAuthUrl: authPath,
		backendAuthStartUrl: authPath + '/start',
	};
}
