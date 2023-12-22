package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientChan      chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barber)

		for {
			// if there are no clients, the barber goes to sleep
			if len(shop.ClientChan) == 0 {
				color.Yellow("There's nothing to do, so %s takes a nap", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes up %s.", client, barber)
					isSleeping = false
				}
				// cut hair
				shop.cutHair(barber, client)
			} else {
				// shop is closed, so send the barber home
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is done cutting %s's hair", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShop() {
	color.Cyan("Closing shop for the day")
	close(shop.ClientChan)
	shop.Open = false

	for a := 0; a < shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)
	color.Green("-------------------------------------------------------")
	color.Green("Closed for the day, all the barbers have left the shop.")
}

func (shop BarberShop) addClient(client string) {
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
