# Proyek belajar golang

Teknologi yang digunakan pada Proyek ini adalah :
- PostgreSQL : SQL Database Server
- GoLang : Bahasa Pemrograman Backend (IDE) versi 1.19.9
- Docker : Container
- Git : Sistem Pengontrol Versi (Kode)
- Postman : Dokumentasi API

## Menjalankan Proyek
### 1. Buat folder log
Buat folder untuk menampung/menyimpan file log lalu ambil path atau address dari folder tersebut.
Pada projek ini file log akan di simpan di ```C:\Users\Public\Logs```.


### 2. Clone Proyek dari Github
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
```shell
cd golang-docker
```

#### Merubah Folder Log: 
Setelah masuk ke folder ```golang-docker``` edit file ```.env``` lalu ubah value pada bagian ```LOG_PATH```.  Paste-kan path atau address dari folder yang telah di buat pada langkah pertama.

### 3. Jalankan proyek ini dengan docker compose
```shell
docker compose up -d
```
tunggu hingga project selesai, dan akan muncul tampilan seperi berikut :
```
....
[+] Running 3/3
 ✔ Network golang-docker_default    Created                  0.9s 
 ✔ Container postgres-local-docker  Started                  2.0s 
 ✔ Container backend                Started                  3.4s 
```
### 4. Database Migrations
Pada project ini menggunakan database migration untuk membantu melakukan tracking perubahan struktur database. 

Langkah untuk menginstall :
- install golang migrate dengan perintah di bawah ini. Saat menginstall Golang Migrate, secara otomatis terdapat executable file di folder `$GOPATH/bin/` dengan nama `migrate`.
```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

- jalankan migration up untuk proses pembuatan table dan user `admin.support`.  
```shell
migrate -database "postgres://emchepe:welcome1@localhost:5432/test_user_management?sslmode=disable" -path ./migrations up
```
note : Agar `migrate` bisa di akses lewat terminal, pastikan `$GOPATH/bin` sudah terpasang di terminal.

### 5. Akses proyek
untuk melakukan testing terhadap Rest API yang tersedia bisa menggunakan Postman dan untuk melihat/membuka database bisa menggunakan DBeaver.
- ### DBeaver
>1. Buka `DBeaver`.
>2. Klik ikon `Connect to a Database` di Pojok-Kiri-Atas.
>3. Pada Kategori `Popular`, pilih `PostgreSQL` lalu klik Next.
>4. Isikan Konfigurasi sesuai dengan File `.env` pada Proyek ini.
>5. Jika berhasil, akan ada Ikon centang hijau pada daftar Database di sebelah kiri.

- ### Postman 
pada proyek ini juga disematkan [Collection Postman (GolangDocker.postman_collection.json)](https://github.com/MCPutro/golang-docker/blob/master/postmanCollection/GolangDocker.postman_collection.json) yang dapat anda import kepostman untuk mencoba API.
>Pre-Build Account : 
>- username : admin.support
>- password : admin123

### 6. Mematikan Proyek
```shell
$ docker compose down -v 
```


<!-- This content will not appear in the rendered Markdown 
```
docker build -t test-go-docker:1.0.1 .   
```

```
docker run --name user-manegement -d -p 9999:9999 -it test-go-docker:1.0.1 
```
-->
