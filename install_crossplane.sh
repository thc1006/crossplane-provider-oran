#!/bin/bash

# This script installs Crossplane into the Kubernetes cluster.
# It uses Helm to add the Crossplane repository, update it,
# and install the chart into the 'crossplane-system' namespace.

set -e

NAMESPACE="crossplane-system"

echo "Checking for Helm installation..."
if ! command -v helm &> /dev/null
then
    echo "Helm could not be found. Please install it first."
    echo "See: https://helm.sh/docs/intro/install/"
    exit 1
fi

echo "Checking for kubectl connection to a cluster..."
if ! kubectl cluster-info &> /dev/null
then
    echo "Could not connect to a Kubernetes cluster."
    echo "Please ensure your KinD cluster is running and kubectl is configured."
    exit 1
fi

echo "Adding Crossplane Helm repository..."
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

echo "Installing Crossplane using Helm..."
helm install crossplane --namespace "${NAMESPACE}" --create-namespace crossplane-stable/crossplane

echo "Waiting for Crossplane pods to be ready..."
kubectl wait --for=condition=Ready pod --all -n "${NAMESPACE}" --timeout=5m

echo "Crossplane installed successfully."
