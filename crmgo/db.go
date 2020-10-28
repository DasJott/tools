package crmgo

import (
	"fmt"
	"strings"

	"cleverreach.com/crtools/crconfig"
	"gopkg.in/mgo.v2"
)

type (
	// DBLogger is the interface to implement to be used as logger in here.
	DBLogger interface {
		Debugln(...interface{})
	}
	defaultLog struct{}
)

func (l *defaultLog) Debugln(args ...interface{}) {
	if Debug {
		fmt.Println(args...)
	}
}

var (
	// Debug can be set to true for retrieving debug information
	Debug bool
	// Logger can be set to any logger implementing the DBLogger interface.
	Logger DBLogger = &defaultLog{}
)

// DB represents and encapsulates the mongo db
type DB struct {
	session *mgo.Session
	dbName  string
}

// Multi is a helper struct to open a suffixed connection
type Multi struct {
	suffix string
}

// WithSuffix opens a connection, where you access the configuration of a certain suffix.
func WithSuffix(suffix string) *Multi {
	if suffix != "" {
		suffix = "_" + strings.ToUpper(suffix)
	}
	return &Multi{suffix}
}

// Open opens the mongo connection
func (m *Multi) Open(dbname string) (*DB, error) {
	host := crconfig.Get("MONGO_HOST"+m.suffix, "localhost")
	port := crconfig.Get("MONGO_PORT"+m.suffix, "27017")
	dialup := "mongodb://" + host + ":" + port

	Logger.Debugln("dialup database", dialup)
	sess, err := mgo.Dial(dialup + "/" + dbname)
	if err == nil {
		err = sess.Ping()
	}
	if err == nil {
		sess.SetMode(mgo.Monotonic, true)
		return &DB{
			session: sess,
			dbName:  dbname,
		}, nil
	}

	return nil, err
}

// MustOpen opens the DB Connection and panics on errors
func (m *Multi) MustOpen(dbname string) *DB {
	d, err := m.Open(dbname)
	if d == nil {
		panic("open db failed " + err.Error())
	}
	return d
}

// Open opens the mongo db
func Open(dbname string) (*DB, error) {
	m := &Multi{}
	return m.Open(dbname)
}

// MustOpen opens the DB Connection and panics on errors
func MustOpen(dbname string) *DB {
	m := &Multi{}
	return m.MustOpen(dbname)
}

// Close closes the db session
func (d *DB) Close() {
	d.session.Close()
	Logger.Debugln("Closed database", d.dbName)
}

// Drop drops the database, be careful!
func (d *DB) Drop() error {
	d.session.Refresh()
	ok := d.session.DB(d.dbName).DropDatabase()
	d.Close()
	Logger.Debugln("Dropped database", d.dbName)
	return ok
}

// C gets the mgo Collection to make the queries on.
func (d *DB) C(name string) *mgo.Collection {
	if err := d.session.Ping(); err != nil {
		d.session.Refresh()
		Logger.Debugln("refreshed mongo session")
	}
	return d.session.DB(d.dbName).C(name)
}
