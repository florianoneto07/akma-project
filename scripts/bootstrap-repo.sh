#!/usr/bin/env bash
# Inicializa o repositório local e publica no GitHub.
#
# Uso:
#   ./scripts/bootstrap-repo.sh <usuario>/<repositorio>
#
# Exemplo:
#   ./scripts/bootstrap-repo.sh florianoneto07/akma-free5gc-docs
#
# Requer: git. O passo de publicação usa a GitHub CLI (gh) se disponível;
# caso contrário, imprime as instruções manuais.

set -euo pipefail

REPO="${1:-}"
if [[ -z "$REPO" ]]; then
  echo "uso: $0 <usuario>/<repositorio>" >&2
  exit 1
fi

cd "$(dirname "$0")/.."

if [[ ! -d .git ]]; then
  git init -b main
fi

git add .
git commit -m "docs: baseline normativo do AKMA Release 19 e decisões de projeto

Consolida as fases F0 e F1:
- baseline de 16 especificações 3GPP Release 19
- análise normativa completa e quatro adendos
- ADR-0001 (topologia), ADR-0002 (codificação), ADR-0003 (identificador Ua*)
- registro de cinco defeitos em especificações
- matriz de rastreabilidade, revisão 3
- implementação de referência da KDF"

if command -v gh >/dev/null 2>&1; then
  gh repo create "$REPO" --public --source=. --remote=origin --push
  echo "publicado em https://github.com/$REPO"
else
  cat <<MSG

git commit criado. A GitHub CLI (gh) não foi encontrada.

Para publicar manualmente:

  1. Crie o repositório vazio em https://github.com/new  (nome: ${REPO#*/})
  2. Execute:

     git remote add origin git@github.com:$REPO.git
     git push -u origin main

MSG
fi
