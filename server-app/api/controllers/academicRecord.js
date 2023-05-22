const academicRecordService = require('../services/academicRecord.js')
const { v4: uuidv4 } = require('uuid')

exports.createAcademicRecord = async(req, res) => {
    try {
        if (req.user.userType != "dosen") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const data = req.body;
        const idKls = data.idKls;
        const idDosen = data.idPtk;
        const idMahasiswa = data.idPd;
        const nilaiAngka = data.nilaiAngka;
        const nilaiHuruf = data.nilaiHuruf;
        const nilaiIndex = data.nilaiIndex;
        var id = data.id;

        // Randomize unique Id if there is no request id given
        if (!id) {
            id = uuidv4()
        }

        const result = await academicRecordService.createAcademicRecord(req.user.username, id, idKls, idDosen, idMahasiswa, nilaiAngka, nilaiHuruf, nilaiIndex)
        res.status(201).send({
            success: true,
            message: "Transaction nilai has been submitted",
        })
    } catch(error){
        res.status(400).send({
            success: false,
            error: error.toString(),
        })    
    }
}

exports.updateAcademicRecord = async(req, res) => {
    try {
        if (req.user.userType != "dosen") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const data = req.body;
        const idKls = data.idKls;
        const idDosen = data.idPtk;
        const idMahasiswa = data.idPd;
        const nilaiAngka = data.nilaiAngka;
        const nilaiHuruf = data.nilaiHuruf;
        const nilaiIndex = data.nilaiIndex;
        const idNilai = req.params.id
       
        const result = await academicRecordService.updateAcademicRecord(req.user.username,  idNilai, idKls, idDosen, idMahasiswa, nilaiAngka, nilaiHuruf, nilaiIndex)
        res.status(200).send({
            success: true,
            message: `Nilai dengan id ${idNilai} telah diubah`, 
        })
    }
    catch(error){
     res.status(400).send({
            success: false,
            error: error.toString(),
        })       
    }
}

exports.signAcademicRecord = async(req, res) => {
    try {
        if (req.user.userType != "dosen") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const idNilai = req.body.idNilai

        const result = await academicRecordService.signAcademicRecord(req.user.username, idNilai)
        res.status(200).send({
            success: true,
            message: `Record nilai dengan id ${idNilai} is signed`,
        })
    }
    catch(error){  
     res.status(400).send({
            success: false,
            error: error.toString(),
        })       
    }
}

exports.deleteAcademicRecord = async(req, res) => {
    try {
        if (req.user.userType != "dosen") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const idNilai = req.params.id
        const result = await academicRecordService.deleteAcademicRecord(req.user.username, idNilai);
        res.status(200).send({
            success: true,
            message: `Record nilai dengan id ${idNilai} telah dihapus`,
        })
    }
    catch(error){
     res.status(400).send({
            success: false,
            error: error.toString(),
        })       
    }
}

exports.getAllAcademicRecord = async(req, res) => {
    try {
        if (req.user.userType != "admin pddikti") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const data = await academicRecordService.getAllAcademicRecord(req.user.username) 
        res.status(200).send({data});
    } catch(error){
        res.status(400).send({
               success: false,
               error: error.toString(),
           })       
       }
}

exports.getAcademicRecordById = async(req, res) => {
    try {
        if (req.user.userType != "dosen") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const id = req.params.id
        const data = await academicRecordService.getAcademicRecordById(req.user.username, id) 
        res.status(200).send({data});
    } catch(error){
        res.status(400).send({
               success: false,
               error: error.toString(),
           })       
       }
}

exports.getAcademicRecordByIdMhsw = async(req, res) => {
    try {
        if (req.user.userType != "mahasiswa") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const id = req.params.id
        const data = await academicRecordService.getAcademicRecordByIdMhsw(req.user.username, id) 
        res.status(200).send({data});
    } catch(error){
        res.status(400).send({
               success: false,
               error: error.toString(),
           })       
       }
}

exports.getAcademicRecordByIdKls = async(req, res) => {
    try {
        if (req.user.userType != "dosen") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const id = req.params.id
        const data = await academicRecordService.getAcademicRecordByIdKls(req.user.username, id) 
        res.status(200).send({data});
    } catch(error){
        res.status(400).send({
               success: false,
               error: error.toString(),
           })       
       }
}

exports.setGrade = async(req, res) => {
    try {
        if (req.user.userType != "dosen") {
            return res.status(403).send({"result":`Forbidden Access for role ${req.user.userType}`})
        }
        const data = await academicRecordService.setGrade(req.user.username)
        res.status(200).send({data});
    } catch(error){
        res.status(400).send({
               success: false,
               error: error.toString(),
           })       
       }
}


