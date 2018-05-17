FROM postgres:latest
ADD omero2cnb /usr/local/bin/
ADD entrypoint.sh /usr/local/bin/
ADD notify.sql /root/
ADD trigger.sql /root/
EXPOSE 8080 8081
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
