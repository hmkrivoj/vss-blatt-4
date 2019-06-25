FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN cd cinema_showing && go build -o main

FROM iron/go
COPY --from=builder /app/cinema_showing/main /app/cinema_showing
EXPOSE 52000-53000
ENTRYPOINT [ "/app/cinema_showing/main" ]
