export interface AppConfig {
    backendSecretsUrl: string
    backendAuthUrl: string
    backendAuthStartUrl: string
    backendSetupUrl: string
    backendSetupPlainUrl: string
    backendSetupPasswordUrl: string
    backendUnsealUrl: string
}

const secretPath = '/api/secrets'
const authPath = '/api/auth'
const setupPath = '/api/setup'

export const config: AppConfig = {
    backendSecretsUrl: secretPath,
    backendAuthUrl: authPath,
    backendAuthStartUrl: authPath + '/start',
    backendSetupUrl: setupPath,
    backendSetupPlainUrl: setupPath + '/plain',
    backendSetupPasswordUrl: setupPath + '/password',
    backendUnsealUrl: setupPath + '/unseal',
};
