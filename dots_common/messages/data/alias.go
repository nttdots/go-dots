package data_messages

import (
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
)

type AliasesRequest struct {
  Aliases types.Aliases `json:"ietf-dots-data-channel:aliases"`
}

type AliasesResponse struct {
  Aliases types.Aliases `json:"ietf-dots-data-channel:aliases"`
}

func (r *AliasesRequest) Validate(method string) bool {
  if len(r.Aliases.Alias) <= 0 {
    log.WithField("len", len(r.Aliases.Alias)).Error("'alias' is not exist.")
    return false
  }
  if len(r.Aliases.Alias) > 1 {
    log.WithField("len", len(r.Aliases.Alias)).Error("multiple 'alias'.")
    return false
  }

  alias := r.Aliases.Alias[0]
  if alias.Name == "" {
    log.Error("Missing required alias 'name' attribute.")
    return false
  }

  if alias.PendingLifetime != nil {
    log.WithField("pending-lifetime", *alias.PendingLifetime).Error("'pending-lifetime' found.")
    return false
  }

  if method == "POST" {
    if len(alias.TargetPrefix) == 0 && len(alias.TargetFQDN) == 0 && len(alias.TargetURI) == 0 {
      log. Error("At least one of the 'target-prefix', 'target-fqdn', or 'target-uri' attributes MUST be present.")
      return false
    }
  }

  for _, pr := range alias.TargetPortRange {
    if pr.UpperPort != nil {
      if *pr.UpperPort < pr.LowerPort {
        log.WithField("lower-port", pr.LowerPort).WithField("upper-port", *pr.UpperPort).Error("'upper-port' must be greater than or equal to 'lower-port'.")
        return false
      }
    }
  }

  return true
}
