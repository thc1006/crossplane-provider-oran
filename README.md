# Crossplane Provider for O-RAN Optical Devices

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/thc1006/crossplane-provider-oran)](https://goreportcard.com/report/github.com/thc1006/crossplane-provider-oran)

A Kubernetes operator and Crossplane provider for managing O-RAN compliant optical network devices with cloud-native GitOps workflows.

## Overview

This project implements a complete Kubernetes operator that manages optical network devices in O-RAN (Open Radio Access Network) environments. It follows the Kubernetes Operator pattern and integrates with Crossplane for infrastructure-as-code capabilities.

### Key Features

- **Kubernetes Native**: Custom Resource Definitions (CRDs) for optical device management
- **Declarative Configuration**: GitOps-ready declarative API for device configuration
- **Status Reporting**: Real-time status updates with conditions and observed state
- **Prometheus Integration**: Built-in metrics for monitoring and observability
- **Hardware Abstraction**: Clean separation between control plane and data plane
- **Crossplane Compatible**: Composable infrastructure with Crossplane XRDs

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                       │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │              Crossplane (Optional)                  │   │
│  │  ┌──────────────┐         ┌──────────────┐        │   │
│  │  │     XRD      │────────▶│ Composition  │        │   │
│  │  └──────────────┘         └──────────────┘        │   │
│  └────────────────────────────────────────────────────┘   │
│                             │                              │
│                             ▼                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │            OpticalDevice CRD (hardware.ran)        │   │
│  └────────────────────────────────────────────────────┘   │
│                             │                              │
│                             ▼                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │         O-RAN Adapter Controller                   │   │
│  │  - Reconciliation Loop                             │   │
│  │  - Status Management                               │   │
│  │  - Prometheus Metrics                              │   │
│  └────────────────────────────────────────────────────┘   │
│                             │                              │
│                             ▼                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │       Hardware Simulator / Actual Device           │   │
│  │       (RESTful API Interface)                      │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

For detailed architecture information, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Quick Start

### Prerequisites

- Go 1.21+
- Docker 20.10+
- Kubernetes 1.24+ (kubectl configured)
- Helm 3.0+ (for Crossplane installation)
- kind or Docker Desktop (for local testing)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/thc1006/crossplane-provider-oran.git
cd crossplane-provider-oran
```

2. **Install CRD**
```bash
kubectl apply -f adapter/config/crd/hardware.ran.example.com_opticaldevices.yaml
```

3. **Deploy Hardware Simulator** (for testing)
```bash
kubectl apply -f k8s/simulator-configmap.yaml
kubectl apply -f k8s/simulator-deployment.yaml
```

4. **Build and Deploy Adapter**
```bash
# Build Docker image
docker build -t oran-adapter:latest .

# Deploy to Kubernetes
kubectl apply -f k8s/adapter-deployment.yaml
```

5. **Verify Installation**
```bash
kubectl get pods
kubectl get crd opticaldevices.hardware.ran.example.com
```

### Create Your First Optical Device

```bash
kubectl apply -f k8s/test-optical-device.yaml
```

**Example OpticalDevice Resource:**
```yaml
apiVersion: hardware.ran.example.com/v1alpha1
kind: OpticalDevice
metadata:
  name: optical-device-001
  namespace: default
spec:
  location: site1
  controllerConfig:
    hostname: hardware-simulator
    port: 8080
  parameters:
    bandwidth: 100Gbps
    laserPower: 15dBm
    channel: 1
```

**Check Status:**
```bash
kubectl get opticaldevices
kubectl describe opticaldevice optical-device-001
```

## Usage

### Managing Optical Devices

**Create a device:**
```bash
kubectl apply -f - <<EOF
apiVersion: hardware.ran.example.com/v1alpha1
kind: OpticalDevice
metadata:
  name: my-device
spec:
  location: datacenter-1
  controllerConfig:
    hostname: 192.168.1.100
    port: 8080
  parameters:
    bandwidth: 200Gbps
    laserPower: 20dBm
    channel: 2
EOF
```

**Update device configuration:**
```bash
kubectl patch opticaldevice my-device --type=merge -p '{"spec":{"parameters":{"bandwidth":"400Gbps"}}}'
```

**View device status:**
```bash
kubectl get opticaldevice my-device -o yaml
```

**Delete device:**
```bash
kubectl delete opticaldevice my-device
```

### Monitoring with Prometheus

The adapter exposes Prometheus metrics on port 8080:

- `opticaldevice_info`: Information about managed optical devices
- `opticaldevice_reconcile_total`: Total number of reconciliation operations

**Example Prometheus query:**
```promql
opticaldevice_info{bandwidth="100Gbps"}
rate(opticaldevice_reconcile_total[5m])
```

A Grafana dashboard is available at `monitoring/grafana_dashboard.json`.

## Development

### Project Structure

```
crossplane-provider-oran/
├── adapter/                    # Adapter controller source code
│   ├── cmd/main.go            # Entry point
│   ├── config/crd/            # Generated CRD manifests
│   ├── internal/
│   │   ├── api/v1alpha1/      # API types
│   │   └── controller/        # Controller logic
│   └── hack/                  # Code generation scripts
├── hardware-simulator/         # Hardware simulator for testing
├── crossplane/                # Crossplane compositions
├── k8s/                       # Kubernetes manifests
├── monitoring/                # Monitoring configs
└── tests/                     # Test resources
```

### Building from Source

**Build adapter binary:**
```bash
cd adapter
go mod download
go build -o bin/manager ./cmd/main.go
```

**Build Docker image:**
```bash
docker build -t oran-adapter:dev .
```

**Run tests:**
```bash
cd adapter
go test ./internal/controller/... -v
```

### Code Generation

The project uses `controller-gen` for generating DeepCopy methods and CRD manifests.

**Install controller-gen:**
```bash
go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
```

**Generate DeepCopy methods:**
```bash
cd adapter
controller-gen object:headerFile=hack/boilerplate.go.txt paths="./internal/api/v1alpha1"
```

**Generate CRD manifests:**
```bash
cd adapter
controller-gen crd paths="./internal/api/v1alpha1" output:crd:artifacts:config=config/crd
```

### Local Development with Kind

**Create kind cluster:**
```bash
bash setup_kind_cluster.sh
```

**Load image to kind:**
```bash
kind load docker-image oran-adapter:dev
```

## Testing

### Unit Tests

```bash
cd adapter
go test ./internal/controller/... -cover
```

### Integration Tests

```bash
# Deploy all components
kubectl apply -f k8s/

# Run test scenarios
kubectl apply -f k8s/test-optical-device.yaml

# Verify reconciliation
kubectl get opticaldevices -w
```

### End-to-End Testing

See [TEST_REPORT.md](TEST_REPORT.md) for comprehensive testing results and scenarios.

## Crossplane Integration

### Install Crossplane

```bash
bash install_crossplane.sh
```

Or manually:
```bash
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update
helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane
```

### Deploy Crossplane Resources

```bash
# Install XRD (Composite Resource Definition)
kubectl apply -f crossplane/definition.yaml

# Install Composition
kubectl apply -f crossplane/composition.yaml

# Create a Claim
kubectl apply -f tests/claim.yaml
```

## API Reference

### OpticalDevice Resource

**Spec Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `location` | string | Physical location of the device |
| `controllerConfig.hostname` | string | Hostname of the device controller |
| `controllerConfig.port` | int | Port number of the device controller |
| `parameters.bandwidth` | string | Optical bandwidth (e.g., "100Gbps") |
| `parameters.laserPower` | string | Laser power level (e.g., "15dBm") |
| `parameters.channel` | int | Channel number |

**Status Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `conditions` | []Condition | Current state conditions |
| `observedBandwidth` | string | Actual bandwidth observed from hardware |
| `observedLaserPower` | string | Actual laser power observed from hardware |
| `lastSyncTime` | string | Last successful synchronization time |

## Troubleshooting

### Pod Not Starting

```bash
kubectl logs -l control-plane=controller-manager
kubectl describe pod -l control-plane=controller-manager
```

### CRD Issues

```bash
kubectl get crd opticaldevices.hardware.ran.example.com -o yaml
kubectl explain opticaldevice.spec
```

### Reconciliation Errors

```bash
kubectl describe opticaldevice <name>
kubectl logs -l control-plane=controller-manager --tail=50
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `chore:` Maintenance tasks
- `test:` Test additions or modifications

## Roadmap

- [ ] Add webhook validation (admission control)
- [ ] Implement Finalizer handling for resource cleanup
- [ ] Add support for multiple device types
- [ ] Enhanced error handling and retry logic
- [ ] Performance optimization for large-scale deployments
- [ ] Helm Chart packaging
- [ ] CI/CD pipeline integration

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Kubebuilder](https://book.kubebuilder.io/)
- Inspired by [Crossplane](https://crossplane.io/)
- Follows [O-RAN Alliance](https://www.o-ran.org/) specifications

## Contact

For questions and support, please open an issue on GitHub.

---

**Documentation:**
- [Architecture](ARCHITECTURE.md)
- [Project Structure](PROJECT_STRUCTURE.md)
- [Test Report](TEST_REPORT.md)
