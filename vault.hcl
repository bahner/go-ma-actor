storage file {
  path = "./vault"
}

listener tcp {
  address = "127.0.0.1:8200"
  tls_disable = 1
}

api_addr = "http://127.0.0.1:8200"

default_lease_ttl = "168h"
max_lease_ttl = "720h"
ui = true
