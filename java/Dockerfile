FROM openjdk:8
EXPOSE 8080:8080
WORKDIR /usr/local/bin
COPY ./target/docker-0.0.1-SNAPSHOT.jar helloworld.jar
ENTRYPOINT [ "java", "-jar", "helloworld.jar"]