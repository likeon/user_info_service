A simple user info service in golang using sqlite and gorm.

## Running
### Prebuilt binaries
Binaries are available as artifacts of [GitHub actions](https://github.com/likeon/user_info_service/actions).

### Containers
Dockerfile is included and has the server as entrypoint.

To build and run:
```bash
docker build . -t user-info-service:latest
poddockerman run --rm -it -p 8080:8080 user-info-service:latest
```
Mount `users.db` from your filesystem if you wish to preserve the database between runs.

## Testing
`go test -v`
