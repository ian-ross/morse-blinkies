package processing

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"

	"github.com/ian-ross/morse-blinkies/mbaas/model"
)

// Make is the main driver function to run the Python blinky creation
// code.
func (bm *BlinkyMaker) Make(text string, rules *model.Rules,
	tmpl *template.Template) (string, string, error) {
	jobDir := path.Join(bm.WorkDir, "blinky-job")
	jsonRules, projName, fullProjName, err := rules.ProjectName(text)
	if err != nil {
		return "", "", err
	}

	if err := os.RemoveAll(jobDir); err != nil {
		return "", "", err
	}
	if err := exec.Command("cp", "-r", bm.TemplateDir, jobDir).Run(); err != nil {
		return "", "", err
	}
	if err := os.Chdir(jobDir); err != nil {
		return "", "", err
	}
	if err := renames("PROJECT_NAME", projName,
		[]string{"pro", "sch", "kicad_pcb"}); err != nil {
		return "", "", err
	}
	rulesfp, err := os.Create("rules.json")
	if err != nil {
		return "", "", err
	}
	if _, err := rulesfp.Write(jsonRules); err != nil {
		return "", "", err
	}
	rulesfp.Close()

	output, err := exec.Command("python", bm.Script, text, "rules.json").CombinedOutput()
	if err != nil {
		return "", string(output), err
	}

	if err := renames("process-mbaas", projName, []string{"net", "info"}); err != nil {
		return "", "", err
	}
	if err := removes("process-mbaas",
		[]string{".log", ".erc", "_lib_sklib.py"}); err != nil {
		return "", "", err
	}

	logfp, err := os.Create(projName + ".log")
	if err != nil {
		return "", "", err
	}
	if _, err := logfp.Write(output); err != nil {
		return "", "", err
	}
	logfp.Close()

	zipFile := path.Join(bm.OutputDir, fullProjName+".zip")
	if err := exec.Command("zip", "-r", zipFile, ".").Run(); err != nil {
		return "", "", err
	}

	htmlfp, err := os.Create(path.Join(bm.OutputDir, fullProjName+".html"))
	if err != nil {
		return "", "", err
	}
	infofp, err := os.Open(projName + ".info")
	if err != nil {
		return "", "", err
	}
	dec := json.NewDecoder(infofp)
	var info map[string]interface{}
	if err := dec.Decode(&info); err != nil {
		return "", "", err
	}
	info["URL"] = "/output/" + fullProjName + ".zip"
	seqlen := len(info["padded_sequence"].([]interface{})[0].(string))
	info["SparklineWidth"] = fmt.Sprintf("%d", 10*seqlen)

	if err := tmpl.ExecuteTemplate(htmlfp, "output", info); err != nil {
		return "", "", err
	}

	return "/output/" + fullProjName + ".html", "", nil
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
