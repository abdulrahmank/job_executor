package dao

import (
	"os"
	"testing"
)

func TestSaveFile(t *testing.T) {
	fileName := "helloworld.sh"

	t.Run("Should save file", func(t *testing.T) {
		contentString := "echo 'hello world'"

		impl := FileDaoImpl{}
		impl.SaveFile(fileName, []byte(contentString))

		file, _ := os.Open(fileName)
		defer file.Close()
		if file == nil {
			t.Error("Expected file to be present")
			t.Fail()
		}
		content := make([]byte, 18)
		if _, e := file.Read(content); e != nil {
			t.Error("Error reading file")
			t.Fail()
		}
		if string(content) != contentString {
			t.Errorf("expected '%s' but was '%s'", contentString, string(content))
			t.Fail()
		}
	})

	_ = os.Remove(fileName)
}
