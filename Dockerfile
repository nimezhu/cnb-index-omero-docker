FROM postgres:latest
ADD omeroExt /usr/local/bin/
ADD pgweb /usr/local/bin/
ADD entrypoint.sh /usr/local/bin/
ADD notify.sql /root/
ADD trigger.sql /root/
EXPOSE 8080 8081
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
