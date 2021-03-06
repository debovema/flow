package instance

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/flow/definition"
	"github.com/project-flogo/flow/model"
	_ "github.com/project-flogo/flow/support/test"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

const defJSON = `
{
    "type": 1,
    "name": "test",
    "model": "test",
    "rootTask": {
      "id": 1,
      "type": 1,
      "activityType": "",
      "name": "root",
      "tasks": [
        {
          "id": 2,
          "activityRef": "test-log",
          "name": "a",
          "attributes": [
            {
              "name": "message",
              "value": "test message",
              "type": "string"
            }
          ]
        },
        {
          "id": 3,
          "activityRef": "test-counter",
          "name": "b",
          "attributes": [
            {
              "name": "counterName",
              "value": "test_counter",
              "type": "string"
            }
          ]
        }
      ],
      "links": [
        { "id": 1, "type": 1,  "name": "","from": 2, "to": 3 }
      ]
    }
  }
`

func TestFullSerialization(t *testing.T) {

	defRep := &definition.DefinitionRep{}
	err := json.Unmarshal([]byte(defJSON), defRep)
	assert.Nil(t, err)

	log.RootLogger().Infof("Def Rep: %v", defRep)

	def, _ := definition.NewDefinition(defRep)
	assert.NotNil(t, def)

	instance := NewIndependentInstance("12345", "uri", def, log.RootLogger())

	instance.Start(nil)

	hasWork := true

	for hasWork && instance.Status() < model.FlowStatusCompleted {
		hasWork = instance.DoStep()

		json, _ := json.Marshal(instance)
		log.RootLogger().Debugf("Snapshot: %s\n", string(json))
	}

}

func TestChangeSerialization(t *testing.T) {

	defRep := &definition.DefinitionRep{}
	err := json.Unmarshal([]byte(defJSON), defRep)
	assert.Nil(t, err)

	log.RootLogger().Infof("Def Rep: %v", defRep)

	def, _ := definition.NewDefinition(defRep)
	assert.NotNil(t, def)

	instance := NewIndependentInstance("12345", "uri", def, log.RootLogger())

	instance.Start(nil)

	hasWork := true

	for hasWork && instance.Status() < model.FlowStatusCompleted {
		hasWork = instance.DoStep()

		json, _ := json.Marshal(instance.ChangeTracker)
		log.RootLogger().Debugf("Change: %s\n", string(json))
	}
}

//func TestIncrementalSerialization(t *testing.T) {
//
//	defRep := &flowdef.DefinitionRep{}
//	json.Unmarshal([]byte(defJSON), defRep)
//
//	idGen, _ := util.NewGenerator()
//	id := idGen.NextAsString()
//
//	def, _ := flowdef.NewDefinition(defRep)
//
//	instance := NewFlowInstance(id, "uri2", def)
//
//	instance.Start(nil)
//
//	hasWork := true
//
//	for hasWork && instance.Status() < StatusCompleted {
//		hasWork = instance.DoStep()
//
//		json, _ := json.Marshal(instance.GetChanges())
//		log.Debugf("Changes: %s\n", string(json))
//	}
//}
