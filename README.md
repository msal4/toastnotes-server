# [Toast Notes API](https://toast.msal.dev) &middot; ![CI](https://github.com/msal4/toastnotes-server/workflows/CI/badge.svg) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/msal4/toastnotes-server/blob/master/LICENSE) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/msal4/toastnotes-server/pulls)


The backend for toast notes app.
### Requirements
- [Go](https://golang.org/)
- [Postgres](https://www.postgresql.org/)
- [Docker](https://www.docker.com/) (optional)

### Setup
- Create a .env file
  ```bash
  cp .env.example .env
  ```
  set the `JWT_SECRET`, and `DATABASE_URL` for example `postgres://<user>:<password>@<host>:<port>/<database-name>` in .env


### Run
- Using docker
  ```bash
  docker build -t <app-tag> .
  docker run -it <app-tag>
  ```
- Without docker
  ```bash
  make dev
  ```
### Deploy
- Set `GIN_MODE=release` in .env
- Docker Compose
  - create a `docker-compose.yml` file and add your dependencies there (e.g. postgres)
  - upload to your server (e.g. using git)
  - build and run using `docker-compose`
- Dokku
  - use the postgres plugin and link it to toastnotes
  - using `dokku config` set your environment variables
  - push to dokku, for more details checkout [dokku docs](http://dokku.viewdocs.io/dokku/)

### License
See [LICENSE](https://github.com/msal4/toastnotes-server/blob/master/LICENSE)
