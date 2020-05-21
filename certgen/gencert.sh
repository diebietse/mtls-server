#!/bin/bash

if [ -z "$1" ]; then
  echo 1>&2 "Usage: $0 <cert-name>"
  echo 1>&2 "Ex: $0 my-phone-cert"
  exit 1
fi

CERT_NAME=$1

TMP_DIR=$(mktemp -d -t cert-XXXXXXXXXX)

echo "Creating '${TMP_DIR}/cert.json' containing cert config"
cat <<EOT > ${TMP_DIR}/cert.json
{
  "CN": "${CERT_NAME}",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "ZA",
      "L": "Cape Town",
      "O": "diebietse.com",
      "OU": "www",
      "ST": "Western Cape"
    }
  ]
}
EOT
echo "Finished creating cert.json"

cp ca.json config.json root.pem root-key.pem ${TMP_DIR}

pushd ${TMP_DIR}

# Create and sign a cert
cfssl genkey ./cert.json | cfssljson -bare ${CERT_NAME}
cfssl sign -config ./config.json -ca root.pem -ca-key root-key.pem ${CERT_NAME}.csr | cfssljson -bare ${CERT_NAME}
# Create p12 file
openssl pkcs12 -export -nodes -clcerts -in ${CERT_NAME}.pem -inkey ${CERT_NAME}-key.pem -out ${CERT_NAME}.p12 -passout pass:mtls

popd

cp ${TMP_DIR}/${CERT_NAME}.p12 .
rm -r ${TMP_DIR}
