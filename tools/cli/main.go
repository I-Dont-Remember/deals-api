package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/I-Dont-Remember/deals-api/tools/cli/deals"

	"github.com/brianvoe/gofakeit"

	"github.com/tucnak/climax"
)

// use lvh instead of localhost to possibly trick any CORS issues out during local development
var (
	basePath = "http://lvh.me:4500"
)

func main() {
	cli := climax.New("cli")
	cli.Brief = "CLI for tasks for Deals On Tap"
	cli.Version = "v0.0.0"

	generateCmd := climax.Command{
		Name:  "generate",
		Brief: "Generate fake development data in the local DB",
		Usage: "go run main.go generate -c 2 -l 1,2 -d 3,5 -a auth-string",
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
			{
				Name:     "authValue",
				Short:    "a",
				Variable: true,
				Usage:    "-a auth-header-string",
			},
		},

		Handle: generateData,
	}

	uploadCmd := climax.Command{
		Name:  "upload",
		Brief: "Upload deals and locations from files or pass a directory",
		Usage: "go run main.go upload -c uw-madison [ -b http://api.com/api/  ] [ -a auth-string ][ -d directory/ | ./file1.toml ./file2.toml ]",
		Flags: []climax.Flag{
			{
				Name:     "campusSlug",
				Short:    "c",
				Variable: true,
				Usage:    "-c uw-madison",
			},
			{
				Name:     "directory",
				Short:    "d",
				Variable: true,
				Usage:    "-d ./directory/",
			},
			{
				Name:     "basePath",
				Short:    "b",
				Variable: true,
				Usage:    "-b http://localhost:7895",
			},
			{Name: "authValue",
				Short:    "a",
				Variable: true,
				Usage:    "-a auth-header-string",
			},
		},
		Handle: upload,
	}

	validateCmd := climax.Command{
		Name:   "validate",
		Brief:  "Validate a TOML file",
		Usage:  "go run main.go validate file.toml",
		Handle: validate,
	}

	// createCmd := climax.Command{
	// 	Name:  "create",
	// 	Brief: "Create one of the models",
	// 	Usage: "",
	// }

	cli.AddCommand(generateCmd)
	cli.AddCommand(uploadCmd)
	cli.AddCommand(validateCmd)
	cli.Run()
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

	var authValue string
	if value, ok := ctx.Get("authValue"); ok {
		authValue = value
	} else {
		authValue = "local"
	}

	generate(numCampuses, maxLocations, minLocations, maxDeals, minDeals, authValue)
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

/// generate a big pile of fake data for using in the local DB
func generate(numCampuses, maxLocations, minLocations, maxDeals, minDeals int, authValue string) {
	rand.Seed(time.Now().UnixNano())
	// test 1 function MVP just to make progress cuz all cli things suck

	client := deals.New(basePath, authValue)

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

// used to upload the toml files of locations and their deals
func upload(ctx climax.Context) int {
	campusSlug, ok := ctx.Get("campusSlug")
	if !ok {
		ctx.Log("need to know which campus to upload to")
		os.Exit(1)
	}

	newBasePath, ok := ctx.Get("basePath")
	if ok {
		basePath = newBasePath
	}

	var authValue string
	if value, ok := ctx.Get("authValue"); ok {
		authValue = value
	} else {
		authValue = "local"
	}

	client := deals.New(basePath, authValue)

	// create the campus if it doesn't exist
	_, err := client.CreateCampus(campusSlug, campusSlug)
	check(err, "error creating campus for upload")

	dir, ok := ctx.Get("directory")
	if ok {
		ctx.Log("using directory " + dir)
		// get all .toml files from directory
		absPath, err := filepath.Abs(dir)
		if err != nil {
			ctx.Log("couldn't figure out correct path with " + dir)
			os.Exit(1)
		}

		err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			fileErr := handleFile(campusSlug, path, client)
			return fileErr

		})
		if err != nil {
			ctx.Log("error while walking the directory " + err.Error())
			os.Exit(1)
		}
	} else {
		files := ctx.Args

		for _, f := range files {
			handleFile(campusSlug, f, client)
		}
	}

	return 0
}

// types of toml unmarshaling
// !! this needs to stay  up to date with the fields in models
type locationInfo struct {
	Name           string
	CampusSlug     string `toml:"campus"`
	DisplayAddress string `toml:"addr"`
	Latitude       string
	Longitude      string
	ImageLink      string
	PhoneNumber    string
	Website        string
	YelpLink       string
	Deals          []dealInfo `toml:"deal"`
}

type dealInfo struct {
	Desc   string
	Days   []string
	AllDay bool
	Start  string
	End    string
	Types  []string
}

func validate(ctx climax.Context) int {
	files := ctx.Args
	fmt.Println("[*] validating ", files)

	if len(files) == 0 {
		fmt.Println("[!] need files to validate")
		return 1
	}

	for _, f := range files {
		contents := &locationInfo{}
		ret, err := toml.DecodeFile(f, &contents)
		if err != nil {
			fmt.Println("Err decoding " + err.Error())
			return 1
		}

		if len(ret.Undecoded()) != 0 {
			fmt.Printf("[!] unable to decode these correctly:\n\t%v\n\n", ret.Undecoded())
		}

		fmt.Printf("[*] data:\n\n%+v\n", contents)
	}
	return 0
}

func handleFile(campusSlug, path string, client *deals.Client) error {
	if filepath.Ext(strings.TrimSpace(path)) == ".toml" {
		fmt.Println("[*] handling " + path)
	} else {
		fmt.Println("[!] skipping " + path)
		return nil
	}

	contents := &locationInfo{}
	ret, err := toml.DecodeFile(path, &contents)
	if err != nil {
		fmt.Println("Err decoding " + err.Error())
		return nil
	}

	if len(ret.Undecoded()) != 0 {
		fmt.Println("[!] unable to decode these correctly", ret.Undecoded())
	}

	fmt.Println(contents)

	// create a location from information
	location, err := client.CreateLocation(campusSlug, contents.Name,
		contents.DisplayAddress,
		contents.Latitude,
		contents.Longitude,
		contents.ImageLink,
		contents.PhoneNumber,
		contents.Website,
		contents.YelpLink)
	check(err, "failed creating location from toml")

	for _, d := range contents.Deals {
		_, err := client.CreateDeal(location.ID, d.Desc,
			d.Start,
			d.End,
			d.Days,
			d.Types,
			d.AllDay)
		check(err, "failed creating deals from toml "+d.Desc)
	}

	return nil
}
