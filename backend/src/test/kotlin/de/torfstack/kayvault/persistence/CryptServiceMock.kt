package de.torfstack.kayvault.persistence

import de.torfstack.kayvault.crypto.CryptService
import org.springframework.context.annotation.Primary
import org.springframework.stereotype.Service

@Service
@Primary
class CryptServiceMock : CryptService {
    var encryptGotCalled = false
    var decryptGotCalled = false

    override fun encrypt(plaintext: ByteArray): ByteArray {
        encryptGotCalled = true
        return plaintext
    }

    override fun decrypt(ciphertext: ByteArray): ByteArray {
        decryptGotCalled = true
        return ciphertext
    }
}
