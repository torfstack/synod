package com.torfstack.kayvault.persistence

data class SecretModel(
    val secretValue: String,
    val secretKey: String,
    val secretUrl: String?
)