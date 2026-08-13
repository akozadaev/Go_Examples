package easyjson_example

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	benchmarkJSON  []byte
	benchmarkUser  User
	benchmarkPlain plainUser
)

// plainUser и plainProfile намеренно не имеют методов MarshalJSON и
// UnmarshalJSON. Иначе encoding/json обнаружит методы, сгенерированные
// easyjson для User, и benchmark фактически сравнит easyjson с самим собой.
type plainUser struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"created_at"`
	Tags      []string      `json:"tags"`
	Profile   *plainProfile `json:"profile,omitempty"`
}

type plainProfile struct {
	Bio      string            `json:"bio"`
	Settings map[string]string `json:"settings"`
	IsActive bool              `json:"is_active"`
	Score    float64           `json:"score"`
}

func newBenchmarkUser() *User {
	return &User{
		ID:        1,
		Name:      "Иван Иванов",
		Email:     "ivan@example.com",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:      []string{"developer", "golang"},
		Profile: &Profile{
			Bio:      "Go разработчик",
			Settings: map[string]string{"theme": "dark", "lang": "ru"},
			IsActive: true,
			Score:    95.5,
		},
	}
}

func newBenchmarkPlainUser() *plainUser {
	return &plainUser{
		ID:        1,
		Name:      "Иван Иванов",
		Email:     "ivan@example.com",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:      []string{"developer", "golang"},
		Profile: &plainProfile{
			Bio:      "Go разработчик",
			Settings: map[string]string{"theme": "dark", "lang": "ru"},
			IsActive: true,
			Score:    95.5,
		},
	}
}

func BenchmarkEasyJSONMarshal(b *testing.B) {
	user := newBenchmarkUser()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		data, err := user.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkJSON = data
	}
}

func BenchmarkStandardJSONMarshal(b *testing.B) {
	user := newBenchmarkPlainUser()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		data, err := json.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkJSON = data
	}
}

func BenchmarkEasyJSONUnmarshal(b *testing.B) {
	data := []byte(`{"id":1,"name":"Иван Иванов","email":"ivan@example.com","created_at":"2024-01-01T00:00:00Z","tags":["developer","golang"],"profile":{"bio":"Go разработчик","settings":{"theme":"dark","lang":"ru"},"is_active":true,"score":95.5}}`)
	var user User
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := user.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
	}
	benchmarkUser = user
}

func BenchmarkStandardJSONUnmarshal(b *testing.B) {
	data := []byte(`{"id":1,"name":"Иван Иванов","email":"ivan@example.com","created_at":"2024-01-01T00:00:00Z","tags":["developer","golang"],"profile":{"bio":"Go разработчик","settings":{"theme":"dark","lang":"ru"},"is_active":true,"score":95.5}}`)
	var user plainUser
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := json.Unmarshal(data, &user); err != nil {
			b.Fatal(err)
		}
	}
	benchmarkPlain = user
}

func TestJSONImplementationsAreEquivalent(t *testing.T) {
	easyData, err := newBenchmarkUser().MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	standardData, err := json.Marshal(newBenchmarkPlainUser())
	if err != nil {
		t.Fatal(err)
	}

	var easyValue, standardValue any
	if err := json.Unmarshal(easyData, &easyValue); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(standardData, &standardValue); err != nil {
		t.Fatal(err)
	}
	if !jsonValuesEqual(easyValue, standardValue) {
		t.Fatalf("different JSON values:\neasyjson: %s\nstandard: %s", easyData, standardData)
	}
}

func jsonValuesEqual(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && string(leftJSON) == string(rightJSON)
}
