FROM golang:1.24 AS build

RUN mkdir /src
COPY . /src
WORKDIR /src

RUN go build

FROM debian:bookworm
COPY --from=build /src/gosshtunnel /bin/gosshtunnel

CMD [ "/bin/gosshtunnel" ]
ENTRYPOINT [ "/bin/gosshtunnel" ]
