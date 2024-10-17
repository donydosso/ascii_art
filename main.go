package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Structure pour stocker les données du formulaire
type FormData struct {
	Phrase string
	Choix  string
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/result", resultHandler)
	http.Handle(
		`/static/`,
		http.StripPrefix("/static/", 
		http.FileServer(http.Dir(`static`))))
	http.HandleFunc("/404", notFoundHandler)
	fmt.Println("Serveur démarré sur le port 2080 : http://localhost:2080")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Gestionnaire de la page d'erreur 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: red; font-size: 40px;\">Error 404 - Page not found</p></div>")
}

// Gestionnaire de la page du formulaire
func formHandler(w http.ResponseWriter, r *http.Request) {

	// Obtenir le chemin de l'URL
	path := r.URL.Path

	// Vérifier si le chemin correspond à une des pages disponibles
	if path != "/" && path != "/result" {
		http.Redirect(w, r, "/404", http.StatusNotFound) // Rediriger vers la page d'erreur 404 si le chemin ne correspond à aucune page disponible
		return
	}

	// Vérifier si le formulaire a été soumis
	if r.Method == "POST" {
		// Récupérer les données du formulaire
		err := r.ParseForm()
		if err != nil {
			// Gérer l'erreur de parsing du formulaire
			http.Error(w, "Erreur lors de la récupération des données du formulaire", http.StatusInternalServerError)
			fmt.Println("Erreur lors de la récupération des données du formulaire : ",err)
			return
		}

		phrase := r.PostFormValue("phrase")
		choix := r.PostFormValue("choix")

		// Encoder le texte saisi
		encodedPhrase := url.QueryEscape(phrase)

		// Rediriger vers la page de résultat avec les données du formulaire
		redirectURL := "/result?phrase=" + encodedPhrase + "&choix=" + choix
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Afficher le formulaire
	tmpl := template.Must(template.ParseFiles("./templates/form.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		// Gérer l'erreur de parsing du formulaire
		http.Error(w, "Erreur lors de la lecture du template", http.StatusInternalServerError)
		fmt.Println("Erreur lors de la lecture du template : ", err)
		return
	}
}

// Gestionnaire de la page de résultat
func resultHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les données du formulaire de la requête précédente
	phrase := r.URL.Query().Get("phrase")
	choix := r.URL.Query().Get("choix")
	Choix := ""

	if choix == "option1" {
		Choix = "shadow"
	} else if choix == "option2" {
		Choix = "standard"
	} else if choix == "option3" {
		Choix = "thinkertoy"
	}

	// Afficher le résultat sur une nouvelle page
	tmpl := template.Must(template.ParseFiles("./templates/result.html"))
	data := FormData{
		Phrase: phrase,
		Choix:  choix,
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		// Gérer l'erreur de parsing du formulaire
		http.Error(w, "Erreur lors de la lecture du template", http.StatusInternalServerError)
		fmt.Println("Erreur lors de la lecture du template : ", err)
		return
	}

	// ASCII ART

	// vérifie que l'entrée est valide
	if phrase == "" {
		http.Error(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: red; font-size: 30px;\">Error 400 - Bad request : Empty input</p></div>", http.StatusBadRequest)
		return
	}

	phrases := strings.Split(phrase, "\r\n")
	if len(phrases) > 2 {
		http.Error(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: red; font-size: 30px;\">Error 400 - Bad request : Too many lines</p></div>", http.StatusBadRequest)
		return
	} else {
		phrase = strings.Join(phrases, "\n")
	}

	// le contenu du fichier template.txt est lu et stocké dans la variable "template"
	// si une erreur se produit lors de la lecture du fichier, le programme affiche une erreur et se termine
	templateFile := "./Banner/" + Choix + ".txt"
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		http.Error(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: blue; font-size: 30px;\">Error 500 - Internal Servor Error : Can't read file</p></div>", http.StatusInternalServerError)
		return
	}

	// Ensuite, le programme traite l'entrée en remplaçant tous les caractères "\n"
	// par des sauts de ligne. Si l'entrée est une chaîne de caractères vide,
	// le programme affiche simplement une nouvelle ligne et se termine

	input := strings.ReplaceAll(phrase, "\\n", "\n")
	if input == "\n" {
		http.Error(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: red; font-size: 30px;\">Error 400 - Bad request</p></div>", http.StatusBadRequest)
		return
	}

	// Si l'entrée contient des caractères
	// qui ne peuvent pas être affichés en ASCII, le programme affiche un message d'erreur et se termine.
	for i := 0; i < len(input); i++ {
		if (input[i] < 32 || input[i] > 127) && input[i] != 10 {
			http.Error(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: red; font-size: 30px;\">Error 400 - Bad request : Invalid character(s)</p></div>", http.StatusBadRequest)
			return
		}
	}

	// Le modèle de texte est divisé en 95 blocs de texte,
	// chacun représentant un caractère ASCII. Ces blocs sont stockés dans la variable "splitted".
	// Si le nombre de blocs n'est pas égal à 95, le modèle est considéré comme incorrect et le programme affiche un message d'erreur.
	splitted := strings.Split(string(template)[1:], "\n\n")
	if len(splitted) != 95 {
		http.Error(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: blue; font-size: 30px;\">Error 500 - Internal Servor Error</p></div>", http.StatusInternalServerError)
		return
	}

	// Le programme divise ensuite l'entrée en lignes et crée une représentation ASCII art de chaque ligne en
	// itérant sur chaque caractère de chaque ligne, en utilisant le modèle approprié stocké dans "splitted"
	lines := strings.Split(input, "\n")
	res := ""
	for _, line := range lines {
		if line == "" && res != "" {
			res += string('\n')
			continue
		}
		// Chaque caractère est représenté par une matrice 8x8 de caractères ASCII,
		// qui est concaténée pour produire une ligne ASCII art complète.
		// Les lignes ASCII art sont ensuite concaténées pour produire la sortie finale, qui est affichée.
		for row := 0; row < 8; row++ {
			for i := 0; i < len(line); i++ {
				temp := strings.Split(splitted[line[i]-32], "\n")[row]
				for j := 0; j < len(temp); j++ {
					res += string(temp[j])
				}
			}
			res += string('\n')
		}
	}

	result := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Résultat</title>
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css">
		<style>
			body {
				background-image: linear-gradient(to right, #8ec5fc, #e0c3fc);
				font-family: "Monotype Corsiva", cursive;
			}

			.container_ascii {
				display: table;
				margin: 40px auto;
				padding: 30px;
				background-color: rgba(255, 255, 255, 0.8);
				border-radius: 10px;
				box-shadow: 0 6px 18px rgba(0, 0, 0, 0.2);
			}

			.ascii-art-title {
				text-align: center;
				margin-bottom: 40px;
				font-size: 50px;
				font-weight: bold;
				color: #6c63ff;
				text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
				background-image: linear-gradient(to right, #8ec5fc, #e0c3fc);
				-webkit-background-clip: text;
				background-clip: text;
				-webkit-text-fill-color: transparent;
			}

			.ascii-art {
				margin-top: 20px;
				padding: 20px;
				border: 3px solid #ccc;
				border-radius: 6px;
				background-color: #f9f9f9;
				font-family: monospace;
				white-space: pre;
				display: inline-block;
				box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
			}

			.export {
				text-align: center;
				margin-top: 30px;
			}

			.export button {
				display: inline-block;
				padding: 10px 20px;
				background-color: #6c63ff;
				color: white;
				text-decoration: none;
				border-radius: 50px;
				font-size: 18px;
				font-weight: bold;
				box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
				transition: background-color 0.3s ease;
				border: none;
				cursor: pointer;
			}

			.export button:hover {
				background-color: #5048e5;
			}
		</style>
	</head>
	<body>
		<div class="container_ascii">
			<h1 class="ascii-art-title">ASCII-ART</h1>
			<div class="ascii-art">%s</div>
			<div class="export">
				<button class="exporter">Exporter</button>
			</div>
		</div>

		<script>
			function exportText() {
				var texteGenere = document.querySelector(".ascii-art").textContent;
				var blob = new Blob([texteGenere], { type: "text/plain" });
				var url = window.URL.createObjectURL(blob);
				var a = document.createElement("a");
				a.href = url;
				a.download = "ascii_art.txt";
				a.style.display = "none"; // Cacher l'élément d'ancre
				document.body.appendChild(a);
				a.click();
				document.body.removeChild(a);
				window.URL.revokeObjectURL(url);
			}

			document.addEventListener("DOMContentLoaded", function() {
				var exportButton = document.querySelector(".exporter");
				exportButton.addEventListener("click", exportText);
			});
		</script>
	</body>
	</html>
  `, res)

	// Afficher le résultat
	fmt.Fprintln(w, "<div style=\"display: flex; justify-content: center;\"><p style=\"font-weight: bold; color: green; font-size: 30px;\">OK 200 : Succeded</p></div>")
	fmt.Fprintln(w, result)
}
