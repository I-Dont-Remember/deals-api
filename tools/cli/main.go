package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"

	"github.com/I-Dont-Remember/deals-api/tools/cli/deals"
	"github.com/tucnak/climax"
)

func main() {
	d := climax.New("d")
	d.Brief = "CLI for tasks for Deals On Tap"
	d.Version = "v0.0.0"

	generateCmd := climax.Command{
		Name:  "generate",
		Brief: "Generate fake development data in the local DB",
		Usage: "go run main.go generate -c 2 -l 1,2 -d 3,5",
		Help: "Make sure to start up both the local API server as well as " +
			"Localstack & run the table creation script. Seems like this " +
			"should be automated.",

		Flags: []climax.Flag{
			{
				Name:     "numCampuses",
				Short:    "c",
				Variable: true,
				Usage:    "-c 5",
			},
			{
				Name:     "locationRange",
				Short:    "l",
				Variable: true,
				Usage:    "-l 3,10",
			},
			{
				Name:     "dealRange",
				Short:    "d",
				Variable: true,
				Usage:    "-d 8,15",
			},
		},

		Handle: generateData,
	}

	d.AddCommand(generateCmd)
	d.Run()
}

func check(e error, msg string) {
	if e != nil {
		fmt.Println(msg)
		fmt.Println(e)
		os.Exit(1)
	}
}

// returns the min and max
func parseRange(r string) (int, int, error) {
	parts := strings.Split(r, ",")

	if len(parts) != 2 {
		return -1, -1, errors.New("unable to parse range from " + r)
	}

	min, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return -1, -1, errors.New("unable to parse range from " + r)
	}

	max, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return -1, -1, errors.New("unable to parse range from " + r)
	}

	return int(min), int(max), nil
}

func generateData(ctx climax.Context) int {
	numCampuses := 3
	locationRange := "3,10"
	dealRange := "5,15"

	//generate()
	if strNum, ok := ctx.Get("numCampuses"); ok {
		num, err := strconv.ParseInt(strNum, 10, 32)
		if err != nil {
			ctx.Log("failed to parse number of campuses")
			return 1
		}
		numCampuses = int(num)
	}
	if lrange, ok := ctx.Get("locationRange"); ok {
		locationRange = lrange
	}
	if drange, ok := ctx.Get("dealRange"); ok {
		dealRange = drange
	}
	fmt.Printf("Making %d campuses. Location range: %s Deal range: %s\n", numCampuses, locationRange, dealRange)

	minLocations, maxLocations, err := parseRange(locationRange)
	if err != nil {
		ctx.Log("failed parsing location range")
		return 1
	}

	minDeals, maxDeals, err := parseRange(dealRange)
	if err != nil {
		ctx.Log("failed parsing deal range")
		return 1
	}

	generate(numCampuses, maxLocations, minLocations, maxDeals, minDeals)
	return 0
}

func getRandomFromRange(min, max int) int {
	numRange := max - min
	// picks [0,n)
	num := rand.Intn(numRange + 1)
	// bring it back to match our min & max
	return num + min
}

func temp_convertStructAddrToString(a *gofakeit.AddressInfo) string {
	return fmt.Sprintf("%s %s, %s, %s\n", a.Address, a.Street, a.City, a.State)
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

// func getRandomDays() []string {
// 	num := getRandomFromRange(1, 7)

// }

func getRandomTypes() []string {
	options := []string{"Drinks", "Food", "Event"}
	num := getRandomFromRange(1, 3)

	return []string{options[num-1]}
}

func generate(numCampuses, maxLocations, minLocations, maxDeals, minDeals int) {
	rand.Seed(time.Now().UnixNano())
	// test 1 function MVP just to make progress cuz all cli things suck
	basePath := "http://localhost:4500"

	client := deals.New(basePath)

	for i := 0; i < numCampuses; i++ {
		slug := gofakeit.Username()
		name := gofakeit.JobTitle()
		campus, err := client.CreateCampus(slug, name)
		check(err, "createCampus failed")

		campusSlug := campus.Slug
		fmt.Println("[*] creating campus ", campusSlug)

		numLocations := getRandomFromRange(minLocations, maxLocations)

		fmt.Printf("  %d locations\n", numLocations)
		for j := 0; j < numLocations; j++ {
			name := gofakeit.Name()
			fmt.Println("    * creating location ", name)
			addr := temp_convertStructAddrToString(gofakeit.Address())
			lat := floatToString(gofakeit.Latitude())
			long := floatToString(gofakeit.Longitude())
			imageLink := gofakeit.URL()
			phoneNumber := gofakeit.Phone()
			website := gofakeit.URL()
			yelpLink := gofakeit.URL()
			location, err := client.CreateLocation(campusSlug, name, addr, lat, long, imageLink, phoneNumber, website, yelpLink)
			check(err, "createLocation failed")

			locationID := location.ID

			numDeals := getRandomFromRange(minDeals, maxDeals)
			fmt.Printf("      - %d deals\n", numDeals)
			for k := 0; k < numDeals; k++ {
				desc := gofakeit.HipsterSentence(8)
				start := "start"
				end := "end"
				days := []string{gofakeit.WeekDay()}
				types := getRandomTypes()
				allDay := false
				_, err := client.CreateDeal(locationID, desc, start, end, days, types, allDay)
				check(err, "createDeal failed")
			}
		}
	}

	fmt.Printf("Created %d campuses\n", numCampuses)
}
