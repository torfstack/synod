package com.torfstack.kayvault.crypto

interface CryptService {
    fun encrypt(plaintext: ByteArray): ByteArray
    fun decrypt(ciphertext: ByteArray): ByteArray
}
