package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseTimeToSeconds(t string) int {
    parts := strings.Split(t, ":")
    h, _ := strconv.Atoi(parts[0])
    m, _ := strconv.Atoi(parts[1])
    s, _ := strconv.Atoi(parts[2])
    return h*3600 + m*60 + s
}
func SecondsToTimeString(seconds int) string {
    hours := (seconds / 3600)
    minutes := (seconds / 60) % 60
    return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func calculateDuration(departure_time string, arrival_time string) string{
        arrivalTime := ParseTimeToSeconds(arrival_time)
        departureTime := ParseTimeToSeconds(departure_time)
        total := arrivalTime-departureTime
        
        return SecondsToTimeString(total)
}
func NewPath(e *Edge) (Path){
        duration := calculateDuration(e.DepartureTime,e.ArrivalTime)
        return Path{
                TotalDuration: duration,
                Edge: *e,
        }
}
