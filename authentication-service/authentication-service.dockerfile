# Use an official Golang image to build the application
FROM alpine:latest

RUN mkdir /app
COPY authApp /app 

CMD [ "/app/authApp" ]
