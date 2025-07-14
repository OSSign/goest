# gost
Go SignTool - A cross-platform tool for code signing based on [Relic from sassoftware](sassoftware/relic).

## Usage
### Azure Key Vault
To sign a file with Azure Key Vault, use the following command:
```bash
# Using CLI Flags
gost sign azurekv --url https://<your-key-vault-name>.vault.azure.net/certificates/your-certificate/certificate-hash/ --tenant <your-tenant-uuid> --client <your-client-id> --secret <your-client-secret> file1.ps1 file2.exe file3.dll

# Using environment variables
export GOST_AZUREKV_URL="https://<your-key-vault-name>.vault.azure.net/certificates/your-certificate/certificate-hash/"
export GOST_AZUREKV_TENANT="<your-tenant-uuid>"
export GOST_AZUREKV_CLIENT="<your-client-id>"
export GOST_AZUREKV_SECRET="<your-client-secret>"

gost sign azurekv file1.ps1 file2.exe file3.dll

# Using config file
echo <<EOF > ./gost-example-config.yaml
azurekv:
  url: https://<your-key-vault-name>.vault.azure.net/certificates/your-certificate/certificate-hash/
  tenant: <your-tenant-uuid>
  client: <your-client-id>
  secret: <your-client-secret>
EOF

gost sign azurekv --config ./gost-example-config.yaml file1.ps1 file2
```
