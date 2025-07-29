package main

import (
	"encoding/csv"
	"os"
    "strings"
	"strconv"
)

func LoadGraph(stopsFile string,  stopTimesFile string, serviceFile string, tripFile string) (*Graph, error) {
    stops, err := LoadStops(stopsFile)
    if err != nil {
        return nil, err
    }

    trips, err := LoadTrips(stopTimesFile)
    if err != nil {
        return nil, err
    }
        
    services, err := LoadServices(serviceFile)
    if err != nil {
        return nil, err
    }

    trip2ServiceMap, err := LoadMap(tripFile)
    if err != nil {
        return nil, err
    }
    //edges, err := LoadEdges(stopTimesFile)
    //if err != nil {
    //    return nil, err
    //}

    graph := &Graph{
        Stops: stops,
        Trips: trips,
        Calendar: services,
        Trip2ServiceMap: trip2ServiceMap,
    }

    return graph, nil
}
func LoadMap(filename string) (map[string]string,error){
        
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1
    
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    tripMap := make(map[string]string)
    for i, record := range records{
            if i == 0 {
                    continue
            }
            //trip_id->idx 1
            //service_id->idx 2
            tripMap[record[1]] = record[2]

    }
    return tripMap,nil

}
func LoadServices(filename string) (map[string]Service, error){
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1

    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    services := make(map[string]Service)
    for i, record := range records {
            if i == 0 {
                    continue
            }
            services[record[0]] = Service{
                Monday: record[1] != "0",
                Tuesday:record[2] != "0",
                Wednesday: record[3] != "0",
                Thursday: record[4] != "0",
                Friday: record[5] != "0",
                Saturday: record[6] != "0",
                Sunday: record[7] != "0",
                StartDate: record[8],
                EndDate: record[9],
            }

    }
    return services, nil
}
func LoadStops(filename string) (map[string]Stop, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1

    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    stops := make(map[string]Stop)

    for i, record := range records {
        if i == 0 {
            continue
        }
        lat, _ := strconv.ParseFloat(record[2], 64)
        lon, _ := strconv.ParseFloat(record[3], 64)

        stops[record[0]] = Stop{
            ID:   record[0],
            Name: record[1],
            Lat:  lat,
            Lon:  lon,
        }
    }

    return stops, nil
}

func LoadTrips(filename string) (map[string]map[string]Stoptime, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1

    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    tripStopTimes := make(map[string]map[string]Stoptime)

    for i, record := range records {
        if i == 0 {
            continue
        }

        tripID := strings.TrimSpace(record[0])
        stopID := strings.TrimSpace(record[1])
        arrivalTimeStr := strings.TrimSpace(record[2])
        departureTimeStr := strings.TrimSpace(record[3])
        stopSequenceStr := strings.TrimSpace(record[5])
        stopSequence, _ := strconv.Atoi(stopSequenceStr)
        stoptime := Stoptime{
                ArrivalTime: arrivalTimeStr,
                DepartureTime: departureTimeStr,
                StopSequence: stopSequence,
        }
        //If this trip exist already
        if _, exists := tripStopTimes[tripID]; !exists{
                //initialize the map with make
                tripStopTimes[tripID] = make(map[string]Stoptime)
        }
        tripStopTimes[tripID][stopID] = stoptime
    }
    return tripStopTimes, nil
}
