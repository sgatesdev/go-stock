package stream

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"samgates.io/go-stock/models"
	"samgates.io/go-stock/utils"
)

var (
	StockToConnectionMap StockToConnectionMapType
	SessionIds           SessionIdsType
	db                   *gorm.DB
)

func init() {
	StockToConnectionMap = StockToConnectionMapType{
		stockToConnection: make(map[string][]ConnNameToWSConn),
	}

	SessionIds = SessionIdsType{
		sessionIds: make([]string, 0),
	}

	db = utils.SetupDB()
}

// map of connection UUID to connection UUID:websocket connection
type StockToConnectionMapType struct {
	stockToConnection map[string][]ConnNameToWSConn
	mutex             sync.Mutex
}

type ConnNameToWSConn struct {
	sessionId string
	conn      *websocket.Conn
}

type SessionIdsType struct {
	sessionIds []string
	mutex      sync.Mutex
}

// crud for connections/map to stocks

func AddConnection(conn *websocket.Conn) string {
	id := GetSessionId()
	SessionIds.mutex.Lock()
	SessionIds.sessionIds = append(SessionIds.sessionIds, id)
	SessionIds.mutex.Unlock()

	// send connection id to client
	msg := fmt.Sprintf("Connected to /ws/prices. Connection id: %s", id)
	conn.WriteMessage(websocket.TextMessage, []byte(msg))

	// get all stocks, stream data
	// TODO: implement new col in stocks table to indicate whether or not to stream
	stocks, err := getStocks()
	if err != nil {
		fmt.Println(err)
	}

	StockToConnectionMap.mutex.Lock()
	for _, s := range stocks {
		if StockToConnectionMap.stockToConnection[s.ID] == nil {
			StockToConnectionMap.stockToConnection[s.ID] = make([]ConnNameToWSConn, 0)
		}
		StockToConnectionMap.stockToConnection[s.ID] = append(StockToConnectionMap.stockToConnection[s.ID], ConnNameToWSConn{
			sessionId: id,
			conn:      conn,
		})
	}
	StockToConnectionMap.mutex.Unlock()

	fmt.Println(StockToConnectionMap.stockToConnection)

	return id
}

func RemoveConnection(sessionId string) {
	// remove connection from map
	deleteStocks := make([]string, 0)
	StockToConnectionMap.mutex.Lock()

	// remove connections from stocks
	for stockId, _ := range StockToConnectionMap.stockToConnection {
		connections := make([]ConnNameToWSConn, 0)
		for _, c := range StockToConnectionMap.stockToConnection[stockId] {
			if c.sessionId != sessionId {
				connections = append(connections, c)
			}
		}

		StockToConnectionMap.stockToConnection[stockId] = connections

		if len(connections) == 0 {
			deleteStocks = append(deleteStocks, stockId)
		}
	}

	// clean up empty stocks with no connections
	for _, s := range deleteStocks {
		delete(StockToConnectionMap.stockToConnection, s)
	}
	StockToConnectionMap.mutex.Unlock()

	// remove connection id from session ids
	SessionIds.mutex.Lock()
	for i, id := range SessionIds.sessionIds {
		if id == sessionId {
			SessionIds.sessionIds = append(SessionIds.sessionIds[:i], SessionIds.sessionIds[i+1:]...)
		}
	}
	SessionIds.mutex.Unlock()

	fmt.Println(StockToConnectionMap.stockToConnection)
}

func SendPriceUpdate(p models.Price) {
	StockToConnectionMap.mutex.Lock()
	for _, c := range StockToConnectionMap.stockToConnection[p.StockID] {
		// check if connection is still open
		if c.conn != nil {
			fmt.Println(time.Now().String()+" SENDING UPDATE FOR ", p.StockID, " TO ", c.sessionId)
			c.conn.WriteJSON(p)
		}
	}
	StockToConnectionMap.mutex.Unlock()
}

func GetSessionId() string {
	return fmt.Sprintf("%s", uuid.New())
}

func getStocks() ([]models.Stock, error) {
	stocks := []models.Stock{}
	err := db.Table("stocks").Find(&stocks).Error
	return stocks, err
}
