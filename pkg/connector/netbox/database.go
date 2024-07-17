//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package netbox

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/cloudflare/octopus/pkg/connector/netbox/model"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func init() {
	orm.SetTableNameInflector(func(s string) string {
		return s
	})
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}

type dbParams struct {
	host         string
	port         uint
	dBname       string
	user         string
	password     string
	useTLS       bool
	caCertPath   string
	logDBQueries bool
}

type database struct {
	params dbParams

	pgdb *pg.DB

	contentTypeDcimDevice                 int32
	contentTypeDcimInterface              int32
	contentTypeIpamIpaddress              int32
	contentTypeIpamPrefix                 int32
	contentTypeCircuitsCircuit            int32
	contentTypeCircuitsCircuittermination int32
	contentTypeFrontPort                  int32
	contentTypeRearPort                   int32
}

func newDB(params dbParams) *database {
	return &database{
		params: params,
	}
}

func (db *database) connect() error {
	tlsConfig, err := db.getTLSConfig()
	if err != nil {
		return fmt.Errorf("error building TLS config: %v", err)
	}

	db.pgdb = pg.Connect(&pg.Options{
		Addr:      fmt.Sprintf("%s:%d", db.params.host, db.params.port),
		User:      db.params.user,
		Password:  db.params.password,
		Database:  db.params.dBname,
		TLSConfig: tlsConfig,
	})

	if db.params.logDBQueries {
		db.pgdb.AddQueryHook(dbLogger{})
	}

	return nil
}

func (db *database) getTLSConfig() (*tls.Config, error) {
	if !db.params.useTLS {
		return nil, nil
	}

	cfg := &tls.Config{
		ServerName:         db.params.host,
		InsecureSkipVerify: false,
	}

	if db.params.caCertPath != "" {
		certPool := x509.NewCertPool()
		pem, err := os.ReadFile(db.params.caCertPath)
		if err != nil {
			return nil, fmt.Errorf("error while reading CA cert from %q: %v", db.params.caCertPath, err)
		}

		ok := certPool.AppendCertsFromPEM(pem)
		if !ok {
			return nil, fmt.Errorf("error while add CA cert to pool")
		}

		cfg.RootCAs = certPool
	}

	return cfg, nil
}

func (db *database) getDevices() ([]*model.DcimDevice, error) {
	dcimDevices := make([]*model.DcimDevice, 0)

	err := db.pgdb.Model(&dcimDevices).Relation("DeviceRole").Relation("Site").Relation("DeviceType").Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	tagsByID, err := db.tagsByID(uint(db.contentTypeDcimDevice))
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %v", err)
	}

	for _, dev := range dcimDevices {
		dev.Tags = tagsByID[dev.ID]
	}

	return dcimDevices, nil
}

func (db *database) getInterfaces() (map[int64]*model.DcimInterface, error) {
	dcimInterfaces := make([]*model.DcimInterface, 0)

	err := db.pgdb.Model(&dcimInterfaces).Relation("Parent").Relation("Device").Relation("LAG").Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	tagsByID, err := db.tagsByID(uint(db.contentTypeDcimInterface))
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %v", err)
	}

	for _, ifa := range dcimInterfaces {
		ifa.Tags = tagsByID[ifa.ID]
	}

	res := make(map[int64]*model.DcimInterface)
	for _, ifa := range dcimInterfaces {
		res[ifa.ID] = ifa
	}

	return res, nil
}

func (db *database) getIPAddresses() ([]*model.IpamIpaddress, error) {
	addrs := make([]*model.IpamIpaddress, 0)

	err := db.pgdb.Model(&addrs).Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	return addrs, nil
}

func (db *database) getCables() ([]*model.DcimCable, error) {
	cables := make([]*model.DcimCable, 0)

	err := db.pgdb.Model(&cables).Relation("Terminations").Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	return cables, nil
}

func (db *database) getPrefixes() ([]*model.IpamPrefix, error) {
	prefixes := make([]*model.IpamPrefix, 0)

	err := db.pgdb.Model(&prefixes).Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	tagsByID, err := db.tagsByID(uint(db.contentTypeIpamPrefix))
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %v", err)
	}

	for _, pfx := range prefixes {
		pfx.Tags = tagsByID[pfx.ID]
	}

	return prefixes, nil
}

func (db *database) getCircuits() ([]*model.CircuitsCircuit, error) {
	circuits := make([]*model.CircuitsCircuit, 0)

	err := db.pgdb.Model(&circuits).Relation("Provider").Relation("Type").Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	tagsByID, err := db.tagsByID(uint(db.contentTypeCircuitsCircuit))
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %v", err)
	}

	for _, c := range circuits {
		c.Tags = tagsByID[c.ID]
	}

	return circuits, nil
}

func (db *database) getCircuitTerminations() ([]*model.CircuitsCircuittermination, error) {
	cts := make([]*model.CircuitsCircuittermination, 0)

	err := db.pgdb.Model(&cts).Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	return cts, nil
}

func (db *database) getFrontports() ([]*model.DcimFrontport, error) {
	fps := make([]*model.DcimFrontport, 0)

	err := db.pgdb.Model(&fps).Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	return fps, nil
}

func (db *database) getRearports() ([]*model.DcimRearport, error) {
	rps := make([]*model.DcimRearport, 0)

	err := db.pgdb.Model(&rps).Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	return rps, nil
}

func (db *database) tagsByID(contentTypeID uint) (map[int64][]string, error) {
	tags, err := db.getTags(contentTypeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get tags: %v", err)
	}

	ret := make(map[int64][]string)
	for _, tag := range tags {
		if _, exists := ret[int64(tag.ObjectID)]; !exists {
			ret[int64(tag.ObjectID)] = make([]string, 0, 1)
		}

		ret[int64(tag.ObjectID)] = append(ret[int64(tag.ObjectID)], tag.Tag.Name)
	}

	return ret, nil
}

func (db *database) getTags(contentTypeID uint) ([]model.ExtrasTaggeditem, error) {
	tagsMapping := make([]model.ExtrasTaggeditem, 0)
	err := db.pgdb.Model(&tagsMapping).Relation("Tag").Where("content_type_id = ?", contentTypeID).Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	return tagsMapping, nil
}

func (db *database) getDjangoContentTypes() ([]model.DjangoContentType, error) {
	contentTypes := make([]model.DjangoContentType, 0)

	err := db.pgdb.Model(&contentTypes).Select()
	if err != nil {
		return nil, fmt.Errorf("select failed: %v", err)
	}

	return contentTypes, nil
}

func (db *database) loadContentTypes() error {
	types, err := db.getDjangoContentTypes()
	if err != nil {
		return err
	}

	for _, t := range types {
		switch t.AppLabel {
		case "dcim":
			switch t.Model {
			case "device":
				db.contentTypeDcimDevice = t.ID
			case "interface":
				db.contentTypeDcimInterface = t.ID
			case "frontport":
				db.contentTypeFrontPort = t.ID
			case "rearport":
				db.contentTypeRearPort = t.ID
			}
		case "ipam":
			{
				switch t.Model {
				case "ipaddress":
					db.contentTypeIpamIpaddress = t.ID
				case "prefix":
					db.contentTypeIpamPrefix = t.ID
				}
			}
		case "circuits":
			{
				switch t.Model {
				case "circuit":
					db.contentTypeCircuitsCircuit = t.ID
				case "circuittermination":
					db.contentTypeCircuitsCircuittermination = t.ID
				}
			}
		}
	}

	return nil
}
