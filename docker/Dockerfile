FROM java:openjdk-8u111-jre-alpine

LABEL author="tri.doan@trustingsocial.com"

ENV app /app
WORKDIR ${app}

COPY ./smppsim* ${app}/smpp/
COPY ./entrypoint.sh /entrypoint.sh

EXPOSE 2775 
EXPOSE 2777
EXPOSE 2779

RUN chmod 700 /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
