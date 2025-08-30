# ImageKit - 이미지 변환 CLI 도구

미리캔버스(MiriCanvas)에 최적화된 이미지 변환 도구입니다.

## 빠른 설치

### macOS

```bash
# Intel Mac
curl -L https://github.com/pyhub-apps/pyhub-imagekit/releases/latest/download/imagekit-darwin-amd64 -o imagekit
chmod +x imagekit
sudo mv imagekit /usr/local/bin/

# Apple Silicon (M1/M2)
curl -L https://github.com/pyhub-apps/pyhub-imagekit/releases/latest/download/imagekit-darwin-arm64 -o imagekit
chmod +x imagekit
sudo mv imagekit /usr/local/bin/
```

### Windows

PowerShell을 관리자 권한으로 실행:

```powershell
# Windows (x64)
Invoke-WebRequest -Uri "https://github.com/pyhub-apps/pyhub-imagekit/releases/latest/download/imagekit-windows-amd64.exe" -OutFile "imagekit.exe"
Move-Item -Path "imagekit.exe" -Destination "C:\Windows\System32\imagekit.exe"

# 또는 사용자 폴더에 설치
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\bin"
Invoke-WebRequest -Uri "https://github.com/pyhub-apps/pyhub-imagekit/releases/latest/download/imagekit-windows-amd64.exe" -OutFile "$env:USERPROFILE\bin\imagekit.exe"
# 환경 변수에 경로 추가 (한 번만 실행)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\bin", [EnvironmentVariableTarget]::User)
```

### 설치 확인

```bash
imagekit --version
```

## 주요 기능

- ✅ **이미지 크기 변환**: 원하는 픽셀 크기나 비율로 이미지 리사이징
- ✅ **DPI 변환**: 72, 96, 150, 300 DPI로 변환
- ✅ **가장자리 크롭**: 이미지 가장자리 제거 (여백 제거용)
- ✅ **배치 처리**: glob 패턴으로 여러 파일 동시 처리
- ✅ **형식 지원**: JPG, PNG 이미지 지원
- ✅ **WebAssembly 버전**: 브라우저에서 직접 실행 가능 (서버 전송 없음)
- ✅ **고품질 변환**: 이미지 품질 손실 최소화

## 소스에서 빌드

### 빌드 방법

```bash
# 저장소 클론
git clone https://github.com/allieus/pyhub-imagekit.git
cd pyhub-imagekit

# 빌드
make build

# 또는 직접 빌드
go build -o imagekit cmd/imagekit/main.go
```

### 크로스 플랫폼 빌드

```bash
# 모든 플랫폼용 빌드
make build-all

# 개별 플랫폼
make build-windows
make build-mac
make build-linux
```

## 사용법

### 이미지 정보 확인

```bash
imagekit info image.jpg
```

### 크기 변환

```bash
# 특정 크기로 변환 (비율 유지)
imagekit convert --width=1920 --height=1080 input.jpg output.jpg

# 너비만 지정 (비율 유지)
imagekit convert --width=800 input.jpg output.jpg

# 높이만 지정 (비율 유지)
imagekit convert --height=600 input.jpg output.jpg

# 배수로 크기 지정 (2배 확대)
imagekit convert --width=2x input.jpg output.jpg    # 2배 너비
imagekit convert --width=x2 input.jpg output.jpg    # x2 형식도 지원
imagekit convert --width=2x --height=2x input.jpg output.jpg  # 전체 2배

# 축소 (0.5배 = 절반 크기)
imagekit convert --width=0.5x input.jpg output.jpg  # 절반 크기
imagekit convert --width=0.25x input.jpg thumbnail.jpg  # 1/4 크기 (썸네일)

# 정확한 크기로 변환 (비율 무시)
imagekit convert --width=800 --height=600 --mode=exact input.jpg output.jpg

# 채우기 모드 (크롭)
imagekit convert --width=800 --height=600 --mode=fill input.jpg output.jpg
```

### DPI 변환

```bash
# DPI를 96으로 변환
imagekit convert --dpi=96 input.jpg output.jpg

# DPI를 72로 변환 (웹용)
imagekit convert --dpi=72 input.jpg output.jpg
```

### 크기와 DPI 동시 변환

```bash
imagekit convert --width=1920 --height=1080 --dpi=96 input.jpg output.jpg
```

### 배치 처리 (여러 파일 동시 변환)

```bash
# 모든 JPG 파일을 1920픽셀 너비로 변환
imagekit convert --width=1920 "*.jpg"

# 디렉토리의 모든 PNG 파일 DPI 변환
imagekit convert --dpi=96 "images/*.png"

# 여러 파일 크기 조정 (결과: image1_converted.jpg, image2_converted.jpg ...)
imagekit convert --width=800 --height=600 "photos/*.jpg"

# 모든 이미지를 2배로 확대
imagekit convert --width=2x --height=2x "*.jpg"

# 썸네일 일괄 생성 (25% 크기)
imagekit convert --width=0.25x --height=0.25x "originals/*.jpg"
```

### 가장자리 크롭

```bash
# 하단 100픽셀 제거 (여백 제거용)
imagekit crop --bottom=100 input.jpg output.jpg

# 상단 10% 제거 (퍼센트 단위)
imagekit crop --top=10% header-logo.jpg clean.jpg

# 모든 가장자리에서 20픽셀씩 제거
imagekit crop --top=20 --bottom=20 --left=20 --right=20 input.jpg output.jpg

# 여러 파일 배치 크롭
imagekit crop --bottom=50 "watermarked/*.jpg"
imagekit crop --top=15% "photos/*.png"
```

## 명령어 옵션

### convert 명령어

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--width` | 목표 너비 (픽셀 또는 배수: 1920, 2x, x2, 0.5x) | - |
| `--height` | 목표 높이 (픽셀 또는 배수: 1080, 2x, x2, 0.5x) | - |
| `--dpi` | 목표 DPI | - |
| `--mode` | 리사이징 모드 (fit, fill, exact) | fit |
| `--quality` | JPEG 품질 (1-100) | 95 |

### crop 명령어

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--top` | 상단에서 제거할 영역 (픽셀 또는 %) | - |
| `--bottom` | 하단에서 제거할 영역 (픽셀 또는 %) | - |
| `--left` | 좌측에서 제거할 영역 (픽셀 또는 %) | - |
| `--right` | 우측에서 제거할 영역 (픽셀 또는 %) | - |

## 리사이징 모드

- **fit**: 지정된 크기 내에서 비율을 유지하며 맞춤
- **fill**: 지정된 크기를 채우며, 필요시 크롭
- **exact**: 정확한 크기로 변환 (비율 변경 가능)

## 개발

### 테스트 실행

```bash
# 모든 테스트 실행
make test

# 커버리지 포함
make test-coverage

# 특정 패키지 테스트
go test ./pkg/transform/...
```

### 코드 포맷팅

```bash
make fmt
```

### 정적 분석

```bash
make lint
```

## 라이브러리로 사용

```go
import "github.com/allieus/pyhub-imagekit/pkg/transform"

// 트랜스포머 생성
transformer := transform.NewTransformer()

// 이미지 리사이징
options := transform.ResizeOptions{
    Width:   1920,
    Height:  1080,
    Mode:    transform.ResizeFit,
    Quality: 95,
}
err := transformer.Resize(input, output, options)

// DPI 설정
err := transformer.SetDPI(input, output, 96)

// 가장자리 크롭
cropOptions := transform.EdgeCropOptions{
    Top:    transform.CropValue{Value: 10, IsPercent: true},
    Bottom: transform.CropValue{Value: 100, IsPercent: false},
}
err := transformer.CropEdges(input, output, cropOptions)
```

## 요구사항

- Go 1.19 이상

## 라이선스

MIT License

## 기여

이슈 및 풀 리퀘스트를 환영합니다!