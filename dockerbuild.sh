#!/bin/zsh

docker buildx build --platform linux/amd64 -t container-bay:latest ./
docker tag container-bay:latest cju.iptime.org:8443/container-bay:latest
docker push cju.iptime.org:8443/container-bay:latest
docker system prune -f


# -v /var/run/docker.sock:/var/run/docker.sock \ 
# --privileged \  <- 컨테이너 실행시 추가 필요