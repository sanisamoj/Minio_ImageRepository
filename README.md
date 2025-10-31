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

    https://www.sanisamojrepository.com/image-repo/files/9cca7e78-955e-438b-bb39-b8c15c93bef0.png

    https://www.sanisamojrepository.com/image-repo/files/e287831f-e409-4dab-ad42-79dd4c3b6437.png
