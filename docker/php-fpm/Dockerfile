FROM php:7.4-fpm

WORKDIR /var/www/task-list/

COPY --from=builder:latest /var/www/task-list /var/www/task-list

RUN apt-get update && apt-get install -y curl openssl unzip libzip-dev
RUN docker-php-ext-install pdo pdo_mysql zip


