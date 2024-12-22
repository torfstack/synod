package com.torfstack.kayvault.crypto

import assertk.assertAll
import assertk.assertThat
import assertk.assertions.*
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest

@SpringBootTest
class CryptServiceTest {

    @Autowired
    lateinit var cryptService: CryptService

    companion object {
        const val TAG_LENGTH = 16
        const val NONCE_LENGTH = 12
    }

    @Test
    fun `encrypt returns ciphertext different from ciphertext of correct length`() {
        val plaintext = "Testing plaintext".toByteArray()
        val ciphertext = cryptService.encrypt(plaintext)
        assertThat(ciphertext).isNotEqualTo(plaintext)
        assertThat(ciphertext).hasSize(plaintext.size + TAG_LENGTH + NONCE_LENGTH)
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