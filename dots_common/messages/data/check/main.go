package main

import "encoding/json"
import "fmt"
import "github.com/davecgh/go-spew/spew"
import "github.com/nttdots/go-dots/dots_common/messages/data"

func check(s string, d interface{}) {
  err := json.Unmarshal([]byte(s), &d)

  fmt.Printf("%s\n => ", s)
  if err != nil {
    fmt.Printf("%v\n", err)
  } else {
    spew.Dump(d)
  }
}

func main() {

  check(`
    {
      "ietf-dots-data-channel:dots-client": [
        {
          "cuid": "dz6pHjaADkaFTbjr0JGBpw",
          "cdid": "7eeaf349529eb55ed50113"
        }
      ]
    }
  `, data_messages.ClientRequest{})

  check(`
            {
              "ietf-dots-data-channel:dots-client": [
                {
                  "cuid": "dz6pHjaADkaFTbjr0JGBpw"
                }
              ]
            }
  `, data_messages.ClientRequest{})

  check(`
   {
     "ietf-dots-data-channel:aliases": {
       "alias": [
         {
           "name": "https1",
           "target-protocol": [
             6
           ],
           "target-prefix": [
             "2001:db8:6401::1/128",
             "2001:db8:6401::2/128"
           ],
           "target-port-range": [
             {
               "lower-port": 443
             }
           ]
         }
       ]
     }
   }
  `, data_messages.AliasesRequest{})

  check(`
   {
     "ietf-dots-data-channel:aliases": {
       "alias": [
         {
           "name": "Server1",
           "target-protocol": [
             6
           ],
           "target-prefix": [
             "2001:db8:6401::1/128",
             "2001:db8:6401::2/128"
           ],
           "target-port-range": [
             {
               "lower-port": 443
             }
           ],
           "pending-lifetime": 3596
         },
         {
           "name": "Server2",
           "target-protocol": [
             6
           ],
           "target-prefix": [
             "2001:db8:6401::10/128",
             "2001:db8:6401::20/128"
           ],
           "target-port-range": [
             {
               "lower-port": 80
             }
           ],
           "pending-lifetime": 9869
         }
       ]
     }
   }
  `, data_messages.AliasesResponse{})

  check(`
    {
     "ietf-dots-data-channel:capabilities": {
       "address-family": ["ipv4", "ipv6"],
       "forwarding-actions": ["drop", "accept"],
       "rate-limit": true,
       "fragment": ["v4-fragment", "v6-fragment"],
       "transport-protocols": [1, 6, 17, 58],
       "ipv4": {
         "length": true,
         "protocol": true,
         "destination-prefix": true,
         "source-prefix": true
       },
       "ipv6": {
         "length": true,
         "protocol": true,
         "destination-prefix": true,
         "source-prefix": true
       },
       "tcp": {
         "flags": true,
         "source-port": true,
         "destination-port": true,
         "port-range": true
       },
       "udp": {
         "length": true,
         "source-port": true,
         "destination-port": true,
         "port-range": true
       },
       "icmp": {
         "type": true,
         "code": true
       }
     }
   }
  `, data_messages.CapabilitiesResponse{})

  check(`
    {
     "ietf-dots-data-channel:acls": {
       "acl": [
         {
           "name": "sample-ipv4-acl",
           "type": "ipv4-acl-type",
           "activation-type": "activate-when-mitigating",
           "aces": {
             "ace": [
               {
                 "name": "rule1",
                 "matches": {
                   "ipv4": {
                     "destination-ipv4-network": "198.51.100.0/24",
                     "source-ipv4-network": "192.0.2.0/24"
                     ,"flags": "reserved more fragment"
                   }
                 },
                 "actions": {
                   "forwarding": "drop"
                 }
               }
             ]
           }
         }
       ]
     }
    }
  `, data_messages.ACLsRequest{})
}
