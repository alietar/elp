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
> Les données topographiques prennet de la place, par département : ~xMo pour 25 m, ~xMo pour 5 m, ~xMo pour 1 m.

```console
// Syntaxe
go run main.go -accuracy-1|-accuracy-5|-accuracy-25 -dl-some <n°depart>|-dl-all

// Exemple pour le rhône en précision 5m
go run main.go -accuracy-5 -dl-some 69
```

# Usage

```console
go run main.go
```

Puis aller sur [http://localhost:8080](htpp://localhost:8080)


# Benchmarks

Outils utilisés :
- [**pprof**](https://github.com/google/pprof) : `go run main.go -perf` puis `go tool pprof -http=:8000 cpu.prof`
- [**plow**](https://github.com/six-ddc/plow) : `plow -c 10 http://127.0.0.1:8080/points --rate 30/1s --body='{"lat":45.7838052,"lng":4.821928,"deniv":1,"accuracy":1}' -m POST -d 10s`

# Pour aller plus loin

Plusieurs éléments pourraient être mis en oeuvre pour étendre la précision et l'utilisation de notre projet :  

- Prise en compte de la Corse et des territoires d'Outre-Mer