package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
    // Set default log level for logrus with default debug
    lvl, ok := os.LookupEnv("LOG_LEVEL")
    if !ok {
        lvl = "debug"
    }
    ll, err := log.ParseLevel(lvl)
    if err != nil {
        ll = log.DebugLevel
    }
    log.SetLevel(ll)

    // var config *pgxpool.Config
    // config, err = pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
    // if err != nil {
    //     fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
    //     os.Exit(1)
    // }
    // ctx := context.Background()
    // db.DB, err = pgxpool.ConnectConfig(ctx, config)
}

func testInsert() error {
	ctx := context.Background()
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback(ctx)

	name := "test6"
	nowtime := time.Now()
	// durationtime := time.Duration(500000)
	durationtime := time.Duration(5000000000)
	// sql := "insert into test1(name, duration, ts) values ('test2', '00:00:00.005', NOW())"
	sqlp := `insert into test1(
		name,
		duration,
		ts)
		values ($1, $2, $3)`
	_, err = tx.Exec( ctx, sqlp,
		name,
		durationtime,
		nowtime)
	if err != nil {
		log.Error(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	msg := fmt.Sprintf("Insert into test1 of name = %v.\n", name)
	log.Debug(msg)
	return err
}

func main() {
	for {
		e := testInsert()
		if e != nil {
			// fmt.Printf("E not nil e: %v\n", e)
		}
		time.Sleep(10 * time.Second)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var duration pgtype.Interval
	var ts pgtype.Time
	err = conn.QueryRow(context.Background(), "select name, duration, ts from test1 where name=$1", "test1").Scan(&name, &duration, &ts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("row: ", name, duration, ts)

	// err = conn.exec(context.Background(), sql)
	// insert into test1(name, duration, ts) values ('test1', '00:00:00.005', NOW());

	// loc, _ := time.LoadLocation("Asia/Shanghai")
	//set timezone,
	// now := time.Now().In(loc)

	// var timeTz pgtype.Timestamptz
	//timeTz = pgtype.Timestamptz{Status: pgtype.Null}
	// timeTz, _ = time.Parse("00:00:00.000005", "00:00:00.000005")
	const layout = "00:00:00.000005"
	timeTz, _ := time.Parse(layout, "00:00:00.000005")
	fmt.Printf("timeTz: %v\n", timeTz)
	// pgTime = pgtype.Timestamptz{Time: timeTz, Status: pgtype.Null}
	var pgTime pgtype.Timestamptz
	pgTime = pgtype.Timestamptz{Time: timeTz}
	fmt.Printf("pgTime: %v\n", pgTime)
	// timeTz = pgtype.Timestamptz{Time: "00:00:00.000005", Status: pgtype.Null}
	var tnow time.Time
	tnow = time.Now()
	fmt.Printf("tnow: %v\n", tnow)
	os.Exit(1)
	var v int32
	v = 1
	//var pv pgtype.Int2
	// pv = pgtype.Int2(1)
	// var pv int16
	// pv := pgtype.Int2(1)
	// var t time.Duration
	// var t pgtype.Interval
	// t = "00:00:00.1"
	t := 1
	fmt.Printf("t: %v\n", t)
	// pt := pgtype.Interval(t/1000)
	// pt := t

	// var pt pgtype.Interval
	// pt = pgtype.Interval("40:00:00")
	pt := 1

	a, b := time.Parse("2006-01-02 3:04PM", "1970-01-01 9:00PM")
	fmt.Printf("a %v\n", a)
	fmt.Printf("b %v\n", b)
	timeVal := pgtype.Timestamptz{Status: pgtype.Null}
	// timeVal := pgtype.Interval{Microseconds:10000, Status: pgtype.Null}
	// timeVal := pgtype.Interval{Microseconds:10000}
	// timeDur := time.Time(timeVal)
	timeDur := timeVal
	fmt.Printf("timeVal %v\n", timeVal)
	fmt.Printf("timeDur %v\n", timeDur)

	tt := time.Now()
	// tunixmicro := time.Now().UnixMicro()
	tunixmicro := tt.UnixMicro()
	// j := pgtype.Interval{Microseconds:tt}
	// j := pgtype.Time{Microseconds:tt}
	// j := pgtype.Time(tt)
	var j time.Time
	j = tt
	fmt.Printf("j: %v\n", j)
	var jt pgtype.Timestamptz
	// jt = pgtype.Timestamptz{Time: j, Status: pgtype.Null}
	jt = pgtype.Timestamptz{Time: j, Status: 1}
	// jt = pgtype.Timestamptz{Time: j}
	fmt.Printf("pt: %v\n", jt)
	fmt.Printf("tunixmicro: %v\n", tunixmicro)
	os.Exit(1)
	// tt := time.Now().UnixNano()
	// pgtype.Timestamp(tt)
	fmt.Println(tt)
	// fmt.Println(t.String())
	// fmt.Println(t.Format("2006-01-02 15:04:05"))
	// fmt.Printf("timeVal %v\n", time.Duration(timeVal))

	// *dst = pgtype.Interval{Microseconds: value / 1000, Status: Present}
	// pt = pgtype.Interval{Microseconds: t / 1000, Status: Present}
	// pt = pgtype.Interval{Microseconds: t / 1000, Status: "Present"}
	fmt.Printf("pt: %v\n", pt)

	pv := 1
	fmt.Printf("v: %v pv: %v\n", v, pv)
	// "2006-01-02 15:04:05.999999999Z07:00:00"

	var n int = 97
	m := int64(n)
	fmt.Printf("m: %v\n", m)

	var badboys int = 1921

	// explicit type conversion
	var badboys2 float64 = float64(badboys)

	var badboys3 int64 = int64(badboys)

	var badboys4 uint = uint(badboys)
	fmt.Printf("v: %v pv: %v, %v", badboys2, badboys3, badboys4)

}
