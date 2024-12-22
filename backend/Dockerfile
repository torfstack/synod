FROM gradle:8.12.0-jdk23-alpine AS builder

WORKDIR /opt/kayvault
COPY . .

RUN ./gradlew clean bootJar

FROM eclipse-temurin:23-alpine AS runner

RUN mkdir /opt/kayvault
COPY --from=builder /opt/kayvault/build/libs/kayvault*SNAPSHOT.jar /opt/kayvault/kayvault.jar
CMD ["java", "-Dspring.profiles.active=dev", "-jar", "/opt/kayvault/kayvault.jar"]

