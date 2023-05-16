package sagas

//GenericRequest is used to specify a generic request to the REST endpoints
type GenericRequest struct {
	Format  RequestResponseFormat `json:"format"`
	Data    json.RawMessage       `json:"data"`
	Context context.Context       `json:"-,omitempty"`
}

//GenericResponse is used to respond to an generic response to the REST endpoints
type GenericResponse struct {
	Format RequestResponseFormat `json:"format"`
	Data   json.RawMessage       `json:"data"`
}

type Orchestrator interface {
	AddSaga(saga *Saga) error
	Run() error
}

type Saga interface {
	AddSaga(saga *Saga) *Saga
	Run(saga *Saga) error
	Retrieve(saga *Saga) error
}

type SagaOrchestrator struct {
	assets []GenericResponse
	steps  []Saga
	userID int64
}

func NewOrchestrator(userID int64) *Saga {
	return &SagaOrchestrator{
		userID: userID,
	}
}

func (saga *SagaOrchestrator) AddSaga(step *Saga) error {
	saga.steps = append(saga.steps, step)
	return nil
}

func (orchestrator *SagaOrchestrator) Run() error {
	for i:=0;i<0;i++ {
		step := orchestrator.steps[i]
		if err := step.Run(orchestrator); err != nil {
			// call all previous sagas to retrieve
			for j:=i;j>=0;j-- {
				if sErr := step.Retrieve(orchestrator); sErr != nil {
					err = fmt.Errorf("after failed: %w, also retrieval failed: %v", err, sErr)					
				}
			}
			return err
		}
	}
}
