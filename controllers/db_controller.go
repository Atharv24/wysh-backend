package controllers

import (
	"context"
	"github.com/mindstand/gogm/v2"
	"wysh-app/models"
)

func ConnectDB() {
	// define your configuration
	config := gogm.Config{
		Host: "0.0.0.0",
		Port: 7687, // deprecated in favor of protocol
		// IsCluster:                 false,
		Protocol: "neo4j", //also supports neo4j+s, neo4j+ssc, bolt, bolt+s and bolt+ssc
		// Specify CA Public Key when using +ssc or +s
		Username:           "neo4j",
		Password:           "passpass",
		PoolSize:           200,
		IndexStrategy:      gogm.VALIDATE_INDEX,       //other options are ASSERT_INDEX and IGNORE_INDEX
		TargetDbs:          []string{"wysh"},          // default logger wraps the go "log" package, implement the Logger interface from gogm to use your own logger
		Logger:             gogm.GetDefaultLogger(),   // define the log level
		LogLevel:           "DEBUG",                   // enable neo4j go driver to log
		EnableDriverLogs:   false,                     // enable gogm to log params in cypher queries. WARNING THIS IS A SECURITY RISK! Only use this when debugging
		EnableLogParams:    false,                     // enable open tracing. Ensure contexts have spans already. GoGM does not make root spans, only child spans
		OpentracingEnabled: false,                     // specify the method gogm will use to generate Load queries
		LoadStrategy:       gogm.SCHEMA_LOAD_STRATEGY, // set to SCHEMA_LOAD_STRATEGY for schema-aware queries which may reduce load on the database
	}

	// register all vertices and edges
	// this is so that GoGM doesn't have to do reflect processing of each edge in real time
	// use nil or gogm.DefaultPrimaryKeyStrategy if you only want graph ids
	// we are using the default key strategy since our vertices are using BaseNode
	_gogm, err := gogm.New(&config, gogm.DefaultPrimaryKeyStrategy, &models.Article{}, &models.Variation{}, &models.Brand{}, &models.VariantEdge{}, &models.Size{}, &models.Color{}, &models.Store{}, &models.Tag{})

	if err != nil {
		panic(err)
	}

	gogm.SetGlobalGogm(_gogm)
}

func getArticleByStoreArticleCode(storeArticleCode string, storeName string) *models.Article {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeRead})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)

	var article models.Article
	query := "MATCH (a:Article) -[r:HAS_VARIANT]-> (:Variation) -[:SOLD_AT] -> (s:Store{name:$storeName})" +
		"WHERE r.variant_store_id STARTS WITH $storeArticleCode RETURN DISTINCT a"
	err = sess.Query(context.Background(), query, map[string]interface{}{
		"storeArticleCode": storeArticleCode,
		"storeName":        storeName,
	}, &article)

	if err != nil {
		return nil
	}

	return &article
}

// TODO: Add getTrendsFunction

func getArticleDetailByVariationId(articleVariationId int) *models.ArticleMini {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)
	a, _, err := sess.QueryRaw(context.Background(),
		"MATCH (a:Article)-[r:HAS_VARIANT]->(v)-->(s), (a)-->(b:Brand) where ID(r)=$id RETURN distinct ID(r) as id, a.name as name, b.name as brand, r.current_price as current_price, r.base_price as base_price, r.article_url as article_url, a.default_image_url as image_url",
		map[string]interface{}{
			"id": articleVariationId,
		},
	)
	row := a[0]
	articleDetail := models.ArticleMini{
		ID:           row[0].(int64),
		Name:         row[1].(string),
		Brand:        row[2].(string),
		CurrentPrice: row[3].(int64),
		BasePrice:    row[4].(int64),
		ArticleUrl:   row[5].(string),
		ImageUrl:     row[6].(string),
	}

	if err != nil {
		return nil
	}
	return &articleDetail
}

func insertArticle(article *models.Article) {
	// Create a session to run transactions in. Sessions are lightweight to
	// create and use. Sessions are NOT thread safe.
	//param is readonly, we're going to make stuff, so we're going to do read write
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)
	err = sess.Save(context.Background(), article)
	if err != nil {
		for i := 0; i < dbRetries; i++ {
			err = sess.Save(context.Background(), article)
			if err == nil {
				break
			}
		}
		if err != nil {
			panic(err)
		}
	}
}
