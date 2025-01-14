package auth

type SessionService interface {
    // CreateSession creates a new session for the given user.
    CreateSession(user int64) (*Session, error)

    // GetSession returns the session for the given token.
    GetSession(token string) (*Session, error)

    // DeleteSession deletes the session for the given token.
    DeleteSession(token string) error
}

type Session string

