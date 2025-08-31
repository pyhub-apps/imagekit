# PWA Icons

이 디렉토리에는 Progressive Web App을 위한 아이콘들이 포함됩니다.

## 필요한 아이콘 크기

- icon-16x16.png
- icon-32x32.png
- icon-72x72.png
- icon-96x96.png
- icon-128x128.png
- icon-144x144.png
- icon-152x152.png
- icon-180x180.png (Apple Touch Icon)
- icon-192x192.png
- icon-384x384.png
- icon-512x512.png

## 아이콘 생성 방법

1. ChatGPT (GPT-5)를 사용하여 512x512 기본 아이콘 생성
2. ImageMagick 또는 온라인 도구를 사용하여 다양한 크기로 변환

### ImageMagick 사용 예시:

```bash
# 512x512 원본에서 모든 크기 생성
convert icon-512x512.png -resize 384x384 icon-384x384.png
convert icon-512x512.png -resize 192x192 icon-192x192.png
convert icon-512x512.png -resize 180x180 icon-180x180.png
convert icon-512x512.png -resize 152x152 icon-152x152.png
convert icon-512x512.png -resize 144x144 icon-144x144.png
convert icon-512x512.png -resize 128x128 icon-128x128.png
convert icon-512x512.png -resize 96x96 icon-96x96.png
convert icon-512x512.png -resize 72x72 icon-72x72.png
convert icon-512x512.png -resize 32x32 icon-32x32.png
convert icon-512x512.png -resize 16x16 icon-16x16.png
```

## 디자인 가이드라인

- 배경: 투명 또는 단색
- 색상: #667eea (보라색 그라데이션)
- 스타일: 모던, 미니멀리스틱
- 요소: 이미지 변환을 상징하는 아이콘