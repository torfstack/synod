package de.torfstack.kayvault.persistence

data class SecretModel(
    val secretValue: String,
    val secretKey: String,
    val secretUrl: String?
) {
    companion object {
        fun fromEntity(entity: SecretEntity): SecretModel {
            return SecretModel(
                secretValue = entity.secretValue,
                secretKey = entity.secretKey,
                secretUrl = entity.secretUrl
            )
        }
    }
}
