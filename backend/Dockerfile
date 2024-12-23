ARG BUILD_IMAGE=gradle:8.12.0-jdk21-alpine
ARG RUN_IMAGE=eclipse-temurin:21-alpine

FROM $BUILD_IMAGE AS pre-build-cache

ARG BUILD_HOME=/opt/kayvault

WORKDIR $BUILD_HOME
COPY build.gradle.kts settings.gradle.kts gradle $BUILD_HOME/

RUN gradle dependencies

FROM pre-build-cache AS builder

ARG BUILD_HOME=/opt/kayvault

WORKDIR $BUILD_HOME
COPY src $BUILD_HOME/src

RUN gradle bootJar

FROM $RUN_IMAGE AS runner

ARG BUILD_HOME=/opt/kayvault

COPY --from=builder $BUILD_HOME/build/libs/kayvault*SNAPSHOT.jar kayvault.jar
CMD ["java", "-Dspring.profiles.active=dev", "-jar", "kayvault.jar"]

