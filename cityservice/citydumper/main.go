package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	elastigo "github.com/dimroc/elastigo/lib"
	citylib "github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"os"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "citydumper"
	app.Version = "0.0.1"
	app.Usage = "Manage a city's Elasticsearch Geoevent store. See subcommands (help dump)."
	app.Commands = []cli.Command{
		addDumpCommand(),
		addImportCommand(),
		addPruneCommand(),
	}

	app.Run(os.Args)
}

func addPruneCommand() cli.Command {
	var elasticsearchUrl string
	command := cli.Command{
		Name:  "prune",
		Usage: "prune <datestring>. Delete geoevents that are before a certain date",
	}

	command.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "http://localhost:9200",
			Usage:       "The URL of the Elasticsearch server to write to",
			EnvVar:      "ELASTICSEARCH_URL",
			Destination: &elasticsearchUrl,
		},
	}

	command.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			fmt.Println("Must pass date string. See help.")
			return
		}

		before := c.Args()[0]
		println("Deleting all Geoevents before " + before)

		elastic := citylib.NewBulkElasticConnection(elasticsearchUrl)
		defer elastic.Close()

		// Set up Elasticsearch reading
		dsl := elastigo.Search(citylib.ES_IndexName).Type(citylib.ES_TypeName).Size("1000").
			Pretty().Filter(
			elastigo.Filter().Range("createdAt", nil, nil, before, nil, ""))

		elastic.ScanAndScrollGeoEvents(dsl, func(geoevent citylib.GeoEvent) {
			elastic.BulkIndexer.Delete(citylib.ES_IndexName, citylib.ES_TypeName, geoevent.Id)
		})
	}

	return command
}

func addImportCommand() cli.Command {
	var elasticsearchUrl string
	command := cli.Command{
		Name:    "import",
		Aliases: []string{"i"},
		Usage:   "import <flags> <filename> Import geoevents from a Cityservice JSONL file.",
	}

	command.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "http://localhost:9200",
			Usage:       "The URL of the Elasticsearch server to write to",
			EnvVar:      "ELASTICSEARCH_URL",
			Destination: &elasticsearchUrl,
		},
	}

	command.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			fmt.Println("Must pass input filename as argument. See help.")
			return
		}

		elastic := citylib.NewBulkElasticConnection(elasticsearchUrl)
		//elastic.SetRequestTracer(RequestTracer)

		filename := c.Args()[0]
		file, err := os.Open(filename)
		Check(err)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			rawMessage := json.RawMessage(scanner.Bytes())
			geoevent := citylib.GeoEventFromRawMessage(&rawMessage)
			geoevent = correctMediaOrigin(geoevent)
			elastic.Write(geoevent)
		}

		if err := scanner.Err(); err != nil {
			Logger.Warning("reading standard input:", err)
		}
	}

	return command
}

func correctMediaOrigin(g citylib.GeoEvent) citylib.GeoEvent {
	if len(g.MediaOrigin) > 0 {
		return g
	}

	if len(g.MediaUrl) > 0 && strings.Contains(g.MediaUrl, "cdninstagram") {
		g.MediaOrigin = "instagram"
	} else {
		g.MediaOrigin = "twitter"
	}

	return g
}

func addDumpCommand() cli.Command {
	var citykey, after, before, neighborhoods, elasticsearchUrl, index string
	command := cli.Command{
		Name:    "dump",
		Aliases: []string{"d"},
		Usage:   "dump <flags> <filename> Export geoevents from a Cityservice Elasticsearch store to JSONL.",
	}

	command.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "citykey, c",
			Usage:       "Required. The key for the city you are trying to dump (aka export)",
			Destination: &citykey,
			EnvVar:      "CITYKEY",
		},
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "http://localhost:9200",
			Usage:       "The URL of the Elasticsearch server to read from",
			EnvVar:      "ELASTICSEARCH_URL",
			Destination: &elasticsearchUrl,
		},
		cli.StringFlag{
			Name:        "after, a",
			Value:       "1980-01-01",
			Usage:       "The date string which geoevents must be after",
			Destination: &after,
		},
		cli.StringFlag{
			Name:        "before, b",
			Value:       time.Now().Format(time.RFC3339),
			Usage:       "The date string which geoevents must be before",
			Destination: &before,
		},
		cli.StringFlag{
			Name:        "neighborhoods, n",
			Usage:       "A comma separated string specifying neighborhoods to filter for",
			Destination: &neighborhoods,
		},
		cli.StringFlag{
			Name:        "index, i",
			Value:       citylib.ES_IndexName,
			Usage:       "The specific index to run against.",
			Destination: &index,
		},
	}

	command.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			fmt.Println("Must pass output filename as argument. See help.")
			return
		}

		filename := c.Args()[0]

		if len(citykey) == 0 {
			fmt.Println("Must set citykey flag. See help")
			return
		}

		fmt.Println("Dumping city " + citykey + " from " + elasticsearchUrl)
		elastic := citylib.NewElasticConnection(elasticsearchUrl)
		elastic.SetRequestTracer(RequestTracer)
		outputfile, err := os.Create(filename)
		Check(err)
		defer outputfile.Close()
		defer elastic.Close()

		// Set up Output File
		writer := bufio.NewWriter(outputfile)
		defer writer.Flush()

		filters := []*elastigo.FilterOp{
			elastigo.Filter().Term("city", citykey),
			elastigo.Filter().Range("createdAt", after, nil, before, nil, ""),
		}

		neighborhoodArray := removeEmptyEntries(strings.Split(neighborhoods, ","))
		if len(neighborhoodArray) > 0 {
			filters = append(filters, elastigo.Filter().Terms("neighborhoods", "", neighborhoodArray))
		}

		// Set up Elasticsearch reading
		dsl := elastigo.Search(index).Type(citylib.ES_TypeName).Size("1000").
			Pretty().Filter(elastigo.Filter().And(filters...))

		elastic.ScanAndScrollGeoEvents(dsl, func(geoevent citylib.GeoEvent) {
			writeGeoevent(writer, geoevent)
		})

		fmt.Println("Output written to " + outputfile.Name())
	}

	return command
}

func writeGeoevent(writer *bufio.Writer, geoevent citylib.GeoEvent) {
	_, err := writer.WriteString(ToJsonStringUnsafe(geoevent) + "\n")
	Check(err)
}

func removeEmptyEntries(array []string) []string {
	cleaned := []string{}
	for _, entry := range array {
		if len(entry) > 0 {
			cleaned = append(cleaned, entry)
		}
	}

	return cleaned
}
