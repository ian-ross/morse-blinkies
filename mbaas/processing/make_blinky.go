package processing

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Make is the main driver function to run the Python blinky creation
// code.
func (bm *BlinkyMaker) Make(text string, tmpl *template.Template) (string, error) {
	jobDir := path.Join(bm.WorkDir, "blinky-job")
	projName := strings.ReplaceAll(strings.ToLower(text), " ", "-")

	if err := os.RemoveAll(jobDir); err != nil {
		return "", err
	}
	if err := exec.Command("cp", "-r", bm.TemplateDir, jobDir).Run(); err != nil {
		return "", err
	}
	if err := os.Chdir(jobDir); err != nil {
		return "", err
	}
	if err := renames("PROJECT_NAME", projName,
		[]string{"pro", "sch", "kicad_pcb"}); err != nil {
		return "", err
	}

	output, err := exec.Command("python", bm.Script, text).CombinedOutput()
	if err != nil {
		return "", err
	}

	if err := renames("process-mbaas", projName, []string{"net", "info"}); err != nil {
		return "", err
	}
	if err := removes("process-mbaas",
		[]string{".log", ".erc", "_lib_sklib.py"}); err != nil {
		return "", err
	}

	logfp, err := os.Create(projName + ".log")
	if err != nil {
		return "", err
	}
	defer logfp.Close()
	if _, err := logfp.Write(output); err != nil {
		return "", err
	}

	zipFile := path.Join(bm.OutputDir, projName+".zip")
	if err := exec.Command("zip", "-r", zipFile, ".").Run(); err != nil {
		return "", err
	}

	htmlfp, err := os.Create(path.Join(bm.OutputDir, projName+".html"))
	if err != nil {
		return "", err
	}
	infofp, err := os.Open(projName + ".info")
	if err != nil {
		return "", err
	}
	dec := json.NewDecoder(infofp)
	var info map[string]interface{}
	if err := dec.Decode(&info); err != nil {
		return "", err
	}
	info["URL"] = "/output/" + projName + ".zip"
	info["SparklineWidth"] = fmt.Sprintf("%d", 10*len(info["padded_sequence"].(string)))

	if err := tmpl.ExecuteTemplate(htmlfp, "output", info); err != nil {
		return "", err
	}

	return "/output/" + projName + ".html", nil
}

func renames(in string, out string, types []string) error {
	for _, t := range types {
		if err := os.Rename(in+"."+t, out+"."+t); err != nil {
			return err
		}
	}
	return nil
}

func removes(in string, types []string) error {
	for _, t := range types {
		if err := os.Remove(in + t); err != nil {
			return err
		}
	}
	return nil
}
