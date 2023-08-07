package de.torfstack.kayvault.controller

import de.torfstack.kayvault.persistence.SecretModel
import de.torfstack.kayvault.persistence.SecretService
import de.torfstack.kayvault.validation.TokenValidator
import mu.KotlinLogging
import org.springframework.web.bind.annotation.*

val Log = KotlinLogging.logger{}

@RestController
@CrossOrigin
class SecretController(val secretService: SecretService, val tokenVerifier: TokenValidator) {

    @GetMapping("secret")
    fun getSecret(@RequestHeader authorization: String): List<SecretRequestEntity> {
        val user = userFromHeader(authorization)
        Log.info { "returning secrets for user $user" }
        return secretsForUser(user)
    }

    @PostMapping("secret")
    fun postSecret(@RequestHeader authorization: String, @RequestBody entity: SecretRequestEntity): List<SecretRequestEntity> {
        val user = userFromHeader(authorization)
        Log.info { "adding a secret for user $user" }
        secretService.addSecretForUser(user, entity.toModel())
        return secretsForUser(user)
    }

    private fun userFromHeader(header: String): String {
        val authorization = header.removePrefix("Bearer").trim()
        when (val result = tokenVerifier.validate(authorization)) {
            is TokenValidator.InvalidVerification -> {
                Log.warn { "token could not be verified" }
                throw IllegalArgumentException("token could not be verified", result.ex)
            }
            is TokenValidator.ValidVerification -> {
                Log.debug { "validated token successfully and got user ${result.user}" }
                return result.user
            }
        }
    }

    private fun secretsForUser(user: String): List<SecretRequestEntity> {
        return secretService.secretsForUser(user)
            .map {
                SecretRequestEntity(
                    key = it.secretKey,
                    value = it.secretValue,
                    url = it.secretUrl,
                    notes = ""
                )
            }
    }

    data class SecretRequestEntity(
        val key: String,
        val value: String,
        val notes: String?,
        val url: String?
    ) {
        fun toModel(): SecretModel {
            return SecretModel(
                secretValue = value,
                secretUrl = url,
                secretKey = key
            )
        }
    }
}
