package main

import "fmt"

func main() {
    a := []int{1, 2, 3, 4, 5, 6}
    b := []int{0, 1, 2, 3, 4, 5, 6, 7}

    commonElements := 0

    for _, x := range a {
        for _, y := range b {
            if x == y {
                commonElements++
                fmt.Printf("%d ", x)
                break // Break out of the inner loop once a match is found.
            }
        }
    }

    fmt.Printf("\nTotal %d common elements\n", commonElements)
}

