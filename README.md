# kubectl-aliases (Go Edition)

This repository provides an experimental utility for programmatically
generating concise and conflict-free shell aliases for `kubectl` commands
using Go.

It was designed to support rapid Kubernetes workflows, reduce keystrokes, and
help operators during debugging and triage. The logic is intentionally
minimal and designed to be easy to audit or extend.

## Overview

This utility uses heuristics to prioritize common CRUD operations like `get`,
`apply`, `describe`, `delete`, and `logs`, along with modifier support such as
`-o yaml`, `--all-namespaces`, or `-f FILE`.

If your cluster defines Custom Resource Definitions (CRDs) with short names,
those will be detected automatically. Aliases will be generated for those
resources as well — especially useful when managing platform extensions or
custom controllers.

The generator outputs a flat alias file with shell functions and basic alias
definitions. It is designed to work well in POSIX-compatible environments.

## What It Produces

Running the generator creates a file called `kubectl_aliases`. You can source
it directly in your terminal session or add it to your shell config.

```sh
source ./kubectl_aliases
```

Or persist it:

```sh
echo "source ~/path/to/kubectl_aliases" >> ~/.bashrc
```

The aliases will include commands like:

```sh
alias kg='kubectl get'                   # get
alias kgj='kubectl get -o json'          # get + output
alias kgjn='kubectl get -o json --all-ns'# get + output + scope
alias ke='kubectl edit'                  # edit
alias kd='kubectl describe'              # describe
alias kdel='kubectl delete'              # delete
function ka() { kubectl apply  "$1"; } # apply + input modifier
function kdelf() { kubectl delete -f "$1"; } # delete + input modifier
```

## Modifiers

| Modifier Type | Suffix | Flag / Effect            |
| ------------- | ------ | ------------------------ |
| Output JSON   | `j`    | `-o json`                |
| Output YAML   | `y`    | `-o yaml`                |
| Output Wide   | `l`    | `-o wide`                |
| Scope (All NS)| `n`    | `--all-namespaces`       |
| Input File    | `f`    | `-f` for apply/delete     |

## Testing and Validation

Basic tests are included and run automatically using:

```sh
make all
```

A GitHub Actions workflow will also run `make all` and ensure that changes to
`kubectl_aliases` are properly reflected in the README and file outputs. A
Kind cluster will be started automatically if needed to verify dynamic CRD
detection.

## Extending or Modifying

To change command priorities or modifier support, edit the `CommandConfig`
map in `generate_aliases.go`. You can also add your own logic to support
non-standard workflows or override how short names are chosen.

This is a first-pass generator. Contributions, bug reports, and experiments
are welcome but may not be prioritized as this is a low maintenance side project.
Also see [CONTRIBUTING.md](./CONTRIBUTING.md).

## All Available Aliases

```
apply => ka
auth => kau
create => kc
debug => kde
delete => kdel
describe => kd
edit => ke
events => kev
get => kg
logs => kl
top => kt
```

## Full Alias Reference

```bash
# Auto-generated kubectl function alias file
alias ka='kubectl apply' # apply
alias kau='kubectl auth' # auth
alias kc='kubectl create' # create
alias kccj='kubectl create cj' # create cronjobs
alias kccm='kubectl create cm' # create configmaps
alias kccrd='kubectl create crd' # create customresourcedefinitions
alias kccs='kubectl create cs' # create componentstatuses
alias kccsr='kubectl create csr' # create certificatesigningrequests
alias kcdeploy='kubectl create deploy' # create deployments
alias kcds='kubectl create ds' # create daemonsets
alias kcep='kubectl create ep' # create endpoints
alias kcev='kubectl create ev' # create events
alias kchpa='kubectl create hpa' # create horizontalpodautoscalers
alias kcing='kubectl create ing' # create ingresses
alias kclimits='kubectl create limits' # create limitranges
alias kcnetpol='kubectl create netpol' # create networkpolicies
alias kcno='kubectl create no' # create nodes
alias kcns='kubectl create ns' # create namespaces
alias kcpc='kubectl create pc' # create priorityclasses
alias kcpdb='kubectl create pdb' # create poddisruptionbudgets
alias kcpo='kubectl create po' # create pods
alias kcpv='kubectl create pv' # create persistentvolumes
alias kcpvc='kubectl create pvc' # create persistentvolumeclaims
alias kcquota='kubectl create quota' # create resourcequotas
alias kcrc='kubectl create rc' # create replicationcontrollers
alias kcrs='kubectl create rs' # create replicasets
alias kcsa='kubectl create sa' # create serviceaccounts
alias kcsc='kubectl create sc' # create storageclasses
alias kcsts='kubectl create sts' # create statefulsets
alias kcsvc='kubectl create svc' # create services
alias kd='kubectl describe' # describe
alias kdcj='kubectl describe cj' # describe cronjobs
alias kdcm='kubectl describe cm' # describe configmaps
alias kdcrd='kubectl describe crd' # describe customresourcedefinitions
alias kdcs='kubectl describe cs' # describe componentstatuses
alias kdcsr='kubectl describe csr' # describe certificatesigningrequests
alias kddeploy='kubectl describe deploy' # describe deployments
alias kdds='kubectl describe ds' # describe daemonsets
alias kde='kubectl debug' # debug
alias kdel='kubectl delete' # delete
alias kdelcj='kubectl delete cj' # delete cronjobs
alias kdelcm='kubectl delete cm' # delete configmaps
alias kdelcrd='kubectl delete crd' # delete customresourcedefinitions
alias kdelcs='kubectl delete cs' # delete componentstatuses
alias kdelcsr='kubectl delete csr' # delete certificatesigningrequests
alias kdeldeploy='kubectl delete deploy' # delete deployments
alias kdelds='kubectl delete ds' # delete daemonsets
alias kdelep='kubectl delete ep' # delete endpoints
alias kdelev='kubectl delete ev' # delete events
alias kdelhpa='kubectl delete hpa' # delete horizontalpodautoscalers
alias kdeling='kubectl delete ing' # delete ingresses
alias kdellimits='kubectl delete limits' # delete limitranges
alias kdelnetpol='kubectl delete netpol' # delete networkpolicies
alias kdelno='kubectl delete no' # delete nodes
alias kdelns='kubectl delete ns' # delete namespaces
alias kdelpc='kubectl delete pc' # delete priorityclasses
alias kdelpdb='kubectl delete pdb' # delete poddisruptionbudgets
alias kdelpo='kubectl delete po' # delete pods
alias kdelpv='kubectl delete pv' # delete persistentvolumes
alias kdelpvc='kubectl delete pvc' # delete persistentvolumeclaims
alias kdelquota='kubectl delete quota' # delete resourcequotas
alias kdelrc='kubectl delete rc' # delete replicationcontrollers
alias kdelrs='kubectl delete rs' # delete replicasets
alias kdelsa='kubectl delete sa' # delete serviceaccounts
alias kdelsc='kubectl delete sc' # delete storageclasses
alias kdelsts='kubectl delete sts' # delete statefulsets
alias kdelsvc='kubectl delete svc' # delete services
alias kdep='kubectl describe ep' # describe endpoints
alias kdev='kubectl describe ev' # describe events
alias kdhpa='kubectl describe hpa' # describe horizontalpodautoscalers
alias kding='kubectl describe ing' # describe ingresses
alias kdlimits='kubectl describe limits' # describe limitranges
alias kdnetpol='kubectl describe netpol' # describe networkpolicies
alias kdno='kubectl describe no' # describe nodes
alias kdns='kubectl describe ns' # describe namespaces
alias kdpc='kubectl describe pc' # describe priorityclasses
alias kdpdb='kubectl describe pdb' # describe poddisruptionbudgets
alias kdpo='kubectl describe po' # describe pods
alias kdpv='kubectl describe pv' # describe persistentvolumes
alias kdpvc='kubectl describe pvc' # describe persistentvolumeclaims
alias kdquota='kubectl describe quota' # describe resourcequotas
alias kdrc='kubectl describe rc' # describe replicationcontrollers
alias kdrs='kubectl describe rs' # describe replicasets
alias kdsa='kubectl describe sa' # describe serviceaccounts
alias kdsc='kubectl describe sc' # describe storageclasses
alias kdsts='kubectl describe sts' # describe statefulsets
alias kdsvc='kubectl describe svc' # describe services
alias ke='kubectl edit' # edit
alias kecj='kubectl edit cj' # edit cronjobs
alias kecm='kubectl edit cm' # edit configmaps
alias kecrd='kubectl edit crd' # edit customresourcedefinitions
alias kecs='kubectl edit cs' # edit componentstatuses
alias kecsr='kubectl edit csr' # edit certificatesigningrequests
alias kedeploy='kubectl edit deploy' # edit deployments
alias keds='kubectl edit ds' # edit daemonsets
alias keep='kubectl edit ep' # edit endpoints
alias keev='kubectl edit ev' # edit events
alias kehpa='kubectl edit hpa' # edit horizontalpodautoscalers
alias keing='kubectl edit ing' # edit ingresses
alias kelimits='kubectl edit limits' # edit limitranges
alias kenetpol='kubectl edit netpol' # edit networkpolicies
alias keno='kubectl edit no' # edit nodes
alias kens='kubectl edit ns' # edit namespaces
alias kepc='kubectl edit pc' # edit priorityclasses
alias kepdb='kubectl edit pdb' # edit poddisruptionbudgets
alias kepo='kubectl edit po' # edit pods
alias kepv='kubectl edit pv' # edit persistentvolumes
alias kepvc='kubectl edit pvc' # edit persistentvolumeclaims
alias kequota='kubectl edit quota' # edit resourcequotas
alias kerc='kubectl edit rc' # edit replicationcontrollers
alias kers='kubectl edit rs' # edit replicasets
alias kesa='kubectl edit sa' # edit serviceaccounts
alias kesc='kubectl edit sc' # edit storageclasses
alias kests='kubectl edit sts' # edit statefulsets
alias kesvc='kubectl edit svc' # edit services
alias kev='kubectl events' # events
alias kg='kubectl get' # get
alias kgcj='kubectl get cj' # get cronjobs
alias kgcm='kubectl get cm' # get configmaps
alias kgcrd='kubectl get crd' # get customresourcedefinitions
alias kgcs='kubectl get cs' # get componentstatuses
alias kgcsr='kubectl get csr' # get certificatesigningrequests
alias kgdeploy='kubectl get deploy' # get deployments
alias kgds='kubectl get ds' # get daemonsets
alias kgep='kubectl get ep' # get endpoints
alias kgev='kubectl get ev' # get events
alias kghpa='kubectl get hpa' # get horizontalpodautoscalers
alias kging='kubectl get ing' # get ingresses
alias kglimits='kubectl get limits' # get limitranges
alias kgnetpol='kubectl get netpol' # get networkpolicies
alias kgno='kubectl get no' # get nodes
alias kgns='kubectl get ns' # get namespaces
alias kgpc='kubectl get pc' # get priorityclasses
alias kgpdb='kubectl get pdb' # get poddisruptionbudgets
alias kgpo='kubectl get po' # get pods
alias kgpv='kubectl get pv' # get persistentvolumes
alias kgpvc='kubectl get pvc' # get persistentvolumeclaims
alias kgquota='kubectl get quota' # get resourcequotas
alias kgrc='kubectl get rc' # get replicationcontrollers
alias kgrs='kubectl get rs' # get replicasets
alias kgsa='kubectl get sa' # get serviceaccounts
alias kgsc='kubectl get sc' # get storageclasses
alias kgsts='kubectl get sts' # get statefulsets
alias kgsvc='kubectl get svc' # get services
alias kl='kubectl logs' # logs
alias kt='kubectl top' # top
function ka() { kubectl apply  "$1"; } # apply + input modifier
function kaf() { kubectl apply -f "$1"; } # apply + input modifier
function kaun() { kubectl auth --all-namespaces "$1"; } # auth + scope modifier
function kdel() { kubectl delete  "$1"; } # delete + input modifier
function kdelf() { kubectl delete -f "$1"; } # delete + input modifier
function kden() { kubectl debug --all-namespaces "$1"; } # debug + scope modifier
function kdn() { kubectl describe --all-namespaces "$1"; } # describe + scope modifier
function kevj() { kubectl events -o json "$1"; } # events + output modifier
function kevl() { kubectl events -o wide "$1"; } # events + output modifier
function kevn() { kubectl events --all-namespaces "$1"; } # events + scope modifier
function kevy() { kubectl events -o yaml "$1"; } # events + output modifier
function kgcjj() { kubectl get cj "$1" -o json; } # get cronjobs + output
function kgcjl() { kubectl get cj "$1" -o wide; } # get cronjobs + output
function kgcjy() { kubectl get cj "$1" -o yaml; } # get cronjobs + output
function kgcmj() { kubectl get cm "$1" -o json; } # get configmaps + output
function kgcml() { kubectl get cm "$1" -o wide; } # get configmaps + output
function kgcmy() { kubectl get cm "$1" -o yaml; } # get configmaps + output
function kgcrdj() { kubectl get crd "$1" -o json; } # get customresourcedefinitions + output
function kgcrdl() { kubectl get crd "$1" -o wide; } # get customresourcedefinitions + output
function kgcrdy() { kubectl get crd "$1" -o yaml; } # get customresourcedefinitions + output
function kgcsj() { kubectl get cs "$1" -o json; } # get componentstatuses + output
function kgcsl() { kubectl get cs "$1" -o wide; } # get componentstatuses + output
function kgcsrj() { kubectl get csr "$1" -o json; } # get certificatesigningrequests + output
function kgcsrl() { kubectl get csr "$1" -o wide; } # get certificatesigningrequests + output
function kgcsry() { kubectl get csr "$1" -o yaml; } # get certificatesigningrequests + output
function kgcsy() { kubectl get cs "$1" -o yaml; } # get componentstatuses + output
function kgdeployj() { kubectl get deploy "$1" -o json; } # get deployments + output
function kgdeployl() { kubectl get deploy "$1" -o wide; } # get deployments + output
function kgdeployy() { kubectl get deploy "$1" -o yaml; } # get deployments + output
function kgdsj() { kubectl get ds "$1" -o json; } # get daemonsets + output
function kgdsl() { kubectl get ds "$1" -o wide; } # get daemonsets + output
function kgdsy() { kubectl get ds "$1" -o yaml; } # get daemonsets + output
function kgepj() { kubectl get ep "$1" -o json; } # get endpoints + output
function kgepl() { kubectl get ep "$1" -o wide; } # get endpoints + output
function kgepy() { kubectl get ep "$1" -o yaml; } # get endpoints + output
function kgevj() { kubectl get ev "$1" -o json; } # get events + output
function kgevl() { kubectl get ev "$1" -o wide; } # get events + output
function kgevy() { kubectl get ev "$1" -o yaml; } # get events + output
function kghpaj() { kubectl get hpa "$1" -o json; } # get horizontalpodautoscalers + output
function kghpal() { kubectl get hpa "$1" -o wide; } # get horizontalpodautoscalers + output
function kghpay() { kubectl get hpa "$1" -o yaml; } # get horizontalpodautoscalers + output
function kgingj() { kubectl get ing "$1" -o json; } # get ingresses + output
function kgingl() { kubectl get ing "$1" -o wide; } # get ingresses + output
function kgingy() { kubectl get ing "$1" -o yaml; } # get ingresses + output
function kgj() { kubectl get -o json "$1"; } # get + output modifier
function kgl() { kubectl get -o wide "$1"; } # get + output modifier
function kglimitsj() { kubectl get limits "$1" -o json; } # get limitranges + output
function kglimitsl() { kubectl get limits "$1" -o wide; } # get limitranges + output
function kglimitsy() { kubectl get limits "$1" -o yaml; } # get limitranges + output
function kgn() { kubectl get --all-namespaces "$1"; } # get + scope modifier
function kgnetpolj() { kubectl get netpol "$1" -o json; } # get networkpolicies + output
function kgnetpoll() { kubectl get netpol "$1" -o wide; } # get networkpolicies + output
function kgnetpoly() { kubectl get netpol "$1" -o yaml; } # get networkpolicies + output
function kgnoj() { kubectl get no "$1" -o json; } # get nodes + output
function kgnol() { kubectl get no "$1" -o wide; } # get nodes + output
function kgnoy() { kubectl get no "$1" -o yaml; } # get nodes + output
function kgnsj() { kubectl get ns "$1" -o json; } # get namespaces + output
function kgnsl() { kubectl get ns "$1" -o wide; } # get namespaces + output
function kgnsy() { kubectl get ns "$1" -o yaml; } # get namespaces + output
function kgpcj() { kubectl get pc "$1" -o json; } # get priorityclasses + output
function kgpcl() { kubectl get pc "$1" -o wide; } # get priorityclasses + output
function kgpcy() { kubectl get pc "$1" -o yaml; } # get priorityclasses + output
function kgpdbj() { kubectl get pdb "$1" -o json; } # get poddisruptionbudgets + output
function kgpdbl() { kubectl get pdb "$1" -o wide; } # get poddisruptionbudgets + output
function kgpdby() { kubectl get pdb "$1" -o yaml; } # get poddisruptionbudgets + output
function kgpoj() { kubectl get po "$1" -o json; } # get pods + output
function kgpol() { kubectl get po "$1" -o wide; } # get pods + output
function kgpoy() { kubectl get po "$1" -o yaml; } # get pods + output
function kgpvcj() { kubectl get pvc "$1" -o json; } # get persistentvolumeclaims + output
function kgpvcl() { kubectl get pvc "$1" -o wide; } # get persistentvolumeclaims + output
function kgpvcy() { kubectl get pvc "$1" -o yaml; } # get persistentvolumeclaims + output
function kgpvj() { kubectl get pv "$1" -o json; } # get persistentvolumes + output
function kgpvl() { kubectl get pv "$1" -o wide; } # get persistentvolumes + output
function kgpvy() { kubectl get pv "$1" -o yaml; } # get persistentvolumes + output
function kgquotaj() { kubectl get quota "$1" -o json; } # get resourcequotas + output
function kgquotal() { kubectl get quota "$1" -o wide; } # get resourcequotas + output
function kgquotay() { kubectl get quota "$1" -o yaml; } # get resourcequotas + output
function kgrcj() { kubectl get rc "$1" -o json; } # get replicationcontrollers + output
function kgrcl() { kubectl get rc "$1" -o wide; } # get replicationcontrollers + output
function kgrcy() { kubectl get rc "$1" -o yaml; } # get replicationcontrollers + output
function kgrsj() { kubectl get rs "$1" -o json; } # get replicasets + output
function kgrsl() { kubectl get rs "$1" -o wide; } # get replicasets + output
function kgrsy() { kubectl get rs "$1" -o yaml; } # get replicasets + output
function kgsaj() { kubectl get sa "$1" -o json; } # get serviceaccounts + output
function kgsal() { kubectl get sa "$1" -o wide; } # get serviceaccounts + output
function kgsay() { kubectl get sa "$1" -o yaml; } # get serviceaccounts + output
function kgscj() { kubectl get sc "$1" -o json; } # get storageclasses + output
function kgscl() { kubectl get sc "$1" -o wide; } # get storageclasses + output
function kgscy() { kubectl get sc "$1" -o yaml; } # get storageclasses + output
function kgstsj() { kubectl get sts "$1" -o json; } # get statefulsets + output
function kgstsl() { kubectl get sts "$1" -o wide; } # get statefulsets + output
function kgstsy() { kubectl get sts "$1" -o yaml; } # get statefulsets + output
function kgsvcj() { kubectl get svc "$1" -o json; } # get services + output
function kgsvcl() { kubectl get svc "$1" -o wide; } # get services + output
function kgsvcy() { kubectl get svc "$1" -o yaml; } # get services + output
function kgy() { kubectl get -o yaml "$1"; } # get + output modifier
function klj() { kubectl logs -o json "$1"; } # logs + output modifier
function kll() { kubectl logs -o wide "$1"; } # logs + output modifier
function kln() { kubectl logs --all-namespaces "$1"; } # logs + scope modifier
function kly() { kubectl logs -o yaml "$1"; } # logs + output modifier
function ktn() { kubectl top --all-namespaces "$1"; } # top + scope modifier
```
