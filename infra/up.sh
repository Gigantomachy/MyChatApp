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

# aws eks update-kubeconfig --region ca-central-1 --name mychatapp

echo "==> Waiting for the nodes to become Ready"
kubectl wait --for=condition=Ready node --all --timeout=5m
kubectl get nodes -o wide

echo "==> Applying k8s manifests"
kubectl apply -k k8s/

echo "==> Waiting for deployments"
kubectl rollout status deployment/backend --timeout=3m
kubectl rollout status deployment/frontend --timeout=3m

ALB_DNS=$(terraform -chdir=cluster output -raw alb_dns_name)
echo
echo "App is up: http://${ALB_DNS}"
echo "Remember to run ./down.sh when you are finished."
