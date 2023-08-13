export const manifest = {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.png"]),
	mimeTypes: {".png":"image/png"},
	_: {
		client: {"start":"_app/immutable/entry/start.0fba647a.js","app":"_app/immutable/entry/app.932dbc24.js","imports":["_app/immutable/entry/start.0fba647a.js","_app/immutable/chunks/index.bed43da4.js","_app/immutable/chunks/singletons.f13267be.js","_app/immutable/entry/app.932dbc24.js","_app/immutable/chunks/index.bed43da4.js"],"stylesheets":[],"fonts":[]},
		nodes: [
			() => import('./nodes/0.js'),
			() => import('./nodes/1.js'),
			() => import('./nodes/2.js')
		],
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			}
		],
		matchers: async () => {
			
			return {  };
		}
	}
};
