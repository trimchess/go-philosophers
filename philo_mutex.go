package main

/**
  Extended version MWI
	Each fork is a mutex, philos have a state 0 (sleeping) or 1 (eating)
	When a philo wakes up, he fetches its forks and starts eating
	After the eat timeout, he releases its forks and goes back sleeping
	(some say thinking or philosophizing)
  Algorithm with counter-meaningful taking of the forks
  The mutex blocking is the synchronisation mechanism
*/

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

//array, each position representing a fork
var forks [5]sync.Mutex

func philo(id int) {
	//We want to log to a file
	f, err := os.OpenFile("philo_mutex.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	// some initialistaions for each philosopher...
	portions := 0
	forkLeft := id
	forkRight := id - 1
	if id == 0 {
		forkRight = 4
	}
	rand.Seed(time.Now().Unix())
	state := 0
	log.Println(",Philosopher,", id, ",initialized")
	//...when done start living!
	for {
		if state == 0 {
			//wakes up from sleeping, waiting for first fork and  fetch it
			//then waiting for second fork and  fetch it too
			//then change to state eating and eat spaghetti for a while
			//even philos first take the left and after the right fork
			//odd philos first take the right and after the left fork
			if id%2 == 0 {
				forks[forkLeft].Lock()
				log.Println(",Philosopher,", id, ",left fetched")
				forks[forkRight].Lock()
				log.Println(",Philosopher,", id, ",right fetched")
			} else {
				forks[forkRight].Lock()
				log.Println(",Philosopher,", id, ",right fetched")
				forks[forkLeft].Lock()
				log.Println(",Philosopher,", id, ",left fetched")
			}
			portions++
			log.Println(",Philosopher,", id, ",eating,", portions)
			state = 1
		} else if state == 1 {
			//enough eaten, release forks and go to sleep
			//do not forget to release both f
			forks[forkLeft].Unlock()
			log.Println(",Philosopher,", id, ",left released")
			forks[forkRight].Unlock()
			log.Println(",Philosopher,", id, ",right released")
			log.Println(",Philosopher,", id, ",sleeping,")
			state = 0
		} else {
			//unexpected
			state = 0
		}
		//wait for state changes (sleeping and eating)
		random := rand.Intn(20) + 1
		t := time.NewTicker(time.Duration(random) * time.Second)
		<-t.C
	}
}

func init() {
	//start your philosopher in conncurence
	for i := 0; i < 5; i++ {
		go philo(i)
	}

}

func main() {
	//as long as the philosophers live
	runtime.Goexit()
}
