FROM amazoncorretto:17-alpine

RUN mkdir /opt/kayvault
COPY build/libs/kayvault*SNAPSHOT.jar /opt/kayvault/kayvault.jar
CMD ["java", "-Dspring.profiles.active=dev", "-jar", "/opt/kayvault/kayvault.jar"]

