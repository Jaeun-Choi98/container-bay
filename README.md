# container-bay

## 도커 데몬 실행 (build-server)
### 1. 명령어로 실행
sudo dockerd -H unix:///var/run/docker.sock -H tcp://0.0.0.0:2375 --tls=false

### 2. 도커 설정 파일 수정한 뒤, 도커 데몬 실행 


