#!/bin/sh
export VAULT_TOKEN=root
export VAULT_ADDR=http://127.0.0.1:8200

vault secrets enable pki
vault secrets tune -max-lease-ttl=8760h pki

vault write -field=certificate pki/root/generate/internal \
    common_name="blueprint.local" \
    ttl=8760h > root_ca.crt

vault write pki/config/urls \
    issuing_certificates="http://openbao-vault-0.openbao.svc:8200/v1/pki/ca" \
    crl_distribution_points="http://openbao-vault-0.openbao.svc:8200/v1/pki/crl"

vault write pki/roles/blueprint-dot-local \
    allowed_domains="blueprint.local" \
    allow_subdomains=true \
    max_ttl="720h"

# Create policy for cert-manager
cat <<'POLICY' > /tmp/cert-manager-policy.hcl
path "pki*" { capabilities = ["read", "list"] }
path "pki/sign/blueprint-dot-local" { capabilities = ["create", "update"] }
path "pki/issue/blueprint-dot-local" { capabilities = ["create"] }
POLICY
vault policy write cert-manager /tmp/cert-manager-policy.hcl

# Enable AppRole for simple token-less auth (or use a direct token)
vault auth enable kubernetes
vault write auth/kubernetes/config \
    kubernetes_host="https://kubernetes.default.svc:443"

vault write auth/kubernetes/role/cert-manager \
    bound_service_account_names="*" \
    bound_service_account_namespaces="*" \
    policies=cert-manager \
    ttl=24h

echo "VAULT SETUP COMPLETE"
