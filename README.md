# apprecert: Local Certificate Authority Management Tool

## Table of Contents

- [apprecert: Local Certificate Authority Management Tool](#apprecert-local-certificate-authority-management-tool)
	- [Table of Contents](#table-of-contents)
	- [Overview](#overview)
	- [Features](#features)
	- [Prerequisites](#prerequisites)
		- [Required](#required)
		- [Supported Operating Systems](#supported-operating-systems)
		- [Optional Integrations](#optional-integrations)
	- [Installation](#installation)
		- [Clone the Repository](#clone-the-repository)
		- [Build the Tool](#build-the-tool)
			- [For Current Platform](#for-current-platform)
			- [Cross-Compilation (e.g., Windows from Linux/macOS)](#cross-compilation-eg-windows-from-linuxmacos)
	- [Usage](#usage)
		- [Certificate Storage Locations](#certificate-storage-locations)
		- [Generating Root CA](#generating-root-ca)
		- [Generating Host Certificates](#generating-host-certificates)
		- [Managing Trust Stores](#managing-trust-stores)
			- [Install CA](#install-ca)
			- [Uninstall CA](#uninstall-ca)
	- [Generating Certificates for Different Platforms and Ecosystems](#generating-certificates-for-different-platforms-and-ecosystems)
		- [macOS Certificate Generation](#macos-certificate-generation)
		- [Linux Certificate Generation](#linux-certificate-generation)
		- [Windows Certificate Generation](#windows-certificate-generation)
		- [Docker Certificate Generation](#docker-certificate-generation)
		- [Git Certificate Generation](#git-certificate-generation)
		- [Java Keystore Certificate Generation](#java-keystore-certificate-generation)
		- [Node.js Certificate Generation](#nodejs-certificate-generation)
		- [NSS (Network Security Services) Certificate Generation](#nss-network-security-services-certificate-generation)
		- [Python Certificate Generation](#python-certificate-generation)
		- [Multi-Environment Certificate Generation](#multi-environment-certificate-generation)
	- [Advanced Certificate Generation](#advanced-certificate-generation)
		- [Subject Alternative Name (SAN) Certificates](#subject-alternative-name-san-certificates)
		- [Specifying Certificate Validity](#specifying-certificate-validity)
	- [Kubernetes Integration](#kubernetes-integration)
		- [How It Works](#how-it-works)
		- [k8s Prerequisites](#k8s-prerequisites)
		- [k8s Installation](#k8s-installation)
		- [Removal](#removal)
		- [Manual Verification and Troubleshooting](#manual-verification-and-troubleshooting)
			- [Checking Existing ConfigMap](#checking-existing-configmap)
			- [Common Error Handling](#common-error-handling)
			- [Resetting Cluster Configuration](#resetting-cluster-configuration)
			- [Restart Cluster Components](#restart-cluster-components)
		- [Advanced Usage](#advanced-usage)
			- [Custom CA Directory](#custom-ca-directory)
			- [Manual Patch Application](#manual-patch-application)
			- [Debugging](#debugging)
		- [Best Practices](#best-practices)
	- [General Best Practices](#general-best-practices)
	- [Troubleshooting Certificate Generation](#troubleshooting-certificate-generation)

## Overview

`apprecert` is a command-line tool designed to simplify the generation and management of a local Certificate Authority (CA) for development environments. It helps developers create and trust locally-issued TLS certificates without manually configuring trust stores.

## Features

- Generate a root CA certificate and private key
- Create host-specific TLS certificates signed by your custom CA
- Automatically install/uninstall the CA into various platform trust stores
- Support for Linux, macOS, and Windows
- Integration with multiple ecosystem tools (Java, Firefox, Node.js, Python, Git)

## Prerequisites

### Required

- Go (version 1.22 or newer recommended)
- Git

### Supported Operating Systems

- Linux
- macOS
- Windows

### Optional Integrations

- Java (requires `JAVA_HOME` to be set)
- Python 3
- Git
- Node.js
- Docker
- NSS (Network Security Services)
  - Mozilla Firefox
  - Thunderbird
  - Google Chrome (Linux/macOS)
  - LibreOffice
  - Postfix
  - Sendmail
  - OpenVPN
  - Curl (NSS builds)
  - Apache HTTP Server (mod_nss)
  - Red Hat/Fedora Linux
  - Evolution Email Client
  - NetworkManager
  - Pidgin

## Installation

### Clone the Repository

```bash
git clone https://github.com/appremon/apprecert.git
cd apprecert
```

### Build the Tool

#### For Current Platform

```bash
go build -o apprecert ./cmd/apprecert
```

#### Cross-Compilation (e.g., Windows from Linux/macOS)

```bash
GOOS=windows GOARCH=amd64 go build -o apprecert.exe ./cmd/apprecert
```

## Usage

### Certificate Storage Locations

By default, apprecert stores certificates in platform-specific directories:

- Linux/Unix: `~/.local/share/apprecert`
- macOS: `~/Library/Application Support/Apprecert`
- Windows: `%LOCALAPPDATA%\Apprecert`

You can override this by setting the `CAROOT` environment variable:

```bash
export CAROOT="/path/to/custom/dir"
```

### Generating Root CA

To generate and install the root CA:

```bash
./apprecert -install
```

This command will:

- Create a root CA certificate (`rootCA.pem`)
- Create a private key (`rootCA-key.pem`)
- Attempt to install the CA into system trust stores

### Generating Host Certificates

To generate a certificate for a specific hostname:

```bash
./apprecert myservice.local
```

This creates:

- Host certificate: `myservice.local-cert.pem`
- Private key: `myservice.local-key.pem`

### Managing Trust Stores

#### Install CA

```bash
./apprecert -install
```

- Adds CA to system trust stores
- Attempts integration with Java, Firefox, and other supported stores

#### Uninstall CA

```bash
./apprecert -uninstall
```

- Removes CA from system trust stores

## Generating Certificates for Different Platforms and Ecosystems

> Note: The following sections provide detailed certificate generation instructions for various platforms and ecosystems.

### macOS Certificate Generation

```bash
# Generate macOS-specific host certificate
./apprecert myservice.mac.local

# For local development servers
./apprecert localhost
./apprecert 127.0.0.1
```

### Linux Certificate Generation

```bash
# Generate Linux-specific host certificates
./apprecert myservice.linux.local
./apprecert localhost
./apprecert 127.0.0.1

# For local development environments
./apprecert dev.local
./apprecert testing.local
```

### Windows Certificate Generation

```powershell
# Generate Windows-specific host certificates
.\apprecert.exe myservice.win.local
.\apprecert.exe localhost
.\apprecert.exe 127.0.0.1

# For local development domains
.\apprecert.exe dev.local
.\apprecert.exe testing.local
```

### Docker Certificate Generation

```bash
# Generate certificates for Docker containers
./apprecert docker.local
./apprecert localhost
./apprecert host.docker.internal

# For Docker Compose services
./apprecert service1.docker.local
./apprecert service2.docker.local
```

### Git Certificate Generation

```bash
# Generate certificates for Git repositories and services
./apprecert git.local
./apprecert github.local
./apprecert gitlab.local

# Configure Git to use the generated CA
git config --global http.sslCAInfo "$CAROOT/rootCA.pem"
```

### Java Keystore Certificate Generation

```bash
# Ensure JAVA_HOME is set
export JAVA_HOME=/path/to/java/jdk

# Generate Java-specific certificates
./apprecert java.local
./apprecert localhost.java

# Manually add to Java keystore if needed
keytool -importcert -alias apprecert-root \
        -file "$CAROOT/rootCA.pem" \
        -keystore "$JAVA_HOME/lib/security/cacerts" \
        -storepass changeit
```

### Node.js Certificate Generation

```bash
# Generate certificates for Node.js applications
./apprecert nodejs.local
./apprecert localhost.node

# Configure Node.js to trust the CA
export NODE_EXTRA_CA_CERTS="$CAROOT/rootCA.pem"
```

### NSS (Network Security Services) Certificate Generation

```bash
# Generate certificates for NSS-based applications (Firefox, Thunderbird)
./apprecert nss.local
./apprecert firefox.local

# Manually add to NSS profiles if needed
certutil -A -d sql:$HOME/.mozilla/firefox/*.default \
         -t "C,," -n "apprecert-rootCA" \
         -i "$CAROOT/rootCA.pem"
```

### Python Certificate Generation

```bash
# Generate certificates for Python applications
./apprecert python.local
./apprecert localhost.py

# Append to Python's certifi bundle
cat "$CAROOT/rootCA.pem" | sudo tee -a "$(python3 -m certifi)"
```

### Multi-Environment Certificate Generation

```bash
# Generate wildcard certificates for multiple environments
./apprecert "*.dev.local"
./apprecert "*.test.local"
./apprecert "*.staging.local"
```

## Advanced Certificate Generation

### Subject Alternative Name (SAN) Certificates

```bash
# Generate certificates with multiple hostnames
./apprecert myservice.local \
            localhost \
            127.0.0.1 \
            host.docker.internal
```

### Specifying Certificate Validity

```bash
# Set custom validity period (if supported by future versions)
./apprecert myservice.local --days 365
```

## Kubernetes Integration

The `apprecert` tool supports integrating a custom Certificate Authority (CA) into a Kubernetes cluster by creating and managing a ConfigMap containing the CA certificate. This ensures that the cluster recognizes and trusts certificates issued by the custom CA.

### How It Works

1. **Certificate Management**:
   - The CA certificate is located at `CAROOT/rootCA.pem` (or as defined by the `CAROOT` environment variable)

2. **Kubernetes ConfigMap**:
   - A ConfigMap named `custom-ca-bundle` is created in the `kube-system` namespace
   - This ConfigMap includes the `rootCA.pem` certificate as `ca.crt`

3. **Cluster Configuration Patch**:
   - The `kubeadm-config` ConfigMap in the `kube-system` namespace is patched to reference the custom CA bundle directory (`/etc/kubernetes/pki/custom-ca-bundle`)

### k8s Prerequisites

- A running Kubernetes cluster with `kubectl` installed and configured
- Administrator permissions to the cluster

### k8s Installation

To install the custom CA into a Kubernetes cluster:

1. Run the installation command:

   ```bash
   ./apprecert -install
   ```

2. This will:
   - Check for the presence of the `rootCA.pem` certificate
   - Create a ConfigMap named `custom-ca-bundle` in the `kube-system` namespace
   - Patch the `kubeadm-config` ConfigMap to include the custom CA bundle path

3. Verify the ConfigMap:

   ```bash
   kubectl get configmap custom-ca-bundle -n kube-system
   ```

4. Confirm the `kubeadm-config` ConfigMap has been patched:

   ```bash
   kubectl get configmap kubeadm-config -n kube-system -o yaml
   ```

   Ensure that the `ClusterConfiguration` includes:

   ```yaml
   certificatesDir: /etc/kubernetes/pki/custom-ca-bundle
   ```

### Removal

To uninstall the custom CA from the Kubernetes cluster:

1. Run the uninstall command:

   ```bash
   ./apprecert -uninstall
   ```

2. This will:
   - Delete the `custom-ca-bundle` ConfigMap in the `kube-system` namespace
   - Ensure no lingering references to the custom CA remain

3. Verify the ConfigMap has been removed:

   ```bash
   kubectl get configmap custom-ca-bundle -n kube-system
   ```

   The output should indicate the ConfigMap does not exist.

### Manual Verification and Troubleshooting

#### Checking Existing ConfigMap

Before creating the ConfigMap, the tool checks for its existence. To manually check:

```bash
kubectl get configmap custom-ca-bundle -n kube-system
```

#### Common Error Handling

1. **"ConfigMap 'custom-ca-bundle' already exists"**:
   Delete the existing ConfigMap if you want to replace it:

   ```bash
   kubectl delete configmap custom-ca-bundle -n kube-system
   ```

2. **"Failed to patch cluster configuration"**:
   - Ensure you have administrator privileges
   - Verify the cluster supports `kubeadm`

#### Resetting Cluster Configuration

To reset the `kubeadm-config` to its original state, manually edit the ConfigMap:

```bash
kubectl edit configmap kubeadm-config -n kube-system
```

#### Restart Cluster Components

After modifying the trust store or configuration, restart cluster components to ensure changes take effect.

### Advanced Usage

#### Custom CA Directory

If your Kubernetes cluster uses a custom directory for CA certificates, ensure the path matches the patched `certificatesDir` in `kubeadm-config`.

#### Manual Patch Application

If automatic patching fails, manually apply the patch:

```bash
kubectl patch configmap kubeadm-config -n kube-system \
  --type=json \
  -p='[{"op": "add", "path": "/data/ClusterConfiguration/certificatesDir", "value": "/etc/kubernetes/pki/custom-ca-bundle"}]'
```

#### Debugging

Use verbose logging to debug issues:

```bash
./apprecert -install | tee debug.log
```

### Best Practices

- Always backup your existing cluster configuration before making changes
- Verify the CA certificate before installation
- Use minimal necessary permissions
- Monitor cluster health after CA integration
-

## General Best Practices

1. Always use unique, local-only domain names
2. Avoid using production domain names
3. Keep your root CA private and secure
4. Regenerate certificates periodically
5. Remove unused certificates

## Troubleshooting Certificate Generation

- Ensure unique hostnames across your development environment
- Check that `CAROOT` is set and accessible
- Verify platform-specific trust store requirements
- Restart applications after certificate generation

---

**Note**: For the most up-to-date information and support, please visit the [project repository](https://github.com/appremon/apprecert).
