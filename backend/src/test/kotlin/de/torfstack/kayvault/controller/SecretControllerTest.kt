package de.torfstack.kayvault.controller

import assertk.assertThat
import assertk.assertions.hasSize
import de.torfstack.kayvault.persistence.SecretService
import de.torfstack.kayvault.validation.TokenValidator
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.Mockito
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.mock.mockito.MockBean
import org.springframework.http.MediaType
import org.springframework.test.annotation.DirtiesContext
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.content
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.status

@SpringBootTest
@AutoConfigureMockMvc
@DirtiesContext(classMode = DirtiesContext.ClassMode.BEFORE_EACH_TEST_METHOD)
class SecretControllerTest {

    @MockBean
    lateinit var tokenValidator: TokenValidator

    @Autowired
    lateinit var mockMvc: MockMvc

    @Autowired
    lateinit var secretService: SecretService

    companion object {
        const val MOCK_USER = "user"
    }

    @BeforeEach
    fun `mock token validator`() {
        Mockito.`when`(tokenValidator.validate(MOCK_USER)).thenAnswer { TokenValidator.ValidVerification(MOCK_USER) }
    }

    @Test
    fun `getting secret returns nothing`() {
        mockMvc.perform(get("/secret").header("Authorization", "Bearer $MOCK_USER"))
            .andExpect(status().isOk)
            .andExpect(content().string("[]"))
    }

    @Test
    fun `secret for one user can be retrieved`() {
        mockMvc.perform(post("/secret").header("Authorization", "Bearer $MOCK_USER")
            .contentType(MediaType.APPLICATION_JSON)
            .content("""
                {
                    "key": "system",
                    "value": "abc-def-ghij"
                }
            """.trimIndent()))
            .andExpect(status().isOk)

        val secrets = secretService.secretsForUser(MOCK_USER)
        assertThat(secrets).hasSize(1)
    }
}