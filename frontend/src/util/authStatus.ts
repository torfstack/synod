export type AuthStatus = {
    isAuthenticated: boolean,
    isSetup: boolean,
    needsToUnseal: boolean,
}

export const EmptyAuthStatus: AuthStatus = {isAuthenticated: false, isSetup: false, needsToUnseal: false}