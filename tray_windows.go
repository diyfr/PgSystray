//go:generate windres -i PgSystray.rc -O coff -o resource.syso
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lxn/walk"
)

var (
	mw                   *walk.MainWindow
	ni                   *walk.NotifyIcon
	statusPgServerAction *walk.Action
	startPgServerAction  *walk.Action
	stopPgServerAction   *walk.Action
	startPgShellAction   *walk.Action
	showSettingsAction   *walk.Action
	startPgAdmin         *walk.Action

	err error

	imgStatus   = "resources/windows/postgresql-16.ico"
	imgStart    = "resources/windows/play-16.ico"
	imgStop     = "resources/windows/stop-16.ico"
	imgConsole  = "resources/windows/console-16.ico"
	imgSettings = "resources/windows/settings-16.ico"
	imgSync     = "resources/windows/sync-16.ico"
	imgExit     = "resources/windows/exit-16.ico"
)

func ShowNotification(msg string) {
	if ni != nil {
		icon, err := walk.NewIconFromResourceId(100)
		if err != nil {
			log.Fatal(err)
		}
		if err := ni.ShowCustom(strTitle, msg, icon); err != nil {
			log.Fatal(err)
		}
	}
}

func SetNotificationToolTip(msg string) {
	if err := ni.SetToolTip(fmt.Sprintf("%s\n%s", strTitle, msg)); err != nil {
		log.Fatal(err)
	}
}

func CreateTray() {
	mw, err = walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}
	mw.Hide()

	icon, err := walk.NewIconFromResourceId(101)
	if err != nil {
		log.Fatal(err)
	}

	ni, err = walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	SetNotificationToolTip("")

	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		ShowNotification(strHelp)
	})

	// Start status item {{{
	statusPgServerAction = walk.NewAction()
	if err := statusPgServerAction.SetText(strStopped); err != nil {
		log.Fatal(err)
	}
	statusPgServerImg, stsImgErr := walk.NewBitmapFromFileForDPI(imgStatus, 96)
	if stsImgErr != nil {
		log.Printf("status resource - %s\n", stsImgErr.Error())
	}
	if err := statusPgServerAction.SetImage(statusPgServerImg); err != nil {
		log.Printf("statusPgServerAction - %v\n", err.Error())
	}
	statusPgServerAction.SetEnabled(false)
	if err := ni.ContextMenu().Actions().Add(statusPgServerAction); err != nil {
		log.Fatal(err)
	}
	// }}} End status item

	// Start separator item {{{
	if err := ni.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
		log.Fatal(err)
	}
	// }}} Start separator item

	// Start start item {{{
	startPgServerAction = walk.NewAction()
	if err := startPgServerAction.SetText(strStart); err != nil {
		log.Fatal(err)
	}
	startPgServerImg, _ := walk.NewBitmapFromFileForDPI(imgStart, 96)
	if err := startPgServerAction.SetImage(startPgServerImg); err != nil {
		log.Fatal(err)
	}
	startPgServerAction.SetEnabled(!serverStatus && len(conf.UsedVersion) > 0)

	startPgServerAction.Triggered().Attach(func() { go startPg() })
	if err := ni.ContextMenu().Actions().Add(startPgServerAction); err != nil {
		log.Fatal(err)
	}
	// }}} End start item

	// Start stop item {{{
	stopPgServerAction = walk.NewAction()
	if err := stopPgServerAction.SetText(strStop); err != nil {
		log.Fatal(err)
	}
	stopImg, _ := walk.NewBitmapFromFileForDPI(imgStop, 96)
	if err := stopPgServerAction.SetImage(stopImg); err != nil {
		log.Fatal(err)
	}
	stopPgServerAction.SetEnabled(serverStatus)

	stopPgServerAction.Triggered().Attach(func() { go stopPg() })
	if err := ni.ContextMenu().Actions().Add(stopPgServerAction); err != nil {
		log.Fatal(err)
	}
	// }}} End stop item

	if err := ni.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
		log.Fatal(err)
	}

	// Start shell item {{{
	startPgShellAction = walk.NewAction()
	if err := startPgShellAction.SetText(strStartShell); err != nil {
		log.Fatal(err)
	}
	startPgShellImg, _ := walk.NewBitmapFromFileForDPI(imgConsole, 96)
	if err := startPgShellAction.SetImage(startPgShellImg); err != nil {
		log.Fatal(err)
	}
	startPgShellAction.SetEnabled(false)
	startPgShellAction.Triggered().Attach(func() { go startShell() })
	if err := ni.ContextMenu().Actions().Add(startPgShellAction); err != nil {
		log.Fatal(err)
	}
	// }}} End shell item

	if err := ni.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
		log.Fatal(err)
	}

	// Start settings item {{{
	showSettingsAction = walk.NewAction()
	if err := showSettingsAction.SetText(strSettings); err != nil {
		log.Fatal(err)
	}
	settingsImg, _ := walk.NewBitmapFromFileForDPI(imgSettings, 96)
	if err := showSettingsAction.SetImage(settingsImg); err != nil {
		log.Fatal(err)
	}
	showSettingsAction.Triggered().Attach(func() {
		ShowSettingsDialog()
	})
	if err := ni.ContextMenu().Actions().Add(showSettingsAction); err != nil {
		log.Fatal(err)
	}
	// }}} End settings item

	// Start checkForUpdates item {{{
	startPgAdmin = walk.NewAction()
	if err := startPgAdmin.SetText(strPgAdmin); err != nil {
		log.Fatal(err)
	}
	cfuImg, _ := walk.NewBitmapFromFileForDPI(imgSync, 96)
	if err := startPgAdmin.SetImage(cfuImg); err != nil {
		log.Fatal(err)
	}
	startPgAdmin.Triggered().Attach(func() {
		checkNewestVersion()
	})
	if err := ni.ContextMenu().Actions().Add(startPgAdmin); err != nil {
		log.Fatal(err)
	}
	// }}} End checkForUpdates item

	// Start exit item {{{
	exitAction := walk.NewAction()
	if err := exitAction.SetText(strExit); err != nil {
		log.Fatal(err)
	}
	exitImg, _ := walk.NewBitmapFromFileForDPI(imgExit, 96)
	if err := exitAction.SetImage(exitImg); err != nil {
		log.Fatal(err)
	}

	exitAction.Triggered().Attach(func() { quit() })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
	// }}} End exit item

	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	if conf.CheckForUpdates {
		log.Println("Check for updates enabled")
		checkNewestVersion()
	}
	mw.Run()
}

func SetStatus(msg string) {
	statusPgServerAction.SetText(msg)
}

func EnableStart(enable bool) {
	startPgServerAction.SetEnabled(enable)
}

func EnableStop(enable bool) {
	stopPgServerAction.SetEnabled(enable)
}

func EnableShell(enable bool) {
	startPgShellAction.SetEnabled(enable)
}

func EnablePgAdmin(enable bool) {
	startPgAdmin.SetEnabled(enable)
}

func ShowMessage(msg string) {
	walk.MsgBox(nil, strTitle, msg, walk.MsgBoxIconInformation)
}

func AskQuestion(msg string) int {
	return walk.MsgBox(nil, strTitle, msg, walk.MsgBoxIconQuestion+walk.MsgBoxOKCancel+walk.MsgBoxDefButton1)
}

func AppQuit() {
	ni.SetVisible(false)
	walk.App().Exit(0)
	os.Exit(0)
}
