package repository_test

import (
	"context"
	"polarite/repository"
	"polarite/resources"
	"testing"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
)

func TestDependency_InsertPaste(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	text := `One morning, when Gregor Samsa woke from troubled dreams, he found himself transformed in his bed into a horrible vermin. He lay on his armour-like back, and if he lifted his head a little he could see his brown belly, slightly domed and divided by arches into stiff sections. The bedding was hardly able to cover it and seemed ready to slide off any moment. His many legs, pitifully thin compared with the size of the rest of him, waved about helplessly as he looked. "What's happened to me? " he thought. It wasn't a dream.

		His room, a proper human room although a little too small, lay peacefully between its four familiar walls. A collection of textile samples lay spread out on the table - Samsa was a travelling salesman - and above it there hung a picture that he had recently cut out of an illustrated magazine and housed in a nice, gilded frame. It showed a lady fitted out with a fur hat and fur boa who sat upright, raising a heavy fur muff that covered the whole of her lower arm towards the viewer. Gregor then turned to look out the window at the dull weather.
		
		Drops of rain could be heard hitting the pane, which made him feel quite sad. "How about if I sleep a little bit longer and forget all this nonsense", he thought, but that was something he was unable to do because he was used to sleeping on his right, and in his present state couldn't get into that position. However hard he threw himself onto his right, he always rolled back to where he was. He must have tried it a hundred times, shut his eyes so that he wouldn't have to look at the floundering legs, and only stopped when he began to feel a mild, dull pain there that he had never felt before.
		
		"Oh, God", he thought, "what a strenuous career it is that I've chosen! Travelling day in and day out. Doing business like this takes much more effort than doing your own business at home, and on top of that there's the curse of travelling, worries about making train connections, bad and irregular food, contact with different people all the time so that you can never get to know anyone or become friendly with them. It can all go to Hell! " He felt a slight itch up on his belly; pushed himself slowly up on his back towards the headboard so that he could lift his head better; found where the itch was, and saw that it was covered with lots of little white spots which he didn't know what to make of; and when he tried to feel the place with one of his legs he drew it quickly back because as soon as he touched it he was overcome by a cold shudder. He slid back into his former position. "Getting up early all the time", he thought, "it makes you stupid.`

	hash, err := resources.Hash([]byte(text))
	if err != nil {
		t.Fatalf("hashing text: %s", err.Error())
	}

	id, err := nanoid.GenerateString(nanoid.DefaultAlphabet, 6)
	if err != nil {
		t.Fatalf("generating id: %s", err.Error())
	}

	item, err := dependency.InsertPaste(
		ctx,
		repository.Item{
			ID:    id,
			Paste: []byte(text),
			Hash:  hash,
			Metadata: repository.Metadata{
				CreatorIP: "127.0.0.1",
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Minute),
			},
		},
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if item.ID != id {
		t.Errorf("expected item.ID to be equivalent with id, got: %s:%s", item.ID, id)
	}

	id2, err := nanoid.GenerateString(nanoid.DefaultAlphabet, 6)
	if err != nil {
		t.Fatalf("generating id2: %s", err.Error())
	}

	item2, err := dependency.InsertPaste(
		ctx,
		repository.Item{
			ID:    id2,
			Paste: []byte(text),
			Hash:  hash,
			Metadata: repository.Metadata{
				CreatorIP: "127.0.0.1",
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Minute),
			},
		},
	)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if item2.ID != item.ID {
		t.Errorf("expecting item2.ID to be equivalent with item.ID, got %s:%s", item.ID, item2.ID)
	}
}
