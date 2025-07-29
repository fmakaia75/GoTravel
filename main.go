package main

import (
        "encoding/json"
        "net/http"
        "fmt"
)
func main() {
    graph, err := LoadGraph("data/stops.txt", "data/stop_times.txt", "data/calendar.txt","data/trips.txt")
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/search_stops", func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query().Get("q")
        matches := graph.SearchStopsByName(query)

        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("["))
        for i, stop := range matches {
            if i > 0 {
                w.Write([]byte(","))
            }
            fmt.Fprintf(w, `{"id":"%s", "name":"%s", "lat":%.6f, "lon":%.6f}`,
                stop.ID, stop.Name, stop.Lat, stop.Lon)
        }
        w.Write([]byte("]"))
    })
    http.HandleFunc("/find_path", func(w http.ResponseWriter, r *http.Request) {
        from := r.URL.Query().Get("from")
        to := r.URL.Query().Get("to")
        departureDate := r.URL.Query().Get("departure_date")
        //result egdes is []Edges for given dep->arr
        resultEdge , err:= graph.SearchResult(from, to)
        if err != nil {
                fmt.Printf("Error during search result : %+v",err)
        }
        //init needed path for the final Result struct
        var paths []Path
        for _, edge := range resultEdge {
                edgeCopy := edge       // Create a fresh copy
                paths = append(paths, NewPath(&edgeCopy)) 
        }
        response := Result{
                Paths: paths,
        }
        //we have a response results but We need to organise the calendar now
        validEdges := response.CalendarFilter(departureDate, graph.Calendar,graph.Trip2ServiceMap)

        filteredPaths := make([]Path, 0, len(validEdges))
        for _, edge := range validEdges {
                edgeCopy := edge
                filteredPaths = append(filteredPaths, NewPath(&edgeCopy))
        }

        finalResponse := Result{
                Paths: filteredPaths,
        }
        
        w.Header().Set("Content-Type", "application/json")
        err = json.NewEncoder(w).Encode(finalResponse)
        if err != nil {
                http.Error(w, "Failed to encode result", http.StatusInternalServerError)
                return
        }
    }) 
    http.Handle("/", http.FileServer(http.Dir("./static")))
    fmt.Printf("Graph loaded: %d stops, %d trips\n", len(graph.Stops), len(graph.Trips))
    fmt.Println("Server started on port 8081...")
    err = http.ListenAndServe(":8081", nil)
    if err != nil {
            fmt.Println("Erreur serveur :", err)
    }
}
