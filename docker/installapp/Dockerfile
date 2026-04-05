FROM tasklistfpm:latest
WORKDIR /var/www/task-list


COPY ./docker/installapp/wait-for-it.sh /wait-for-it.sh
COPY ./docker/installapp/cmd.sh /cmd.sh
RUN chmod +x /wait-for-it.sh /cmd.sh
ENTRYPOINT ["/cmd.sh"]