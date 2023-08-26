export class Secret {
    constructor(
        readonly key: string,
        readonly value: string,
        readonly url: string | undefined,
        readonly tags: string[] = []
    ) {}
}
