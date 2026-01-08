package algo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// prend en paramètres le numéro de département et télécharge la base de données correspondant
func downloadUnzipDB(depart int) string {

	err := os.MkdirAll("bd", 0755) // Assure que le dossier "bd" existe (le crée si nécessaire) ; err est nil si le dossier est prêt

	if err != nil {
		fmt.Println("Erreur création dossier BD :", err)
		return ""
	}

	motif := fmt.Sprintf("D%03d", depart) // Motif spécifique au département dans l'URL (ex: D069 pour le rhône), 03d signifie 3 chiffres (zéro padding devant si nécessaire)

	apiURL := fmt.Sprintf("https://data.geopf.fr/telechargement/resource/BDALTI?zone=%s", motif)

	// Pour récupérer les headers
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Accept", "application/atom+xml") // On force le format XML
	req.Header.Set("User-Agent", "Mozilla/5.0")      // On simule un navigateur

	rep, err := client.Do(req)

	if err != nil {
		fmt.Println("impossible de joindre l'URL :", err)
		return ""
	}
	defer rep.Body.Close()

	data, err := io.ReadAll(rep.Body) // data contient le contenu du fichier sous forme de byte
	if err != nil {
		fmt.Println("erreur lecture catalogue :", err)
		return ""
	}

	contenu := string(data)

	// Dans le XML, on cherche le contenu entre <title>...</title> situé à l'intérieur de la balise <entry>

	// On cherche l'entrée
	posEntry := strings.Index(contenu, "<entry>")
	if posEntry == -1 {
		fmt.Println("Aucune entrée trouvée pour ce département")
		return ""
	}

	reste := contenu[posEntry:] // On travaille à partir de la position de l'entrée
	startTag := "<title>"
	endTag := "</title>"

	// On cherche le titre à l'intérieur de cette entrée
	posStart := strings.Index(reste, startTag) + len(startTag)
	posEnd := strings.Index(reste, endTag)

	if posStart == -1 || posEnd == -1 {
		fmt.Println("Impossible d'extraire le nom de la ressource")
		return ""
	}

	subResource := reste[posStart:posEnd] // Nom de la ressource à télécharger
	nomFichier := subResource + ".7z"
	finalURL := fmt.Sprintf("https://data.geopf.fr/telechargement/download/BDALTI/%s/%s", subResource, nomFichier)

	out, err := os.Create(nomFichier) // Création du fichier local pour écrire les données téléchargées
	if err != nil {
		fmt.Println("erreur création fichier :", err)
		return ""
	}

	defer out.Close() // Assure la fermeture du fichier à la fin de la fonction (uniquement par sécurité car on le fait manuellement après)

	fmt.Printf("Téléchargement lancé : %s\n", nomFichier)

	// On réutilise le client
	clientDownload := &http.Client{}
	reqDownload, _ := http.NewRequest("GET", finalURL, nil)

	// On simule un navigateur pour passer l'erreur 403
	reqDownload.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	reqDownload.Header.Set("Accept", "*/*")

	fichier, err := clientDownload.Do(reqDownload) // Téléchargement du fichier
	if err != nil {
		fmt.Println("erreur téléchargement fichier :", err)
		return ""
	}
	defer fichier.Body.Close()

	if fichier.StatusCode != http.StatusOK { // Vérification du code de statut HTTP
		fmt.Println("erreur serveur :", fichier.Status)
		return ""
	}

	_, err = io.Copy(out, fichier.Body) // Utilisation de io.Copy pour ne pas saturer la RAM (flux direct vers disque)
	if err != nil {
		fmt.Println("erreur pendant l'écriture :", err)
		return ""
	}

	fmt.Printf("Base de données téléchargée\n")

	out.Close() // On ferme le fichier avant de le décompresser

	// Décompression du fichier .7z
	fmt.Println("Décompression du fichier...")

	cmd := exec.Command("7z", "e", nomFichier, "-obd", "-y", "-r", "*.asc") // commande pour extraire uniquement les fichiers .asc dans le dossier bd (e pour extraire, o pour choisir le répertoire, y pour répondre oui à toutes les questions si besoin, r pour parcourir tous les sous-dossiers du .7z)

	// On lance la commande et on attend la fin
	err = cmd.Run()
	if err != nil {
		fmt.Println("Erreur lors de la décompression du fichier :", err)
		return ""
	}
	os.Remove(nomFichier) // On supprime le fichier .7z après extraction
	fmt.Println("Extraction terminée")

	return ""

}

func downloadAllDepartements() {
	for compteur := 1; compteur < 96; compteur++ { // On ne s'occupe que des départements métropolitains (1 à 95)
		if compteur == 20 { // On ignore le 20 car c'est la Corse (2A et 2B)
			continue
		}
		fmt.Printf("Département %02d en cours de téléchargement\n", compteur)
		downloadUnzipDB(compteur)
	}
}
