# GO

Idées d'algo :
 - Calcul de chemin à +-10m sur une carte
 - Calcul de la planéité du terrain/ville

To-Do List :
 - [x] Programme qui lis toutes les cases de la bd et récupère les limites de la case pour les condenser dans un tableau
 - [ ] Prendre en argument du main.go les coordonnées initiales (et +- de hauteur)
 - [x] Trouver le fichier correspondant à ces coordonnées
 - [x] Exécuter l'algo dans ce fichier en partant des coordonnées
 - [x] Retourner le résultat
 - [ ] Aller chercher dans les cases d'à côté pour éxécuter l'algo
 - [ ] Multithreading
 - [ ] Serveur TCP


## Algo :

On ne peut se déplacer que de gauche à droit, et de bas en haut, pas en diagonale.