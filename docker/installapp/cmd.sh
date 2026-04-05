#!/bin/bash
/wait-for-it.sh mysql:3306 
/wait-for-it.sh php-fpm:9000

php bin/console doctrine:database:create --if-not-exists 

php bin/console --no-interaction doctrine:migrations:migrate
