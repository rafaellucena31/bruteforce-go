package bfworkergeneric

import (
	"bruteforce-generic/models"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"
)

//StartBruteForce Inicia o processo de tentar descobrir a password
func StartBruteForce(job models.JobGeneric, chQuit chan struct{}, chProgresso chan models.ProgressoGeneric) {
	giroMaximo := len(job.Alfabeto) - 1
	giroMaximoPascal := len(job.Alfabeto)
	sairDoLoop := false
	dtStart := time.Now()
	temporizador := time.Now()
	totalChaves := (TransformArrayToNumber(job.ArrayFinal, giroMaximoPascal) - TransformArrayToNumber(job.ArrayInicial, giroMaximoPascal)) + 1
	if totalChaves > 500000000 {
		fmt.Printf("AA: %v AF: %v T: %f\n", job.ArrayAtual, job.ArrayFinal, totalChaves)
	}
	//fmt.Printf("%f\n", totalChaves)
	//time.Sleep(4 * time.Second)
	//fmt.Printf("AA: %v AF: %v\n", job.ArrayAtual, job.ArrayFinal)
loopCrack:
	for !IntsEquals(job.ArrayAtual, job.ArrayFinal) {

		//time.Sleep(1 * time.Second)
		select {
		case <-chQuit:
			sairDoLoop = true
			writeJobChannel(chProgresso, &temporizador, 10, job.ArrayAtual, job.ArrayInicial, giroMaximoPascal, totalChaves, dtStart, job.NumberJob, "", "", "Abortado")
			break loopCrack
		default:
			palavraTemp := arrayToString(job.ArrayAtual, job.Alfabeto)
			//fmt.Printf("\n%s\n", palavraTemp)
			payloadFull := strings.Replace(job.Payload, "##hash##", palavraTemp, 1)
			//fmt.Println(payloadFull)
			hash := md5.Sum([]byte(payloadFull))

			if BytesEquals(hash, job.Md5BytesHope) {
				writeJobChannel(chProgresso, &temporizador, -1, job.ArrayAtual, job.ArrayInicial, giroMaximoPascal, totalChaves, dtStart, job.NumberJob, palavraTemp, palavraTemp, "Password Found")
				sairDoLoop = true
				close(chQuit)
				break loopCrack
			}
			writeJobChannel(chProgresso, &temporizador, 10, job.ArrayAtual, job.ArrayInicial, giroMaximoPascal, totalChaves, dtStart, job.NumberJob, palavraTemp, "", "Executando")
			//fmt.Printf("Array: %v - %f \n", job.ArrayAtual, TransformArrayToNumber(job.ArrayAtual, len(job.Alfabeto)))
			RotacionarArray(job.ArrayAtual, 1, giroMaximo)
		}
	}

	if !sairDoLoop && IntsEquals(job.ArrayAtual, job.ArrayFinal) {
		//fmt.Printf("AA: %v AF: %v\n", job.ArrayAtual, job.ArrayFinal)
		palavraTemp := arrayToString(job.ArrayAtual, job.Alfabeto)
		//fmt.Printf("\n%s\n", palavraTemp)
		payloadFull := strings.Replace(job.Payload, "##hash##", palavraTemp, 1)
		hash := md5.Sum([]byte(payloadFull))
		if BytesEquals(hash, job.Md5BytesHope) {
			writeJobChannel(chProgresso, &temporizador, -1, job.ArrayAtual, job.ArrayInicial, giroMaximoPascal, totalChaves, dtStart, job.NumberJob, palavraTemp, palavraTemp, "Password Found")
			close(chQuit)
		}
		writeJobChannel(chProgresso, &temporizador, 10, job.ArrayAtual, job.ArrayInicial, giroMaximoPascal, totalChaves, dtStart, job.NumberJob, palavraTemp, "", "Finalizado")
	}
}

func writeJobChannel(chProgresso chan models.ProgressoGeneric, temporizador *time.Time, TempoSegundosDisparo float64, ArrayAtual []int, ArrayInicial []int, giroMaximoPascal int, totalChaves float64, dtStart time.Time, NumberJob float64, palavraTemp string, passwordFound string, Status string) {
	//feitos := (TransformArrayToNumber(ArrayAtual, giroMaximoPascal)) - (TransformArrayToNumber(ArrayInicial, giroMaximoPascal)) + 1
	totalElapsed := time.Since(*temporizador).Seconds()
	if TempoSegundosDisparo > 0 {
		if math.Ceil(totalElapsed) >= TempoSegundosDisparo {
			*temporizador = time.Now()
			//fmt.Printf("aqui")
			//time.Sleep(40 * time.Second)
			progresso := createProgressJob(ArrayAtual, ArrayInicial, giroMaximoPascal, totalChaves, dtStart, NumberJob, palavraTemp, passwordFound, Status)
			chProgresso <- progresso
		}
	} else {
		progresso := createProgressJob(ArrayAtual, ArrayInicial, giroMaximoPascal, totalChaves, dtStart, NumberJob, palavraTemp, passwordFound, Status)
		chProgresso <- progresso
	}
}

func createProgressJob(ArrayAtual []int, ArrayInicial []int, giroMaximoPascal int, totalChaves float64, dtStart time.Time, NumberJob float64, palavraTemp string, passwordFound string, Status string) models.ProgressoGeneric {

	feitos := (TransformArrayToNumber(ArrayAtual, giroMaximoPascal)) - (TransformArrayToNumber(ArrayInicial, giroMaximoPascal)) + 1

	progressoRelativo := feitos / totalChaves
	progressoPercentual := math.Round(progressoRelativo*100) / 100
	progressoJob := models.ProgressoGeneric{}
	progressoJob.DtStart = dtStart
	progressoJob.NumberJob = NumberJob
	progressoJob.UltimaPalavra = palavraTemp
	progressoJob.Progresso = progressoPercentual
	progressoJob.KeysPerSecond = math.Ceil(feitos / math.Ceil(time.Since(dtStart).Seconds()))
	progressoJob.TotalTestada = feitos
	progressoJob.TotalChaves = totalChaves
	progressoJob.PasswordFound = passwordFound
	progressoJob.Status = Status
	return progressoJob
}

//IntsEquals Verifica se dois arrays são iguais
func IntsEquals(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

//BytesEquals compara 2 arrays de bytes
func BytesEquals(a [16]byte, b [16]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

//TransformArrayToNumber Transforma um array na base do tamanho do alfabeto para a base decimal
func TransformArrayToNumber(array []int, tamanhoGiro int) float64 {
	NumeroBase10 := float64(0)
	j := len(array) - 1
	for i := (len(array) - 1); i >= 0; i-- {
		NumeroBase10 = NumeroBase10 + float64(array[i])*math.Pow(float64(tamanhoGiro), float64((j-i)))
	}
	return NumeroBase10
}

//TransformNumberToArray Dado um número na base 10, descobre a posicao correta do cilindro baseado no tamanho do alfabeto
func TransformNumberToArray(posicao float64, tamanhoGiro int, slicerOutArray *[]int) {
	quociente := math.Floor(posicao / float64(tamanhoGiro))
	resto := math.Ceil(math.Mod(posicao, float64(tamanhoGiro)))
	if quociente >= float64(tamanhoGiro) {
		TransformNumberToArray(quociente, tamanhoGiro, slicerOutArray)
		*slicerOutArray = append(*slicerOutArray, int(resto))
	} else {
		*slicerOutArray = append(*slicerOutArray, int(quociente), int(resto))
	}
}

//PaddingArray Adiciona 0 a esquerda do array
func PaddingArray(tamanhoArray int, slicerOutArray []int) []int {
	if tamanhoArray > len(slicerOutArray) {
		tmp := make([]int, tamanhoArray)
		pos := tamanhoArray - len(slicerOutArray)
		copy(tmp[pos:], (slicerOutArray))
		return tmp
	} else if tamanhoArray < len(slicerOutArray) {
		return slicerOutArray[0:tamanhoArray]
	}
	return slicerOutArray
}

func arrayToString(array []int, alfabeto string) string {
	retString := ""
	for i := 0; i < len(array); i++ {
		retString = retString + string(alfabeto[array[i]])
	}
	return retString
}

//RotacionarArray Gira o ArrayAtual x vezes a partir de uma posição.
func RotacionarArray(ArrayAtual []int, qtdGiros, indiceMaximo int) []int {
	ArrayAGirar := ArrayAtual
	intLenArray := len(ArrayAGirar) - 1
	for i := 0; i < qtdGiros; i++ {
		rotorArray(&ArrayAGirar, intLenArray, indiceMaximo)
	}
	return ArrayAGirar
}

//RotacionarArraySafe não deixa girar mais que o máximo do cilindro, e se isso acontecer, informa na variável booleana
func RotacionarArraySafe(ArrayAtual []int, qtdGiros, indiceMaximo int) ([]int, bool) {
	ArrayAGirar := make([]int, len(ArrayAtual))
	copy(ArrayAGirar, ArrayAtual)

	ArrayLimite := make([]int, len(ArrayAtual))
	for i := 0; i < len(ArrayLimite); i++ {
		ArrayLimite[i] = indiceMaximo
	}
	intLenArray := len(ArrayAGirar) - 1
	for i := 0; i < qtdGiros; i++ {
		if IntsEquals(ArrayAGirar, ArrayLimite) {
			return ArrayAGirar, true
		}
		rotorArray(&ArrayAGirar, intLenArray, indiceMaximo)

	}
	return ArrayAGirar, false
}

func rotorArray(ptrArrayAtual *[]int, rotorMaster, giroMaximo int) {
	if (*ptrArrayAtual)[rotorMaster] < giroMaximo {
		(*ptrArrayAtual)[rotorMaster] = (*ptrArrayAtual)[rotorMaster] + 1
	} else {
		(*ptrArrayAtual)[rotorMaster] = 0
		if rotorMaster > 0 {
			rotorArray(ptrArrayAtual, rotorMaster-1, giroMaximo)
		}
	}
}

func getMD5Hash2(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
