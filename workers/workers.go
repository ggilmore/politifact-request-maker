package workers

import (
	"log"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/ggilmore/politifact-request-maker/scraper"
	"github.com/zabawaba99/firego"
)

type Connection struct{}

//SendSubjectStatements - routine that endlessly pushes statements to Firebase
func SendSubjectStatements(c chan Connection, peopleEndpoint, subjectsEndpoint string, src chan scraper.Statement, f *firego.Firebase) {

	for s := range src {

		for _, sub := range s.Subject {

			//send to ~/person/subject/date=stmt
			for _, p := range s.Target {

				<-c

				go func(name, sub, date string, s scraper.Statement) {
					err := f.Child(peopleEndpoint).Child(name).Child(sub).Child(date).Set(s)
					if err != nil {
						log.Printf("Error when sending statement %v to endpoint: %v/%v/%v/%v \n. Retrying... \n err: %v",
							s, peopleEndpoint, name, sub, date, err)
						//try again, probably EOF
						src <- s
					}
					c <- Connection{}
				}(p.NameSlug, sub.SubjectSlug, s.RulingDate, s)
			}

			//send to ~/subject/date=stmt
			<-c
			go func(sub, date string, s scraper.Statement) {
				err := f.Child(subjectsEndpoint).Child(sub).Child(date).Set(s)
				if err != nil {
					log.Printf("Error when sending statement %v to endpoint: %v/%v/%v \n  err: %v", s,
						subjectsEndpoint, sub, date, err)
					//try again, probably EOF
					src <- s
				}
				c <- Connection{}
			}(sub.SubjectSlug, s.RulingDate, s)
		}
	}
}

//DiffStatements - diff old statements w/ new statements, put them on the out channel
//we only want to send updated statements to Firebase
func DiffStatements(in <-chan []scraper.Statement, out chan<- scraper.Statement) {
	var old []scraper.Statement

	for new := range in {
		diff := scraper.DiffStmts(old, new)
		for _, s := range diff {
			out <- s
		}

		old = diff
	}

}

func MakePolitifactRequests(rate, size int, out chan<- []scraper.Statement) {
	limiter := time.NewTicker(time.Second * time.Duration(rate))
	for _ = range limiter.C {
		out <- scraper.StatementsByDate(size)
	}
}

type rssPkg struct {
	date string
	name string
	item *rss.Item
}

func pollRSS(rate int, name, ep string, out chan<- rssPkg) {
	limiter := time.NewTicker(time.Second * time.Duration(rate))
	log.Print(ep)
	feed, err := rss.Fetch(ep)

	var seen = make(map[string]bool)

	if err != nil {
		log.Printf("PollRSS: %v", err)
	}

	for _ = range limiter.C {
		for _, itm := range feed.Items {
			date := itm.Date.String()
			if !seen[date] {
				log.Print("New " + name + "!")
				out <- rssPkg{date, name, itm}
				seen[date] = true
			}
			log.Print("done: " + name)
		}
		err = feed.Update()
		if err != nil {
			log.Printf("PollRSS: %v", err)
		}
	}
}

func sendRSS(c chan Connection, ps chan rssPkg, storyChildName string, f *firego.Firebase) {

	for p := range ps {
		var p = p
		<-c
		go func() {
			log.Printf("sending rss to endpoint: %v/%v/%v/ ", storyChildName, p.name, p.date)
			err := f.Child(storyChildName).Child(p.name).Child(p.date).Set(p.item)
			if err != nil {
				log.Printf("SendRSS: Error when sending item %v \n to endpoint: %v/%v/%v/ \n  err: %v", p.item, storyChildName, p.name, p.date, err)
				ps <- p
			}
			c <- Connection{}
		}()

	}
}

func SetupRSS(rate int, names []string, c chan Connection, storyChildName string, f *firego.Firebase) {
	ps := make(chan rssPkg)

	for _, n := range names {
		go pollRSS(rate, n, scraper.RSSEndpoint+n, ps)
	}

	go sendRSS(c, ps, storyChildName, f)

}
