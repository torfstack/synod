package com.torfstack.kayvault.persistence

import jakarta.persistence.*

@Entity(name = "secrets")
class SecretEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    var secretId: Long = 0

    @Column(name = "secretValue")
    @Lob
    var secretValue: ByteArray = ByteArray(0)

    @Column(name = "secretKey")
    var secretKey: String = ""

    @Column(name = "secretUrl")
    var secretUrl: String? = null

    @Column(name = "forUser")
    var forUser: String = ""
}