# tiernow

A modern tier list creator with real-time collaboration and no ads.

## Developing

Run `docker compose up`. This should start the object storage, backend and frontend.

You can then access the frontend via `localhost`.
The backend runs at `localhost:5452`.

### First time setup

When starting the project for the first time, the object storage must be set up for development.
This will eventually be automated.

1. Run `docker exec -it tiernow-garage /garage status`. Copy the ID of the node.
2. Run `docker exec -it tiernow-garage /garage layout assign -z local1 -c 1G <node ID>`. This will assign the node to the zone `local1` and give it 1GB of capacity.
3. Run `docker exec -it tiernow-garage /garage layout apply --version 1`. This commits the assignment.
4. Run `docker exec -it tiernow-garage /garage bucket create tiernow`. This creates our bucket `tiernow`.
5. Copy the `.env.example` file to `.env`. 
6. Run `docker exec -it tiernow-garage /garage key create tiernow-app-key`. This creates a key to which we will grant access to `tiernow`. Copy the "Key ID" and the "Secret key" to the respective variables in `.env`.
7. Run `docker exec -it tiernow-garage /garage bucket allow --read --write --owner tiernow --key tiernow-app-key`. This grants access to the bucket with our key.
8. Now, stop the Docker Compose and start it again to read the new env variables.

## Building

To create a production version:

```bash
bun run build
```

You can preview the production build with `bun run preview`.
