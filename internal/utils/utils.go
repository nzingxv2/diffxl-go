// File: internal/utils/utils.go
package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func HandleError(err error, message string) {
	if err != nil {
		logrus.WithError(err).Error(message)
		os.Exit(1)
	}
}

func GetSheetsToCompare(beforeSheets, afterSheets []string, specifiedSheet string) []string {
	// Sheet selection logic...
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}