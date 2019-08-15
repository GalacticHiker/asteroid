FROM alpine
RUN apk add --no-cache bash
COPY logmill-linux /
CMD [ "/logmill-linux", "--protocol", "tcp", "--destAddr", "192.168.0.29:514", "--logsPerTick=1", "--nLogsToSend=10", "--tick=1s", "--template=defaultKVP" ]
