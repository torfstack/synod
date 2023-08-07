package de.torfstack.kayvault.validation

interface TokenValidator {
    fun validate(token: String): VerificationResult

    sealed class VerificationResult

    class ValidVerification(val user: String): VerificationResult()
    class InvalidVerification(val ex: Exception): VerificationResult()
}