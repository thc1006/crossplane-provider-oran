#!/bin/bash

# This script sets up a Kubernetes cluster using KinD (Kubernetes in Docker).
# The cluster will have one control-plane node and two worker nodes.

set -e

CLUSTER_NAME="o-ran-dev"

echo "Checking for KinD installation..."
if ! command -v kind &> /dev/null
then
    echo "KinD could not be found. Please install it first."
    echo "See: https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    exit 1
fi

echo "Creating KinD cluster with name: ${CLUSTER_NAME}"

kind create cluster --name "${CLUSTER_NAME}" --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
EOF

echo "KinD cluster '${CLUSTER_NAME}' created successfully."
echo "You can now use kubectl to interact with your cluster."
echo "For example: kubectl cluster-info --context kind-${CLUSTER_NAME}"
