const certificateRouter = require('express').Router()
const certificateController = require('../controllers/academicCertificate.js')
const auth = require('../middleware/auth.js')

certificateRouter.use(auth)

certificateRouter.post('/', certificateController.createAcademicCertificate)
certificateRouter.post('/sign/:id', certificateController.signIjazah)
certificateRouter.post('/identifier/', certificateController.generateIdentifier)
certificateRouter.get('/identifier/', certificateController.getIdentifier)
certificateRouter.post('/verify/', certificateController.verify)
certificateRouter.post('/signer/', certificateController.addSigner)

// Ijazah
certificateRouter.get('/ijazah/', certificateController.getAllIjazah)
certificateRouter.get('/ijazah/:id', certificateController.getIjazahById)
certificateRouter.get('/ijazah/pt/:id', certificateController.getIjazahByIdPt)
certificateRouter.get('/ijazah/prodi/:id', certificateController.getIjazahByIdProdi)
certificateRouter.get('/ijazah/mahasiswa/:id', certificateController.getIjazahByIdMahasiswa)
certificateRouter.put('/ijazah/:id', certificateController.updateIjazah)

//transkrip
certificateRouter.get('/transkrip/', certificateController.getAllTranskrip)
certificateRouter.get('/transkrip/:id', certificateController.getTranskripById)
certificateRouter.get('/transkrip/pt/:id', certificateController.getTranskripByIdPt)
certificateRouter.get('/transkrip/prodi/:id', certificateController.getTranskripByIdProdi)
certificateRouter.get('/transkrip/mahasiswa/:id', certificateController.getTranskripByIdMahasiswa)
certificateRouter.put('/transkrip/:id', certificateController.updateTranskrip)

module.exports = certificateRouter;