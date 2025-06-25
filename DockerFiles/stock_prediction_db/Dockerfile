# syntax=docker/dockerfile:1
FROM mysql:8.0

ENV MYSQL_DATABASE=stock_prediction_db
ENV MYSQL_USER=ethan
ENV MYSQL_PASSWORD=040323
ENV MYSQL_ROOT_PASSWORD=Front2Back!

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 3306
