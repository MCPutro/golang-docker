# Proyek belajar golang

Teknologi yang digunakan pada Proyek ini adalah :
- PostgreSQL : SQL Database Server
- GoLang : Bahasa Pemrograman Backend (IDE) versi 1.19.9
- Docker : Container
- Git : Sistem Pengontrol Versi (Kode)
- Postman : Dokumentasi API

## Menjalankan Proyek
### 1. Clone Proyek dari Github
```
$ mkdir go-project
$ cd go-project
$ git clone https://github.com/MCPutro/golang-docker.git
```

tunggu hingga porses selesai, dan akan muncul seperti berikut :
```
Cloning into 'golang-docker'...
remote: Enumerating objects: 261, done.
remote: Counting objects: 100% (261/261), done.
Receiving objects: 100% (261/261), 53.14 KiB | 1.18 MiB/s, done.
Resolving deltas:  22% (27/122)0% (172/172), done.
Resolving deltas:  27% (33/122)reused 208 (delta 78), pack-reused 0
Resolving deltas: 100% (122/122), done.
```
lalu masuk kedalam folder golang-docker
```
$ cd golang-docker
```

### 2. Jalankan proyek ini dengan docker-compose
```
$ docker compose up -d
```
tunggu hingga project selesai, dan akan muncul tampilan seperi berikut :
```
....
[+] Running 3/3
 ✔ Network golang-docker_default    Created                  0.9s 
 ✔ Container postgres-local-docker  Started                  2.0s 
 ✔ Container backend                Started                  3.4s 
```
### 3. Akses proyek
untuk melakukan testing terhadap Rest API yang tersedia bisa menggunakan Postman dan untuk melihat/membuka database bisa menggunakan DBeaver.
- ### DBeaver
>1. Buka `DBeaver`.
>2. Klik ikon `Connect to a Database` di Pojok-Kiri-Atas.
>3. Pada Kategori `Popular`, pilih `PostgreSQL` lalu klik Next.
>4. Isikan Konfigurasi sesuai dengan File `.env` pada Proyek ini.
>5. Jika berhasil, akan ada Ikon centang hijau pada daftar Database di sebelah kiri.

- ### Postman
pada proyek ini juga disematkan Collection Postman yang dapat anda Gunakan untuk mencoba API.
>Pre-Build Account : 
>- username : admin.support
>- password : admin123

### 4. Mematikan Proyek
```
$ docker-compose down -v 
```


<!-- This content will not appear in the rendered Markdown 
```
docker build -t test-go-docker:1.0.1 .   
```

```
docker run --name user-manegement -d -p 9999:9999 -it test-go-docker:1.0.1 
```

```dockerfile
#FROM golang:1.19.9 AS builder
FROM --platform=$BUILDPLATFORM golang:1.19.9-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o bin/backend main.go

CMD ["/app/bin/backend"]

FROM builder AS dev-envs

RUN <<EOF
apk update
apk add git
EOF
#
RUN <<EOF
addgroup -S docker
adduser -S --shell /bin/bash --ingroup docker vscode
adduser -S --shell /bin/bash --ingroup docker vscode
EOF

# install Docker tools (cli, buildx, compose)
COPY --from=gloursdocker/docker / /

CMD ["go", "run", "main.go"]

FROM scratch
COPY --from=builder /app/bin/backend /usr/local/bin/backend
COPY --from=builder /app/.env .

#COPY .env /usr/local/bin/

CMD ["/usr/local/bin/backend"]

```


```
docker compose up -d
```

```
docker-compose down --volumes
```
-->