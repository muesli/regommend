/*
 * Simple recommendation engine
 *     Copyright (c) 2014, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package regommend

import (
	_ "errors"
	_ "log"
	"fmt"
	"math"
)

func cosineSim(t1, t2 map[interface{}]float64) float64 {
	sum_xy := 0.0
	sum_x2 := 0.0
	sum_y2 := 0.0

	for key, x := range t1 {
		y, ok := t2[key]
		if ok {
			fmt.Println("Found shared:", key, x, y)

			sum_xy += x * y
			sum_x2 += math.Pow(x, 2)
			sum_y2 += math.Pow(y, 2)
		}
	}

	denominator := math.Sqrt(sum_x2) * math.Sqrt(sum_y2)
	if denominator == 0 {
		return 0
	}

	return sum_xy / denominator
}
