package de.torfstack.kayvault.persistence

interface CryptService {
    fun encrypt(plaintext: String): String
    fun decrypt(ciphertext: String): String
}
