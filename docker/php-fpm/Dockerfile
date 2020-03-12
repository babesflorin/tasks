FROM php:7.4-fpm

WORKDIR /var/www/task-list/

RUN apt-get update && apt-get install -y curl openssl unzip libzip-dev
RUN docker-php-ext-install pdo pdo_mysql zip
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x wait-for-it.sh
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
RUN composer install
