# UE (UERANSIM): K_AKMA e A-KID no ME

**MVP:** MVP-3a (derivação e log) · **Data:** 18 de outubro de 2026

O MVP-4 completa este documento com a K_AF, o comando do `nr-cli` e o cliente Ua\*.

## 1. O que o UE precisa fazer

`[3GPP]` TS 33.535 §6.1: *"The UE shall generate the AKMA Anchor Key (KAKMA) and the A-KID from the KAUSF before initiating communication with an AKMA Application Function."* O §6.2.1, passo 1, repete a exigência e acrescenta que o UE envia o A-KID ao AF.

| Valor | Derivação | Referência |
|---|---|---|
| K_AKMA | KDF(K_AUSF), FC `0x80`, P0 `"AKMA"`, P1 SUPI | TS 33.535 A.2 |
| A-TID | KDF(K_AUSF), FC `0x81`, P0 `"A-TID"`, P1 SUPI | TS 33.535 A.3 |
| A-KID | `rid<RID>.atid<A-TID>@5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org` | TS 33.535 §6.1; TS 29.522 Tab. 5.14.5.4.2-1; ADR-0002 |

- **SUPI:** o IMSI, só os dígitos (TS 33.501 A.7.0).
- **RID:** o do SUCI do UE (TS 23.003 §2.2B).
- **Domínio:** o da rede de origem, com o MNC em três dígitos (TS 23.003 §28.2).
- **A-TID:** em base64 com preenchimento, 44 caracteres (RFC 4648 §4), conforme o ADR-0002.

`[FREE5GC]` O UERANSIM (a versão instalada na VM do free5GC, *commit* `85a0fbfd`) calcula a K_AUSF nos dois métodos de autenticação primária, mas não tem nada de AKMA.

## 2. Quando o UE deriva

`[INTERPRETAÇÃO]` O UE deriva K_AKMA e A-KID quando passa a usar o contexto de segurança NAS criado pela autenticação primária bem-sucedida. É o momento em que o `nr-ue` aceita o *Security Mode Command* que ativa esse contexto:

- **5G AKA:** o UE não recebe confirmação explícita de sucesso. O *Security Mode Command* que ativa o contexto novo é a primeira prova de que a rede aceitou o RES\*.
- **EAP-AKA':** o EAP-Success chega no mesmo *Security Mode Command*.

Isso atende ao "antes de comunicar com um AF" do §6.1: nenhuma aplicação fala com um AF antes de o UE se registrar. Cada autenticação primária nova gera K_AKMA e A-KID novos, como acontece no AUSF.

O material fica no contexto de segurança NAS (`UeKeys`) e é descartado com ele.

## 3. Configuração

`[MVP]` Opção `akma` na configuração do UE, prevista no plano (MVP-4). O padrão é `false`, que preserva o comportamento original do UERANSIM. O *patch* `0002` liga a opção no `free5gc-ue.yaml`:

```yaml
# AKMA (3GPP TS 33.535): derive K_AKMA and the A-KID after each successful primary authentication
akma: true
```

Os demais dados vêm do que o UE já tem:

- `supi`: só IMSI, como na rede (DV-3); com outro tipo de SUPI, o UE registra um aviso e não deriva;
- `routingIndicator`: o mesmo que vai no SUCI. `[FREE5GC]` Sem o campo, o UERANSIM põe `"0000"` no SUCI, e o A-KID usa o mesmo valor;
- `mcc` e `mnc`: o comprimento do MNC vem da configuração, e não do IMSI.

## 4. Log

O UERANSIM imprime cada linha como `[data hora] [componente] [nível] mensagem`, com nível a partir de `debug`. A derivação acontece na camada NAS e usa o componente `nas`:

```
[2026-10-18 10:15:42.118] [nas] [info] AKMA key material derived, A-KID[rid0000.atid…=@5gc.mnc093.mcc208.3gppnetwork.org]
```

| Nível | Mensagem | Quando |
|---|---|---|
| `info` | `AKMA key material derived, A-KID[<A-KID>]` | contexto novo ativado e `akma: true` |
| `warn` | `AKMA key material not derived, the SUPI is not an IMSI` | `akma: true` com SUPI de outro tipo |

O formato `texto[valor]` é o mesmo das mensagens do UERANSIM, como `Selected integrity[2] ciphering[0]`. K_AUSF e K_AKMA nunca aparecem. O A-KID impresso deve ser idêntico ao que o AUSF e o AAnF registram (`docs/free5gc-side/logs-akma.md`).

## 5. Arquivos tocados

Árvore: `~/UERANSIM`, na VM do free5GC (D5). Binário: `build/nr-ue`, com o original em `build/nr-ue.orig`. *Patches*: `af/ue/ueransim/patches/`.

| Arquivo | Mudança | *Patch* |
|---|---|---|
| `src/ue/types.hpp` | `UeConfig::akma`; `UeKeys::kAkma` e `UeKeys::aKid`, copiados em `deepCopy` | `0001-akma-key-derivation.patch` |
| `src/ue.cpp` | lê `akma`; copia a opção para cada UE (`-n`) | idem |
| `src/ue/nas/keys.{hpp,cpp}` | `CalculateKAkma`, `CalculateATid`, `ConstructAKid`, `DeriveAkmaKeys`; codificação base64 | idem |
| `src/ue/nas/mm/mm.hpp`, `security.cpp` | `NasMm::deriveAkmaKeys`, chamada quando o contexto novo vira o corrente | idem |
| `config/free5gc-ue.yaml` | `akma: true` | `0002-free5gc-ue-akma.patch` |

A KDF é a do próprio UERANSIM (`crypto::CalculateKdfKey`), a mesma que ele usa para K_SEAF e K_AMF. O UERANSIM não tem base64, então o *patch* traz um codificador pequeno.

## 6. Como reproduzir numa instalação limpa

Na máquina do UERANSIM, com o repositório em `~/akma-free5gc`:

```bash
~/akma-free5gc/af/scripts/patch-ueransim.sh apply
```

O *script* faz quatro coisas:

- confere e aplica os *patches* em ordem;
- guarda `build/nr-ue.orig`;
- recompila só o `nr-ue` em `cmake-build-release`, aproveitando a compilação já feita (sem ela, roda o `make` completo);
- copia o binário novo para `build/`.

Alterações locais em arquivos que os *patches* não tocam, como o `free5gc-gnb.yaml`, não bloqueiam a aplicação. Para desfazer:

```bash
~/akma-free5gc/af/scripts/patch-ueransim.sh revert
```

## 7. Como verificar

Registre o UE com um assinante AKMA (`docs/free5gc-side/assinatura-akma.md`). A linha `AKMA key material derived` deve aparecer na tela do UE com o mesmo A-KID que o AAnF registrou. O procedimento e o resultado estão em `docs/free5gc-side/verificacao.md`, seção 4.

## 8. Limites

- **K_AKMA ainda não conferida.** A igualdade dos A-KID confere a KDF do UERANSIM, a A-TID (FC `0x81`), o base64 e o formato. A K_AKMA (FC `0x80`) passa pela mesma KDF, mas só é conferida no MVP-4: pelos vetores do MVP-0 e pela K_AF que autentica o Ua\*.
- **Sem NFKC.** O C++ do UERANSIM não aplica NFKC; isso só importa para a K_AF (FQDN do AF, MVP-4).
