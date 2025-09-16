export interface AppConfig {
  backendSecretsUrl: string;
  backendAuthUrl: string;
  backendAuthStartUrl: string;
}

const secretPath = '/api/secrets';
const authPath = '/api/auth';

export const config: AppConfig = {
  backendSecretsUrl: secretPath,
  backendAuthUrl: authPath,
  backendAuthStartUrl: authPath + '/start',
};
