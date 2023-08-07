package de.torfstack.kayvault.persistence

import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface SecretRepository : JpaRepository<SecretEntity, Long> {
    fun findByForUser(forUser: String): List<SecretEntity>
}