# docker-recreate

I always wanted to have the functionality of re-creating docker containers based on new images, coming from, for example, CI builds. Docker does not provide this functionality out of the box, and it does involve some operations, like transferring environment variables, copying links, etc.

It's like using the [ES7 spread operator](http://redux.js.org/docs/recipes/UsingObjectSpreadOperator.html): `const newContainer = {...oldContainer};` (At least, the mental model is.)

Luckily, [@lanrat](https://github.com/lanrat) wrote a small Gist doing exactly this: https://gist.github.com/lanrat/8a8b385959627a7b29f1. This is the origin of this small application. `docker-machine` is fine, and with `docker-recreate`, CI builds (or, simply, different image versions) are easily applied.


## Usage

You'll need [Go](https://golang.org/) and [`dep`](https://github.com/golang/dep). To install, check out the code via `go get github.com/fallafeljan/docker-recreate`. Install via `make`. Then, run it:

```
docker-recreate [-p] [-d] id [tag]
```

- `-p` Pull the image from the repository.
- `-d` Delete the then-old container when the new one is running.
- `id` Container ID of the container to-be-recreated.
- `tag` A different tag than the currently selected may be specified. (`staging`, for example.)


#### Private Repositories

You may provide credentials for private repositories by specifying them in `~/.recreate.json`:

```json
{
  "registries": [
    {
      "host": "registry.acme.corp",
      "user": "dane",
      "password": "joe",
    }
  ]
}
```
