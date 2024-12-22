package com.torfstack.kayvault

import assertk.assertThat
import assertk.assertions.isNotNull
import com.torfstack.kayvault.controller.SecretController
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest

@SpringBootTest
class KayVaultApplicationTests {

	@Autowired
	val secretController: SecretController? = null

	@Test
	fun contextLoads() {
		assertThat(secretController).isNotNull()
	}
}
