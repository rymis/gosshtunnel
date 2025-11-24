FROM golang:1.24

RUN mkdir /src
COPY . /src

FROM scratch
COPY --from=0 /src/gosshtunnel /bin/gosshtunnel

CMD [ "/bin/gosshtunnel" ]
