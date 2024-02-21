import {useState} from 'react';
import './App.css';
import {SelectFolder, OpenSourceURL, StartConversion, SaveConfigVar} from "../wailsjs/go/main/App";
import {EventsOn} from "../wailsjs/runtime/runtime";

let firstTime:boolean = false;

//gets called twice? im genuinely
function App() {
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    let logs:HTMLTextAreaElement;
    let isWindows = false;
    if (window.navigator.userAgent.includes("Windows")) {
        isWindows = true;
    }

    function selectEngineFolder() {
        SelectFolder("Select Psych Engine path").then((path) => {
            (document.getElementById("path-engine") as HTMLInputElement).value = path;
        });
    }

    function selectModFolder() {
        SelectFolder("Select Mod path").then((path) => {
            (document.getElementById("path-mod") as HTMLInputElement).value = path;
            resetModName();
        });
    }

    function resetModName() {
        let splitted = [];
        if (isWindows) {
            splitted = (document.getElementById("path-mod") as HTMLInputElement).value.split("\\");
        }
        else {
            splitted = (document.getElementById("path-mod") as HTMLInputElement).value.split("/");
        }
        (document.getElementById("mod-name") as HTMLInputElement).value = splitted[splitted.length - 1];
    }

    function convert() {
        SaveConfigVar("enginePath", (document.getElementById("path-engine") as HTMLInputElement).value);
        while (logs == null) {
            logs = (document.getElementById("logs") as HTMLTextAreaElement);
        }
        logs.value = "";
        (document.getElementById("convert") as HTMLButtonElement).disabled = true;
        (document.getElementById("convert") as HTMLButtonElement).style.opacity = "0.5";
        StartConversion((document.getElementById("path-mod") as HTMLInputElement).value, (document.getElementById("mod-name") as HTMLInputElement).value).then(() => {
            (document.getElementById("convert") as HTMLButtonElement).disabled = false;
            (document.getElementById("convert") as HTMLButtonElement).style.opacity = "1";
        });
    }

    (document.getElementById("source-link") as HTMLParagraphElement).onclick = (e) => {
        OpenSourceURL("https://github.com/Snirozu");
    }

    if (!firstTime) {
        firstTime = true;
        EventsOn("addLog", (msg) => {
            if (logs == null) {
                logs = (document.getElementById("logs") as HTMLTextAreaElement);
            }
            logs.value += (logs.value != "" ? "\n" : "") + " " + msg;
            logs.scrollTop = logs.scrollHeight;
        });

        EventsOn("loadData", (msg) => {
            (document.getElementById("path-engine") as HTMLInputElement).value = msg.EnginePath;
        });
    }

    return (        
        <div id="App">
            <h1> Configuration </h1>
            <div className='item'>
                <p> Select Psych Engine path here: </p>
                <div id="input" className="input-box">
                    <input id="path-engine" className="input" onChange={updateName} name="input" type="text"/>
                    <button className="btn" onClick={selectEngineFolder}> Select </button>
                </div>
            </div>

            <div className='item'>
                <p> Select Mod path here: </p>
                <div id="input" className="input-box">
                    <input id="path-mod" className="input" onChange={updateName} name="input" type="text"/>
                    <button className="btn" onClick={selectModFolder}> Select </button>
                </div>
            </div>

            <br/>

            <h1> Import Settings </h1>

            <div className='item'>
                <p> Choose the name of this mod: </p>
                <div id="input" className="input-box">
                    <input id="mod-name" className="input" onChange={updateName} name="input" type="text"/>
                    <button id="reset" className="btn" onClick={resetModName}> Reset </button>
                </div>
            </div>

            <br/>

            <button id="convert" className="big-btn" onClick={convert}> Convert </button>

            <br/><br/>

            <textarea id="logs" disabled placeholder='Logs will be displayed here.'></textarea>
        </div>
    )
}

export default App
