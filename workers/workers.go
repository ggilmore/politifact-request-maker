package workers

import (
	"log"
	"time"

	"github.com/ggilmore/politifact-request-maker/scraper"
	"github.com/zabawaba99/firego"
)

type token struct{}

//SendSubjectStatements - routine that endlessly pushes statements to Firebase
func SendSubjectStatements(maxReq int, peopleEndpoint, subjectsEndpoint string, src <-chan scraper.Statement, f *firego.Firebase) {

	//simple semaphore to control the maxmium number of concurrent requests
	//that we can have at any given time
	resources := make(chan token, maxReq)

	//preload all the resource tokens
	for i := 0; i < maxReq; i++ {
		resources <- token{}
	}

	for s := range src {

		for _, sub := range s.Subject {

			//send to ~/person/subject/date=stmt
			for _, p := range s.Target {

				<-resources

				go func(name, sub, date string, s scraper.Statement) {
					err := f.Child(peopleEndpoint).Child(name).Child(sub).Child(date).Set(s)
					if err != nil {
						log.Fatalf("Error when sending statement %v to endpoint: %v/%v/%v/%v \n err: %v",
							s, peopleEndpoint, name, sub, date, err)
					}
					resources <- token{}
				}(p.NameSlug, sub.SubjectSlug, s.RulingDate, s)
			}

			//send to ~/subject/date=stmt
			<-resources
			go func(sub, date string, s scraper.Statement) {
				err := f.Child(subjectsEndpoint).Child(sub).Child(date).Set(s)
				if err != nil {
					log.Fatalf("Error when sending statement %v to endpoint: %v/%v/%v \n  err: %v", s,
						subjectsEndpoint, sub, date, err)
				}
				resources <- token{}
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
