package service

import (
	"bytes"
	"context"
	"github.com/gen2brain/go-fitz"
	uuid "github.com/satori/go.uuid"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service/stru"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/handler"
	"gitlab.yoyiit.com/youyi/go-core/store"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"image"
	"image/draw"
	"image/png"
	"io"
)

type PdfToImageService interface {
	UploadPdfPng(imageUrl string) (pngUrl string, err error)
	UploadPdfToImageJsdk(ctx context.Context, pdfUrl string) (string, error)
}

type pdfToImageService struct {
	ossConfig *store.OSSConfig
}

func (s *pdfToImageService) UploadPdfPng(pdfUrl string) (pngUrl string, err error) {
	body, err := store.GetOSSFile(pdfUrl, s.ossConfig)
	if err != nil {
		return "", handler.HandleError(err)
	}
	defer body.Close()
	objectBytes, err := io.ReadAll(body)
	if err != nil {
		return "", handler.HandleError(err)
	}
	doc, err := fitz.NewFromMemory(objectBytes)
	if err != nil {
		return "", handler.HandleError(err)
	}
	var images []image.Image
	for n := 0; n < doc.NumPage(); n++ {
		pngImage, err := doc.ImagePNG(n, 300)
		if err != nil {
			return "", handler.HandleError(err)
		}
		img, err := png.Decode(bytes.NewReader(pngImage))
		if err != nil {
			return "", handler.HandleError(err)
		}
		images = append(images, img)
	}
	// 创建新的合并后的图像
	bounds := image.Rect(0, 0, images[0].Bounds().Dx(), images[0].Bounds().Dy()*len(images))
	merged := image.NewRGBA(bounds)
	// 将所有图片依次绘制到新图像中
	offset := 0
	for _, img := range images {
		draw.Draw(merged, img.Bounds().Add(image.Point{0, offset}), img, image.Point{0, 0}, draw.Src)
		offset += img.Bounds().Dy()
	}
	mergedBytes, err := encodePNG(merged)
	if err != nil {
		return "", handler.HandleError(err)
	}
	electronicReceiptFile, err := store.UploadOSSFileBytes("png", ".png", mergedBytes, s.ossConfig, false)
	return electronicReceiptFile, nil
}

func encodePNG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (s *pdfToImageService) UploadPdfToImageJsdk(ctx context.Context, pdfUrl string) (string, error) {
	newFileName := uuid.NewV4().String()
	requestBody := stru.PdfToImageJsdkRequest{
		OssPath:    newFileName,
		PdfUrl:     pdfUrl,
		SourceType: 1,
	}
	var responseData stru.PdfToImageJsdkResponse
	host := config.GetString("jsdk.pdfToImg.url", "")

	err := util.PostHttpResult(ctx, host, &requestBody, &responseData)
	if err != nil {
		return "", handler.HandleError(err)
	}
	return responseData.OssPath, nil
}
