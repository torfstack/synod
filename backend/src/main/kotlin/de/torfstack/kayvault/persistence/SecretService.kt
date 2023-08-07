package de.torfstack.kayvault.persistence

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.stereotype.Service

@Service
class SecretService @Autowired constructor(val repo: SecretRepository) {
    fun secretsForUser(user: String): List<SecretModel> {
        return repo.findByForUser(user).map {
            SecretModel.fromEntity(it)
        }
    }

    fun addSecretForUser(user: String, model: SecretModel) {
        repo.save(
            SecretEntity().also {
                it.secretValue = model.secretValue
                it.secretKey = model.secretKey
                it.secretUrl = model.secretUrl
                it.forUser = user
            }
        )
    }
}