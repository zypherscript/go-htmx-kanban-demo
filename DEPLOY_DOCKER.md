# Kanban Board Docker Deployment

## Build the Docker image

```bash
docker build -t kanban-board .
```

## Run the container (recommended persistent data)

```bash
docker run -p 8080:8080 -v $PWD/data:/app/data -e KANBAN_DATA_FILE=/app/data/tasks.json kanban-board
```
- This mounts your local data directory to /app/data in the container.
- The app will create and manage tasks.json automatically inside data.
- The app will be available at http://localhost:8080

## Custom Data Location

You can use a different data file or directory:
```bash
docker run -p 8080:8080 -v /your/path:/app/data -e KANBAN_DATA_FILE=/app/data/tasks.json kanban-board
```

## Updating the Container
- Stop the old container
- Rebuild the image if you change code/templates
- Start a new container (your data is safe in the mapped file)

## Notes
- If you want to pre-seed tasks, copy a `tasks.json` into your project before building/running.
- The `templates/` folder is copied into the image at build time.
- For production, use a reverse proxy (nginx, Caddy) for HTTPS.

---

Happy deploying! 🚀
