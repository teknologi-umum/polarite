package repository_test

import (
	"context"
	"errors"
	"polarite/repository"
	"polarite/resources"
	"testing"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
)

func TestDependency_GetItemById(t *testing.T) {
	t.Run("NotExists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		_, err := dependency.GetItemById(ctx, "this-shall-not-exists")
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expecting an error of not found, got %s", err.Error())
		}
	})

	t.Run("Normal", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		text := `But I must explain to you how all this mistaken idea of denouncing pleasure and praising pain was born and I will give you a complete account of the system, and expound the actual teachings of the great explorer of the truth, the master-builder of human happiness. No one rejects, dislikes, or avoids pleasure itself, because it is pleasure, but because those who do not know how to pursue pleasure rationally encounter consequences that are extremely painful. Nor again is there anyone who loves or pursues or desires to obtain pain of itself, because it is pain, but because occasionally circumstances occur in which toil and pain can procure him some great pleasure. To take a trivial example, which of us ever undertakes laborious physical exercise, except to obtain some advantage from it? But who has any right to find fault with a man who chooses to enjoy a pleasure that has no annoying consequences, or one who avoids a pain that produces no resultant pleasure? On the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammelled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided.`
		id, err := nanoid.GenerateString(nanoid.DefaultAlphabet, 6)
		if err != nil {
			t.Fatalf("generating id: %s", err.Error())
		}

		hash, err := resources.Hash([]byte(text))
		if err != nil {
			t.Fatalf("generating hash: %s", err.Error())
		}

		paste := repository.Item{
			ID:    id,
			Paste: []byte(text),
			Hash:  hash,
			Metadata: repository.Metadata{
				CreatorIP: "127.0.0.1",
				ExpiresAt: time.Now().Add(time.Minute * 5),
				CreatedAt: time.Now(),
			},
		}

		_, err = dependency.InsertPaste(ctx, paste)
		if err != nil {
			t.Fatalf("inserting paste: %s", err.Error())
		}

		item, err := dependency.GetItemById(ctx, id)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if item.ID != id {
			t.Errorf("expected item.ID to be equal with id: %s", item.ID)
		}

		if string(item.Paste) != text {
			t.Errorf("expected item.Paste to be equal with text, got: %s", string(item.Paste))
		}
	})
}

func TestDependency_ReadHash(t *testing.T) {
	t.Run("NotExists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		exists, id, err := dependency.ReadHash(ctx, "this-shall-not-exists")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if id != "" {
			t.Errorf("id should be empty, got %s", id)
		}

		if exists {
			t.Errorf("expecting not exists, got exists true")
		}
	})

	t.Run("Normal", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		text := `The European languages are members of the same family. Their separate existence is a myth. For science, music, sport, etc, Europe uses the same vocabulary. The languages only differ in their grammar, their pronunciation and their most common words. Everyone realizes why a new common language would be desirable: one could refuse to pay expensive translators. To achieve this, it would be necessary to have uniform grammar, pronunciation and more common words. If several languages coalesce, the grammar of the resulting language is more simple and regular than that of the individual languages. The new common language will be more simple and regular than the existing European languages. It will be as simple as Occidental; in fact, it will be Occidental. To an English person, it will seem like simplified English, as a skeptical Cambridge friend of mine told me what Occidental is. The European languages are members of the same family. Their separate existence is a myth. For science, music, sport, etc, Europe uses the same vocabulary. The languages only differ in their grammar, their pronunciation and their most common words. Everyone realizes why a new common language would be desirable: one could refuse to pay expensive translators. To achieve this, it would be necessary to have uniform grammar, pronunciation and more common words. If several languages coalesce, the grammar of the resulting language is more simple and regular than that of the individual languages. The new common language will be more simple and regular than the existing European languages. It will be as simple as Occidental; in fact, it will be Occidental. To an English person, it will seem like simplified English, as a skeptical Cambridge friend of mine told me what Occidental is.The European languages are members of the same family.`
		id, err := nanoid.GenerateString(nanoid.DefaultAlphabet, 6)
		if err != nil {
			t.Fatalf("generating id: %s", err.Error())
		}

		hash, err := resources.Hash([]byte(text))
		if err != nil {
			t.Fatalf("generating hash: %s", err.Error())
		}

		paste := repository.Item{
			ID:    id,
			Paste: []byte(text),
			Hash:  hash,
			Metadata: repository.Metadata{
				CreatorIP: "127.0.0.1",
				ExpiresAt: time.Now().Add(time.Minute * 5),
				CreatedAt: time.Now(),
			},
		}

		_, err = dependency.InsertPaste(ctx, paste)
		if err != nil {
			t.Fatalf("inserting paste: %s", err.Error())
		}

		exists, itemId, err := dependency.ReadHash(ctx, hash)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if !exists {
			t.Errorf("expecting exists to be true, got %t", exists)
		}

		if itemId != id {
			t.Errorf("expected itemId to be equal with id: %s", itemId)
		}
	})
}
