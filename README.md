<h1 align="center" style="font-weight: bold;">Authentication-Servicе 🔒</h1>

<p align="center">
<a href="#tech">Technologies</a>
<a href="#started">Getting Started</a>
</p>

<p align="center">Simple Authentication-Servicе with HTTP & gRPC</p>
 
<h2 id="technologies">💻 Technologies</h2>

- Golang
- HTTP router [go-chi](https://github.com/go-chi/chi)
- [gRPC](https://grpc.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
 
<h2 id="started">🚀 Getting started</h2>
 
<h3>Cloning</h3>

```bash
git clone https://github.com/D1mitrii/Authentication-Service.git
```
 
<h3>Config .env variables</h2>

Use the `.env.example` as reference to create your configuration file `.env`
 
<h3>Starting</h3>

Run Authentication-Service in Container
```bash
cd Authentication-Service
docker compose up
```
At the first launch, you will need to perform migrations from directory `migrations/`. For example, using the goose/migrate utility.