# Course Setup Verification

## Goal

Verify that your development environment is ready for this curriculum. You will confirm that your terminal, Git, and a text editor are installed and configured, and establish your personal development workspace directory for course projects.

## Prerequisites

- Lessons 01 through 10 of Module 00 (Orientation) completed
- A terminal application open (Bash, Zsh, PowerShell, etc.)
- Git installed on your system
- A text editor or IDE installed (VS Code, GoLand, Neovim, etc.)

## Tasks

Execute the following environment verification steps in your terminal:

1. **Terminal and Shell**: Verify that your terminal shell is working and inspect your current path:
   ```bash
   pwd    # On Windows PowerShell: Get-Location
   ```
2. **Git Installation & Identity**: Confirm Git is installed and that your committer identity is configured:
   ```bash
   git --version
   git config user.name
   git config user.email
   ```
   If your name or email returns blank, configure them:
   ```bash
   git config --global user.name "Your Name"
   git config --global user.email "your.email@example.com"
   ```
3. **Text Editor**: Verify your editor from the command line (e.g., `code --version` for VS Code, or confirm your IDE is launched).
4. **Initialize Your Learning Workspace**: Create a dedicated directory outside the curriculum repository where you will store your personal code, experiment scripts, and project solutions:
   ```bash
   mkdir -p ~/the-go-engineer-workspace
   ```
5. **Record Your Setup Evidence**: In your learning workspace directory, create an environment summary file named `setup.txt`:
   ```bash
   git --version > setup.txt
   git config user.name >> setup.txt
   git config user.email >> setup.txt
   ```

## Verification

Review your `setup.txt` file and verify:
- Git reports a valid version (e.g., `git version 2.43.0` or higher).
- Your committer name and email are recorded.
- Your personal workspace directory exists and is ready for project work.

When you are done, note this milestone in your personal portfolio plan from Lesson 10.

## NEXT UP

Advance to the setup assessment:

- [Assessment: Course Setup](../../assessments/course-setup/README.md)
