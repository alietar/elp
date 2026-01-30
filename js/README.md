
<div align="center">
  <img src="https://www.espritjeu.com/upload/image/flip-7-p-image-98429-grande.jpg" alt="IGN Logo" width="100px" />
</div>

<p align="center">JS - 3TC Groupe 32</p>

<div align="center">
  <img alt="JavaScript" src="https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E">
  <img alt="React" src="https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB">
</div>

# Flip 7 made in JS

## Présentation du projet

Nous travaillons sur la recréation du jeu Flip 7 en JavaScript, en divisant notre projet en deux modules. Le premier module se charge de la logique du jeu, incluant la gestion de la main du joueur, le lancement des parties et des manches, ainsi qu'une classe dédiée à la résolution des actions. Le deuxième module est une interface construite à l’aide du paquet ink, permettant de visualiser la carte du jeu et de faciliter les choix du joueur.

## Utilisation

Pour construire le projet, il suffit d'exécuter la commande suivante dans le répertoire /JS :

```bash
npm run build 
```

Ensuite, pour lancer le programme, utilisez la commande suivante :

```bash
node dist/cli.js
```

## Comment y jouer ?

Au démarrage, le jeu vous invite à : 
- Choisir le nombre de joueurs 
- Nommer les joueurs

Une fois la partie lancée, chaque joueur peut, à son tour :
- Tirer une carte
- Voir votre main
- Utiliser le helper pour voir le pourcentage de chances de tirer un doublon 
- Choisir de s'arrêter


