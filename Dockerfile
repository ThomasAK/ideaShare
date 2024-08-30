FROM golang:1.23.0-bookworm as server
COPY server /server
WORKDIR /server
RUN CGO_ENABLED=0 go mod download && \
go build -o server main.go


FROM node:20-bookworm-slim as ui
COPY ui /ui
WORKDIR /ui
RUN npm i && npm run build

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=server /server/server /app/server
COPY --from=ui /ui/dist /app/public
CMD ["/app/server"]