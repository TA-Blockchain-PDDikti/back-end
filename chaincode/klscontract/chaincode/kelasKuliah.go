package chaincode

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ============================================================================================================================
// Logger
// ============================================================================================================================

var logger = flogging.MustGetLogger("KLSContract")


// ============================================================================================================================
// Contract Definitions
// ============================================================================================================================

type KLSContract struct {
	contractapi.Contract
}


// ============================================================================================================================
// Asset Definitions - The ledger will store Kelas Kuliah (KLS) data
// ============================================================================================================================

type KelasKuliah struct {
	ID      			string 		`json:"id"`
	IdSMS				string 		`json:"idSms"`
	IdMK				string 		`json:"idMk"`
	NamaKLS				string 		`json:"namaKls"`
	Semester			string 		`json:"semester"`
	SKS					int 		`json:"sks"`
	ListPTK				[]string 	`json:"listPtk"`
	ListPD				[]string 	`json:"listPd"`
}


// ============================================================================================================================
// Struct Definitions - Pendidik Tenaga Kependidikan (PTK)
// ============================================================================================================================

type PendidikTenagaKependidikan struct {
	ID      			string `json:"id"`
	IdSP				string `json:"idSp"`
	IdSMS				string `json:"idSms"`
	NamaPTK				string `json:"namaPtk"`
	NIDN				string `json:"nidn"`
}


// ============================================================================================================================
// Struct Definitions - Peserta Didik (PD)
// ============================================================================================================================

type PesertaDidik struct {
	ID      			string 	`json:"id"`
	IdSP				string 	`json:"idSp"`
	IdSMS				string 	`json:"idSms"`
	NamaPD				string 	`json:"namaPd"`
	NIPD				string 	`json:"nipd"`
}


// ============================================================================================================================
// Struct Definitions - Full Kelas Kuliah (KLS) data
// ============================================================================================================================

type KelasKuliahFull struct {
	ID      			string 							`json:"id"`
	IdSMS				string 							`json:"idSms"`
	IdMK				string 							`json:"idMk"`
	NamaKLS				string 							`json:"namaKls"`
	Semester			string 							`json:"semester"`
	SKS					int 							`json:"sks"`
	ListPTK				[]PendidikTenagaKependidikan 	`json:"listPtk"`
	ListPD				[]PesertaDidik 					`json:"listPd"`
}


// ============================================================================================================================
// Error Messages
// ============================================================================================================================

const (
	ER11 string = "ER11-Incorrect number of arguments. Required %d arguments, but you have %d arguments."
	ER12        = "ER12-KelasKuliah with id '%s' already exists."
	ER13        = "ER13-KelasKuliah with id '%s' doesn't exist."
	ER31        = "ER31-Failed to change to world state: %v."
	ER32        = "ER32-Failed to read from world state: %v."
	ER33        = "ER33-Failed to get result from iterator: %v."
	ER34        = "ER34-Failed unmarshaling JSON: %v."
	ER35        = "ER35-Failed parsing string to integer: %v."
	ER36        = "ER36-Failed parsing string to float: %v."
	ER37        = "ER37-Failed to query another chaincode (%s): %v."
	ER41        = "ER41-Access is not permitted with MSDPID '%s'."
	ER42        = "ER42-Unknown MSPID: '%s'."
)


// ============================================================================================================================
// Channel Name & Contract Name In The Channel
// ============================================================================================================================

const (
	AcademicChannel	string = "academicchannel"
	PTKContract 	string = "ptkcontract"
	PDContract 		string = "pdcontract"
	QSCC			string = "qscc"
)


// ============================================================================================================================
// CreateKls - Issues a new Kelas Kuliah (KLS) to the world state with given details.
// Inputs - ID, Id SMS, Id MK, Nama KLS, Semester, SKS
// ============================================================================================================================

func (s *KLSContract) CreateKls(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run CreateKls function with args: %+q.", args)

	if len(args) != 6 {
		logger.Errorf(ER11, 6, len(args))
		return fmt.Errorf(ER11, 6, len(args))
	}

	id:= args[0]
	idSms:= args[1]
	idMk:= args[2]
	namaKls:= args[3]
	semester:= args[4]
	sksStr:= args[5]

	exists, err := isKlsExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		logger.Errorf(ER12, id)
		return fmt.Errorf(ER12, id)
	}

	sks, err := strconv.Atoi(sksStr)
	if err != nil {
		logger.Errorf(ER35, id)
		return fmt.Errorf(ER35, id)
	}

	kls := KelasKuliah{
		ID:      			id,
		IdSMS:				idSms,
		IdMK:				idMk,
		NamaKLS:			namaKls,
		Semester:			semester,
		SKS:				sks,
		ListPTK:			[]string{},
		ListPD:				[]string{},
	}

	klsJSON, err := json.Marshal(kls)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, klsJSON)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}


// ============================================================================================================================
// UpdateKls - Updates an existing Kelas Kuliah (KLS) in the world state with provided parameters.
// Inputs - ID, Id SMS, Id MK, Nama KLS, Semester, SKS
// ============================================================================================================================

func (s *KLSContract) UpdateKls(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run UpdateKls function with args: %+q.", args)

	if len(args) != 6 {
		logger.Errorf(ER11, 6, len(args))
		return fmt.Errorf(ER11, 6, len(args))
	}

	id:= args[0]
	idSms:= args[1]
	idMk:= args[2]
	namaKls:= args[3]
	semester:= args[4]
	sksStr:= args[5]

	kls, err := getKlsStateById(ctx, id)
	if err != nil {
		return err
	}

	sks, err := strconv.Atoi(sksStr)
	if err != nil {
		logger.Errorf(ER35, id)
		return fmt.Errorf(ER35, id)
	}

	kls.IdSMS = idSms
	kls.IdMK = idMk
	kls.NamaKLS = namaKls
	kls.Semester = semester
	kls.SKS = sks

	klsJSON, err := json.Marshal(kls)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, klsJSON)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}


// ============================================================================================================================
// UpdateKlsListPtk - Updates List PTK of an existing Kelas Kuliah (KLS) in the world state.
// Arguments - ID, List PTK
// ============================================================================================================================

func (s *KLSContract) UpdateKlsListPtk(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run UpdateKlsListPtk function with args: %+q.", args)

	if len(args) != 2 {
		logger.Errorf(ER11, 2, len(args))
		return fmt.Errorf(ER11, 2, len(args))
	}

	id:= args[0]
	listPtkStr:= args[1]

	kls, err := getKlsStateById(ctx, id)
	if err != nil {
		return err
	}

	listPtkStr = strings.Replace(listPtkStr, "[", "", -1)
	listPtkStr = strings.Replace(listPtkStr, "]", "", -1)
	splitter := regexp.MustCompile(` *, *`)
	listPtk :=  splitter.Split(listPtkStr, -1)

	kls.ListPTK = listPtk

	klsJSON, err := json.Marshal(kls)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, klsJSON)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}


// ============================================================================================================================
// UpdateKlsListPd - Updates List PD of an existing Kelas Kuliah (KLS) in the world state.
// Arguments - ID, List PD
// ============================================================================================================================

func (s *KLSContract) UpdateKlsListPd(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run UpdateKlsListPd function with args: %+q.", args)

	if len(args) != 2 {
		logger.Errorf(ER11, 2, len(args))
		return fmt.Errorf(ER11, 2, len(args))
	}

	id:= args[0]
	listPdStr:= args[1]

	kls, err := getKlsStateById(ctx, id)
	if err != nil {
		return err
	}

	listPdStr = strings.Replace(listPdStr, "[", "", -1)
	listPdStr = strings.Replace(listPdStr, "]", "", -1)
	splitter := regexp.MustCompile(` *, *`)
	listPd :=  splitter.Split(listPdStr, -1)

	kls.ListPD = listPd

	klsJSON, err := json.Marshal(kls)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, klsJSON)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}


// ============================================================================================================================
// DeleteKls - Deletes an given Kelas Kuliah (KLS) from the world state.
// Arguments - ID
// ============================================================================================================================

func (s *KLSContract) DeleteKls(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run DeleteKls function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return fmt.Errorf(ER11, 1, len(args))
	}

	id:= args[0]

	exists, err := isKlsExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ER13, id)
	}

	err = ctx.GetStub().DelState(id)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}


// ============================================================================================================================
// GetAllKls - Returns all Kelas Kuliah (KLS) found in world state.
// No Arguments
// ============================================================================================================================

func (s *KLSContract) GetAllKls(ctx contractapi.TransactionContextInterface) ([]*KelasKuliah, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetAllKls function with args: %+q.", args)

	if len(args) != 0 {
		logger.Errorf(ER11, 0, len(args))
		return nil, fmt.Errorf(ER11, 0, len(args))
	}

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf(ER32, err)
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}


// ============================================================================================================================
// GetKlsById - Get the Kelas Kuliah (KLS) stored in the world state with given id.
// Arguments - ID
// ============================================================================================================================

func (s *KLSContract) GetKlsById(ctx contractapi.TransactionContextInterface) (*KelasKuliah, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetKlsById function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	id:= args[0]

	kls, err := getKlsStateById(ctx, id)
	if err != nil {
		return nil, err
	}

	return kls, nil
}


// ============================================================================================================================
// GetFullKlsById - Get the Kelas Kuliah (KLS) stored in the world state with given id.
// Arguments - ID
// ============================================================================================================================

func (s *KLSContract) GetFullKlsById(ctx contractapi.TransactionContextInterface) (*KelasKuliahFull, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetFullKlsById function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	id:= args[0]

	kls, err := getKlsStateById(ctx, id)
	if err != nil {
		return nil, err
	}

	var klsFull KelasKuliahFull

	klsFull.ID = kls.ID
	klsFull.IdSMS = kls.IdSMS
	klsFull.IdMK = kls.IdMK
	klsFull.NamaKLS = kls.NamaKLS
	klsFull.Semester = kls.Semester
	klsFull.SKS = kls.SKS

	if len(kls.ListPTK) != 0 {
		listPtk, err := getPtkByIds(ctx, "[" + strings.Join(kls.ListPTK, ",") + "]")
			if err != nil {
				return nil, err
			}
		klsFull.ListPTK = listPtk
	} else {
		klsFull.ListPTK = []PendidikTenagaKependidikan{}
	}

	if len(kls.ListPD) != 0 {
		listPd, err := getPdByIds(ctx, "[" + strings.Join(kls.ListPD, ",") + "]")
			if err != nil {
				return nil, err
			}
		klsFull.ListPD = listPd
	} else {
		klsFull.ListPD = []PesertaDidik{}
	}

	return &klsFull, nil
}


// ============================================================================================================================
// GetKlsByIdMk - Get the Kelas Kuliah (KLS) stored in the world state with given IdMk.
// Arguments - idMk
// ============================================================================================================================

func (t *KLSContract) GetKlsByIdMk(ctx contractapi.TransactionContextInterface) ([]*KelasKuliah, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetKlsByIdMk function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	idMk:= args[0]

	queryString := fmt.Sprintf(`{"selector":{"idMk":"%s"}}`, idMk)
	return getQueryResultForQueryString(ctx, queryString)
}


// ============================================================================================================================
// GetKlsByIdPtk - Get the Kelas Kuliah (KLS) stored in the world state with given IdPtk.
// Arguments - idPtk
// ============================================================================================================================

func (t *KLSContract) GetKlsByIdPtk(ctx contractapi.TransactionContextInterface) ([]*KelasKuliah, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetKlsByIdPtk function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	idPtk:= args[0]

	queryString := fmt.Sprintf(`{"selector":{"listPtk":{"$elemMatch":{"$eq":"%s"}}}}`, idPtk)
	return getQueryResultForQueryString(ctx, queryString)
}


// ============================================================================================================================
// isKlsExists - Returns true when Kelas Kuliah (KLS) with given ID exists in world state.
// ============================================================================================================================

func isKlsExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	logger.Infof("Run isKlsExists function with id: '%s'.", id)

	klsJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		logger.Errorf(ER32, err)
		return false, fmt.Errorf(ER32, err)
	}

	return klsJSON != nil, nil
}


// ============================================================================================================================
// getKlsStateById - Get KLS state with given id.
// ============================================================================================================================

func getKlsStateById(ctx contractapi.TransactionContextInterface, id string) (*KelasKuliah, error) {
	logger.Infof("Run getKlsStateById function with id: '%s'.", id)

	klsJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf(ER32, err)
	}
	if klsJSON == nil {
		return nil, fmt.Errorf(ER13, id)
	}

	var kls KelasKuliah
	err = json.Unmarshal(klsJSON, &kls)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return &kls, nil
}


// ============================================================================================================================
// getPtkByIds - Get PTK with given idsPtk.
// ============================================================================================================================

func getPtkByIds(ctx contractapi.TransactionContextInterface, idsPtk string) ([]PendidikTenagaKependidikan, error) {
	logger.Infof("Run getPtkByIds function with idsPtk: '%s'.", idsPtk)

	params := []string{"GetPtkByIds", idsPtk}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(PTKContract, queryArgs, AcademicChannel)
	logger.Infof("Response Payload: '%+q'.", response.Payload)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(ER37, PTKContract, response.Message)
	}

	var listPtk []PendidikTenagaKependidikan
	err := json.Unmarshal([]byte(response.Payload), &listPtk)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return listPtk, nil
}


// ============================================================================================================================
// getPdByIds - Get PD with given idsPd.
// ============================================================================================================================

func getPdByIds(ctx contractapi.TransactionContextInterface, idsPd string) ([]PesertaDidik, error) {
	logger.Infof("Run getPdByIds function with idsPd: '%s'.", idsPd)

	params := []string{"GetPdByIds", idsPd}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(PDContract, queryArgs, AcademicChannel)
	logger.Infof("Response Payload: '%+q'.", response.Payload)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(ER37, PDContract, response.Message)
	}

	var listPd []PesertaDidik
	err := json.Unmarshal([]byte(response.Payload), &listPd)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return listPd, nil
}


// ============================================================================================================================
// constructQueryResponseFromIterator - Constructs a slice of assets from the resultsIterator.
// ============================================================================================================================

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*KelasKuliah, error) {
	logger.Infof("Run constructQueryResponseFromIterator function.")

	var klsList []*KelasKuliah

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf(ER33, err)
		}

		var kls KelasKuliah
		err = json.Unmarshal(queryResult.Value, &kls)
		if err != nil {
			return nil, fmt.Errorf(ER34, err)
		}
		klsList = append(klsList, &kls)
	}

	return klsList, nil
}


// ============================================================================================================================
// getQueryResultForQueryString - Get a query result from query string
// ============================================================================================================================

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*KelasKuliah, error) {
	logger.Infof("Run getQueryResultForQueryString function with queryString: '%s'.", queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf(ER32, err)
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}
