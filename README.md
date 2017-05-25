# docker-recreate

I always wanted to have the functionality of re-creating docker containers based on new images, coming from, for example, CI builds. Docker does not provide this functionality out of the box, and it does involve some operations, like transferring environment variables, copying links, etc.

Luckily, [@lanrat](https://github.com/lanrat) wrote a small Gist doing exactly this: https://gist.github.com/lanrat/8a8b385959627a7b29f1. This is the origin of this small application. `docker-machine` is fine, and with `docker-recreate`, CI builds (or, simply, switching image versions) are easily incorporated.


## Usage

```
docker-recreate [-p] id
```

`id` specifies the container ID of the container to-be-recreated. Specify `-p` to make sure the newest image is pulled from the repository.
