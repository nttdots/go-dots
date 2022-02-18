package data_controllers

import (
	"fmt"
	"strconv"
	"net/http"
  
	"github.com/julienschmidt/httprouter"
	"github.com/nttdots/go-dots/dots_server/db"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/dots_server/models/data"
	"github.com/nttdots/go-dots/dots_server/db_models/data"
	log "github.com/sirupsen/logrus"
	types    "github.com/nttdots/go-dots/dots_common/types/data"
	messages "github.com/nttdots/go-dots/dots_common/messages/data"
  )

type VendorMappingController struct {}

// Get vendor-mapping
func (v *VendorMappingController) Get(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
	cuid := p.ByName("cuid")
	log.WithField("cuid", cuid).Info("[VendorMappingController] GET")
	isAfterTransaction := false
  
	// Check missing 'cuid'
	if cuid == "" {
		errMsg := "Missing required path 'cuid' value."
	    log.Error(errMsg)
	    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, errMsg, isAfterTransaction)
	}
	return WithTransaction(func (tx *db.Tx) (Response, error) {
		return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
			return findVendorMapping(tx, &client.Id, client.Cuid, r, true)
		})
	})
}

// Get vendor-mapping of sever
func (v *VendorMappingController) GetVendorMappingOfServer(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
	log.Info("[VendorMappingController] GET")
	isAfterTransaction := false
	capabilities := getCapabilities()
	if *capabilities.Capabilities.VendorMappingEnabled == false {
		errMsg := "'vendor-mapping-enabled' is set to 'false'. Failed to Get the Dots server's vendor attack mapping details."
		log.Error(errMsg)
	    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
	}
	return WithTransaction(func (tx *db.Tx) (Response, error) {
		return findVendorMapping(tx, nil, "", r, true)
	})
}

// Put vendor-mapping
func (vc *VendorMappingController) Put(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
	var errMsg string
	cuid := p.ByName("cuid")
	vendorId := p.ByName("vendorId")
	log.WithField("cuid", cuid).Info("[VendorMappingController] PUT")
	isAfterTransaction := false
	// Check missing 'cuid'
	if cuid == "" {
		errMsg = "Missing required path 'cuid' value."
		log.Error(errMsg)
		return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, errMsg, isAfterTransaction)
	}
	if vendorId == "" {
		errMsg = "Missing required path 'vendor-id' value."
		log.Error(errMsg)
	    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, errMsg, isAfterTransaction)
	}
	req := messages.VendorMappingRequest{}
	err := Unmarshal(r, &req)
	if err != nil {
		errMsg = fmt.Sprintf("Invalid body data format: %+v", err)
		log.Error(errMsg)
		return ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
	}
	// Validate body data
	vId, err := strconv.Atoi(vendorId)
	if err != nil {
		errMsg := "Failed to convert 'vendor-id' to int"
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
	}
	errMsg = messages.ValidateWithVendorId(vId, &req)
	if errMsg != "" {
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
	}
	errMsg = messages.ValidateVendorMapping(&req)
	if errMsg != "" {
		return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, errMsg, isAfterTransaction)
	}
	return WithTransaction(func (tx *db.Tx) (Response, error) {
		return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
			isAfterTransaction = true
			// Find vendor-mapping by vendor-id
			e, err := data_models.FindVendorByVendorId(tx, client.Id, vId)
			if err != nil {
				errMsg = fmt.Sprintf("Failed to get vendor with 'vendor-id' = %+v. Error: %+v", vId, err)
				log.Errorf(errMsg)
				return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
			}
			if e.Id == 0 {
				errMsg := fmt.Sprintf("Not Found vendor-mapping by specified vendor-id = %+v", vId)
				log.Errorf(errMsg)
				return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
			}
			// Save attack-detail
			e.VendorName  = *req.VendorMapping.Vendor[0].VendorName
			if req.VendorMapping.Vendor[0].DescriptionLang != nil {
				e.DescriptionLang = *req.VendorMapping.Vendor[0].DescriptionLang
			} else {
				e.DescriptionLang = "en-US"
			}
			lastUpdated, _ := strconv.ParseUint(*req.VendorMapping.Vendor[0].LastUpdated, 10, 64)
			e.LastUpdated = lastUpdated
			e.AttackMapping = nil
			for _, am := range req.VendorMapping.Vendor[0].AttackMapping {
				attackMapping := data_models.AttackMapping{}
				attackMapping.AttackId = int(*am.AttackId)
				attackMapping.AttackDescription = *am.AttackDescription
				e.AttackMapping = append(e.AttackMapping, attackMapping)
			}
			err = e.Save(tx)
			if err != nil {
				errMsg = fmt.Sprintf("Failed to save vendor-mapping with vendor-id = %+v", vId)
				log.Errorf(errMsg)
				return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
			}
			return EmptyResponse(http.StatusNoContent)
		})
	})
}

// Delete one vendor-mapping
func (v *VendorMappingController) DeleteOne(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
	var errMsg string
	cuid := p.ByName("cuid")
	vendorId := p.ByName("vendorId")
	log.WithField("cuid", cuid).WithField("vendor-id", vendorId).Info("[VendorMappingController] DELETE One")
	isAfterTransaction := false
	// Check missing 'cuid'
	if cuid == "" {
		errMsg = "Missing required path 'cuid' value."
		log.Error(errMsg)
		return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, errMsg, isAfterTransaction)
	}
	if vendorId == "" {
		errMsg = "Missing required path 'vendor-id' value."
		log.Error(errMsg)
	    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, errMsg, isAfterTransaction)
	}
	vId, err := strconv.Atoi(vendorId)
	if err != nil {
		errMsg := "Failed to convert 'vendor-id' to int"
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
	}
	capabilities := getCapabilities()
	if *capabilities.Capabilities.VendorMappingEnabled == false {
		errMsg := "'vendor-mapping-enabled' is set to 'false'. Failed to Delete the Dots server's vendor attack mapping details."
		log.Error(errMsg)
	    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
	}
	return WithTransaction(func (tx *db.Tx) (Response, error) {
		return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
			// Delete vendor-mapping
			return deleteVendorAttackMapping(tx, client.Id, vId, true)
		})
	})
}

// Delete all vendor-mapping
func (v *VendorMappingController) DeleteAll(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
	var errMsg string
	cuid := p.ByName("cuid")
	log.WithField("cuid", cuid).Info("[VendorMappingController] DELETE All")
	isAfterTransaction := false
	// Check missing 'cuid'
	if cuid == "" {
		errMsg = "Missing required path 'cuid' value."
		log.Error(errMsg)
		return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, errMsg, isAfterTransaction)
	}
	capabilities := getCapabilities()
	if *capabilities.Capabilities.VendorMappingEnabled == false {
		errMsg := "'vendor-mapping-enabled' is set to 'false'. Failed to Delete the Dots server's vendor attack mapping details."
		log.Error(errMsg)
	    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
	}
	return WithTransaction(func (tx *db.Tx) (Response, error) {
		return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
			isAfterTransaction = true
			vendorList, err := data_models.FindVendorMappingByClientId(tx, &client.Id)
			if err != nil {
				return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, err.Error(), isAfterTransaction)
			}
			// If vendor-mapping doesn't exist, Dots server return 200
			if len(vendorList) <= 0 {
				return EmptyResponse(http.StatusNoContent)
			}
			// Delete vendor-mapping
			for _, vendor := range vendorList {
				res, err := deleteVendorAttackMapping(tx, client.Id, vendor.VendorId, isAfterTransaction)
				if res.Code != http.StatusNoContent {
					return res, err
				}
			}
			return EmptyResponse(http.StatusNoContent)
		})
	})
}

// Find vendor-mapping
func findVendorMapping(tx *db.Tx, clientId *int64, cuid string, r *http.Request, isAfterTransaction bool) (Response, error) {
	// Find vendor-mapping by client_id
	vendorList, err := data_models.FindVendorMappingByClientId(tx, clientId)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, err.Error(), isAfterTransaction)
	}
	if len(vendorList) < 1 {
		errMsg := fmt.Sprintf("Not Found vendor-mapping by specified dots-client = %+v", cuid)
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
	}
	q := r.URL.Query()
	var depth *int
	  if a, ok := q["depth"]; ok {
		tmpDepth, err := strconv.Atoi(a[0])
		if err != nil {
			errMsg := "Failed to convert 'depth' to int"
			log.Error(errMsg)
			return ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
		}
		depth = & tmpDepth
	} else {
		depth = nil
	}
	tv := vendorList.ToTypesVendorMapping(depth)
	s := messages.VendorMappingResponse{}
	s.VendorMapping = *tv
	cp, err := messages.ContentFromRequest(r)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, err.Error(), isAfterTransaction)
	}
	m, err := messages.ToMap(s, cp)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, err.Error(), isAfterTransaction)
	}
	return YangJsonResponse(m)
}

// Get vendor-mapping by cuid
func GetVendorMappingByCuid(customer *models.Customer, cuid string) (res *types.VendorMapping, err error) {
	_, err = WithTransaction(func (tx *db.Tx) (Response, error) {
		// Find client by cuid
		client, err := data_models.FindClientByCuid(tx, customer, cuid)
		if err != nil {
		  return Response{}, err
		}
		if client == nil {
			return Response{}, nil
		}
		// Find vendor-mapping by client_id
		vendorList, err := data_models.FindVendorMappingByClientId(tx, &client.Id)
		if err != nil {
			return Response{}, err
		}
		res = vendorList.ToTypesVendorMapping(nil)
		return Response{}, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete a vendor attack mapping
func deleteVendorAttackMapping(tx *db.Tx, clientId int64, vId int, isAfterTransaction bool) (Response, error) {
	errMsg := ""
	// Find vendor-mapping by vendor-id
	e, err := data_models.FindVendorByVendorId(tx, clientId, vId)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to get vendor with 'vendor-id' = %+v. Error: %+v", vId, err)
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
	}
	if e.Id == 0 {
		errMsg := fmt.Sprintf("Not Found vendor-mapping by specified vendor-id = %+v", vId)
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
	}
	// Delete vendor-mapping by id
	err = data_db_models.DeleteVendorMappingById(tx.Session, e.Id)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to delete vendor with 'vendor-id' = %+v. Error: %+v", vId, err)
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
	}
	// Delete attack-mapping by vendor-mapping id
	err = data_db_models.DeleteAttackMappingByVendorMappingId(tx.Session, e.Id)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to delete attack-mapping with 'vendor-id' = %+v. Error: %+v", vId, err)
		log.Errorf(errMsg)
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
	}
	return EmptyResponse(http.StatusNoContent)
}