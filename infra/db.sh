#!/usr/bin/env bash
# creates the long-lived Cassandra DB layer (3 EC2 instances).
# ONLY RUN THIS IN WSL2



set -euo pipefail

cd "$(dirname "$0")"

REGION=$(terraform -chdir=network output -raw aws_region 2>/dev/null || echo "ca-central-1")
SECRETS_DIR="k8s/secrets"

gen_secrets() {
  mkdir -p "$SECRETS_DIR"
  [ -f "$SECRETS_DIR/cassandra-app-username.txt" ] || printf 'mychatapp' > "$SECRETS_DIR/cassandra-app-username.txt"
  [ -f "$SECRETS_DIR/cassandra-app-password.txt" ] || openssl rand -hex 24 | tr -d '\n' > "$SECRETS_DIR/cassandra-app-password.txt"
  [ -f "$SECRETS_DIR/cassandra-superuser.txt" ]    || openssl rand -hex 24 | tr -d '\n' > "$SECRETS_DIR/cassandra-superuser.txt"
}

run_ansible() {
  local aws_dir="${USERPROFILE:-/mnt/c/Users/$USER}/.aws"
  mkdir -p ~/.ssh
  cp db/mychatapp-key.pem ~/.ssh/mychatapp.pem
  chmod 600 ~/.ssh/mychatapp.pem
  ( cd ansible &&
    ANSIBLE_CONFIG="$(pwd)/ansible.cfg"
    AWS_SHARED_CREDENTIALS_FILE="$aws_dir/credentials" \
    AWS_CONFIG_FILE="$aws_dir/config" \
    ansible-playbook playbook.yaml )
}

case "${1:-apply}" in
  apply|up)
    gen_secrets
    terraform -chdir=db init -input=false
    terraform -chdir=db apply -auto-approve
    run_ansible
    echo "Cassandra DB is up. Verify with: ./db.sh status"
    ;;
  reconfigure)
    # re-run ansible only (config / schema changes)
    run_ansible
    ;;
  stop)
    aws ec2 stop-instances --region "$REGION" \
      --filters "Name=tag:Role,Values=cassandra" "Name=instance-state-name,Values=running"
    ;;
  start)
    aws ec2 start-instances --region "$REGION" \
      --filters "Name=tag:Role,Values=cassandra" "Name=instance-state-name,Values=stopped"
    ;;
  status)
    terraform -chdir=db output
    aws ec2 describe-instances --region "$REGION" \
      --filters "Name=tag:Role,Values=cassandra" \
      --query "Reservations[].Instances[].[Tags[?Key=='Name'].Value|[0],InstanceId,State.Name,PrivateIpAddress]" \
      --output table
    ;;
  destroy)
    read -r -p "Terminate the DB instances and DELETE all Cassandra data? Type 'yes': " ok
    if [ "$ok" = "yes" ]; then
      terraform -chdir=db destroy -auto-approve
    else
      echo "aborted"
    fi
    ;;
  *)
    echo "usage: $0 {apply|reconfigure|stop|start|status|destroy}"
    exit 1
    ;;
esac
