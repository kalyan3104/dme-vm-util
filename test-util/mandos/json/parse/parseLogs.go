package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/model"
	oj "github.com/kalyan3104/dme-vm-util/test-util/orderedjson"
)

func (p *Parser) processLogList(logsRaw oj.OJsonObject) ([]*mj.LogEntry, error) {
	logList, isList := logsRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("unmarshalled logs list is not a list")
	}
	var logEntries []*mj.LogEntry
	var err error
	for _, logRaw := range logList.AsList() {
		logMap, isMap := logRaw.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("unmarshalled log entry is not a map")
		}
		logEntry := mj.LogEntry{}
		for _, kvp := range logMap.OrderedKV {
			switch kvp.Key {
			case "address":
				accountStr, err := p.parseString(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("unmarshalled log entry address is not a json string: %w", err)
				}
				logEntry.Address, err = p.parseAccountAddress(accountStr)
				if err != nil {
					return nil, err
				}
			case "identifier":
				strVal, err := p.parseString(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid log identifier: %w", err)
				}
				var identifierValue []byte
				identifierValue, err = p.parseAnyValueAsByteArray(strVal)
				if err != nil {
					return nil, fmt.Errorf("invalid log identifier: %w", err)
				}
				if len(identifierValue) != 32 {
					return nil, fmt.Errorf("invalid log identifier - should be 32 bytes in length")
				}
				logEntry.Identifier = mj.JSONBytes{Value: identifierValue, Original: strVal}
			case "topics":
				logEntry.Topics, err = p.parseByteArrayList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("unmarshalled log entry topics is not big int list: %w", err)
				}
			case "data":
				logEntry.Data, err = p.processAnyValueAsByteArray(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("cannot parse log entry data: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown log field: %s", kvp.Key)
			}
		}
		logEntries = append(logEntries, &logEntry)
	}

	return logEntries, nil
}
