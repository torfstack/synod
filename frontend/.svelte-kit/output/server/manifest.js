export const manifest = {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.png"]),
	mimeTypes: {".png":"image/png"},
	_: {
		client: {"start":"_app/immutable/entry/start.47bd6137.js","app":"_app/immutable/entry/app.74ad2698.js","imports":["_app/immutable/entry/start.47bd6137.js","_app/immutable/chunks/index.bed43da4.js","_app/immutable/chunks/singletons.6acc8227.js","_app/immutable/entry/app.74ad2698.js","_app/immutable/chunks/index.bed43da4.js"],"stylesheets":[],"fonts":[]},
		nodes: [
			() => import('./nodes/0.js'),
			() => import('./nodes/1.js')
		],
		routes: [
			
		],
		matchers: async () => {
			
			return {  };
		}
	}
};
