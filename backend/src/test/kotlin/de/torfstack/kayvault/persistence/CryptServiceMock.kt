package de.torfstack.kayvault.persistence

import org.springframework.context.annotation.Primary
import org.springframework.stereotype.Service

@Service
@Primary
class CryptServiceMock : CryptService {
    var encryptGotCalled = false
    var decryptGotCalled = false

    override fun encrypt(plaintext: String): String {
        encryptGotCalled = true
        return plaintext
    }

    override fun decrypt(ciphertext: String): String {
        decryptGotCalled = true
        return ciphertext
    }
}
