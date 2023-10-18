package de.torfstack.kayvault.crypto

interface KeyRetriever {
    fun key(): ByteArray
}