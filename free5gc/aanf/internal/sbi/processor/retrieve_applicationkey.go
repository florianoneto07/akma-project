package processor

import (
	"errors"
	"net/http"
	"net/netip"
	"time"

	"github.com/gin-gonic/gin"

	aanf_context "github.com/florianoneto07/akma-free5gc/free5gc/aanf/internal/context"
	"github.com/florianoneto07/akma-free5gc/free5gc/aanf/internal/logger"
	"github.com/florianoneto07/akma-free5gc/free5gc/aanf/internal/util"
	"github.com/florianoneto07/akma-free5gc/shared/akma/afid"
	"github.com/florianoneto07/akma-free5gc/shared/akma/akid"
	"github.com/florianoneto07/akma-free5gc/shared/akma/codec"
	"github.com/free5gc/openapi/models"
)

// RetrieveApplicationKey returns the K_AF of the requesting AF
// (Naanf_AKMA_ApplicationKey_Get and _AnonUser_Get, TS 29.535 cl. 4.2.2.3.2;
// TS 33.535 cl. 6.2.1 and 6.2.2). remote is the source address of the AF.
func (p *Processor) RetrieveApplicationKey(c *gin.Context, req *models.AkmaAfKeyRequest, remote netip.Addr) {
	aanfCtx := p.Context()

	if pd := validateAkmaAfKeyRequest(req); pd != nil {
		util.WriteProblem(c, pd)
		return
	}
	af, err := afid.Parse(req.AfId)
	if err != nil {
		util.WriteProblem(c, util.Problem(http.StatusBadRequest, util.CauseMandatoryIEIncorrect, "",
			util.Param("/afId", "expected \"<FQDN>.<Ua* security protocol id>\"")))
		return
	}
	id, err := akid.Parse(req.AKId, aanfCtx.AKidOptions)
	if err != nil {
		util.WriteProblem(c, util.Problem(http.StatusBadRequest, util.CauseMandatoryIEIncorrect, "",
			util.Param("/aKId", "not an A-KID \"rid<RID>.<label><A-TID>@<realm>\"")))
		return
	}

	// TS 33.535 cl. 6.2.1 step 2: the AAnF checks, with its local policy,
	// whether it can provide the service to the AF, and otherwise rejects
	// the procedure. No application error is defined for this case
	// (TS 29.535 table 5.1.7.3-1): 403 without cause.
	if !aanfCtx.AfPolicy.Authorize(af, remote) {
		logger.AkmaLog.Warnf("AF not authorized for the requested AF_ID (source %s)", remote)
		logger.AkmaLog.Debugf("AF_ID %s rejected for source %s", af.Transport(), remote)
		util.WriteProblem(c, util.Problem(http.StatusForbidden, "", "AF not authorized by local policy"))
		return
	}

	// The roaming status is not subscribed (deviation DV-2), so
	// ROAMING_AKMA_SERVICE_DENIED is never returned. A K_AF is always
	// derivable when K_AKMA is present, so "204 No Content" is never returned
	// either (deviation DV-7).
	supi, kaf, expiry, err := aanfCtx.Store.ApplicationKey(id, af, aanfCtx.KafLifetime)
	if errors.Is(err, aanf_context.ErrKAkmaNotPresent) {
		logger.AkmaLog.Infof("K_AKMA not present for the A-KID")
		logger.AkmaLog.Debugf("K_AKMA not present for A-KID %s", req.AKId)
		util.WriteProblem(c, util.Problem(http.StatusForbidden, util.CauseKAkmaNotPresent, ""))
		return
	}
	if err != nil {
		logger.AkmaLog.Errorf("K_AF derivation failed: %v", err)
		util.WriteProblem(c, util.Problem(http.StatusInternalServerError, util.CauseSystemFailure, ""))
		return
	}

	exp := expiry.UTC()
	rsp := models.AkmaAfKeyData{
		Kaf:      codec.HexKey(kaf),
		Expiry:   &exp,
		SuppFeat: negotiatedFeatures(req.SuppFeat),
	}
	// TS 29.535 table 5.1.4.3.2-2 NOTE 2: "supi" is included except in the
	// anonymous user access procedure (TS 33.535 cl. 6.2.2).
	if !req.AnonInd {
		rsp.Supi = supi
	}
	logger.AkmaLog.Infof("K_AF provided, expires %s", exp.Format(time.RFC3339))
	logger.AkmaLog.Debugf("K_AF provided to %s for %s", af.Transport(), supi)
	util.WriteJSON(c, http.StatusOK, rsp)
}

// validateAkmaAfKeyRequest checks the mandatory attributes of TS 29.522
// table 5.14.5.3.2-1. "notifUri" belongs to the RoamingRestriction feature,
// which is not supported, and is ignored.
func validateAkmaAfKeyRequest(req *models.AkmaAfKeyRequest) *models.ProblemDetails {
	var missing []models.InvalidParam
	if req.AfId == "" {
		missing = append(missing, util.Param("/afId", "required"))
	}
	if req.AKId == "" {
		missing = append(missing, util.Param("/aKId", "required"))
	}
	if len(missing) > 0 {
		return util.Problem(http.StatusBadRequest, util.CauseMandatoryIEMissing, "", missing...)
	}
	if !suppFeatPattern.MatchString(req.SuppFeat) {
		return util.Problem(http.StatusBadRequest, util.CauseOptionalIEIncorrect, "",
			util.Param("/suppFeat", "not a SupportedFeatures string"))
	}
	return nil
}
