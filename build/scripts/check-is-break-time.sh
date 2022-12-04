#!/usr/bin/env bash

break_times=("12:40:00-13:59:59" "18:30:00-18:59:59" "20:30:00-09:59:59")

function in_time_range()
{
    cur_ts=$(date -j "+%s")
    time_range=$1
    start_day_time=$(echo $time_range | awk -F '-' '{print $1}')
    end_day_time=$(echo $time_range | awk -F '-' '{print $2}')
    start_time=$(date -j "+%Y-%m-%d "$start_day_time)
    end_time=$(date -j "+%Y-%m-%d "$end_day_time)
    start_ts=$(date -j -f "%Y-%m-%d %H:%M:%S" "$start_time" "+%s")
    end_ts=$(date -j -f "%Y-%m-%d %H:%M:%S" "$end_time" "+%s")
    if [ $start_ts -lt $end_ts ]; then
        if [ $cur_ts -gt $start_ts -a $cur_ts -lt $end_ts ]; then
            return 1
        fi
    fi
    if [ $start_ts -gt $end_ts ]; then
        if [ $cur_ts -lt $end_ts -o $cur_ts -gt $start_ts ]; then
            return 1
        fi
    fi
    return 0
}

weekday=$(date -j "+%w")
if [ $weekday -eq 0 -o $weekday -eq 6 ]; then
  exit 0
fi

for break_time in ${break_times[@]}; do
    in_time_range $break_time
    if [ $? -eq 1 ]; then
      exit 0
    fi
done

echo "please commit in break time!"

exit 1