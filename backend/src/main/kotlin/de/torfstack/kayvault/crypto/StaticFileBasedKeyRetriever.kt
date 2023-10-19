package de.torfstack.kayvault.crypto

import org.springframework.stereotype.Service
import java.nio.file.Files
import java.nio.file.OpenOption
import java.nio.file.StandardOpenOption
import java.security.SecureRandom
import kotlin.io.path.Path

@Service
class StaticFileBasedKeyRetriever: KeyRetriever {

    override fun key(): ByteArray {
        if (!Files.exists(Path("key.file"))) {
            val key = ByteArray(32)
            SecureRandom().nextBytes(key)
            Files.write(Path("key.file"), key, StandardOpenOption.CREATE_NEW)
        }
        return Files.readAllBytes(Path("key.file"))
    }
}