package gpsfiles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type MapAccuracy uint8
type DepartmentNb uint8

const (
	ACCURACY_1 MapAccuracy = iota
	ACCURACY_5
	ACCURACY_25
)

type Social struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Url     string `xml:"id"`
	Content string `xml:"content"`
}

func FetchingDepartmentUrl(nb DepartmentNb, accuracy MapAccuracy) {

	// var users Users
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	// xml.Unmarshal(byteValue, &users)
}

func computeFirstURL(nb DepartmentNb, accuracy MapAccuracy) (url string) {
	// Motif spécifique au département dans l'URL (ex: D069 pour le rhône)
	depart := fmt.Sprintf("D%03d", nb)

	if accuracy == ACCURACY_25 {
		url = fmt.Sprintf("https://data.geopf.fr/telechargement/resource/BDALTI?zone=%s", depart)
	} else {
		url = fmt.Sprintf("https://data.geopf.fr/telechargement/resource/RGEALTI?zone=%s", depart)
	}

	return
}

// func fetchDataURL(url string) (dataUrl string) {
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Accept", "application/atom+xml") // On force le format XML
// 	req.Header.Set("User-Agent", "Mozilla/5.0")      // On simule un navigateur
//
// 	rep, err := client.Do(req)
//
// 	if err != nil {
// 		fmt.Println("impossible de joindre l'URL :", err)
// 		return ""
// 	}
// 	defer rep.Body.Close()
//
// 	data, err := io.ReadAll(rep.Body) // data contient le contenu du fichier sous forme de byte
// 	if err != nil {
// 		fmt.Println("erreur lecture catalogue :", err)
// 		return ""
// 	}
//
// 	// contenu := string(data)
//
// }

func downloadUnizpFile(path string) {

}

// prend en paramètres le numéro de département et télécharge la base de données correspondant
func DownloadUnzipDB(depart int) string {
	fmt.Printf("Starting download procedure for department %02d\n", depart)

	err := os.MkdirAll("db", 0755) // Assure que le dossier "bd" existe (le crée si nécessaire) ; err est nil si le dossier est prêt

	if err != nil {
		fmt.Println("Error while creating the DB folder:", err)
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

	fmt.Printf("Download started %s\n", nomFichier)

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

	fmt.Printf("Download finished\n")

	out.Close() // On ferme le fichier avant de le décompresser

	// Décompression du fichier .7z
	fmt.Println("Decompressing the file...")

	cmd := exec.Command("7zz", "e", nomFichier, "-odb", "-y", "-r", "*.asc") // commande pour extraire uniquement les fichiers .asc dans le dossier bd (e pour extraire, o pour choisir le répertoire, y pour répondre oui à toutes les questions si besoin, r pour parcourir tous les sous-dossiers du .7z)

	// On lance la commande et on attend la fin
	err = cmd.Run()
	if err != nil {
		fmt.Println("Erreur lors de la décompression du fichier :", err)
		return ""
	}
	os.Remove(nomFichier) // On supprime le fichier .7z après extraction
	fmt.Println("Extraction finished")

	return ""

}

func DownloadAllDepartements() {
	for compteur := 1; compteur < 96; compteur++ { // On ne s'occupe que des départements métropolitains (1 à 95)
		if compteur == 20 { // On ignore le 20 car c'est la Corse (2A et 2B)
			continue
		}

		DownloadUnzipDB(compteur)
	}
}

// On définit la structure attendue du JSON
type Depart struct {
	CodeDepartement string `json:"codeDepartement"`
}

// prend en paramètre les coordonnées x et y en wsg84 et retourne le département correspondant
func GetDepartement(x, y float64) int {
	url := fmt.Sprintf("https://geo.api.gouv.fr/communes?lat=%f&lon=%f&fields=codeDepartement", x, y)

	client := &http.Client{Timeout: 10 * time.Second} // Requête avec un Timeout
	rep, err := client.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return -1
	}
	defer rep.Body.Close() //rep.Body est le flux de données

	var data []Depart                             // car le JSON retourné est une liste
	err = json.NewDecoder(rep.Body).Decode(&data) // On décode le JSON qui est de la forme [{codeDepartement: "XX"}] en bytes
	if err != nil {
		fmt.Println("Erreur de décodage :", err)
		return -1
	}
	num, err := strconv.Atoi(data[0].CodeDepartement)
	if err != nil {
		fmt.Println("Erreur de conversion en entier :", err)
		return -1
	}
	return num
}
