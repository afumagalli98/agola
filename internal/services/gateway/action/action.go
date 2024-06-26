// Copyright 2019 Sorint.lab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"github.com/rs/zerolog"

	scommon "agola.io/agola/internal/services/common"
	"agola.io/agola/internal/util"
	csclient "agola.io/agola/services/configstore/client"
	nsclient "agola.io/agola/services/notification/client"
	rsclient "agola.io/agola/services/runservice/client"
)

func APIErrorFromRemoteError(err error, options ...util.APIErrorOption) error {
	opts := []util.APIErrorOption{util.WithAPIErrorCallerDepth(2)}
	opts = append(opts, options...)

	rerr, ok := util.AsRemoteError(err)
	if !ok {
		return util.NewAPIErrorWrap(util.ErrInternal, err, opts...)
	}

	// assume remote error from internal services can be propagated from the gateway
	for _, detailedError := range rerr.DetailedErrors {
		opts = append(opts, util.WithAPIErrorDetailedError(util.NewAPIDetailedError(util.ErrorCode(detailedError.Code), util.WithAPIDetailedErrorDetails(detailedError.Details))))
	}

	return util.NewAPIErrorWrap(util.KindFromRemoteError(err), err, opts...)
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type ActionHandler struct {
	log                          zerolog.Logger
	sd                           *scommon.TokenSigningData
	sc                           *scommon.CookieSigningData
	configstoreClient            *csclient.Client
	runserviceClient             *rsclient.Client
	notificationClient           *nsclient.Client
	agolaID                      string
	apiExposedURL                string
	webExposedURL                string
	unsecureCookies              bool
	organizationMemberAddingMode OrganizationMemberAddingMode
}

type OrganizationMemberAddingMode string

const (
	OrganizationMemberAddingModeDirect     OrganizationMemberAddingMode = "direct"
	OrganizationMemberAddingModeInvitation OrganizationMemberAddingMode = "invitation"
)

func NewActionHandler(log zerolog.Logger, sd *scommon.TokenSigningData, sc *scommon.CookieSigningData, configstoreClient *csclient.Client, runserviceClient *rsclient.Client, notificationClient *nsclient.Client, agolaID, apiExposedURL, webExposedURL string, unsecureCookies bool, organizationMemberAddingMode OrganizationMemberAddingMode) *ActionHandler {
	return &ActionHandler{
		log:                          log,
		sd:                           sd,
		sc:                           sc,
		configstoreClient:            configstoreClient,
		runserviceClient:             runserviceClient,
		notificationClient:           notificationClient,
		agolaID:                      agolaID,
		apiExposedURL:                apiExposedURL,
		webExposedURL:                webExposedURL,
		unsecureCookies:              unsecureCookies,
		organizationMemberAddingMode: organizationMemberAddingMode,
	}
}
