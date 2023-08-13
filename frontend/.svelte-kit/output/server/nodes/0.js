import * as universal from '../entries/pages/_layout.js';

export const index = 0;
export const component = async () => (await import('../entries/pages/_layout.svelte.js')).default;
export { universal };
export const universal_id = "src/routes/+layout.js";
export const imports = ["_app/immutable/nodes/0.3100100b.js","_app/immutable/chunks/index.bed43da4.js"];
export const stylesheets = ["_app/immutable/assets/0.7dd25a13.css"];
export const fonts = [];
