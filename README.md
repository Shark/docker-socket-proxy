# docker-socket-proxy

`docker-socket-proxy` is a small utility which listens on a UNIX domain socket and proxies every HTTP request to `/var/run/docker.sock`. It filters incoming requests and currently only allows GET requests. It will respond with `403 Forbidden` for any other method. This allows you to use the socket provided by `docker-socket-proxy` as a drop-in where you normally would mount `/var/run/docker.sock`. Exposing this socket to any container is [inherently insecure](https://www.lvh.io/posts/dont-expose-the-docker-socket-not-even-to-a-container.html), but often [very useful](https://docs.traefik.io/#docker). This utility should make most applications work (they normally listen for events from the Docker engine) and provide you a decent(:tm:) amount of security.

## Caveats

- Commands such as `docker attach` or `docker run` will currently not work (the proxy only supports traditional request-reponse style use)
- Only `GET` requests will be allowed. I probably will add the ability to whitelist specific endpoints in the future.

## Building

- Clone the repository in your `$GOPATH/src/github.com/Shark/docker-socket-proxy`
- Run `go build .`

## Usage

Run `./docker-socket-proxy`

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request! :)

## History

- v0.0.1 (2016-11-03): initial version

## License

This project is licensed under the MIT License. See LICENSE for details.
