package de.torfstack.kayvault.crypto

import org.springframework.stereotype.Service
import java.nio.ByteBuffer
import java.security.SecureRandom
import javax.crypto.Cipher
import javax.crypto.CipherOutputStream
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec
import javax.crypto.spec.SecretKeySpec

@Service
class CryptServiceImpl(private val keyRetriever: KeyRetriever) : CryptService {

    override fun encrypt(plaintext: ByteArray): ByteArray {
        val iv = ByteArray(12)
        SecureRandom().nextBytes(iv)
        val encrypted = cipher(iv, Cipher.ENCRYPT_MODE).doFinal(plaintext)
        val complete = ByteArray(12+encrypted.size)
        iv.copyInto(complete)
        encrypted.copyInto(complete, 12)
        return complete
    }

    override fun decrypt(ciphertext: ByteArray): ByteArray {
        val iv = ciphertext.copyOfRange(0, 12)
        val toDecrypt = ciphertext.copyOfRange(12, ciphertext.size)
        return cipher(iv, Cipher.DECRYPT_MODE).doFinal(toDecrypt)
    }

    private fun cipher(iv: ByteArray, mode: Int): Cipher {
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        val spec = GCMParameterSpec(128, iv)
        cipher.init(mode, key(), spec)
        return cipher
    }

    private fun key(): SecretKey {
        return SecretKeySpec(keyRetriever.key(), "AES")
    }
}