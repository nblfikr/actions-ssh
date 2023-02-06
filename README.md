[![Audit](https://github.com/nblfikr/actions-ssh/actions/workflows/audit.yml/badge.svg)](https://github.com/nblfikr/actions-ssh/actions/workflows/audit.yml)

# About

Actions ini digunakan untuk menjalankan sebuah script pada server host tujuan melalui protokol jaringan aman atau SSH (Secure Shell).

# Preparation

Untuk menggunakannya cukup mudah, Anda harus menyiapkan **private_key** dan **known_hosts** dan mengubahnya menjadi base64 encode yang kemudian disimpan di Github Secret.

## Private Key


`ssh-keygen -t rsa -b 4096`

Perintah di atas menghasilkan dua buah file `id_rsa` dan `id_rsa.pub` pada direktori `~/.ssh`. Untuk menggunakannya salin file `id_rsa.pub` ke server target menggunakan perintah berikut

`ssh-copy-id -i ~/.ssh/id_rsa.pub user@host`


****NOTE: Private Key harus memilik passphrase.***

[See reference](https://www.ssh.com/academy/ssh/keygen)

## Known Hosts

`ssh-keyscan -R [target-host]`

Perintah tersebut menghasilkan file `~/.ssh/known_hosts`.


## Change to base64 encode

Setelah berhasil membuat private key dan known hosts. Langkah selanjutnya adalah mengubah keduannya menjadi base64 encode.

**Known Hosts**

`echo $(cat $HOME/.ssh/known_hosts | base64 >> /tmp/known_hosts && tr -d '\n' < /tmp/known_hosts)`


**Private Key**

`echo $(cat $HOME/.ssh/id_rsa | base64 >> /tmp/private_key && tr -d '\n' < /tmp/private_key)`

[See reference](https://docs.github.com/en/actions/security-guides/encrypted-secrets#storing-base64-binary-blobs-as-secrets)

# Requirements

`host`: *(Required)* Alamat IP host target

`port`: *(Optional)* Default 22

`user`: *(Required)* User host target

`private_key`: *(Required)* Path file private key

`passphrase`: *(Required)* Passphrase private key

`known_hosts`: *(Required)* Path file known_hosts

`command`: *(Required)* Perintah yang akan dijalankan