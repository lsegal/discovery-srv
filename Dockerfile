FROM alpine:3.2
ADD slack-srv /slack-srv
ENTRYPOINT [ "/slack-srv" ]
