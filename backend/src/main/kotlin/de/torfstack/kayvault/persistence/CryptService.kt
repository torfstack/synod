package de.torfstack.kayvault.persistence

interface CryptService {
    fun encrypt(plaintext: ByteArray): ByteArray
    fun decrypt(ciphertext: ByteArray): ByteArray
}
