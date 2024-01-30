package main

import (
	"github.com/charmbracelet/bubbles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	todo status = iota
	inProgress
	done
)


type Task struct {
	title string
	description string
	Status status
}

func (t *Task) FilterValue() string {
	return t.title
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

type Model struct {
	list list.Model
	err error
}

func (m *Model) initList() {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate())
}