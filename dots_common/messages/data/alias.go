package data_messages

import (
  "fmt"
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
)

type AliasesRequest struct {
  Aliases types.Aliases `json:"ietf-dots-data-channel:aliases"`
}

type AliasesResponse struct {
  Aliases types.Aliases `json:"ietf-dots-data-channel:aliases"`
}

func (r *AliasesRequest) Validate(method string) (bool, string) {
  errorMsg := ""

  if len(r.Aliases.Alias) <= 0 {
    log.WithField("len", len(r.Aliases.Alias)).Error("'alias' is not exist.")
    errorMsg = fmt.Sprintf("Body Data Error : 'alias' is not exist")
    return false, errorMsg
  }
  if len(r.Aliases.Alias) > 1 {
    log.WithField("len", len(r.Aliases.Alias)).Error("multiple 'alias'.")
    errorMsg = fmt.Sprintf("Body Data Error : Have multiple 'alias' (%d)", len(r.Aliases.Alias))
    return false, errorMsg
  }

  alias := r.Aliases.Alias[0]
  if alias.Name == "" {
    log.Error("Missing required alias 'name' attribute.")
    errorMsg = fmt.Sprintf("Body Data Error : Missing alias 'name'")
    return false, errorMsg
  }

  if alias.PendingLifetime != nil {
    log.WithField("pending-lifetime", *alias.PendingLifetime).Error("'pending-lifetime' found.")
    errorMsg = fmt.Sprintf("Body Data Error : Found NoConfig Attribute 'pending-lifetime' (%v)", *alias.PendingLifetime)
    return false, errorMsg
  }

  if method == "POST" && (len(alias.TargetPrefix) == 0 && len(alias.TargetFQDN) == 0 && len(alias.TargetURI) == 0) {
    log. Error("At least one of the 'target-prefix', 'target-fqdn', or 'target-uri' attributes MUST be present.")
    errorMsg = fmt.Sprintf("Body Data Error : At least one of the 'target-prefix', 'target-fqdn', or 'target-uri' attributes MUST be present")
    return false, errorMsg
  } else if method == "PUT" && len(alias.TargetPrefix) == 0 {
    log. Error("Missing required 'target-prefix' attribute.")
    errorMsg = fmt.Sprintf("Body Data Error : Missing 'target-prefix'")
    return false, errorMsg
  }


  for _, pr := range alias.TargetPortRange {
    if pr.UpperPort != nil {
      if *pr.UpperPort < pr.LowerPort {
        log.WithField("lower-port", pr.LowerPort).WithField("upper-port", *pr.UpperPort).Error("'upper-port' must be greater than or equal to 'lower-port'.")
        errorMsg = fmt.Sprintf("Body Data Error : 'upper-port' must be greater than or equal to 'lower-port'")
        return false, errorMsg
      }
    }
  }

  return true, ""
}

func (r *AliasesRequest) ValidateWithName(name string, method string) (bool, string) {
  bValid, errorMsg := r.Validate(method)
  if !bValid {
    return false, errorMsg
  }

  alias := r.Aliases.Alias[0]
  if alias.Name != name {
    log.WithField("name(req)", alias.Name).WithField("name(URI)", name).Error("request/URI name mismatch.")
    errorMsg = fmt.Sprintf("Request/URI name mismatch : (%v) / (%v)", alias.Name, name)
    return false, errorMsg
  }

  return true, ""
}