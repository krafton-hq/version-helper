# version-maker

# 커맨드
```shell
$ versionmaker client raise --ci-hint=(tc|jenkins|gha|azp) --counter=(fox|local)
# 주어진 조건에서 값 가져와 build count 올림

$ versionmaker client set <semver>
# 지정한 버전으로 설정

# versionmaker client 공통 파라미터
# --chdir=<path>: WorkingDirectory 변경하기 
# --gen-dir=<path>
# --gen-header=header.h
```

```shell
$ versionmaker server --ci-hint=(tc|jenkins|gha|azp) --counter=(fox|local) --counter-file=<path>

# versionmaker server 파라미터
# --gen-dir=<path>
# --gen-json=version.json
# --gen-txt=version.txt
```

```shell
$ versionmaker devops --ci-hint=(tc|jenkins|gha|azp)

# versionmaker devops 파라미터
# --gen-dir=<path>
# --gen-json=version.json
# --gen-txt=version.txt
```

1. Version Maker: Version Generate
2. 빌드 + 패키징
3. 버전 메타데이터 업로드
