package gpsfiles

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type MapAccuracy string

const (
	ACCURACY_1  MapAccuracy = "1M"
	ACCURACY_5  MapAccuracy = "5M"
	ACCURACY_25 MapAccuracy = "25M"
)

type Link struct {
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
}

type Entry struct {
	Link Link `xml:"link"`
}

type Feed struct {
	Entries []Entry `xml:"entry"`
}

func computeCapabilitiesURL(nb int, accuracy MapAccuracy) (url string) {
	// Motif spécifique au département dans l'URL (ex: D069 pour le rhône)
	depart := fmt.Sprintf("D%03d", nb)

	if accuracy == ACCURACY_25 {
		url = fmt.Sprintf("https://data.geopf.fr/telechargement/resource/BDALTI?zone=%s", depart)
	} else {
		url = fmt.Sprintf("https://data.geopf.fr/telechargement/resource/RGEALTI?zone=%s", depart)
	}

	return
}

func fetchXMLFromUrl(url string) (feed Feed) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/atom+xml") // On force le format XML
	req.Header.Set("User-Agent", "Mozilla/5.0")      // On simule un navigateur

	rep, err := client.Do(req)

	if err != nil {
		fmt.Println("impossible de joindre l'URL :")
		log.Fatal(err)
	}

	defer rep.Body.Close()

	decoder := xml.NewDecoder(rep.Body)

	if err := decoder.Decode(&feed); err != nil {
		log.Fatalf("Erreur lors du parsing XML : %v", err)
	}

	return
}

func downloadFromURL(url string, filename string) {
	out, err := os.Create(filename) // Création du fichier local pour écrire les données téléchargées
	if err != nil {
		fmt.Println("erreur création fichier :")
		log.Fatal(err)
	}

	defer out.Close() // Assure la fermeture du fichier à la fin de la fonction (uniquement par sécurité car on le fait manuellement après)

	// On réutilise le client
	clientDownload := &http.Client{}
	reqDownload, _ := http.NewRequest("GET", url, nil)

	// On simule un navigateur pour passer l'erreur 403
	reqDownload.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	reqDownload.Header.Set("Accept", "*/*")

	fichier, err := clientDownload.Do(reqDownload) // Téléchargement du fichier
	if err != nil {
		fmt.Println("erreur téléchargement fichier :")
		log.Fatal(err)
	}

	defer fichier.Body.Close()

	if fichier.StatusCode != http.StatusOK { // Vérification du code de statut HTTP
		fmt.Println("erreur serveur :", fichier.Status)
		log.Fatal("Wrong status code")
	}

	counter := &WriteCounter{
		Size: uint64(fichier.ContentLength),
	}

	// 4. Copier les données
	// io.TeeReader divise le flux: il lit resp.Body, écrit dans 'counter',
	// et retourne les données pour qu'elles soient écrites dans 'out' par io.Copy.
	if _, err = io.Copy(out, io.TeeReader(fichier.Body, counter)); err != nil {
		log.Fatal(err)
	}

	fmt.Println(" | \033[32mDownload finished\033[0m")

	out.Close() // On ferme le fichier avant de le décompresser
}

func unzip(zipPath string, outputFolder string) {
	// Décompression du fichier .7z
	fmt.Print("└─> Decompressing the file...")

	cmd := exec.Command("7zz", "e", zipPath, "-o"+outputFolder, "-y", "-r", "*.asc") // commande pour extraire uniquement les fichiers .asc dans le dossier bd (e pour extraire, o pour choisir le répertoire, y pour répondre oui à toutes les questions si besoin, r pour parcourir tous les sous-dossiers du .7z)

	// On lance la commande et on attend la fin
	err := cmd.Run()
	if err != nil {
		fmt.Println("Erreur lors de la décompression du fichier :")
		log.Fatal(err)
	}

	os.Remove(zipPath) // On supprime le fichier .7z après extraction
	fmt.Println(" | \033[32mExtraction finished\033[0m")
}

func DownloadUnzipDepartment(nb int, accuracy MapAccuracy) {
	fmt.Printf("Starting download procedure for department n°\033[1m%02d\033[0m at \033[1m%s\033[0m\n", nb, string(accuracy))

	// Finding the ressource URL
	capabilitiesURL := computeCapabilitiesURL(nb, accuracy)

	var resourcesURL string

	for _, entry := range fetchXMLFromUrl(capabilitiesURL).Entries {
		if strings.Contains(entry.Link.Href, "1M") && accuracy == ACCURACY_1 {
			resourcesURL = entry.Link.Href
			break
		} else if strings.Contains(entry.Link.Href, "5M") && accuracy == ACCURACY_5 {
			resourcesURL = entry.Link.Href
			break
		} else if strings.Contains(entry.Link.Href, "25M") && accuracy == ACCURACY_25 {
			resourcesURL = entry.Link.Href
			break
		}
	}

	if resourcesURL == "" {
		log.Fatal("Didn't find the accuracy in the capabilities")
	}

	// Finding the download URLs
	var downloadURLs []string

	for _, entry := range fetchXMLFromUrl(resourcesURL).Entries {
		if entry.Link.Type == "application/x-7z-compressed" {
			downloadURLs = append(downloadURLs, entry.Link.Href)
		}
	}

	if len(downloadURLs) == 0 {
		log.Fatal("No ressources available for this accuracy")
	}

	err := os.MkdirAll("db", 0755) // Assure que le dossier "bd" existe (le crée si nécessaire) ; err est nil si le dossier est prêt

	if err != nil {
		log.Fatalf("Error while creating the DB folder: %v", err)
	}

	for _, url := range downloadURLs {
		filename := path.Base(url)
		downloadFromURL(url, filename)
		unzip(filename, "./db/"+string(accuracy))
		// unzip(filename, "./db/"+string(accuracy)+"/"+strconv.Itoa(nb))
	}
}

func DownloadAllDepartements(accuracy MapAccuracy) {
	for compteur := 1; compteur < 96; compteur++ { // On ne s'occupe que des départements métropolitains (1 à 95)
		if compteur == 20 { // On ignore le 20 car c'est la Corse (2A et 2B)
			continue
		}

		DownloadUnzipDepartment(compteur, accuracy)
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

// // Used to see progress of download
type WriteCounter struct {
	Total uint64 // Bytes already downloaded
	Size  uint64 // Total amount to download
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	// Deletes line and go to begining
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	percent := float64(wc.Total) / float64(wc.Size) * 100

	fmt.Printf("\r└─> Downloading... %.2f%% (%.2f MB / %.2f MB)",
		percent,
		float64(wc.Total)/1024/1024,
		float64(wc.Size)/1024/1024)
}
