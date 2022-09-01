package controllers

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"os"
)

var driver neo4j.Driver

func ConnectDB() {
	uri := os.Getenv("NEO4J_URI")
	auth := neo4j.BasicAuth(os.Getenv("NEO4J_USERNAME"), os.Getenv("NEO4J_PASSWORD"), "")
	// You typically have one driver instance for the entire application. The
	// driver maintains a pool of instance connections to be used by the sessions.
	// The driver is thread safe.
	var err error
	driver, err = neo4j.NewDriver(uri, auth)
	if err != nil {
		panic(err)
	}
	// Don't forget to close the driver connection when you are finished with it
	defer func(driver neo4j.Driver) {
		err := driver.Close()
		if err != nil {
			panic(err)
		}
	}(driver)
	// Create a session to run transactions in. Sessions are lightweight to
	// create and use. Sessions are NOT thread safe.
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	//records, err := session.WriteTransaction(
	//	func(tx neo4j.Transaction) (interface{}, error) {
	//		// To learn more about the Cypher syntax, see https://neo4j.com/docs/cypher-manual/current/
	//		// The Reference Card is also a good resource for keywords https://neo4j.com/docs/cypher-refcard/current/
	//		createRelationshipBetweenPeopleQuery := `
	//			MERGE (p1:Article { name: $article_name, brand: $brand_name, $variant })
	//			MERGE (p2:Tag { name: $person2_name })
	//			MERGE (p1)-[:KNOWS]->(p2)
	//			RETURN p1, p2`
	//		result, err := tx.Run(createRelationshipBetweenPeopleQuery, map[string]interface{}{
	//			"person1_name": "Alice",
	//			"person2_name": "David",
	//		})
	//		if err != nil {
	//			// Return the error received from driver here to indicate rollback,
	//			// the error is analyzed by the driver to determine if it should try again.
	//			return nil, err
	//		}
	//		// Collects all records and commits the transaction (as long as
	//		// Collect doesn't return an error).
	//		// Beware that Collect will buffer the records in memory.
	//		return result.Collect()
	//	})
	//if err != nil {
	//	panic(err)
	//}
	//for _, record := range records.([]*neo4j.Record) {
	//	firstPerson := record.Values[0].(neo4j.Node)
	//	fmt.Printf("First: '%s'\n", firstPerson.Props["name"].(string))
	//	secondPerson := record.Values[1].(neo4j.Node)
	//	fmt.Printf("Second: '%s'\n", secondPerson.Props["name"].(string))
	//}
}
