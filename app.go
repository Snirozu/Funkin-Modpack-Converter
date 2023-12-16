package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	copy "github.com/otiai10/copy"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/slices"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	config = loadConfig();
}

func (a *App) ready(ctx context.Context) {
	fmt.Println(config);
	runtime.EventsEmit(a.ctx, "loadData", config);
}

func (a *App) SelectFolder(title string) string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
        Title: title,
    })
	if err != nil {
		log.Fatal(err);
		return "null";
	}
	return selection;
}

func (a *App) OpenSourceURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url);
}

func addLog(a *App, message string) {
	runtime.EventsEmit(a.ctx, "addLog", message);
}

func (a *App) SaveConfigVar(key string, value any) {
	switch key {
	case "enginePath": config.EnginePath = value.(string);
	}
	saveConfig();
}

func saveConfig() {
	appdata, err_ := os.UserConfigDir();
	if err_ != nil {
		return;
	}
	configPath := filepath.Join(appdata, "Funkin Modpack Converter", "config.json");
	jsonConfig, err := json.Marshal(config);
	if err != nil {
		return;
	}
	err2 := os.WriteFile(configPath, jsonConfig, 0644);
	if err2 != nil {
		fmt.Println(err2);
	}
}

func loadConfig() Config {
	appdata, err1 := os.UserConfigDir();
	if err1 != nil {
		return newConfig();
	}
	configPath := filepath.Join(appdata, "Funkin Modpack Converter", "config.json");
	data, err2 := os.ReadFile(configPath);
	if os.IsNotExist(err2) {
		config, err5 := json.Marshal(newConfig());
		if err5 != nil {
			return newConfig();
		}
		os.MkdirAll(filepath.Join(appdata, "Funkin Modpack Converter"), 0755);
		os.Create(configPath);
		os.WriteFile(configPath, config, 0644)
		return newConfig();
	} else if err2 != nil {
		return newConfig();
	}

	var config Config = newConfig();
	err4 := json.Unmarshal(data, &config)
	if err4 != nil {
		return newConfig();
	}

	return config;
}

var config Config;
type Config struct {
    EnginePath string
}

func newConfig() Config {
	config := Config{};
	config.EnginePath = "";
	return config;
}

func (a *App) StartConversion(modPath string, modID string) {
	// verify if variables are not empty
	modPath = strings.Trim(modPath, " ");
	modID = strings.Trim(modID, " ");
	if strings.Trim(config.EnginePath, " ") == "" {
		addLog(a, "Psych Engine path is empty!");
		return;
	}
	if modPath == "" {
		addLog(a, "Mod path is empty!");
		return;
	}
	if modID == "" {
		addLog(a, "Mod name is empty!");
		return;
	}

	//check if proper paths exists
	_, err1 := os.Stat(modPath); // check if mod path exists
	if os.IsNotExist(err1) {
		addLog(a, "Mod path: " + modPath + " doesn't exist!");
		return;
	}
	_, err2 := os.Stat(config.EnginePath); // check if engine path exists
	if os.IsNotExist(err2) {
		addLog(a, "Psych Engine path: " + modPath + " doesn't exist!");
		return;
	}
	info1, err3 := os.Stat(filepath.Join(config.EnginePath, "mods", "")); // check if engine mods dir exists
	if os.IsNotExist(err3) || !info1.IsDir() {
		addLog(a, "Psych Engine mods directory doesn't exist!");
		return;
	}
	info2, err4 := os.Stat(filepath.Join(modPath, "assets", "")); // check if mod has assets dir
	if os.IsNotExist(err4) || !info2.IsDir() {
		addLog(a, "Mod's assets directory doesn't exist!");
		return;
	}

	//convert
	modAssetsPath := filepath.Join(modPath, "assets", "");
	modPackPath := filepath.Join(config.EnginePath, "mods", modID);
	err5 := os.MkdirAll(modPackPath, 0755);
	if err5 != nil {
		addLog(a, "Couldn't create modpack location!");
		addLog(a, err5.Error());
		return;
	}

	//copy dirs from assets and put them in the modpack
	_, err6 := os.Stat(filepath.Join(modAssetsPath, "characters")); // copies characters
	if !os.IsNotExist(err6) {
		addLog(a, "Copying characters...");
		err := os.MkdirAll(filepath.Join(modPackPath, "characters"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "characters"), 
			filepath.Join(modPackPath, "characters"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	_, err7 := os.Stat(filepath.Join(modAssetsPath, "data"));
	if !os.IsNotExist(err7) {
		addLog(a, "Copying data...");
		err := os.MkdirAll(filepath.Join(modPackPath, "data"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "data"), 
			filepath.Join(modPackPath, "data"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	_, err8 := os.Stat(filepath.Join(modAssetsPath, "images"));
	if !os.IsNotExist(err8) {
		addLog(a, "Copying images...");
		err := os.MkdirAll(filepath.Join(modPackPath, "images"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "images"), 
			filepath.Join(modPackPath, "images"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}
	os.Remove(filepath.Join(modPackPath, "images", "alphabet.png"));
	os.Remove(filepath.Join(modPackPath, "images", "alphabet.xml"));

	_, err9 := os.Stat(filepath.Join(modAssetsPath, "music"));
	if !os.IsNotExist(err9) {
		addLog(a, "Copying music...");
		err := os.MkdirAll(filepath.Join(modPackPath, "music"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "music"), 
			filepath.Join(modPackPath, "music"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	_, errr1 := os.Stat(filepath.Join(modAssetsPath, "songs"));
	if !os.IsNotExist(errr1) {
		addLog(a, "Copying songs...");
		err := os.MkdirAll(filepath.Join(modPackPath, "songs"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "songs"), 
			filepath.Join(modPackPath, "songs"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	_, errr2 := os.Stat(filepath.Join(modAssetsPath, "sounds"));
	if !os.IsNotExist(errr2) {
		addLog(a, "Copying sounds...");
		err := os.MkdirAll(filepath.Join(modPackPath, "sounds"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "sounds"), 
			filepath.Join(modPackPath, "sounds"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	_, errr3 := os.Stat(filepath.Join(modAssetsPath, "stages"));
	if !os.IsNotExist(errr3) {
		addLog(a, "Copying stages...");
		err := os.MkdirAll(filepath.Join(modPackPath, "stages"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "stages"), 
			filepath.Join(modPackPath, "stages"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	_, errr4 := os.Stat(filepath.Join(modAssetsPath, "videos"));
	if !os.IsNotExist(errr4) {
		addLog(a, "Copying videos...");
		err := os.MkdirAll(filepath.Join(modPackPath, "videos"), 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "videos"), 
			filepath.Join(modPackPath, "videos"),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	_, errr5 := os.Stat(filepath.Join(modAssetsPath, "shared"));
	if !os.IsNotExist(errr5) {
		addLog(a, "Copying shared...");

		err_ := copy.Copy(
			filepath.Join(modAssetsPath, "shared"), 
			filepath.Join(modPackPath),
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	//copy mods/ outside of assets directory
	_, errr6 := os.Stat(filepath.Join(modPath, "mods"));
	if !os.IsNotExist(errr6) {
		addLog(a, "Copying mods...");
		err := os.MkdirAll(modPackPath, 0755);
		if err != nil { addLog(a, err.Error()); return; }

		err_ := copy.Copy(
			filepath.Join(modPath, "mods"),
			modPackPath,
		);
		if err_ != nil { addLog(a, err_.Error()); return; }
	}

	//create freeplay week
	if !os.IsNotExist(err7) && !os.IsNotExist(errr1) { //check if data and songs dirs exist
		addLog(a, "Creating Freeplay Week...");

		rre := os.MkdirAll(filepath.Join(modPackPath, "weeks"), 0755);
		if rre != nil { addLog(a, rre.Error()); return; }

		songs, err := os.ReadDir(filepath.Join(modPackPath, "songs"));
		if err != nil { addLog(a, err.Error()); return; }

		skipSongs := []string{};

		weeks, errl := os.ReadDir(filepath.Join(modAssetsPath, "weeks"));
		if errl != nil { addLog(a, errl.Error()); return; }
		for _, week := range weeks {
			data, er1 := os.ReadFile(filepath.Join(modAssetsPath, "weeks", week.Name()));
			if er1 != nil { addLog(a, er1.Error()); continue; }

			var weekData WeekData;
			er2 := json.Unmarshal(data, &weekData)
			if er2 != nil { addLog(a, er2.Error()); continue; }

			for _, song := range weekData.Songs {
				skipSongs = append(skipSongs, strings.ToLower(song[0].(string)));
			}
		}

		addLog(a, "Skipping songs: " + strings.Join(skipSongs, ", "));

		for _, song := range songs {
			if (song.IsDir() && !slices.Contains(skipSongs, strings.ToLower(song.Name()))) {
				diffs, err_ := os.ReadDir(filepath.Join(modPackPath, "data", song.Name()));
				if err_ != nil { addLog(a, err_.Error()); continue; }

				daDiffs := "";
				daSong := []any{};

				for _, diff := range diffs {
					if (strings.HasPrefix(diff.Name(), song.Name()) && strings.HasSuffix(diff.Name(), ".json")) {
						daDifficulty := diff.Name()[len(song.Name()):(len(diff.Name()) - 5)];
						if daDifficulty == "" {
							daDiffs += "Normal, ";
						} else {
							daDiffs += daDifficulty[1:len(daDifficulty)] + ", ";
						}

						data, er1 := os.ReadFile(filepath.Join(modPackPath, "data", song.Name(), diff.Name()));
						if er1 != nil { addLog(a, er1.Error()); continue; }

						var swagSong SongJson;
						er2 := json.Unmarshal(data, &swagSong)
						if er2 != nil { addLog(a, er2.Error()); continue; }

						daSong = []any{
							song.Name(),
							swagSong.Song.Player2,
							[]int{0, 0, 0},
						}

					}
				}

				if len(daDiffs) >= 2 {
					daDiffs = daDiffs[0:len(daDiffs) - 2];
				}

				weekData := WeekData{};
				weekData.StoryName = modID;
				weekData.Difficulties = daDiffs;
				weekData.HideFreeplay = false;
				weekData.WeekBackground = "";
				weekData.FreeplayColor = []int{
					255,
					255,
					255,
				};
				weekData.WeekBefore = "tutorial";
				weekData.StartUnlocked = true;
				weekData.WeekCharacters = []string{
					"dad",
					"bf",
					"gf",
				};
				weekData.Songs = [][]any{daSong};
				weekData.HideStoryMode = true;
				weekData.WeekName = "";
				weekData.HiddenUntilUnlocked = false;

				jsonWeek, erridk := json.Marshal(weekData);
				if erridk != nil { addLog(a, erridk.Error()); continue; }

				err3 := os.WriteFile(filepath.Join(modPackPath, "weeks", song.Name() + "-week.json"), jsonWeek, 0644);
				if err3 != nil { addLog(a, err3.Error()); continue; }

				addLog(a, "Saved week for song: " + song.Name());
			}
		}
	}

	addLog(a, "Done!");
}

type WeekData struct {
    StoryName string `json:"storyName"`
    Difficulties string `json:"difficulties"`
    HideFreeplay bool `json:"hideFreeplay"`
    WeekBackground string `json:"weekBackground"`
    FreeplayColor []int `json:"freeplayColor"`
	WeekBefore string `json:"weekBefore"`
	StartUnlocked bool `json:"startUnlocked"`
	WeekCharacters []string `json:"weekCharacters"`
	Songs [][]any `json:"songs"`
	HideStoryMode bool `json:"hideStoryMode"`
	WeekName string `json:"weekName"`
	HiddenUntilUnlocked bool `json:"hiddenUntilUnlocked"`
}

type SongJson struct {
    Song Song `json:"song"`
}

type Song struct {
    Player2 string `json:"player2"`
}