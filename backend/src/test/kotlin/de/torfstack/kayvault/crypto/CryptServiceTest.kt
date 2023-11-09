package de.torfstack.kayvault.crypto

import assertk.assertAll
import assertk.assertThat
import assertk.assertions.isEqualTo
import assertk.assertions.isFailure
import assertk.assertions.isNotEqualTo
import assertk.assertions.isSuccess
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest

@SpringBootTest
class CryptServiceTest {

    @Autowired
    lateinit var cryptService: CryptService

    @Test
    fun `encrypt returns ciphertext different from ciphertext`() {
        val plaintext = "Testing plaintext".toByteArray()
        val ciphertext = cryptService.encrypt(plaintext)
        assertThat(ciphertext).isNotEqualTo(plaintext)
    }

    @Test
    fun `decrypt returns plaintext that was encrypted`() {
        val plaintext = "Testing plaintext".toByteArray()
        val ciphertext = cryptService.encrypt(plaintext)
        val decrypted = cryptService.decrypt(ciphertext)
        assertThat(decrypted).isEqualTo(plaintext)
    }

    @Test
    fun `modified ciphertext can not be decrypted`() {
        val plaintext = "Testing plaintext".toByteArray()
        val ciphertext = cryptService.encrypt(plaintext)
        assertAll {
            for (i in ciphertext.indices) {
                val modifiedCiphertext = ciphertext.copyOf()
                modifiedCiphertext[i] = modifiedCiphertext[i].inc()
                assertThat {
                    cryptService.decrypt(modifiedCiphertext)
                }.isFailure()
            }
        }
    }
}