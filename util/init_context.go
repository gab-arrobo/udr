// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package util

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/omec-project/openapi/models"
	"github.com/omec-project/udr/context"
	"github.com/omec-project/udr/factory"
	"github.com/omec-project/udr/logger"
)

func InitUdrContext(context *context.UDRContext) {
	config := factory.UdrConfig
	logger.UtilLog.Infof("udrconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	context.RegisterIPv4 = factory.UDR_DEFAULT_IPV4 // default localhost
	context.SBIPort = factory.UDR_DEFAULT_PORT_INT  // default port
	if sbi := configuration.Sbi; sbi != nil {
		context.UriScheme = models.UriScheme(sbi.Scheme)
		if sbi.RegisterIPv4 != "" {
			context.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			context.SBIPort = sbi.Port
		}
		if tls := sbi.Tls; tls != nil {
			if tls.Key != "" {
				context.Key = tls.Key
			}
			if tls.Pem != "" {
				context.PEM = tls.Pem
			}
		}

		context.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if context.BindingIPv4 != "" {
			logger.UtilLog.Infoln("parsing ServerIPv4 address from ENV variable")
		} else {
			context.BindingIPv4 = sbi.BindingIPv4
			if context.BindingIPv4 == "" {
				logger.UtilLog.Warnln("error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default")
				context.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	if configuration.NrfUri != "" {
		context.NrfUri = configuration.NrfUri
	} else {
		logger.UtilLog.Warnln("NRF Uri is empty. Using localhost as NRF IPv4 address")
		context.NrfUri = fmt.Sprintf("%s://%s:%d", context.UriScheme, "127.0.0.1", 29510)
	}
}
