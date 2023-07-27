package goprivatemodules

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type Payload struct {
	SvcCode string `json:"svcCode" validate:"required"`
	Env     string `json:"env" validate:"required"`
	Tag     string `json:"tag" validate:"required"`
}

// HandleAPIRequest is a handler function for updating helm repository
func HandleAPIRequest(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only POST requests are allowed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload Payload
	err := decoder.Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid payload")
		return
	}

	// Validate the payload
	err = validator.New().Struct(payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Some Required fields are missing. Please add all required fields...")
		return
	}

	// Process the payload
	fmt.Printf("Received payload: %+v\n", payload)
	// fmt.Printf("Received payload: %+v\n", payload.Env)

	// var env = string(payload.Env)

	// 		GIT PULL ##################################

	fmt.Println("Running command")
	dir := "helm-multiple-branch/"
	gitpullcmd := exec.Command("git", "pull")
	gitpullcmd.Dir = dir
	err = gitpullcmd.Run()
	if err != nil {
		// handle error
		fmt.Println("Error:", err)
		log.Fatalln("Error Pull changes from bitbucket", err) // Use to exit the process and terminate api call
		return
	}
	fmt.Println("Completed command")

	// 		GIT UPDATE ###############################

	// Update the YAML value
	// Update the YAML value
	err = updateYAMLValue("helm-multiple-branch/values.yaml", payload.SvcCode, payload.Tag, payload.Env)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to update YAML value")
		return
	}

	err = gitAddCommitPush("./", payload.SvcCode, payload.Tag, payload.Env)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to perform git add, commit, and push")
		zap.S().Warnw("Failed to perform git add, commit, and push", "err", err)
		return
	}

	err = helmApply("helm-multiple-branch/", payload.SvcCode, payload.Tag, payload.Env)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		zap.S().Errorw("Failed to perform helm apply", "err", err)
		log.Println(err)
		fmt.Fprint(w, "Failed to apply helm")
		return
	}

	// fmt.Println("Running command")
	// dir = "./helm-multiple-branch"
	// cmd = exec.Command("git", "pull")
	// cmd.Dir = dir
	// err = cmd.Run()
	// if err != nil {
	// 	// handle error
	// 	fmt.Println("Error:", err)
	// 	log.Fatalln("Error Pull changes from bitbucket", err) // Use to exit the process and terminate api call
	// 	return
	// }
	// fmt.Println("Completed command")

	// Send a response
	response := map[string]string{"status": "success"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func updateYAMLValue(filePath string, svcCode string, tag string, env string) error {
	updateYAMLValueCmd := exec.Command("sed", "-i", fmt.Sprintf("s/tag:.*/tag: %s/", tag), filePath)
	err := updateYAMLValueCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func gitAddCommitPush(dir string, svcCode string, tag string, env string) error {
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = dir
	err := addCmd.Run()
	if err != nil {
		return err
	}
	commitMsg := fmt.Sprintf("Updated %s %s to %s", env, svcCode, tag)

	commitCmd := exec.Command("git", "commit", "-m", commitMsg)
	commitCmd.Dir = dir
	err = commitCmd.Run()
	if err != nil {
		return err
	}

	pushCmd := exec.Command("git", "push")
	pushCmd.Dir = dir
	err = pushCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func helmApply(dir string, svcCode string, tag string, env string) error {
	updateYAMLValueCmd := exec.Command("helm", "upgrade", "app", ".")
	updateYAMLValueCmd.Dir = dir
	err := updateYAMLValueCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
