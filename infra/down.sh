#!/usr/bin/env bash
# Tear the app cluster down. Deliberately does NOT touch the network layer or the
# Cassandra DB layer (./db.sh) - the DB keeps its data and stays reachable.
set -euo pipefail

cd "$(dirname "$0")"

REGION=$(terraform -chdir=network output -raw aws_region)
VPC_ID=$(terraform -chdir=network output -raw vpc_id)

echo "==> Destroying cluster layer"
terraform -chdir=cluster destroy -auto-approve

# The 3 Cassandra instances keep their ENIs, so wait until the only ENIs left in
# the VPC are the DB ones (tagged Role=cassandra).
DB_ENI_COUNT=$(aws ec2 describe-network-interfaces --region "$REGION" \
  --filters "Name=vpc-id,Values=${VPC_ID}" "Name=tag:Role,Values=cassandra" \
  --query "length(NetworkInterfaces)" --output text)

for i in $(seq 1 30); do
  TOTAL_ENI_COUNT=$(aws ec2 describe-network-interfaces --region "$REGION" \
    --filters "Name=vpc-id,Values=${VPC_ID}" \
    --query "length(NetworkInterfaces)" --output text)
  if [ "$TOTAL_ENI_COUNT" -eq "$DB_ENI_COUNT" ]; then
    echo "    VPC is free of EKS ENIs"
    break
  fi
  echo "    waiting for ENIs to release... (${i}/30)"
  sleep 10
done

if [ "$TOTAL_ENI_COUNT" -ne "$DB_ENI_COUNT" ]; then
  echo "    WARNING: ${TOTAL_ENI_COUNT} ENI(s) remain besides the ${DB_ENI_COUNT} Cassandra ENI(s)."
fi

echo
echo "Cluster destroyed. Cassandra DB is untouched (./db.sh status)."
echo "VPC and ECR left in place. To tear the VPC down: terraform -chdir=network destroy"