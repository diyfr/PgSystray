package main

import (
	"fmt"
	"log"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

var (
	settingsDlg = &SettingsDialogWindow{model: NewAVModel()}

	db *walk.DataBinder
	SD Dialog
)

type SettingsDialogWindow struct {
	*walk.Dialog
	model               *AVModel
	availableVersionsLB *walk.ListBox
	existingVersionsCMB *walk.ComboBox
	portPG              *walk.NumberEdit
	usernameLE          *walk.LineEdit
	localeLE            *walk.LineEdit
	acceptPB            *walk.PushButton
	cancelPB            *walk.PushButton
}

type AVItem struct {
	name  string
	value string
}

type AVModel struct {
	walk.ListModelBase
	items []AVItem
}

func NewAVModel() *AVModel {
	av := checkAvailableVersions()

	m := &AVModel{items: make([]AVItem, len(av))}

	for i, v := range av {
		name := v
		value := v
		m.items[i] = AVItem{name, value}
	}
	return m
}

func (m *AVModel) ItemCount() int {
	return len(m.items)
}

func (m *AVModel) Value(index int) interface{} {
	return m.items[index].name
}

func (sd *SettingsDialogWindow) evItemActivated() {
	value := settingsDlg.model.items[settingsDlg.availableVersionsLB.CurrentIndex()].value
	answer := AskQuestion(fmt.Sprintf(strAUSI, value))
	if answer == 1 {
		log.Printf(strInstalling, value)
		conf.UsedVersion = value
		setPaths()
		go install(value)
	}
}

func ShowSettingsDialog() {
	SD = Dialog{
		AssignTo: &settingsDlg.Dialog,
		Title:    strSettings,
		DataBinder: DataBinder{
			AssignTo:   &db,
			DataSource: conf,
		},
		Size:      Size{Width: 250, Height: 250},
		MinSize:   Size{Width: 250, Height: 250},
		MaxSize:   Size{Width: 250, Height: 250},
		FixedSize: true,
		Layout:    VBox{Margins: Margins{Left: 5, Top: 5, Right: 5, Bottom: 5}, Spacing: 5, MarginsZero: false, SpacingZero: false},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					Label{
						Text: "Versions existantes:",
					},
					ComboBox{
						AssignTo: &settingsDlg.existingVersionsCMB,
						Value:    Bind("UsedVersion"),
						Model:    checkExistingVersions(),
					},

					Label{
						Text:        "Version disponibles:",
						ToolTipText: "Double clic pour installer",
					},
					ListBox{
						AssignTo:        &settingsDlg.availableVersionsLB,
						MinSize:         Size{Width: 50, Height: 70},
						MaxSize:         Size{Width: 50, Height: 70},
						Model:           settingsDlg.model,
						OnItemActivated: settingsDlg.evItemActivated,
						ToolTipText:     "Double clic pour installer",
					},
				},
			},
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Compte PostgreSQL:",
					},
					LineEdit{
						AssignTo:    &settingsDlg.usernameLE,
						Text:        Bind("Username"),
						ToolTipText: "postgres par défaut",
						CueBanner:   "postgres par défaut",
					},
					Label{
						Text: "Port",
					},
					NumberEdit{
						AssignTo:  &settingsDlg.portPG,
						Value:     float64(conf.Port),
						Decimals:  0,
						Increment: 1,
						OnValueChanged: func() {
							conf.Port = int(settingsDlg.portPG.Value())
						},
					},
					Label{
						Text: "Locale:",
					},
					LineEdit{
						AssignTo:    &settingsDlg.localeLE,
						Text:        Bind("Locale"),
						ToolTipText: "fr_FR par défaut",
						CueBanner:   "fr_FR par défaut",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &settingsDlg.acceptPB,
						Text:     strOK,
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Println(err)
								return
							}
							saveConfig()
							setPaths()
							checkServerStatus()
							settingsDlg.Hide()
						},
					},
					PushButton{
						AssignTo: &settingsDlg.cancelPB,
						Text:     strCancel,
						OnClicked: func() {
							settingsDlg.Hide()
						},
					},
				},
			},
		},
	}
	SD.Run(nil)
}
