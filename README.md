# Requirements

- OS Linux (or specify `GOOS=linux`)
- `go1.19`
- `docker-compose`

# Deploy

1. Compile executable file: run `go build -o locky_bot` from `app/` directory
2. Build app Docker image: run `docker-compose build` from root repo directory
3. Run PostgreSQL: `docker-compose up -d db`
4. Run app: `docker-compose up bot`