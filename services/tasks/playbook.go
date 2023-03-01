package tasks

import (
	"errors"
	"github.com/ansible-semaphore/semaphore/util"
	"io/ioutil"
	"os"
)

const (
	ShellPlaybookTemplate = `---
- hosts: all
  gather_facts: no
  tasks:
    - name: run shell task
      shell: "{{ semaphore_template_command_line }}"
      register: command_result
      ignore_errors: true
    - debug:
        msg: "{{command_result.stdout.splitlines()}}"
`
	CommandPlaybookTemplate = `---
- hosts: all
  gather_facts: no
  tasks:
    - name: run command task
      command: "{{ semaphore_template_command_line }}"
      register: command_result
      ignore_errors: true
    - debug:
        msg: "{{command_result.stdout.splitlines()}}"
`
	WinShellPlaybookTemplate = `---
- hosts: all
  gather_facts: no
  tasks:
    - name: run win_shell task
      win_shell: "{{ semaphore_template_command_line }}"
      register: command_result
      ignore_errors: true
    - debug:
        msg: "{{command_result.stdout.splitlines()}}"
`
	WinCommandPlaybookTemplate = `---
- hosts: all
  gather_facts: no
  tasks:
    - name: run win_command task
      win_command: "{{ semaphore_template_command_line }}"
      register: command_result
      ignore_errors: true
    - debug:
        msg: "{{command_result.stdout.splitlines()}}"
`
	WinPowershellPlaybookTemplate = `---
- hosts: all
  gather_facts: no
  tasks:
    - name: run win_powershell task
      win_powershell: "{{ semaphore_template_command_line }}"
      register: command_result
      ignore_errors: true
    - debug:
        msg: "{{command_result.stdout.splitlines()}}"
`
)

func (t *TaskRunner) installCommandPlaybook() error {
	if t.template.Type != "command" {
		return nil
	}

	t.Log("installing command playbook")

	path := t.GetCommandPlaybookPath()

	_, err := os.Stat(path)

	//有错误且不是文件不存在错误
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	//没有错误，说明文件存在
	if err == nil {
		return nil
	}

	//文件不存在，创建文件
	var content []byte
	switch t.template.Module {
	case "shell":
		content = []byte(ShellPlaybookTemplate)
	case "command":
		content = []byte(CommandPlaybookTemplate)
	case "win_shell":
		content = []byte(WinShellPlaybookTemplate)
	case "win_command":
		content = []byte(WinCommandPlaybookTemplate)
	case "win_powershell":
		content = []byte(WinPowershellPlaybookTemplate)
	default:
		return errors.New("无匹配的任务模板")
	}

	// create inventory file
	return ioutil.WriteFile(path, content, 0664)
}

func (t *TaskRunner) GetCommandPlaybookPath() string {
	return util.Config.TmpPath + "/playbook_" + t.template.Module + ".yml"
}
