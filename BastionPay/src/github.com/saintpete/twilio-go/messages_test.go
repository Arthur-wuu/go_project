package twilio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := "SM26b3b00f8def53be77c5697183bfe95e"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msg, err := envClient.Messages.Get(ctx, sid)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Sid != sid {
		t.Errorf("expected Sid to equal %s, got %s", sid, msg.Sid)
	}
}

func TestGetPage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	page, err := envClient.Messages.GetPage(ctx, url.Values{"PageSize": []string{"5"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Messages) != 5 {
		t.Fatalf("expected len(messages) to be 5, got %d", len(page.Messages))
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	client, s := getServer(sendMessageResponse)
	defer s.Close()
	msg, err := client.Messages.SendMessage(from, to, "twilio-go testing!", nil)
	if err != nil {
		t.Fatal(err)
	}
	if msg.From != from {
		t.Errorf("expected From to be from, got error")
	}
	if msg.Body != "twilio-go testing!" {
		t.Errorf("expected Body to be twilio-go testing, got %s", msg.Body)
	}
	if msg.NumSegments != 1 {
		t.Errorf("expected NumSegments to be 1, got %d", msg.NumSegments)
	}
	if msg.Status != StatusQueued {
		t.Errorf("expected Status to be StatusQueued, got %s", msg.Status)
	}
}

func TestGetMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msg, err := envClient.Messages.Get(ctx, "SM5a52bc49b2354703bfdea7e92b44b385")
	if err != nil {
		t.Fatal(err)
	}
	if msg.ErrorCode != CodeUnknownDestination {
		t.Errorf("expected Code to be %d, got %d", CodeUnknownDestination, msg.ErrorCode)
	}
	if msg.ErrorMessage == "" {
		t.Errorf(`expected ErrorMessage to be non-empty, got ""`)
	}
}

func TestIterateAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	iter := envClient.Messages.GetPageIterator(url.Values{"PageSize": []string{"500"}})
	count := 0
	start := uint(0)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	for {
		page, err := iter.Next(ctx)
		if err == NoMoreResults {
			break
		}
		if count > 0 && (page.Start <= start || page.Start-start > 500) {
			t.Fatalf("expected page.Start to be greater than previous, got %d, previous %d", page.Start, start)
			return
		} else {
			start = page.Start
		}
		if err != nil {
			t.Fatal(err)
			break
		}
		count++
		if count > 15 {
			fmt.Println("count > 15")
			t.Fail()
			break
		}
	}
	if count < 10 {
		t.Errorf("Too small of a count - expected at least 10, got %d", count)
	}
}

func TestGetMediaURLs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}
	t.Parallel()
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	urls, err := envClient.Messages.GetMediaURLs(ctx, "MM89a8c4a6891c53054e9cd604922bfb61", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 1 {
		t.Errorf("Wrong number of URLs returned: %d", len(urls))
	}
	if !strings.HasPrefix(urls[0].String(), "https://s3-external-1.amazonaws.com/media.twiliocdn.com/"+sid) {
		t.Errorf("wrong url: %s", urls[0].String())
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()
	msg := new(Message)
	if err := json.Unmarshal(getMessageResponse, &msg); err != nil {
		t.Fatal(err)
	}
	if msg.Sid != "SM26b3b00f8def53be77c5697183bfe95e" {
		t.Errorf("wrong sid")
	}
	got := msg.DateCreated.Time.Format(time.RFC3339)
	want := "2016-09-20T22:59:57Z"
	if got != want {
		t.Errorf("msg.DateCreated: got %s, want %s", got, want)
	}
	if msg.Direction != DirectionOutboundReply {
		t.Errorf("wrong direction")
	}
	if msg.Status != StatusDelivered {
		t.Errorf("wrong status")
	}
	if msg.Body != "Welcome to ZomboCom." {
		t.Errorf("wrong body")
	}
	if msg.From != PhoneNumber("+19253920364") {
		t.Errorf("wrong from")
	}
	if msg.FriendlyPrice() != "$0.0075" {
		t.Errorf("wrong friendly price %v, want %v", msg.FriendlyPrice(), "$0.00750")
	}
}

func TestStatusFriendly(t *testing.T) {
	t.Parallel()
	if StatusQueued.Friendly() != "Queued" {
		t.Errorf("expected StatusQueued.Friendly to equal Queued, got %s", StatusQueued.Friendly())
	}
	s := Status("in-progress")
	if f := s.Friendly(); f != "In Progress" {
		t.Errorf("expected In Progress.Friendly to equal In Progress, got %s", f)
	}
}
