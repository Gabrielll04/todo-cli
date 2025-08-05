package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Storage interface {
	Initialize() error
	Create(schedule Schedule) error
	ReadAll() ([]Schedule, error)
	Update(id int, schedule Schedule) error
	Delete(id int) error 
}

type JSONStorage struct {
	filename string
}

type Schedule struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Time    string `json:"time"`
	Details string `json:"details,omitempty"`
}

var ErrScheduleNotFound = errors.New("agendamento não encontrado")
var ErrCannotReadID = errors.New("não foi possível ler ID")

func newJSONStorage(filename string) *JSONStorage {
	return &JSONStorage{filename: filename}
} 

func main() {
	storage := newJSONStorage("schedules.json")

	if err := storage.Initialize(); err != nil {
		fmt.Printf("Erro ao inicializar arquivo: %v\n", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Uso: scheduler [comando]")
		fmt.Println("Comandos disponíveis:")
		fmt.Println("  add - Adiciona novo agendamento")
		fmt.Println("  list - Lista todos os agendamentos")
		fmt.Println("  delete - Remove um agendamento")
		fmt.Println("  edit - Edita um agendamento")
		return
	}

	cmd := os.Args[1]
	switch cmd {

	case "add":
		if len(os.Args) < 5 {
			fmt.Fprintln(os.Stderr, "Uso: add <titulo> <hora> <detalhes>")
			os.Exit(1)
		}

		var newSchedule = Schedule{
			Title:   os.Args[2],
			Time:    os.Args[3],
			Details: os.Args[4],
		}

		if err := storage.Create(newSchedule); err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao adicionar: %v\n", err)
            os.Exit(1)
		}

	case "list":
		schedules, err := storage.ReadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Falha ao listar agendamentos: %v\n", err)
			os.Exit(1)
		}

		for _, s := range schedules {
			fmt.Printf("%d: %s (%s)\n", s.ID, s.Title, s.Time)
		}

	case "delete":
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Falha ao ler ID: %v\n", err)
			os.Exit(1)
		}

		if err := storage.Delete(id); err != nil {
			fmt.Fprintf(os.Stderr, "Falha ao deletar agendamento: %v", err)
			os.Exit(1)
		}

	case "edit":
		if len(os.Args) < 6 {
			fmt.Fprintln(os.Stderr, "Uso: add <ID> <titulo> <hora> <detalhes>")
			os.Exit(1)
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Falha ao ler ID: %v\n", err)
			os.Exit(1)
		}


		var newSchedule = Schedule{
			Title:   os.Args[3],
			Time:    os.Args[4],
			Details: os.Args[5],
		}

		if err := storage.Update(id, newSchedule); err != nil {
			fmt.Fprintf(os.Stderr, "Falha ao editar agendamento: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Comando desconhecido")
	}
}

func (js *JSONStorage) Initialize() error {
	if _, err := os.Stat(js.filename); os.IsNotExist(err) {
		return os.WriteFile(js.filename, []byte("[]"), 0644)
	}

	return nil
}

func (js *JSONStorage) ReadAll() ([]Schedule, error) {
	data, err := os.ReadFile(js.filename)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler agendamento: %w", err)
	}

	var schedules []Schedule
	if err := json.Unmarshal(data, &schedules); err != nil {
		return nil, fmt.Errorf("f    alha ao serializar agendamento: %w", err)
	}

	return schedules, err
}

func (js *JSONStorage) SaveAll(schedules []Schedule) error {
	data, err := json.MarshalIndent(schedules, "", " ")
	if err != nil {
		return fmt.Errorf("falha ao serializar agendamentos: %w", err)
	}
	if err := os.WriteFile(js.filename, data, 0644); err != nil {
		return fmt.Errorf("falha ao criara arquivo: %w", err)
	}

	return nil
}


func (js *JSONStorage) Create(newSchedule Schedule) error {
	schedules, err := js.ReadAll()
	if err != nil {
		return fmt.Errorf("falha ao ler agendamento: %w", err)
	}

	newSchedule.ID = 1
	if len(schedules) > 0 {
		newSchedule.ID = schedules[len(schedules)-1].ID + 1
	}

	schedules = append(schedules, newSchedule)

	if err := js.SaveAll(schedules); err != nil {
		return fmt.Errorf("falha ao salvar alterações: %w", err)
	}

	return nil
}

func (js *JSONStorage) Delete(scheduleID int) error {
	data, err := js.ReadAll()
	if err != nil {
		return fmt.Errorf("falha ao ler agendamento: %w", err)
	}

	found := false
	for i, v := range data {
		if v.ID == scheduleID {
			found = true
			data = append(data[:i], data[i+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("%w: ID %d", ErrScheduleNotFound, scheduleID)
	}

	if err := js.SaveAll(data); err != nil {
		return fmt.Errorf("falha ao salvar alterações: %w", err)
	}

	return nil
}

func (js *JSONStorage) Update(scheduleID int, updatedSchedule Schedule) error {
	data, err := js.ReadAll()
	if err != nil {
		return fmt.Errorf("falha ao ler agendamento: %w", err)
	}

	found := false
	for i, v := range data {
		if v.ID == scheduleID {
			found = true
			data[i] = updatedSchedule
			data[i].ID = scheduleID
			break
		}
	}

	if !found {
		return fmt.Errorf("%w: ID %d", ErrScheduleNotFound, scheduleID)
	}

	if err := js.SaveAll(data); err != nil {
		return fmt.Errorf("falha ao salvar alterações: %w", err)
	}

	return nil
}