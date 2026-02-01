*Groupe 32 - Constant Léopold, Corbard Solène (gr 2), Liétar Anantole*

<div align="center">
  <img src="https://www.tonton-outdoor.com/media/cache/sylius_large/de/92/a187da7d97bb6b3d338052c775ab.png" alt="IGN Logo" width="100px" />
</div>

<p align="center">ELP, GO - 3TC Groupe 32</p>

<div align="center">
  <img alt="Golang" src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white">
  <img alt="HTTP" src="https://img.shields.io/badge/HTTP-green?style=for-the-badge">
</div>

# Sommaire

- [Présentation du projet](#presentation-du-projet)
- [Installtion et mise en place](#installation-et-mise-en-place)
- [Usage](#usage)
- [Benchmarks](#benchmarks)
- [Pour aller plus loin](#pour-aller-plus-loin)

# Présentation du projet

Notre projet a pour objectif de déterminer la zone atteignable autour d'un point selon une contrainte de dénivelé. L'utilisateur indique la position initiale souhaitée en coordonnées GPS WGS84 et le dénivelé maximal accepté (on cherche la zone avec +- ce dénivelé). Notre programme affiche alors une carte interactive affichant la zone en question. Nous utilisons pour cela les données topographiques de l'IGN fournissant l'altitude d'un point tous les 25m. Notre algorithme se concentre uniquement sur les départements français métropolitains (hors Corse).

L'interface a été réalisée dans le cadre du [projet ELM](../elm/README.md).

Le projet est accessible en ligne sur [reachable.lietar.net](reachable.lietar.net) (uniquement le Rhône en 25 m et 5 m).

# Installation et mise en place

### Prérequis

- **Go** v1.25.5
- **7z**

### Installation

```console
git clone https://github.com/alietar/elp.git
cd elp/go
```

### Téléchargement des données topographiques de l'IGN

> [!WARNING]
> Les données topographiques prennent de la place, par département : ~ 40 Mo pour 25 m, ~ 300 Mo pour 5 m, ~ 4 Go pour 1 m.

```console
// Syntaxe
go run main.go -7z=<chemin vers l'éxécutable de 7z, (défaut)7z> -accuracy=<1|5|(défaut)25> -dl-all|-dl-some <n°depart> <n°depart> ... 

// Exemple pour le Rhône et l'Isère en précision 5 m
go run main.go -accuracy=5 -dl-some 69 38

// Exemple pour tous les départements en précision 25 m avec 7z non conventionnel
go run main.go -7z=7zz -dl-all
```

# Usage

```console
go run main.go -port=<port serveur HTTP (défaut)8080>
```

Puis aller sur [http://localhost:8080](htpp://localhost:8080)

# Benchmarks

Outils utilisés :
- [**pprof**](https://github.com/google/pprof) : `go run main.go -perf` puis `go tool pprof -http=:8000 <chemin vers xxx.prof>`
- [**plow**](https://github.com/six-ddc/plow) : `plow -c 10 http://127.0.0.1:8080/points --rate 30/1s --body='{"lat":45.7838052,"lng":4.821928,"deniv":1,"accuracy":1}' -m POST -d 10s`

# Pour aller plus loin

Plusieurs éléments pourraient être mis en oeuvre pour étendre la précision et l'utilisation de notre projet :  

- Prise en compte de la Corse et des territoires d'Outre-Mer