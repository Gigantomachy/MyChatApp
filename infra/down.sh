#!/usr/bin/env bash

# Tear the cluster down.
# Deliberately does NOT touch the network layer because the VPC costs nothing while idle (no NAT gateway).
set -euo pipefail

cd "$(dirname "$0")"

# Any AWS load balancer created by kubernetes is invisible to Terraform.
# If one is left behind it holds ENIs in the subnets and blocks VPC deletion
# later, with a very unhelpful error. Harmless to run when there are none.
if kubectl cluster-info >/dev/null 2>&1; then
  echo "==> Deleting Cassandra StatefulSet + PVCs (frees the EBS volumes)"
  kubectl delete statefulset cassandra --ignore-not-found --timeout=2m || true
  kubectl delete pvc --all --ignore-not-found --timeout=2m || true

  echo "==> Removing Kubernetes-managed load balancers"
  kubectl delete ingress --all --all-namespaces --ignore-not-found --timeout=2m || true
  kubectl delete svc --all-namespaces --field-selector spec.type=LoadBalancer \
    --ignore-not-found --timeout=2m || true
  echo "    waiting 60s for AWS to actually release them"
  sleep 60
else
  echo "==> kubectl cannot reach a cluster; skipping load balancer cleanup"
fi

echo "==> Destroying cluster layer"
terraform -chdir=cluster destroy -auto-approve

echo "==> Waiting for ENIs to be released from the VPC"
VPC_ID=$(terraform -chdir=network output -raw vpc_id)
for i in $(seq 1 30); do
  ENI_COUNT=$(aws ec2 describe-network-interfaces \
    --filters "Name=vpc-id,Values=${VPC_ID}" \
    --query "length(NetworkInterfaces)" \
    --output text)
  if [ "$ENI_COUNT" -eq 0 ]; then
    echo "    VPC is free of ENIs"
    break
  fi
  echo "    waiting for ${ENI_COUNT} ENI(s) to release... (${i}/30)"
  sleep 10
done
if [ "$ENI_COUNT" -ne 0 ]; then
  echo "    WARNING: ${ENI_COUNT} ENI(s) still in the VPC. A network destroy may hang."
fi

echo
echo "Cluster destroyed. VPC and ECR left in place."
echo "To tear those down as well: terraform -chdir=network destroy"
