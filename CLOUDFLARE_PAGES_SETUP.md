# Cloudflare Pages 배포 가이드

## 배포 준비 완료 항목

✅ WebAssembly 소스 코드 (`cmd/wasm/main.go`)  
✅ 웹 인터페이스 (`web/` 디렉토리)  
✅ 빌드 스크립트 (`build-wasm.sh`)  
✅ 헤더 설정 (`web/_headers`)  
✅ Go 런타임 (`web/static/wasm_exec.js`)  

## Cloudflare Pages 설정 방법

### 1. Cloudflare Pages 프로젝트 생성

1. [Cloudflare Dashboard](https://dash.cloudflare.com/) 로그인
2. Pages 섹션으로 이동
3. "Create a project" 클릭
4. "Connect to Git" 선택

### 2. GitHub 저장소 연결

1. GitHub 계정 연결 (처음인 경우)
2. `pyhub-apps/pyhub-imagekit` 저장소 선택
3. "Begin setup" 클릭

### 3. 빌드 설정

**중요**: Cloudflare Pages에서 Go를 사용한 동적 빌드 설정

```
Project name: pyhub-imagekit (또는 원하는 이름)
Production branch: main
Build settings:
  - Framework preset: None
  - Build command: ./build-wasm.sh
  - Build output directory: web
  - Root directory: / (비워두기)
```

### 4. 환경 변수 (필요시)

```
GO_VERSION=1.21
```

### 5. 배포

1. "Save and Deploy" 클릭
2. 첫 빌드가 시작됨
3. 빌드 로그에서 진행 상황 확인

## 빌드 프로세스

Cloudflare Pages는 다음 순서로 빌드합니다:

1. GitHub 저장소 클론
2. Go 환경 설정
3. `build-wasm.sh` 실행
   - Go 모듈 다운로드
   - WebAssembly 컴파일 (`imagekit.wasm`)
4. `web/` 디렉토리 배포

## 도메인 설정

### 기본 도메인
- `pyhub-imagekit.pages.dev`

### 커스텀 도메인 (선택사항)
1. Pages 프로젝트 설정으로 이동
2. "Custom domains" 탭
3. "Add a custom domain" 클릭
4. 도메인 입력 및 DNS 설정

## 자동 배포

- `main` 브랜치에 푸시하면 자동으로 재배포
- Pull Request 생성 시 프리뷰 URL 생성

## 문제 해결

### 빌드 실패 시

1. **Go 버전 확인**
   ```
   환경 변수에 GO_VERSION=1.21 추가
   ```

2. **빌드 명령 확인**
   ```
   Build command: ./build-wasm.sh
   또는
   Build command: bash build-wasm.sh
   ```

3. **디렉토리 구조 확인**
   ```
   web/
   ├── index.html
   ├── static/
   │   ├── app.js
   │   ├── style.css
   │   └── wasm_exec.js
   └── _headers
   ```

### WASM 로드 실패 시

1. 브라우저 개발자 도구 확인
2. Network 탭에서 `imagekit.wasm` 로드 확인
3. Console 탭에서 에러 메시지 확인

## 로컬 테스트

배포 전 로컬에서 테스트:

```bash
# WebAssembly 빌드
./build-wasm.sh

# 로컬 서버 실행
cd web && python3 -m http.server 8080
```

브라우저에서 http://localhost:8080 접속

## 모니터링

### Analytics
Cloudflare Pages는 기본 분석 제공:
- 방문자 수
- 대역폭 사용량
- 요청 수

### Web Analytics (선택사항)
더 자세한 분석을 위해 Cloudflare Web Analytics 설정 가능

## 보안

### 자동 제공 기능
- HTTPS 자동 적용
- DDoS 보호
- 글로벌 CDN

### 추가 설정 (`web/_headers`)
- CSP (Content Security Policy)
- X-Frame-Options
- Cache-Control

## 업데이트 워크플로우

1. 로컬에서 개발 및 테스트
2. `main` 브랜치에 푸시
3. Cloudflare Pages 자동 빌드 및 배포
4. 배포 완료 알림 확인

## 유용한 링크

- [Cloudflare Pages 문서](https://developers.cloudflare.com/pages/)
- [Go on Cloudflare Pages](https://developers.cloudflare.com/pages/framework-guides/deploy-a-go-site/)
- [WebAssembly 가이드](https://developer.mozilla.org/docs/WebAssembly)