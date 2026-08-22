#!/usr/bin/env bash
# Bring the cluster up for a session, the network layer is applied too but is a no-op once it already exists.
set -euo pipefail

cd "$(dirname "$0")"

echo "==> Network layer (VPC)"
terraform -chdir=network init -input=false
terraform -chdir=network apply -auto-approve

echo "==> Cluster layer (EKS + node group). This takes 10-15 minutes."
terraform -chdir=cluster init -input=false
terraform -chdir=cluster apply -auto-approve

echo "==> Configuring kubectl"
eval "$(terraform -chdir=cluster output -raw configure_kubectl)"

echo "==> Waiting for the nodes to become Ready"
kubectl wait --for=condition=Ready node --all --timeout=5m
kubectl get nodes -o wide

helm repo add jetstack https://charts.jetstack.io --force-update >/dev/null
helm repo add k8ssandra https://helm.k8ssandra.io/stable --force-update >/dev/null

helm upgrade --install cert-manager jetstack/cert-manager -n cert-manager --create-namespace --set crds.enabled=true
# cert-manager's webhook must be serving before k8ssandra-operator installs (it issues the operator's webhook certs).
kubectl rollout status deployment/cert-manager-webhook -n cert-manager --timeout=3m

helm upgrade --install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace --set global.clusterScoped=true
kubectl rollout status deployment/k8ssandra-operator -n k8ssandra-operator --timeout=3m
kubectl wait --for=condition=Established crd/k8ssandraclusters.k8ssandra.io --timeout=3m

echo "==> Applying k8s manifests"
kubectl apply -k k8s/

echo "==> Waiting for deployments"
kubectl rollout status deployment/redis --timeout=3m
kubectl rollout status deployment/backend --timeout=3m
kubectl rollout status deployment/frontend --timeout=3m

ALB_DNS=$(terraform -chdir=cluster output -raw alb_dns_name)
echo
echo "App is up: http://${ALB_DNS}"
echo "Remember to run ./down.sh when you are finished."