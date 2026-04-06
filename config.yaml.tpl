db:
  host: db-host
  port: db-port
  user: dbuser
  password: dbuser-password
  dbname: dbname
  sslMode: require  # require | verify-full | disable (disable only for local dev)

auth:
  issuer: issuer-url
  clientId: oidc-client-id
  clientSecret: oidc-client-secret
  redirectUrl: oidc-redirect-url

server:
  port: 8080
  baseUrl: base-url
