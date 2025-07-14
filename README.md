# goest
Go Easy SignTool - A cross-platform tool for code signing based on [Relic from sassoftware](sassoftware/relic).

## Installation
### Binary Releases
You can download pre-built binaries from the [Releases page](github.com/ossign/goest/releases).

### From Source
To build from source, you need to have Go 1.24+ installed. Then run:

```bash
go mod tidy
go build -o goest ./cmd/
```

### From package repository
For Debian, Ubuntu, RPM-based distributions and Alpine Linux, you can use the provided package repositories:

#### Debian/Ubuntu
```bash
sudo curl https://pkg.ossign.org/debian/repository.key -o /etc/apt/keyrings/gitea-ossign.asc
echo "deb [signed-by=/etc/apt/keyrings/gitea-ossign.asc] https://pkg.ossign.org/debian all main" | sudo tee -a /etc/apt/sources.list.d/ossign.list
sudo apt update

sudo apt install goest
```

#### RPM-based distributions
```bash
# on RedHat based distributions
dnf config-manager --add-repo https://pkg.ossign.org/rpm.repo

# on SUSE based distributions
zypper addrepo https://pkg.ossign.org/rpm.repo

```

#### Alpine Linux
```bash
echo "https://pkg.ossign.org/alpine/all/repository" | sudo tee /etc/apk/repositories
curl -JO https://pkg.ossign.org/alpine/key

apk add goest
```

## Usage
### Azure Key Vault
To sign a file with Azure Key Vault, use the following command:
```bash
# Using CLI Flags
goest sign azurekv --url https://<your-key-vault-name>.vault.azure.net/certificates/your-certificate/certificate-hash/ --tenant <your-tenant-uuid> --client <your-client-id> --secret <your-client-secret> file1.ps1 file2.exe file3.dll

# Using environment variables
export goest_AZUREKV_URL="https://<your-key-vault-name>.vault.azure.net/certificates/your-certificate/certificate-hash/"
export goest_AZUREKV_TENANT="<your-tenant-uuid>"
export goest_AZUREKV_CLIENT="<your-client-id>"
export goest_AZUREKV_SECRET="<your-client-secret>"

goest sign azurekv file1.ps1 file2.exe file3.dll

# Using config file
echo <<EOF > ./goest-example-config.yaml
azurekv:
  url: https://<your-key-vault-name>.vault.azure.net/certificates/your-certificate/certificate-hash/
  tenant: <your-tenant-uuid>
  client: <your-client-id>
  secret: <your-client-secret>
EOF

goest sign azurekv --config ./goest-example-config.yaml file1.ps1 file2
```
