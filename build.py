import os
import shutil as morshu
import sys
import PyInstaller.__main__

VERSION = "1.0"
NAME = "Funkin Mod Converter"

PyInstaller.__main__.run([
    'main.py',
    '--onefile',
    '--windowed',
    '--icon=icon.ico',
    "--name=" + NAME
])
morshu.copyfile("window.ui", "dist/window.ui")
morshu.copyfile("icon.ico", "dist/icon.ico")
if os.path.exists("dist/ui_window.py"):
    os.remove("dist/ui_window.py")
if sys.argv[1] == "package":
    morshu.make_archive(NAME + "-" + VERSION, 'zip', "dist/")
