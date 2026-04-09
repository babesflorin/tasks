FROM php:7.4-cli

WORKDIR /var/www/task-list

RUN apt-get update && apt-get install -y \
    curl openssl unzip libzip-dev libsqlite3-dev \
    && docker-php-ext-install pdo pdo_mysql pdo_sqlite zip \
    && rm -rf /var/lib/apt/lists/*

# Install composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

COPY --from=builder:latest /var/www/task-list /var/www/task-list

# Install phpunit bridge
RUN php vendor/bin/simple-phpunit install 2>/dev/null || true

CMD ["php", "vendor/bin/simple-phpunit"]
