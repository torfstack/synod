# Security model

Synod trusts the application server. This applies whether Synod is operated as a
hosted service or installed and operated by a customer.

The server may access plaintext secrets and decrypted key material while it is
performing an authorized operation. Password-protected private keys are encrypted
at rest and held decrypted in server memory only while the user's session is
unsealed. Forgetting the password permanently prevents recovery of that private
key; Synod does not currently provide a recovery mechanism.

Secrets and private keys are encrypted at rest to protect data held in the
database and its backups. TLS protects requests and responses in transit. Synod's
server-side encryption is not end-to-end encryption and does not protect against
a compromised or malicious application server.

For individually shared secrets, Synod encrypts the secret payload with a random
data-encryption key. It encrypts a copy of that key for each authorized user's
public key. Data-encryption keys and decrypted private keys are not sent to the
browser. The server unwraps keys and decrypts secrets after checking the current
user's authorization.

Removing access prevents future retrieval through Synod. It cannot erase a
plaintext secret that a recipient has already viewed or copied.

Future threshold sharing will split a data-encryption key into encrypted shares.
The trusted server will collect the required number of explicitly approved
shares, reconstruct the key in memory, and decrypt the secret. Threshold approval
is an application authorization rule, not a defense against the trusted server.

Sensitive plaintext, decrypted keys, and key shares must not be logged, included
in errors, or persisted in caches, sessions, audit records, or job payloads beyond
the private key intentionally retained by an unsealed session. Temporary key
material must be scoped to the authorized operation and discarded promptly.
