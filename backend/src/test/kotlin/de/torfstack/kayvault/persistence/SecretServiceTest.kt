package de.torfstack.kayvault.persistence

import assertk.assertThat
import assertk.assertions.containsOnly
import assertk.assertions.hasSize
import assertk.assertions.isEmpty
import org.junit.jupiter.api.Test
import org.junit.jupiter.params.ParameterizedTest
import org.junit.jupiter.params.provider.ValueSource
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.annotation.DirtiesContext

@SpringBootTest
@DirtiesContext(classMode = DirtiesContext.ClassMode.BEFORE_EACH_TEST_METHOD)
class SecretServiceTest {

    @Autowired
    lateinit var secretService: SecretService

    @ParameterizedTest
    @ValueSource(ints = [0, 1, 2, 3, 400])
    fun `returns secrets put into the service for same user`(numberOfSecrets: Int) {
        val expected = (1..numberOfSecrets).map {
            SecretModel(secretKey = "key$it", secretUrl = "url$it", secretValue = "secret$it")
        }

        expected.forEach {
            secretService.addSecretForUser("user1", it)
        }

        val actual = secretService.secretsForUser("user1")
        assertThat(actual).hasSize(numberOfSecrets)
        assertThat(actual).containsOnly(*expected.toTypedArray())
    }

    @Test
    fun `does not return secrets of user A to user B`() {
        secretService.addSecretForUser("A", SecretModel(secretValue = "secret", secretKey = "key", secretUrl = "url"))
        val actual = secretService.secretsForUser("B")
        assertThat(actual).isEmpty()
    }
}