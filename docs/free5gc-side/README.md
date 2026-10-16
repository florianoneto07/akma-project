# Lado free5GC

Tudo o que é feito **dentro** de uma instalação do free5GC para habilitar AKMA: o AAnF, os *patches* de AUSF, UDM e `free5gc/openapi`, configuração, registro no NRF e dados de assinatura.

Cada documento desta pasta registra:

1. qual arquivo do free5GC foi tocado (caminho completo na instalação);
2. o que mudou e por quê, com a cláusula 3GPP;
3. como reproduzir numa instalação limpa;
4. como verificar que funcionou.

Base de código: free5GC v4.2.3 (`docs/baseline/ambiente-verificado.md`).

| Documento | MVP | Estado |
|---|---|---|
| `aanf.md`, `openapi.md`, `registro-nrf.md` | MVP-1 | a escrever |
| `udm.md`, `ausf.md`, `assinatura-akma.md`, `verificacao.md` | MVP-2 | a escrever |

O que é comum aos dois lados — KDF, formatos de identificador, codificação, vetores — fica em `docs/shared/`.
