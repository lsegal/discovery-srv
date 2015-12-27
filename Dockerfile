FROM alpine:3.2
ADD discovery-srv /discovery-srv
ENTRYPOINT [ "/discovery-srv" ]
