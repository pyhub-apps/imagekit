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

## Cloudflare Pages 배포

### 자동 배포 (GitHub Actions)

1. Cloudflare Pages 프로젝트 생성
2. GitHub Secrets 설정:
   - `CLOUDFLARE_API_TOKEN`: Cloudflare API 토큰
   - `CLOUDFLARE_ACCOUNT_ID`: Cloudflare 계정 ID

3. main 브랜치에 푸시하면 자동 배포

### 수동 배포

1. WebAssembly 빌드:
```bash
GOOS=js GOARCH=wasm go build -o web/static/imagekit.wasm cmd/wasm/main.go
```

2. Cloudflare Pages 대시보드에서:
   - 새 프로젝트 생성
   - `web` 디렉토리 업로드
   - 빌드 설정 불필요 (정적 파일)

### Wrangler CLI 사용

```bash
# Wrangler 설치
npm install -g wrangler

# 로그인
wrangler login

# 배포
wrangler pages deploy web --project-name=pyhub-imagekit
```

## 기술 스택

- Go WebAssembly
- HTML5/CSS3/JavaScript
- 이미지 처리: disintegration/imaging
- 호스팅: Cloudflare Pages

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