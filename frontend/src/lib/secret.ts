export class Secret {
	constructor(
		readonly id: number | undefined,
		readonly key: string,
		readonly value: string,
		readonly url: string | undefined,
		readonly tags: string[] = []
	) {}
}
