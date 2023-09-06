import { PUBLIC_API_AUTHORITY } from '$env/static/public';

const secretPath = '/secret';

const backendUrl = PUBLIC_API_AUTHORITY;
const backendSecretsUrl = backendUrl + secretPath;
export default backendSecretsUrl;
