From alpine:3.10

ADD tapp-update-check /usr/local/bin
RUN mkdir /etc/certs
ADD ca.crt /etc/certs
ADD tls.crt /etc/certs
ADD tls.key /etc/certs

ENTRYPOINT ["/usr/local/bin/tapp-update-check"]
