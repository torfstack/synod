FROM eclipse-temurin:17

RUN mkdir /opt/kayvault
COPY build/libs/kayvault*SNAPSHOT.jar /opt/kayvault/kayvault.jar
CMD ["java", "-Dspring.profiles.active=prod", "-jar", "/opt/kayvault/kayvault.jar"]

