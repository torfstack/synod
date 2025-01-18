const localhost = 'http://localhost:4000';
const ip = 'http://192.168.178.52:8080';
const deployed = 'https://kayvault.torfstack.com/api';

const secretPath = '/secrets';
const authPath = '/auth';

const backendUrl = localhost;
const backendSecretsUrl = backendUrl + secretPath;
const backendAuthUrl = backendUrl + authPath;

export default {
	backendSecretsUrl,
	backendAuthUrl
};
