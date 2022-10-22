from cgitb import text
import traceback
from PyQt6.QtGui import QKeySequence, QPalette, QColor
from PyQt6.QtWidgets import *
from PyQt6 import uic
from PyQt6 import QtCore
from PyQt6 import QtGui
from appdata import appdata
from PyQt6.QtCore import Qt
import os
import shutil
import json

appdata.set_file_store("")
appdata.set_key_value_store("save")

engineDirs = ['characters', 'custom_events', 'custom_notetypes', 'data', 'fonts', 'images', 'music', 'scripts', 'shaders', 'songs', 'sounds', 'stages', 'videos', 'weeks']

class Window(QDialog):
    def __init__(self):
        super().__init__()
        uic.loadUi("window.ui", self)

        self.setWindowTitle("Psych Engine Mod to Mod Pack Converter")
        self.setWindowIcon(QtGui.QIcon('icon.ico'))

        self.update_palette()
        
        if appdata.get("lastPsychPath") == None:
            appdata["lastPsychPath"] = ""

        self.psychPathPlain.setText(appdata["lastPsychPath"])

        self.modPathBrowse.clicked.connect(self.on_browseModPath)
        self.psychPathBrowse.clicked.connect(self.on_browsePsychopath)
        self.convert.clicked.connect(self.on_convert)
        self.modPathPlain.textChanged.connect(self.on_modPathChange)

    def on_convert(self):
        self.log("clear")
        self.convert.setEnabled(False)

        try:
            if not os.path.exists(self.psychPathPlain.text()):
                self.log("Psych Engine path doesn't exist: '" + self.psychPathPlain.text() + "'")
                self.convert.setEnabled(True)
                return
            if not os.path.exists(self.psychPathPlain.text() + "/mods/"):
                self.log("Psych Engine mods directory doesn't exist: '" + self.psychPathPlain.text() + "/mods/" + "'")
                self.convert.setEnabled(True)
                return

            if not os.path.exists(self.modPathPlain.text()):
                self.log("Mod path doesn't exist: '" + self.modPathPlain.text() + "'")
                self.convert.setEnabled(True)
                return
            if not os.path.exists(self.modPathPlain.text() + "/assets/"):
                self.log("Mod's assets directory doesn't exist: '" + self.modPathPlain.text() + "/assets/" + "'")
                self.convert.setEnabled(True)
                return

            modAssetsLocation = self.modPathPlain.text() + "/assets/"
            modPackLocation = self.psychPathPlain.text() + "/mods/" + self.modNamePlain.text() + "/"
            self.createDir(modPackLocation)

            if (os.path.exists(modAssetsLocation + "characters/")):
                self.log("Copying characters...")
                self.createDir(modPackLocation + "characters/")
                for char in os.listdir(modAssetsLocation + "characters/"):
                    self.copy(modAssetsLocation + "characters/" + char, modPackLocation + "characters/" + char)

            if (os.path.exists(modAssetsLocation + "data/")):
                self.log("Copying data...")
                self.createDir(modPackLocation + "data/")
                for data in os.listdir(modAssetsLocation + "data/"):
                    if not self.pathExistsInEngine("assets/data/" + data):
                        self.copy(modAssetsLocation + "data/" + data, modPackLocation + "data/" + data)

            if (os.path.exists(modAssetsLocation + "images/")):
                self.log("Copying images...")
                self.copy(modAssetsLocation + "images/", modPackLocation + "images/")

            if (os.path.exists(modAssetsLocation + "music/")):
                self.log("Copying music...")
                self.copy(modAssetsLocation + "music/", modPackLocation + "music/")

            if (os.path.exists(modAssetsLocation + "songs/")):
                self.log("Copying songs...")
                self.createDir(modPackLocation + "songs/")
                for song in os.listdir(modAssetsLocation + "songs/"):
                    if not self.pathExistsInEngine("assets/songs/" + song):
                        self.copy(modAssetsLocation + "songs/" + song, modPackLocation + "songs/" + song)

            if (os.path.exists(modAssetsLocation + "sounds/")):
                self.log("Copying sounds...")
                self.copy(modAssetsLocation + "sounds/", modPackLocation + "sounds/")

            if (os.path.exists(modAssetsLocation + "stages/")):
                self.log("Copying stages...")
                self.createDir(modPackLocation + "stages/")
                for stage in os.listdir(modAssetsLocation + "stages/"):
                    self.copy(modAssetsLocation + "stages/" + stage, modPackLocation + "stages/" + stage)

            if (os.path.exists(modAssetsLocation + "videos/")):
                self.log("Copying videos...")
                self.copy(modAssetsLocation + "videos/", modPackLocation + "videos/")

            if (os.path.exists(modAssetsLocation + "shared/")):
                for folder in os.listdir(modAssetsLocation + "shared/"):
                    self.log("Copying shared " + folder + "...")
                    for folder2 in os.listdir(modAssetsLocation + "shared/" + folder + "/"):
                        self.copy(modAssetsLocation + "shared/" + folder + "/" + folder2, modPackLocation + folder + "/" + folder2)
            
            if (os.path.exists(self.modPathPlain.text() + "/mods/")):
                self.log("Copying mods...")
                for folder in os.listdir(self.modPathPlain.text() + "/mods/"):
                    if folder in engineDirs:
                        self.copy(self.modPathPlain.text() + "/mods/" + folder, modPackLocation + folder)
                    else:
                        self.copy(self.modPathPlain.text() + "/mods/" + folder, modPackLocation)

            if os.path.exists(modAssetsLocation + "songs/") and os.path.exists(modAssetsLocation + "data/"):
                self.log("Creating Freeplay Week")
                self.createDir(modPackLocation + "weeks/")
                for song in os.listdir(modAssetsLocation + "songs/"):
                    if (os.path.isdir(modAssetsLocation + "songs/" + song)):
                        diffics = []
                        jsonSong = []
                        for diff in os.listdir(modAssetsLocation + "data/" + song):
                            if not diff.endswith(".json") or not diff.startswith(song):
                                continue
                            
                            daDiff = None
                            dif = diff[0:(len(diff) - 5)]
                            if "-" not in dif or dif.lower() == song.lower():
                                daDiff = "Normal"
                            else:
                                difL = dif.split("-")
                                dif = difL[len(difL) - 1]
                                daDiff = dif

                            if daDiff is not None and daDiff not in diffics:
                                diffics.append(daDiff)

                            if (jsonSong == []):
                                with open(modAssetsLocation + "data/" + song + "/" + diff, 'r', encoding="utf8") as f:
                                    data = json.load(f)
                                    jsonSong = [
                                        song,
                                        data["song"]["player2"],
                                        [0, 0, 0]
                                    ]
                        diffString = ""
                        i = 0
                        for d in diffics:
                            diffString += d
                            if (i != len(diffics) - 1):
                                diffString += ", "
                            i+=1

                        x = {
                            "storyName": self.modNamePlain.text(),
                            "difficulties": diffString,
                            "hideFreeplay": False,
                            "weekBackground": "",
                            "freeplayColor": [
                                255,
                                255,
                                255
                            ],
                            "weekBefore": "tutorial",
                            "startUnlocked": True,
                            "weekCharacters": [
                                "dad",
                                "bf",
                                "gf"
                            ],
                            "songs": [jsonSong],
                            "hideStoryMode": True,
                            "weekName": "",
                            "hiddenUntilUnlocked": False
                        }
                        with open(modPackLocation + "weeks/" + song + "week.json", 'w', encoding="utf8") as f:
                            json.dump(x, f)

            self.log("Done")
            self.convert.setEnabled(True)

        except Exception as exc:
            self.log(str(traceback.format_exc()))

    def log(self, str):
        if str == "clear":
            self.info.setPlainText("")
            return
        self.info.setPlainText(self.info.toPlainText() + str + "\n")

    def copy(self, fro, to):        
        if (os.path.isdir(fro)):
            shutil.copytree(fro, to, dirs_exist_ok=True)
        else:
            if os.path.exists(fro) and os.path.exists(to):
                self.log("WARNING: Ignoring path '" + fro + "'")
                return
            shutil.copyfile(fro, to)

    def createDir(self, path):
        if not os.path.exists(path):
            os.mkdir(path)
        else:
            self.log("WARNING: '" + path + "' already exists")

    def pathExistsInEngine(self, path):
        return os.path.exists(self.psychPathPlain.text() + "/" + path) and os.path.exists(self.modPathPlain.text() + "/" + path)
        #spltPsychPath = str(self.psychPathPlain.text()).split("/")
        #return self.psychPathPlain.text()[0:(len(str(self.psychPathPlain.text())) - len(spltPsychPath[len(spltPsychPath) - 1]))]

    def on_modPathChange(self, text):
        spltxt = str(text).split("/")
        self.modNamePlain.setText(spltxt[len(spltxt) - 1])

    def on_browseModPath(self):
        dlg = QFileDialog()
        dlg.setFileMode(QFileDialog.FileMode.Directory)
        dlg.setFilter(QtCore.QDir.Filter.Dirs)

        filenames = []

        if dlg.exec():
            filenames = dlg.selectedFiles()

            #flist = ""
            #for dir in filenames:
            #    flist += dir + "* "
                
            self.modPathPlain.setText(filenames[0])

    def on_browsePsychopath(self):
        dlg = QFileDialog()
        dlg.setFileMode(QFileDialog.FileMode.Directory)
        dlg.setFilter(QtCore.QDir.Filter.Dirs)

        filenames = []

        if dlg.exec():
            filenames = dlg.selectedFiles()

            appdata["lastPsychPath"] = filenames[0]
            self.psychPathPlain.setText(filenames[0])

    def update_palette(self):
        palette = QPalette()
        palette.setColor(QPalette.ColorRole.Window, QColor(53, 53, 53))
        palette.setColor(QPalette.ColorRole.WindowText, Qt.GlobalColor.white)
        palette.setColor(QPalette.ColorRole.Base, QColor(25, 25, 25))
        palette.setColor(QPalette.ColorRole.AlternateBase, QColor(53, 53, 53))
        palette.setColor(QPalette.ColorRole.ToolTipBase, Qt.GlobalColor.white)
        palette.setColor(QPalette.ColorRole.ToolTipText, Qt.GlobalColor.black)
        palette.setColor(QPalette.ColorRole.Text, Qt.GlobalColor.black)
        palette.setColor(QPalette.ColorRole.Button, QColor(53, 53, 53))
        palette.setColor(QPalette.ColorRole.ButtonText, Qt.GlobalColor.black)
        palette.setColor(QPalette.ColorRole.BrightText, Qt.GlobalColor.red)
        palette.setColor(QPalette.ColorRole.Link, QColor(42, 130, 218))
        palette.setColor(QPalette.ColorRole.Highlight, QColor(42, 130, 218))
        palette.setColor(QPalette.ColorRole.HighlightedText, Qt.GlobalColor.black)
        app.setPalette(palette)

        infoPalette = QPalette()
        infoPalette.setColor(QPalette.ColorRole.Text, Qt.GlobalColor.white)
        self.info.setPalette(infoPalette)

app = QApplication([])
dialog = Window()
dialog.show()
app.exec()
