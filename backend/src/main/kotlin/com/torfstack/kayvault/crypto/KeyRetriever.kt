package com.torfstack.kayvault.crypto

interface KeyRetriever {
    fun key(): ByteArray
}