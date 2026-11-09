package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type Task struct {
	Name     string
	Priority int
	Done     bool
}

const taskTmpl = `{{range .}}[{{if .Done}}X{{else}} {{end}}] {{.Name}} {{.Priority | label}}
{{end}}`

func label(priority int) string {
	switch priority {
	case 1:
		return "[HIGH]"
	case 2:
		return "[MED]"
	case 3:
		return "[LOW]"
	default:
		return "[UNKNOWN]"
	}
}

func RenderTasks(tasks []Task) (string, error) {
	tmpl := template.Must(template.New("tasks").Funcs(template.FuncMap{
		"label": label,
	}).Parse(taskTmpl))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tasks); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func main() {
	tasks := []Task{
		{Name: "Implement login", Priority: 1, Done: true},
		{Name: "Write tests", Priority: 2, Done: false},
		{Name: "Fix bug #42", Priority: 3, Done: false},
	}

	output, err := RenderTasks(tasks)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print(output)
}
