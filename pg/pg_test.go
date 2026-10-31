package pg_test

import (
	"slices"
	"testing"

	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/pg/pgtest"
)

func TestListSort(t *testing.T) {
	db := pgtest.Open(t)
	type fruit struct{ Name string }
	if err := db.Exec("CREATE TABLE fruits (name text); INSERT INTO fruits VALUES ('apple'), ('cherry'), ('banana'), ('date')").Error; err != nil {
		t.Fatal(err)
	}
	sort := pg.Sort(map[string]string{"name": "name"}, "name", true)
	r, err := pg.List[fruit](db.Scopes(sort), pg.Page{Offset: 1, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if want := []fruit{{"cherry"}, {"banana"}}; r.Count != 4 || !slices.Equal(r.Items, want) {
		t.Errorf("List = %d %v, want 4 %v", r.Count, r.Items, want)
	}
}

func TestZeroID(t *testing.T) {
	db := pgtest.Open(t)
	type ref struct {
		ID  int
		Ref pg.ID
	}
	if err := db.Exec("CREATE TABLE refs (id int PRIMARY KEY, ref uuid)").Error; err != nil {
		t.Fatal(err)
	}
	want := []ref{{ID: 1}, {ID: 2, Ref: pg.NewID()}}
	if err := db.Create(&want).Error; err != nil {
		t.Fatal(err)
	}
	var null bool
	if err := db.Raw("SELECT ref IS NULL FROM refs WHERE id = 1").Scan(&null).Error; err != nil || !null {
		t.Fatalf("zero ID stored as non-NULL (err %v)", err)
	}
	var got []ref
	if err := db.Order("id").Find(&got).Error; err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("read back %v, want %v", got, want)
	}
}

func TestEnumCreate(t *testing.T) {
	db := pgtest.Open(t)
	type status string
	labels := func() []status {
		var got []status
		if err := db.Raw("SELECT enumlabel FROM pg_enum WHERE enumtypid = 'status'::regtype ORDER BY enumsortorder").Scan(&got).Error; err != nil {
			t.Fatal(err)
		}
		return got
	}

	if err := (pg.Enum[status]{Name: "status", Values: []status{"draft", "live"}}).Create(db); err != nil {
		t.Fatal(err)
	}
	if got := labels(); !slices.Equal(got, []status{"draft", "live"}) {
		t.Errorf("after create: %v", got)
	}

	grown := pg.Enum[status]{Name: "status", Values: []status{"draft", "live", "isn't"}}
	for range 2 {
		if err := grown.Create(db); err != nil {
			t.Fatal(err)
		}
	}
	if got := labels(); !slices.Equal(got, grown.Values) {
		t.Errorf("after growing: %v", got)
	}
}
