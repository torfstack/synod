import {describe, expect, test} from "bun:test";
import {
    clearPostLoginHint,
    getInitialAuthStatusFromPostLoginHint,
    setPostLoginHint,
} from "./postLoginHint.ts";

const createStorage = (): Storage => {
    const values = new Map<string, string>();

    return {
        get length() {
            return values.size;
        },
        clear() {
            values.clear();
        },
        getItem(key: string) {
            return values.get(key) ?? null;
        },
        key(index: number) {
            return Array.from(values.keys())[index] ?? null;
        },
        removeItem(key: string) {
            values.delete(key);
        },
        setItem(key: string, value: string) {
            values.set(key, value);
        },
    };
};

describe("postLoginHint", () => {
    test("returns null when no hint is set", () => {
        const storage = createStorage();

        expect(getInitialAuthStatusFromPostLoginHint(storage)).toBeNull();
    });

    test("returns an unseal auth status when the login hint is present", () => {
        const storage = createStorage();

        setPostLoginHint(storage);

        expect(getInitialAuthStatusFromPostLoginHint(storage)).toEqual({
            isAuthenticated: true,
            isSetup: true,
            needsToUnseal: true,
        });
    });

    test("clears the login hint", () => {
        const storage = createStorage();

        setPostLoginHint(storage);
        clearPostLoginHint(storage);

        expect(getInitialAuthStatusFromPostLoginHint(storage)).toBeNull();
    });
});
