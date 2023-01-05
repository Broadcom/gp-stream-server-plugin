package common

import "github.com/sirupsen/logrus"

type BaseContext interface {
	// GetLogger returns the logger used by gpss
	GetLogger() logrus.FieldLogger

	// GetProperties returns the custom properties associated with the plugin
	GetProperties() map[string]string
}
