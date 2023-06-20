package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ============================================================================================================================
// Logger
// ============================================================================================================================

var logger = flogging.MustGetLogger("IJZContract")


// ============================================================================================================================
// Contract Definitions
// ============================================================================================================================

type IJZContract struct {
	contractapi.Contract
}


// ============================================================================================================================
// Asset Definitions - The ledger will store Ijazah Mahasiswa (IJZ) data
// ============================================================================================================================

type Ijazah struct {
	ID      			string 			`json:"id"`
	IdSP				string 			`json:"idSp"`
	IdSMS				string 			`json:"idSms"`
	IdPD				string 			`json:"idPd"`
	JenjangPendidikan	string 			`json:"jenjangPendidikan"`
	NomorIjazah			string 			`json:"nomorIjazah"`
	TanggalLulus		string 			`json:"tanggalLulus"`
	RemainingApprover	int 			`json:"remainingApprover"`
	Approvers			[]string 		`json:"Approvers"`
}


// ============================================================================================================================
// Struct Definitions - Satuan Pendidikan (SP)
// ============================================================================================================================

type SatuanPendidikan struct {
	ID      		string `json:"id"`
	NamaSP			string `json:"namaSp"`
}


// ============================================================================================================================
// Struct Definitions - Satuan Managemen Sumberdaya (SMS)
// ============================================================================================================================

type SatuanManagemenSumberdaya struct {
	ID      			string 		`json:"id"`
	NamaSMS				string 		`json:"namaSms"`
	JenjangPendidikan	string 		`json:"jenjangPendidikan"`
	ApproversTSK		[]string 	`json:"approversTsk"`
	ApproversIJZ		[]string 	`json:"approversIjz"`
}


// ============================================================================================================================
// Struct Definitions - Pendidik Tenaga Kependidikan (PTK)
// ============================================================================================================================

type PendidikTenagaKependidikan struct {
	ID      			string `json:"id"`
	NamaPTK				string `json:"namaPtk"`
	NIDN				string `json:"nidn"`
	Jabatan				string `json:"jabatan"`
}


// ============================================================================================================================
// Struct Definitions - Peserta Didik (PD)
// ============================================================================================================================

type PesertaDidik struct {
	ID      			string 	`json:"id"`
	NamaPD				string 	`json:"namaPd"`
	NIPD				string 	`json:"nipd"`
}


// ============================================================================================================================
// Struct Definitions - Result of Query Ijazah Mahasiswa (IJZ) data
// ============================================================================================================================

type IjazahResult struct {
	ID      			string 							`json:"id"`
	SP					*SatuanPendidikan 				`json:"sp"`
	SMS					*SatuanManagemenSumberdaya 		`json:"sms"`
	PD					*PesertaDidik					`json:"pd"`
	JenjangPendidikan	string 							`json:"jenjangPendidikan"`
	NomorIjazah			string 							`json:"nomorIjazah"`
	TanggalLulus		string 							`json:"tanggalLulus"`
	RemainingApprover	int 							`json:"remainingApprover"`
	Approvers			[]*PendidikTenagaKependidikan	`json:"Approvers"`
	ApprovalTxId		[]string						`json:"approvalTxId"`
}


// ============================================================================================================================
// Error Messages
// ============================================================================================================================

const (
	ER11 string = "ER11-Incorrect number of arguments. Required %d arguments, but you have %d arguments."
	ER12        = "ER12-Ijazah with id '%s' already exists."
	ER13        = "ER13-Ijazah with id '%s' doesn't exist."
	ER14        = "ER14-Ijazah with id '%s' no longer require approval."
	ER15        = "ER15-Ijazah with id '%s' already approved by PTK with id '%s'."
	ER16        = "ER16-Ijazah with id '%s' cannot be approved by PTK with id '%s' in this step."
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
	SPContract 		string = "spcontract"
	SMSContract 	string = "smscontract"
	PTKContract 	string = "ptkcontract"
	PDContract 		string = "pdcontract"
)


// ============================================================================================================================
// CreateIjz - Issues a new Ijazah Mahasiswa (IJZ) to the world state with given details.
// Arguments - ID, Id SP, Id SMS, Id PD, Jenjang Pendidikan, Nomor Ijazah, Tanggal Lulus
// ============================================================================================================================

func (s *IJZContract) CreateIjz (ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run CreateIjz function with args: %+q.", args)

	if len(args) != 7 {
		logger.Errorf(ER11, 7, len(args))
		return fmt.Errorf(ER11, 7, len(args))
	}

	id:= args[0]
	idSp:= args[1]
	idSms:= args[2]
	idPd:= args[3]
	jenjangPendidikan:= args[4]
	nomorIjazah:= args[5]
	tanggalLulus:= args[6]

	exists, err := isIjzExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		logger.Errorf(ER12, id)
		return fmt.Errorf(ER12, id)
	}

	smsApproverIjz, err := getSmsApproverIjz(ctx, idSms)
	if err != nil {
		return err
	}

	ijz := Ijazah{
		ID:      			id,
		IdSP:				idSp,
		IdSMS:				idSms,
		IdPD:				idPd,
		JenjangPendidikan:	jenjangPendidikan,
		NomorIjazah:		nomorIjazah,
		TanggalLulus:		tanggalLulus,
		RemainingApprover:	len(smsApproverIjz),
		Approvers:			[]string{},
	}

	ijzJSON, err := json.Marshal(ijz)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, ijzJSON)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}


// ============================================================================================================================
// UpdateIjz - Updates an existing Ijazah Mahasiswa (IJZ) in the world state with provided parameters.
// Arguments - ID, Id SP, Id SMS, Id PD, Jenjang Pendidikan, Nomor Ijazah, Tanggal Lulus
// ============================================================================================================================

func (s *IJZContract) UpdateIjz (ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run UpdateIjz function with args: %+q.", args)

	if len(args) != 7 {
		logger.Errorf(ER11, 7, len(args))
		return fmt.Errorf(ER11, 7, len(args))
	}

	id:= args[0]
	idSp:= args[1]
	idSms:= args[2]
	idPd:= args[3]
	jenjangPendidikan:= args[4]
	nomorIjazah:= args[5]
	tanggalLulus:= args[6]

	exists, err := isIjzExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ER13, id)
	}

	smsApproverIjz, err := getSmsApproverIjz(ctx, idSms)
	if err != nil {
		return err
	}

	ijz := Ijazah{
		ID:      			id,
		IdSP:				idSp,
		IdSMS:				idSms,
		IdPD:				idPd,
		JenjangPendidikan:	jenjangPendidikan,
		NomorIjazah:		nomorIjazah,
		TanggalLulus:		tanggalLulus,
		RemainingApprover:	len(smsApproverIjz),
		Approvers:			[]string{},
	}

	ijzJSON, err := json.Marshal(ijz)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, ijzJSON)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}

// ============================================================================================================================
// AddIjzApproval - Add Approval for an existing Ijazah Mahasiswa (IJZ) in the world state.
// Arguments - ID, Approver Id
// ============================================================================================================================

func (s *IJZContract) AddIjzApproval (ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run AddIjzApproval function with args: %+q.", args)

	if len(args) != 2 {
		logger.Errorf(ER11, 2, len(args))
		return fmt.Errorf(ER11, 2, len(args))
	}

	id:= args[0]
	approver:= args[1]

	exists, err := isIjzExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ER13, id)
	}

	ijz, err := getIjzStateById(ctx, id)
	if err != nil {
		return err
	}

	if ijz.RemainingApprover == 0 {
		return fmt.Errorf(ER14, id)
	}

	if contains(ijz.Approvers, approver) {
		return fmt.Errorf(ER15, id, approver)
	}

	smsApproverIjz, err := getSmsApproverIjz(ctx, ijz.IdSMS)
	if err != nil {
		return err
	}

	approvalStep := len(ijz.Approvers)
	if smsApproverIjz[approvalStep] != approver {
		return fmt.Errorf(ER16, id, approver)
	}

	ijz.Approvers = append(ijz.Approvers, approver)
	ijz.RemainingApprover = ijz.RemainingApprover - 1

	ijzJSON, err := json.Marshal(ijz)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, ijzJSON)
	if err != nil {
		logger.Errorf(ER31, err)
	}

	return err
}


// ============================================================================================================================
// DeleteIjz - Deletes an given Ijazah Mahasiswa (IJZ) from the world state.
// Arguments - ID
// ============================================================================================================================

func (s *IJZContract) DeleteIjz(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run DeleteIjz function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return fmt.Errorf(ER11, 1, len(args))
	}

	id:= args[0]

	exists, err := isIjzExists(ctx, id)
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
// GetAllIjz - Returns all Ijazah Mahasiswa (IJZ) found in world state.
// No Arguments
// ============================================================================================================================

func (s *IJZContract) GetAllIjz(ctx contractapi.TransactionContextInterface) ([]*Ijazah, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetAllIjz function with args: %+q.", args)

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
// GetIjzById - Get the Ijazah Mahasiswa (IJZ) stored in the world state with given id.
// Arguments - ID
// ============================================================================================================================

func (s *IJZContract) GetIjzById (ctx contractapi.TransactionContextInterface) (*IjazahResult, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetIjzById function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	id:= args[0]

	ijz, err := getIjzStateById(ctx, id)
	if err != nil {
		return nil, err
	}

	ijzResult, err := getCompleteDataIjz(ctx, ijz)
	if err != nil {
		return nil, err
	}

	return ijzResult, nil
}


// ============================================================================================================================
// GetIjzByIdSp - Get the Ijazah Mahasiswa (IJZ) stored in the world state with given IdSp.
// Arguments - idSp
// ============================================================================================================================

func (t *IJZContract) GetIjzByIdSp(ctx contractapi.TransactionContextInterface) ([]*IjazahResult, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetIjzByIdSp function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	idSp:= args[0]

	queryString := fmt.Sprintf(`{"selector":{"idSp":"%s"}}`, idSp)
	queryResult, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}

	var ijzList []*IjazahResult

	for _, ijz := range queryResult {
		ijzResult, err := getCompleteDataIjz(ctx, ijz)
		if err != nil {
			return nil, err
		}

		ijzList = append(ijzList, ijzResult)
	}

	return ijzList, nil
}

// ============================================================================================================================
// GetIjzByIdSms - Get the Ijazah Mahasiswa (IJZ) stored in the world state with given IdSms.
// Arguments - idSms
// ============================================================================================================================

func (t *IJZContract) GetIjzByIdSms(ctx contractapi.TransactionContextInterface) ([]*IjazahResult, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetIjzByIdSms function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	idSms:= args[0]

	queryString := fmt.Sprintf(`{"selector":{"idSms":"%s"}}`, idSms)
	queryResult, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}

	var ijzList []*IjazahResult

	for _, ijz := range queryResult {
		ijzResult, err := getCompleteDataIjz(ctx, ijz)
		if err != nil {
			return nil, err
		}

		ijzList = append(ijzList, ijzResult)
	}

	return ijzList, nil
}


// ============================================================================================================================
// GetIjzByIdPd - Get the Ijazah Mahasiswa (IJZ) stored in the world state with given IdPd.
// Arguments - idPd
// ============================================================================================================================

func (t *IJZContract) GetIjzByIdPd(ctx contractapi.TransactionContextInterface) ([]*IjazahResult, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetIjzByIdPd function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return nil, fmt.Errorf(ER11, 1, len(args))
	}

	idPd:= args[0]

	queryString := fmt.Sprintf(`{"selector":{"idPd":"%s"}}`, idPd)
	queryResult, err := getQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}

	var ijzList []*IjazahResult

	for _, ijz := range queryResult {
		ijzResult, err := getCompleteDataIjz(ctx, ijz)
		if err != nil {
			return nil, err
		}

		ijzList = append(ijzList, ijzResult)
	}

	return ijzList, nil
}


// ============================================================================================================================
// GetIjzAddApprovalTxIdById - Get the Ijazah Mahasiswa (IJZ) stored in the world state with given IdPd.
// Arguments - idPd
// ============================================================================================================================

func (t *IJZContract) GetIjzAddApprovalTxIdById(ctx contractapi.TransactionContextInterface) ([]string, error) {
	args := ctx.GetStub().GetStringArgs()[1:]

	logger.Infof("Run GetIjzAddApprovalTxIdById function with args: %+q.", args)

	if len(args) != 1 {
		logger.Errorf(ER11, 1, len(args))
		return []string{}, fmt.Errorf(ER11, 1, len(args))
	}

	id:= args[0]

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return []string{}, fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	txIdList := []string{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return []string{}, fmt.Errorf(err.Error())
		}

		var ijz Ijazah
		err = json.Unmarshal([]byte(response.Value), &ijz)
		if err != nil {
			return nil, fmt.Errorf(ER34, err)
		}

		if (len(ijz.Approvers) == 0) {
			break
		}

		txIdList = append([]string{response.TxId}, txIdList[0:]...)
	}

	return txIdList, nil
}


// ============================================================================================================================
// isIjzExists - Returns true when Ijazah Mahasiswa (IJZ) with given ID exists in world state.
// ============================================================================================================================

func isIjzExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	logger.Infof("Run isIjzExists function with id: '%s'.", id)

	ijzJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		logger.Errorf(ER32, err)
		return false, fmt.Errorf(ER32, err)
	}

	return ijzJSON != nil, nil
}


// ============================================================================================================================
// getIjzStateById - Get IJZ state with given id.
// ============================================================================================================================

func getIjzStateById(ctx contractapi.TransactionContextInterface, id string) (*Ijazah, error) {
	logger.Infof("Run getIjzStateById function with id: '%s'.", id)

	ijzJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf(ER32, err)
	}
	if ijzJSON == nil {
		return nil, fmt.Errorf(ER13, id)
	}

	var ijz Ijazah
	err = json.Unmarshal(ijzJSON, &ijz)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return &ijz, nil
}


// ============================================================================================================================
// getSmsApproverIjz - Get SMS Approver IJZ with given idSms.
// ============================================================================================================================

func getSmsApproverIjz(ctx contractapi.TransactionContextInterface, idSms string) ([]string, error) {
	logger.Infof("Run getSmsApproverIjz function with idSms: '%s'.", idSms)

	params := []string{"GetSmsById", idSms}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(SMSContract, queryArgs, AcademicChannel)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(ER37, SMSContract, response.Payload)
	}

	var sms SatuanManagemenSumberdaya
	err := json.Unmarshal([]byte(response.Payload), &sms)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return sms.ApproversIJZ, nil
}


// ============================================================================================================================
// getSpById - Get SP with given idSp.
// ============================================================================================================================

func getSpById(ctx contractapi.TransactionContextInterface, idSp string) (*SatuanPendidikan, error) {
	logger.Infof("Run getSpById function with idSp: '%s'.", idSp)

	params := []string{"GetSpById", idSp}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(SPContract, queryArgs, AcademicChannel)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(ER37, SPContract, response.Message)
	}

	var sp SatuanPendidikan
	err := json.Unmarshal([]byte(response.Payload), &sp)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return &sp, nil
}


// ============================================================================================================================
// getSmsById - Get SMS with given idSms.
// ============================================================================================================================

func getSmsById(ctx contractapi.TransactionContextInterface, idSms string) (*SatuanManagemenSumberdaya, error) {
	logger.Infof("Run getSmsById function with idSms: '%s'.", idSms)

	params := []string{"GetSmsById", idSms}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(SMSContract, queryArgs, AcademicChannel)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(ER37, SMSContract, response.Message)
	}

	var sms SatuanManagemenSumberdaya
	err := json.Unmarshal([]byte(response.Payload), &sms)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return &sms, nil
}


// ============================================================================================================================
// getPtkById - Get PTK with given idPtk.
// ============================================================================================================================

func getPtkById(ctx contractapi.TransactionContextInterface, idPtk string) (*PendidikTenagaKependidikan, error) {
	logger.Infof("Run getPtkById function with idPtk: '%s'.", idPtk)

	params := []string{"GetPtkById", idPtk}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(PTKContract, queryArgs, AcademicChannel)
	logger.Infof("Response Payload: '%+q'.", response.Payload)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(ER37, PTKContract, response.Message)
	}

	var ptk PendidikTenagaKependidikan
	err := json.Unmarshal([]byte(response.Payload), &ptk)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return &ptk, nil
}


// ============================================================================================================================
// getPdById - Get PD with given idPd.
// ============================================================================================================================

func getPdById(ctx contractapi.TransactionContextInterface, idPd string) (*PesertaDidik, error) {
	logger.Infof("Run getPdById function with idPd: '%s'.", idPd)

	params := []string{"GetPdById", idPd}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(PDContract, queryArgs, AcademicChannel)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(ER37, PDContract, response.Message)
	}

	var pd PesertaDidik
	err := json.Unmarshal([]byte(response.Payload), &pd)
	if err != nil {
		return nil, fmt.Errorf(ER34, err)
	}

	return &pd, nil
}


// ============================================================================================================================
// getIjzAddApprovalTxIdById - Get the Ijazah Mahasiswa (IJZ) Approval Tx Id with given ID.
// ============================================================================================================================

func getIjzAddApprovalTxIdById(ctx contractapi.TransactionContextInterface, id string) ([]string, error) {
	logger.Infof("Run getIjzAddApprovalTxIdById function with id: %s.", id)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return []string{}, fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	txIdList := []string{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return []string{}, fmt.Errorf(err.Error())
		}

		var ijz Ijazah
		err = json.Unmarshal([]byte(response.Value), &ijz)
		if err != nil {
			return nil, fmt.Errorf(ER34, err)
		}

		if (len(ijz.Approvers) == 0) {
			break
		}

		txIdList = append([]string{response.TxId}, txIdList[0:]...)
	}

	return txIdList, nil
}


// ============================================================================================================================
// getCompleteDataIjz - Get Complete the Ijazah Mahasiswa (IJZ).
// ============================================================================================================================

func getCompleteDataIjz(ctx contractapi.TransactionContextInterface, ijz *Ijazah) (*IjazahResult, error) {
	logger.Infof("Run getCompleteDataIjz function with ijz id: '%s'.", ijz.ID)

	var ijzResult IjazahResult

	ijzResult.ID 				= ijz.ID
	ijzResult.JenjangPendidikan = ijz.JenjangPendidikan
	ijzResult.NomorIjazah 		= ijz.NomorIjazah
	ijzResult.TanggalLulus 		= ijz.TanggalLulus
	ijzResult.RemainingApprover = ijz.RemainingApprover

	sp, err := getSpById(ctx, ijz.IdSP)
	if err != nil {
		return nil, err
	}
	ijzResult.SP = sp

	sms, err := getSmsById(ctx, ijz.IdSMS)
	if err != nil {
		return nil, err
	}
	sms.ApproversTSK = []string{}
	sms.ApproversIJZ = []string{}
	ijzResult.SMS = sms

	pd, err := getPdById(ctx, ijz.IdPD)
	if err != nil {
		return nil, err
	}
	ijzResult.PD = pd

	if len(ijz.Approvers) != 0 {
		var listPtk []*PendidikTenagaKependidikan

		for _, ptkId := range ijz.Approvers {
			ptk, err := getPtkById(ctx, ptkId)
			if err != nil {
				return nil, err
			}

			listPtk = append(listPtk, ptk)
		}

		ijzResult.Approvers = listPtk
	} else {
		ijzResult.Approvers = []*PendidikTenagaKependidikan{}
	}

	approvalTxIds, err := getIjzAddApprovalTxIdById(ctx, ijz.ID)
	if err != nil {
		return nil, err
	}
	ijzResult.ApprovalTxId = approvalTxIds

	return &ijzResult, nil
}


// ============================================================================================================================
// constructQueryResponseFromIterator - Constructs a slice of assets from the resultsIterator.
// ============================================================================================================================

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Ijazah, error) {
	logger.Infof("Run constructQueryResponseFromIterator function.")

	var ijzList []*Ijazah

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf(ER33, err)
		}

		var ijz Ijazah
		err = json.Unmarshal(queryResult.Value, &ijz)
		if err != nil {
			return nil, fmt.Errorf(ER34, err)
		}
		ijzList = append(ijzList, &ijz)
	}

	return ijzList, nil
}


// ============================================================================================================================
// getQueryResultForQueryString - Get a query result from query string
// ============================================================================================================================

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Ijazah, error) {
	logger.Infof("Run getQueryResultForQueryString function with queryString: '%s'.", queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf(ER32, err)
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}


// ============================================================================================================================
// contains - Check if the slice contains the given value
// ============================================================================================================================

func contains(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
}
