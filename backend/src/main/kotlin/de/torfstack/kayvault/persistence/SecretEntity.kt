package de.torfstack.kayvault.persistence

import jakarta.persistence.*

@Entity(name = "secrets")
class SecretEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    var secretId: Long = 0

    @Column(name = "secretValue")
    var secretValue: String = ""

    @Column(name = "secretKey")
    var secretKey: String = ""

    @Column(name = "secretUrl")
    var secretUrl: String? = null

    @Column(name = "forUser")
    var forUser: String = ""
}