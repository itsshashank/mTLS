# Create CA private key and self-signed certificate
# adding -nodes to not encrypt the private key
openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca-key2.pem -out ca-cert2.pem

echo "CA's self-signed certificate"
openssl x509 -in ca-cert2.pem -noout -text 

# Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-key2.pem -out client-req2.pem

#  Sign the Client Certificate Request (CSR)
openssl x509 -req -in client-req2.pem -days 60 -CA ca-cert2.pem -CAkey ca-key2.pem -CAcreateserial -out client-cert2.pem -extfile client-ext.conf

echo "Client's signed certificate"
openssl x509 -in client-cert2.pem -noout -text
