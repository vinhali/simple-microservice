FROM golang:1.21.5-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/web-banking-backend

FROM scratch
COPY --from=build /bin/web-banking-backend /bin/web-banking-backend
ENTRYPOINT ["/bin/web-banking-backend"]
