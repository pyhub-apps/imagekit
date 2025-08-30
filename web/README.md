# ImageKit WebAssembly 버전

브라우저에서 직접 실행되는 이미지 변환 도구입니다.

## 특징

- ✅ 서버 전송 없음 (100% 프라이버시)
- ✅ WebAssembly로 빠른 처리
- ✅ 다중 파일 동시 처리
- ✅ 설정 자동 저장 (localStorage)

## 로컬 개발

### 빌드 및 실행

```bash
# WebAssembly 빌드
GOOS=js GOARCH=wasm go build -o web/static/imagekit.wasm cmd/wasm/main.go

# 로컬 서버 실행
cd web && python3 -m http.server 8080

# 또는 Makefile 사용
make -f Makefile.wasm wasm-dev
```

브라우저에서 http://localhost:8080 접속

### 디버깅

테스트 페이지: http://localhost:8080/test.html

## 배포

WebAssembly 파일을 빌드한 후 `web` 디렉토리를 정적 웹 호스팅 서비스에 배포할 수 있습니다.

### 빌드

```bash
# WebAssembly 빌드
GOOS=js GOARCH=wasm go build -o web/static/imagekit.wasm cmd/wasm/main.go

# 또는 Makefile 사용
make -f Makefile.wasm wasm-build
```

## 기술 스택

- Go WebAssembly
- HTML5/CSS3/JavaScript
- 이미지 처리: disintegration/imaging
- 호스팅: 정적 웹 호스팅 서비스

## 브라우저 지원

- Chrome 57+
- Firefox 52+
- Safari 11+
- Edge 16+

## 파일 구조

```
web/
├── index.html          # 메인 애플리케이션
├── test.html          # WebAssembly 테스트 페이지
├── static/
│   ├── app.js         # 애플리케이션 로직
│   ├── style.css      # 스타일시트
│   ├── wasm_exec.js   # Go WebAssembly 런타임
│   └── imagekit.wasm  # 컴파일된 WebAssembly
└── build.sh           # 빌드 스크립트
```

## 주요 기능

- **이미지 크기 변환**: 픽셀 또는 배수(2x, 0.5x)
- **가장자리 크롭**: 픽셀 또는 퍼센트
- **DPI 변경**: 72, 96, 150, 300 DPI
- **배치 처리**: 여러 파일 동시 처리
- **설정 저장**: 사용자 설정 자동 저장