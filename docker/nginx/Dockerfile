FROM nginx:stable-alpine

WORKDIR /var/www/task-list
COPY ./docker/nginx/conf.d/default.conf /etc/nginx/conf.d
COPY ./ /var/www/task-list