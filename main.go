//go:build (windows || linux) && (386 || amd64)
// +build windows linux
// +build 386 amd64

package main

import (
	"os"
	"path/filepath"
)

const (
	configFile      = "PgSystray.json"
	strOK           = "OK"
	strCancel       = "Annuler"
	strExit         = "Quitter"
	strPgAdmin      = "Lancer PgAdmin"
	strHelp         = "Use the context menu to manage."
	strInit         = "Initialisation..."
	strInitFinished = "Initialisation terminée"
	strPSVF         = "Sélectionnez d'abord une version"
	strExtrV        = "Installation de la version %s"
	strExtrVF       = "Installation de la version %s terminée"
	strInstalling   = "Installation PostgreSQL %s\n"
	strInstallation = "Installation terminée"
	strSettings     = "Paramètres"
	strStart        = "Démarrer PostgreSQL"
	strStarted      = "Serveur démarré"
	strStarting     = "Démarrage shell psql"
	strStartShell   = "Démarrer shell psql"
	strStop         = "Arrêter PostgreSQL"
	strStopped      = "Serveur arrêté"
	strTitle        = "PostgreSQL"
	strSNR          = "Serveur non démarré!"
	strDNI          = "Base de donnée non initialisée!"
	strNIV          = "Aucune version installée"
	strDVPW         = "Téléchargement de la version %s. Patientez"
	strAUSI         = "Etes-vous sur de vouloir installer PostgreSQL %s?"
	strIIF          = "La version %s n'est pas installée!\n"
	strStopErr      = "Arrêt erreur: %s\n"
	strStartErr     = "Démarrage erreur: %s\n"
	strFNEErr       = "Les fichiers de la version %s n'existent pas!\n"
)

var (
	dir, _       = filepath.Abs(filepath.Dir(os.Args[0]))
	pgsqlBaseDir = filepath.Join(dir, "pgsql")
	downloadDir  = filepath.Join(dir, "pgsql-downloads")
	logsDir      = filepath.Join(dir, "logs")

	pgInitdb, pgCtl, pgShell, dataDir                                          string
	cmdInitDbArgs, cmdStartArgs, cmdStopArgs, cmdStatusArgs, cmdStartShellArgs []string

	serverStatus = false
	serverPid    int

	osName, osArch string
	archiveType    string
	conf           *Configuration
)

func init() {
	checkOs()
	checkArch()
	checkArchiveType()
}

func main() {
	conf = NewConfiguration()
	loadConfig()
	setPaths()
	CreateTray()
}
