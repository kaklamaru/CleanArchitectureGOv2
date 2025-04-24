# # ✅ ใส่ ARG ก่อน FROM
# ARG GO_VERSION=1.23.3
# FROM golang:${GO_VERSION}-bookworm AS builder

# WORKDIR /usr/src/app

# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

# COPY . .
# RUN go build -v -o /run-app ./cmd

# # ---

# FROM debian:bookworm
# COPY --from=builder /run-app /usr/local/bin/
# CMD ["run-app"]

# ---------- STAGE 1 ----------
    ARG GO_VERSION=1.23.3
    FROM golang:${GO_VERSION}-bookworm AS builder
    WORKDIR /usr/src/app
    COPY go.mod go.sum ./
    RUN go mod download && go mod verify
    COPY . .
    RUN go build -v -o /run-app ./cmd
    
    # ---------- STAGE 2 ----------
    FROM debian:bookworm
    WORKDIR /app
    COPY --from=builder /run-app /usr/local/bin/run-app
    COPY --from=builder /usr/src/app/pkg/utility/filesystem/assets ./pkg/utility/filesystem/assets
    
    # (ไม่ต้อง copy uploads)
    # uploads จะมาจาก volume mount ข้างนอก
    
    CMD ["run-app"]
    