## KayVault Backed

### Postgres

#### Local development

```bash
docker run --name kayvault-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres -p 5432:5432
```

create database

```bash
docker exec kayvault-postgres /bin/sh -c "psql -U postgres -c 'CREATE DATABASE kayvault;'"
```
