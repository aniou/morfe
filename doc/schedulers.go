
	/* 
	 * first version - adjust loop per 1 second

	go func() {
		var time_before     time.Time
		var time_wait       time.Duration = 300 * time.Microsecond
		var cycles          uint32
		var cycles_per_sec  uint32

		time_before = time.Now()
		for {
			cycles = 0
			for cycles < 14318 {
				cycles+=p.CPU0.Execute()
			}
			cycles_per_sec+=cycles

			if time.Since(time_before) > 1*time.Second {
				fmt.Printf("> second %d wait %d\n",  cycles_per_sec, time_wait)
				if cycles_per_sec < 14318000 {
					time_wait = time_wait -  4*time.Microsecond
				}
				if cycles_per_sec > 14400000 {
					time_wait = time_wait +  2*time.Microsecond
				}
				cycles_per_sec=0
				time_before = time.Now()
			}
			time.Sleep(time_wait)


		}
	}()
	*/


	/*
	// second version of adaptive speed loop - CPU.Execute is called every 1/5 of second
	// with two threshold counters - ie sleep timer is changed for every two times,
	// when cpu is too low or too fast
	//
	// if number of cycles is greater than desired number of cycles+4% 
	//    then trigger counter (thresh_max) is decreased 
	//    if trigger counter == 0 then wait loop is increased by 2 microseconds
	//
	//  the same behaviour is if number of cpu cycles is lower than...
	//
	// XXX - just testing, change static values
	go func() {
		var time_before     time.Time
		var time_wait       time.Duration = 300 * time.Microsecond
		var cycles          uint32
		var all_cycles      uint32
		var thresh_min	    byte  = 2
		var thresh_max      byte  = 2
		var low_thresh      uint32 = (14318000 - (14318000/25)) / 5
		var top_thresh      uint32 = (14318000 + (14318000/25)) / 5

		time_before = time.Now()
		for {
			cycles = 0
			for cycles < 14318 {
				cycles+=p.CPU0.Execute()
			}
			all_cycles+=cycles

			if time.Since(time_before) > 200*time.Millisecond {
				//fmt.Printf("cpu0> low_thresh %d cycles %d top_thresh %d cycles*5 %d wait %d\n", 
				//               low_thresh, top_thresh, all_cycles, all_cycles*5, time_wait)
				all_cycles=0
				time_before = time.Now()

				if all_cycles < low_thresh {
					thresh_min-=1
					if thresh_min == 0 {
						time_wait = time_wait - 4*time.Microsecond
						thresh_min = 2
					}
				} else {
					thresh_min=2
				}
				if all_cycles > top_thresh {
					thresh_max-=1
					if thresh_max == 0 {
						time_wait = time_wait + 2*time.Microsecond
						thresh_max = 2
					}
				} else {
					thresh_max = 2
				}

			}
			time.Sleep(time_wait)


		}
	}()
	*/

	/*
	go func() {
		for {
			p.CPU1.Execute()
		}
	}()
	*/
