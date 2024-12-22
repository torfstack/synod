import org.gradle.api.tasks.testing.logging.TestLogEvent
import org.jetbrains.kotlin.gradle.dsl.JvmTarget

plugins {
    id("org.springframework.boot") version "3.4.1"
    id("io.spring.dependency-management") version "1.1.7"
    kotlin("jvm") version "2.1.0"
    kotlin("plugin.spring") version "2.1.0"
}

group = "com.torfstack"
version = "0.0.1-SNAPSHOT"
java.sourceCompatibility = JavaVersion.VERSION_23
java.targetCompatibility = JavaVersion.VERSION_23

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-data-rest")
    //implementation("org.springframework.boot:spring-boot-starter-oauth2-client")
    implementation("org.springframework.boot:spring-boot-starter-logging")
    implementation("org.springframework.boot:spring-boot-starter-web")

    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation("org.springframework.boot:spring-boot-starter-data-jdbc")
    //implementation("org.hibernate:hibernate-core:6.1.7.Final")

    implementation("com.fasterxml.jackson.module:jackson-module-kotlin")

    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")

    implementation("io.github.microutils:kotlin-logging:3.0.5")

    implementation("com.nimbusds:nimbus-jose-jwt:9.37.2")
    implementation("com.google.firebase:firebase-admin:9.4.2")

    runtimeOnly("com.h2database:h2:2.2.220")
    runtimeOnly("org.mariadb.jdbc:mariadb-java-client:3.1.2")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("com.willowtreeapps.assertk:assertk-jvm:0.25")
}

kotlin {
    compilerOptions {
        jvmTarget.set(JvmTarget.JVM_23)
        freeCompilerArgs = listOf("-Xjsr305=strict")
    }

}

tasks.withType<Test> {
    useJUnitPlatform()
    testLogging {
        events.add(TestLogEvent.FAILED)
        events.add(TestLogEvent.PASSED)
        events.add(TestLogEvent.SKIPPED)
    }
}

tasks.named("bootRun") {
    (this as JavaExec).jvmArgs = listOf("-Dspring.profiles.active=dev")
    environment("FIREBASE_SECRET", file("./kayvault.json").readText())
}
