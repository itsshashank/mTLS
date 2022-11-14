openssl ecparam -genkey -name prime256v1 -out server.key

openssl req -new -SHA384 -key server.key -nodes -out server.csr

openssl req -in server.csr -noout -text

openssl x509 -req -SHA384 -extfile v.ext -days 365 -in server.csr -CA /workspaces/mtlscert/cert/ca-cert.pem -CAkey /workspaces/mtlscert/cert/ca-key.pem -CAcreateserial -out server.pem