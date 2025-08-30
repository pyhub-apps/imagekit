# ImageKit - 이미지 변환 CLI 도구

미리캔버스(MiriCanvas)에 최적화된 이미지 변환 도구입니다.

## 주요 기능

- ✅ **이미지 크기 변환**: 원하는 픽셀 크기나 비율로 이미지 리사이징
- ✅ **DPI 변환**: 72, 96, 150, 300 DPI로 변환
- ✅ **워터마크 제거**: 지정 영역 블러/채우기 처리
- ✅ **형식 지원**: JPG, PNG 이미지 지원
- ✅ **고품질 변환**: 이미지 품질 손실 최소화

## 설치

### 소스에서 빌드

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

### 워터마크 제거

```bash
# 특정 영역 블러 처리 (x, y, width, height)
imagekit watermark --area=100,100,200,50 input.jpg output.jpg

# 제거 방법 지정
imagekit watermark --area=100,100,200,50 --method=blur input.jpg output.jpg
imagekit watermark --area=100,100,200,50 --method=fill input.jpg output.jpg
imagekit watermark --area=100,100,200,50 --method=inpaint input.jpg output.jpg
```

## 명령어 옵션

### convert 명령어

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--width` | 목표 너비 (픽셀) | - |
| `--height` | 목표 높이 (픽셀) | - |
| `--dpi` | 목표 DPI | - |
| `--mode` | 리사이징 모드 (fit, fill, exact) | fit |
| `--quality` | JPEG 품질 (1-100) | 95 |

### watermark 명령어

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--area` | 워터마크 영역 (x,y,width,height) | 필수 |
| `--method` | 제거 방법 (blur, fill, inpaint) | blur |

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

// 워터마크 제거
area := transform.Rectangle{
    X:      100,
    Y:      100,
    Width:  200,
    Height: 50,
}
err := transformer.RemoveWatermark(input, output, area)
```

## 요구사항

- Go 1.19 이상

## 라이선스

MIT License

## 기여

이슈 및 풀 리퀘스트를 환영합니다!