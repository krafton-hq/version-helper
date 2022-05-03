# version-maker-net

input, output은 아래 세 가지 파일에서 관리됨
```
Version.txt
- 설명: input, output 버전 파일
- 예시: 0.9.0.2
CommitHash.txt
- 설명: input, output commit 해시 파일 (string 값 비교만 함)
- 예시: 2cdae3298f...
VersionLast.txt
- 설명: output 기존 버전 파일
- 예시: 0.9.0.1
```

## 필수 파라미터
```sh
version-maker -raise=3 # raise번 째 버전 up

version-maker -set=0.9.0.4 # 강제 설정

version-maker -raise=3 -commithash=23decab... #커밋 해시도 같이 저장

version-maker -raise=3 -commithash=23decab... -commitver=3 #commitver번 째 버전도 up

version-maker -raise=3 -generate=header.h -datadir=<path> -gameini=DefaultGame.ini #header랑 ini generate

version-maker -raise=3 -commithash=23decab... -commitver=3 -generate=header.h -datadir=<path> -gameini=DefaultGame.ini #header랑 ini generate
```

# 옵션 파라미터
```sh
-chdir=<path> #workdir 변경
-isquiet=true|false #결과 출력할 때 버전만 STDOUT으로 출력
-last=true|false #true면 이전 버전만 출력하고 종료
-redis=<addr> # raise하거나 가져오거나 저장할 때 사용할 redis 주소
```
