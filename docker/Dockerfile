FROM alpine
RUN apk add --no-cache bash
RUN mkdir -p /logmill/bin
COPY bin/logmill-linux /logmill/bin
RUN mkdir -p /logmill/conf
COPY  conf/example.template /logmill/conf
ENTRYPOINT [ "/logmill/bin/logmill-linux" ]
#CMD [ "/logmill/bin/logmill-linux" ]
#CMD [ "--protocol", "tcp", "--destAddr", "192.168.0.29:514", "--logsPerTick=1", "--nLogsToSend=10", "--tick=1s", "--template=defaultKVP", "--tag", "dockerTag"]
