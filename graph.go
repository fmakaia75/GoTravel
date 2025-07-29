package main

//type Transfert struct {
//	FromStopID: string,
//	FromRouteID: string,
//	ToStopID: string,
//	ToRouteID: string,
//}

type Stoptime struct{
        //StopID string
        ArrivalTime string
        DepartureTime string
        StopSequence int
}

type Stop struct {
        ID   string
        Name string
        Lat  float64
        Lon  float64
}

type Graph struct {
        Stops map[string]Stop
        Trips  map[string]map[string]Stoptime
        Calendar map[string]Service
        Trip2ServiceMap map[string]string
        //	Transfers map[string]Transfer
}

type Result struct {
        Paths []Path `json:"path"`
}

type Path struct {
        TotalDuration string `json:"total_duration"`
        Edge Edge `json:"edge"`
}

type Edge struct {
        DepartureStop Stop `json:"departure_stop"`
        ArrivalStop Stop `json:"arrival_stop"`
        DepartureTime string `json:"departure_time"`
        ArrivalTime   string `json:"arrival_time"`
        TripID        string `json:"trip_id"`
}
type Service struct {
        Monday bool
        Tuesday bool
        Wednesday bool
        Thursday bool
        Friday bool
        Saturday bool
        Sunday bool
        StartDate string
        EndDate string

}
