// search.go - Complete file with station grouping and route finding
package main

import (
	"fmt"
	"strings"
    "time"
)

func (graph *Graph) SearchResult(departure string, destination string) ([]Edge, error) {
        //given a departure/ arrival stop look for route first        

        dep_stop, dep_ok := graph.Stops[departure] 
        arr_stop, arr_ok := graph.Stops[destination] 
        //Check if those stops exists 

        if !arr_ok ||!dep_ok {
                return nil, fmt.Errorf("Critical error one stop does not exist") 
        }

        fmt.Printf("Stops does exist: %s ;  %s\n", dep_stop.Name, arr_stop.Name)
        fmt.Printf("Stops does exist: %s ;  %s\n", dep_stop.ID, arr_stop.ID)
        
        //Check if there is a direct route from dep_stop.id to arr_stop.id in Edge 
        routes, err:= graph.SearchDirectRouteById(dep_stop.ID, arr_stop.ID)
        if err != nil {
                return nil, err  // retourne l'erreur directement
        }

        fmt.Printf("Total of %d match(s) found\n",len(routes))
        return routes, nil
}

func (graph *Graph) SearchStopsByName(query string) []Stop {
    query = strings.ToLower(query)
    var matches []Stop

    for _, stop := range graph.Stops {
        if strings.Contains(strings.ToLower(stop.Name), query) {
            matches = append(matches, stop)
        }
    }

    return matches
}

func (graph *Graph) SearchDirectRouteById(fromId string, toId string) ([]Edge, error) {
    var edges []Edge
        
    for tripID, stops := range graph.Trips {
            from, okFrom := stops[fromId]
            to, okTo := stops[toId]
            if okFrom && okTo {
                    //fmt.Printf("Found both stops in trip %s\n", tripID)
                    if from.StopSequence < to.StopSequence {
                            //get stop name from graph too 
                            edges = append(edges, Edge{
                                DepartureStop: graph.Stops[fromId],
                                ArrivalStop: graph.Stops[toId],
                                DepartureTime: from.DepartureTime,
                                ArrivalTime: to.DepartureTime,
                                TripID: tripID,
                            })
                    }
            } else {
                    //fmt.Printf("Stops not both in trip %s: from id=%v, to id=%v\n", tripID, from, to)
            }
    }

    return edges, nil
}
func (r *Result) CalendarFilter(departureDate string, calendar map[string]Service, tripToService map[string]string) []Edge {
    fmt.Printf("This is how the date has been entered: %s\n", departureDate)

    // Parse date
    date, err := time.Parse("2006-01-02", departureDate)
    if err != nil {
        fmt.Printf("Invalid date format: %v\n", err)
        return []Edge{}
    }

    weekday := date.Weekday() // Sunday == 0
    validEdges := []Edge{}
    for _, path := range r.Paths {
        tripID := path.Edge.TripID
        serviceID, ok := tripToService[tripID]
        if !ok {
            fmt.Printf("No service_id for trip %s\n", tripID)
            continue
        }

        service, ok := calendar[serviceID]
        if !ok {
            fmt.Printf("No service found for service_id %s\n", serviceID)
            continue
        }

        // Check if date is within start and end
        start, err1 := time.Parse("20060102", service.StartDate)
        end, err2 := time.Parse("20060102", service.EndDate)
        if err1 != nil || err2 != nil {
            fmt.Printf("Bad start/end date for service %s\n", serviceID)
            continue
        }

        if date.Before(start) || date.After(end) {
            continue
        }

        // Check if weekday is allowed
        if isServiceRunningOn(service, weekday) {
            validEdges = append(validEdges,path.Edge)
        }
    }

    // No match
    return validEdges
}

// Helper to check weekday
func isServiceRunningOn(service Service, day time.Weekday) bool {
    switch day {
    case time.Monday:
        return service.Monday
    case time.Tuesday:
        return service.Tuesday
    case time.Wednesday:
        return service.Wednesday
    case time.Thursday:
        return service.Thursday
    case time.Friday:
        return service.Friday
    case time.Saturday:
        return service.Saturday
    case time.Sunday:
        return service.Sunday
    default:
        return false
    }
}

