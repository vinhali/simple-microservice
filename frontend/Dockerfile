FROM golang:1.21.5-alpine AS build

WORKDIR /src/
COPY main.go go.* ./
COPY files/ files/
RUN CGO_ENABLED=0 go build -o /bin/web-banking-frontend

# Use a temporary image with a shell for inspection
FROM alpine:latest as inspection
COPY --from=build /bin /bin

# Inspect the contents of /bin
RUN ls -R /bin

# Switch to scratch for the final image
FROM scratch
COPY --from=build /bin/web-banking-frontend /bin/web-banking-frontend
COPY --from=build /src/files/ /bin/files/
ENTRYPOINT ["/bin/web-banking-frontend"]
