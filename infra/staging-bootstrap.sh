#!/usr/bin/env bash
# Bootstrap скрипт для первого запуска staging на Oracle Cloud (Ubuntu 22.04 ARM64)
# Запуск: ssh ubuntu@<IP> 'bash -s' < infra/staging-bootstrap.sh
set -euo pipefail

DOMAIN="${DOMAIN:-staging.pkt.kz}"
GITHUB_OWNER="${GITHUB_OWNER:-KAZmake}"

echo "=== [1/6] Обновление системы ==="
sudo apt-get update -qq && sudo apt-get upgrade -y -qq

echo "=== [2/6] Установка k3s ==="
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="server --disable=traefik" sh -
# Traefik отключаем и устанавливаем вручную — чтобы контролировать версию

# Ждём k3s
until kubectl get nodes 2>/dev/null | grep -q Ready; do sleep 2; done
echo "k3s готов"

# Kubeconfig для пользователя
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown "$USER" ~/.kube/config
export KUBECONFIG=~/.kube/config

echo "=== [3/6] Установка Traefik через Helm ==="
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik \
  --namespace kube-system \
  --set ports.web.redirectTo.port=websecure \
  --set ports.websecure.tls.enabled=true \
  --wait

echo "=== [4/6] Создание namespace staging ==="
kubectl apply -f - <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: staging
EOF

echo "=== [5/6] TLS Secret (Cloudflare Origin Certificate) ==="
echo "Загрузи сертификат вручную:"
echo "  kubectl create secret tls pkt-staging-tls -n staging \\"
echo "    --cert=origin.crt --key=origin.key"

echo "=== [6/6] Keycloak realm config ==="
kubectl create configmap keycloak-realm-config -n staging \
  --from-file=realm-pkt.json=infra/keycloak/realm-pkt.json \
  --dry-run=client -o yaml | kubectl apply -f -

echo ""
echo "=== Bootstrap завершён ==="
echo "Следующие шаги:"
echo "  1. Скопируй kubeconfig:"
echo "     cat ~/.kube/config | base64 → GitHub Secret KUBECONFIG_STAGING"
echo "  2. Создай secrets (заполни реальными значениями):"
echo "     kubectl apply -f k8s/staging/secrets.yaml  # (не в git!)"
echo "  3. Загрузи TLS сертификат (шаг 5 выше)"
echo "  4. Запусти первый деплой: git push origin develop"
echo "     или вручную: kubectl apply -k k8s/staging/"
echo ""
echo "Kubeconfig для GitHub Secrets:"
cat ~/.kube/config | base64 -w0
echo ""
