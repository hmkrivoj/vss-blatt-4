FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o user/main user/main.go

FROM iron/go
COPY --from=builder /app/user/main /app/user
EXPOSE 52000-53000
ENTRYPOINT [ "/app/user/main" ]
