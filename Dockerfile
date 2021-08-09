FROM ubuntu:20.04

COPY ./output/ui /opt/ui/

COPY ./secret/ui_rsa_pri.pem /opt/ui/secret/
COPY ./secret/ui_rsa_pub.pem /opt/ui/secret/

COPY ./pg/users.sql /opt/ui/pg/

WORKDIR /opt/ui
ENTRYPOINT ["/opt/ui/ui"]