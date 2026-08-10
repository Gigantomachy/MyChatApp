# MyChatApp — dev infrastructure

Two Terraform layers with separate state, so the expensive half can be destroyed between sessions while the cheap half persists.

```
infra/
  network/    VPC (public subnets, no NAT) — long-lived, ~$0/hr idle
  cluster/    EKS control plane + 1 managed node + addons — ~$0.20/hr, destroy when idle
  up.sh       apply both layers, configure kubectl
  down.sh     destroy the cluster layer only
```

## First run

```bash
chmod +x up.sh down.sh

terraform -chdir=network init
terraform -chdir=network apply

terraform -chdir=cluster init
terraform -chdir=cluster apply        # 10–15 minutes

eval "$(terraform -chdir=cluster output -raw configure_kubectl)"
kubectl get nodes
```

## Daily rhythm

```bash
./up.sh      # ~15 min
# ... work ...
./down.sh    # ~10 min
```

## Prerequisites

- Terraform >= 1.10
- AWS CLI v2, authenticated (`aws sts get-caller-identity` should succeed)
- kubectl
