package com.torfstack.kayvault.persistence

import org.springframework.data.jpa.repository.JpaRepository

interface SecretRepository : JpaRepository<SecretEntity, Long> {
    fun findByForUser(forUser: String): List<SecretEntity>
}