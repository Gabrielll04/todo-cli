package main

import (
	"os"
	"testing"
)
func TestStorage(t *testing.T) {
	testFileName := "storage_test.json"
	defer os.Remove(testFileName)

	storage := newJSONStorage(testFileName)

	t.Run("Initialize", func(t *testing.T) {
		err := storage.Initialize()
		if err != nil {
			t.Fatalf("Initialize() falhou: %v", err)
		}

		if _, err := os.Stat(storage.filename); os.IsNotExist(err) {
			t.Error("Arquivo não foi criado")
		}
	})

	t.Run("Create", func(t *testing.T) {
		testSchedule := Schedule{
			Title:   "Title",
			Time:    "20:20",
			Details: "AAAAAAAAAAAA",
		}

		err := storage.Create(testSchedule)
		if err != nil {
			t.Fatalf("Create() falhou: %v", err)
		}

		schedules, err := storage.ReadAll()
		if err != nil {
			t.Fatalf("ReadAll() não funcionou")
		}

		if len(schedules) != 1 {
			t.Errorf("Esperado 1 agendamento, obtidos %d", len(schedules))
		}

		if schedules[0].Title != testSchedule.Title {
			t.Errorf("Título incorreto. Esperado: %s, Obtido: %s", testSchedule.Title, schedules[0].Title)
		}
	})

	t.Run("Update", func(t *testing.T) {
		updated := Schedule{
			Title:   "TitleUpdated",
			Time:    "66:66",
			Details: "BBBBBBBUPDATED",
		}

		err := storage.Update(1, updated)
		if err != nil {
			t.Fatalf("Update() falhou: %v", err)
		}

		schedules, err := storage.ReadAll()
		if err != nil {
			t.Fatalf("ReadAll() não funcionou")
		}

		if schedules[0].Title != updated.Title {
			t.Errorf("Update não funcionou")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := storage.Delete(1)
		if err != nil {
			t.Fatalf("Delete() falhou: %v", err)
		}

		schedules, err := storage.ReadAll()
		if err != nil {
			t.Fatalf("ReadAll() não funcionou")
		}

		if len(schedules) != 0 {
			t.Error("Delete não removeu o agendamento")
		}
	})
}

func TestStorageErrors(t *testing.T) {
	t.Run("Arquivo inválido", func(t *testing.T) {
		invalidStorage := newJSONStorage("/invalid/path/storage_test.json")

		err := invalidStorage.Initialize()
		if err != nil {
			t.Error("Esperado erro com caminho inválido")
		}
	})

	t.Run("Schedule não encontrado", func(t *testing.T) {
		testFilename := "test_notfound.json"
		defer os.Remove(testFilename)

		storage := newJSONStorage(testFilename)
		storage.Initialize()

		if err := storage.Update(999, Schedule{}); err != nil {
			t.Error("Esperado erro ao atualizar agendamento inexistente")
		}

		if err := storage.Delete(999); err != nil {
			t.Error("Esperado erro ao deletar agendamento inexistente")
		}
	})
}
