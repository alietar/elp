*Groupe 32 - Constant Léopold, Corbard Solène (gr 2), Liétar Anantole*

<div align="center">
  <img src="https://www.tonton-outdoor.com/media/cache/sylius_large/de/92/a187da7d97bb6b3d338052c775ab.png" alt="IGN Logo" width="100px" />
</div>

<p align="center">ELP, ELM - 3TC Groupe 33</p>

<div align="center">
  <img alt="ELM" src="https://img.shields.io/badge/Elm-60B5CC?style=for-the-badge&logo=elm&logoColor=white">
  <img alt="HTTP" src="https://img.shields.io/badge/HTTP-green?style=for-the-badge">
</div>

# Sommaire

- [Présentation du projet](#presentation-du-projet)
- [Installtion et mise en place](#installation-et-mise-en-place)
- [Usage](#usage)
- [Benchmarks](#benchmarks)
- [Pour aller plus loin](#pour-aller-plus-loin)

# Présentation du projet

Pour notre projet, nous n'avons pas choisi un des sujets proposés, car nous voulions lier notre projet Go au projet ELM. En effet notre projet Go permet de calculer les points atteignables depuis des coordonnées fournies, en descendant ou en montant au maximum de plus ou moins x mètres. Nous voulions une visualisation graphique sur carte du résultat, d'où la complicité avec le projet ELM.
Ainsi sur notre site il est possible de cliquer sur la carte ou de renseigner manuellement les coordonnées, puis de choisir un delta d'altitude maxium et une précision, et finalement d'appuyer sur le bouton de calcul. Les données sont envoyées avec une requête POST au serveur, qui les traite puis renvoie une liste de carrés. Le site décode la liste et affiche le résultat sur la carte.

# Architecture des fichiers


## map.js

Pour l'affichage de la carte et des carrés nous avons utilisé la bilbiothèque [leaflet.js](https://leafletjs.com/).

## Carte.elm

La communication avec le code de map.js se fait via des ports définis dans Carte.elm.

## DrawSquare.elm

Ce fichier permet de calculer les quatres coins en coordonnées GPS de chaque carré.

## Interface.elm

Il s'agit de l'interface pour la saisie des coordonnées, du delta d'altitude et de la précision, ainsi que le bouton d'envoie des données qui sert aussi d'indicateur de statut. 

## UserApi.elm

Ce fichier définit les schémas de données pour les requêtes HTTP et gère l'envoie des coordonnées de départ.

## Main.elm

Finalement le Main coordonne les différents modules et gère l'update de chacun.

# Installation et mise en place

### Prérequis

- **ELM** v1.25.5

### Installation

```bash
git clone https://github.com/alietar/elp.git
cd elp/elm
elm make src/Main.elm --output elm.js

```

# Usage

```bash
# Ouvrir dans un naviguateur le fichier index.html
# Exemple:
firefox index.html&
```