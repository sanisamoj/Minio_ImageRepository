# IMAGE_REPO - Servidor de repositório de mídias.

O repositório de mídias é um servidor que armazena vários tipos de mídias, como imagens, vídeos gifs e documentos, podendo ser utilizado como um microservice de armazamento integrando a alguma arquitetura maior.

## Recursos
- ***Armazenamento***: Armazenamento com o Minio.

## Para instalação

#### Executar com docker-compose.

    docker compose -p image_repo up --build -d

### Crie a rede interna com as seguintes informações

    docker network create image-repo-network

### Para rodar o ecosistema da ADA no Docker-compose:

    docker compose -p image-repo up --build -d

```.env
MINIO_ENDPOINT=
MINIO_ACCESS_KEY_ID=
MINIO_SECRET_ACCESS_KEY=
MINIO_BUCKET=uploads
JWT_STORAGE_UPLOAD_SECRET=
API_PORT=6868
SELF_HOST=https://www.sanisamojrepository.com/image-repo #Não é obrigatório passar, apenas se usar algum proxy

EMAIL_AUTH_USER=
EMAIL_AUTH_PASS=
EMAIL_HOST=
EMAIL_PORT=

REDIS_HOST=
REDIS_PORT=
REDIS_PASSWORD=
```

### Para acessar um dos arquivos do Storage:

> https://www.sanisamojrepository.com/image-repo/files/dbd3baa4-5872-49a1-b394-92f45e63a4b8.png

> https://www.sanisamojrepository.com/image-repo/files/e6ca49ca-61e0-4e76-b92e-51ef213ec38b.webp

> https://www.sanisamojrepository.com/image-repo/files/27f43ae4-190d-454d-80ca-67f1dab4e1e0.mp4

> https://www.sanisamojrepository.com/image-repo/files/c45bf718-34e8-4b6b-a07f-8025cf074ab1.gif
