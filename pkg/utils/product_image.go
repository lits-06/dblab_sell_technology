package utils

import (
	"fmt"
	"regexp"
)

func generateImageFilename(productName string, productID string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleanedProductName := re.ReplaceAllString(productName, "_")

	filename := fmt.Sprintf("%s_%s.png", cleanedProductName, productID)

	return filename
}