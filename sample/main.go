package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var workingHours = 10 * time.Second
var arrivalRate = 100

func main() {
	// seed random channel
	rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// print out our welcome message
	color.Green("======Opening for the day========")

	// create channels
	doneChan := make(chan bool)
	clientChan := make(chan string, 10)

	// create a shop
	shop := BarbingSalon{
		NumberOfSeats:   10,
		NumberOfBarbers: 0,
		HairCutDuration: time.Second * 2,
		DoneChan:        doneChan,
		ClientChan:      clientChan,
		Open:            true,
	}

	// add barbers
	shop.addBarber("James")

	shopClosing := make(chan bool)
	closed := make(chan bool)

	// open shop
	go func() {
		<-time.After(workingHours)
		shopClosing <- true
		shop.close()
		closed <- true
	}()
	// create clients
	i := 1
	go func() {
		for {
			randomArrivalRate := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomArrivalRate)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	// finish for the day

}

type BarbingSalon struct {
	NumberOfSeats   int
	NumberOfBarbers int
	HairCutDuration time.Duration
	DoneChan        chan bool
	ClientChan      chan string
	Open            bool
}

func (shop *BarbingSalon) addBarber(barber string) {
	shop.NumberOfBarbers++
	color.Yellow("%s has resumed duty...")
	isSleeping := false

	go func() {
		for {
			if len(shop.ClientChan) == 0 {
				color.Yellow("No clients, %s takes a nap")
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientChan

			if shopOpen {
				// cut hair
				if isSleeping {
					color.Yellow("%s wakes up %s.", client, barber)
					isSleeping = false
				}
				shop.cutHair(barber, client)
			} else {
				// send barber home
				color.Cyan("%s is going home.", barber)
				shop.DoneChan <- true
				return
			}
		}
	}()
}

func (shop *BarbingSalon) cutHair(barber, client string) {
	color.Magenta("%s is cutting %s's hair", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.GreenString("%s is finished cutting %s's hair", barber, client)
}

func (shop *BarbingSalon) close() {
	color.Magenta("We're closing shop")
	close(shop.ClientChan)
	for a := 0; a < shop.NumberOfBarbers; a++ {
		<-shop.DoneChan
	}
	shop.Open = false
	close(shop.DoneChan)
	color.Green("-------------------------------------------------------")
	color.Green("Closed for the day, all the barbers have left the shop.")
}
func (shop *BarbingSalon) addClient(client string) {
	// print message
	color.Green("***** %s arrives!", client)
	if shop.Open {
		select {
		case shop.ClientChan <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves", client)
		}
	} else {
		color.Red("The shop is already closed, so %s leaves")
	}
}