package models

import (
	"time"
)

//JobGeneric -> Representacao completa para um JOB de BruteForce
type JobGeneric struct {
	ArrayInicial []int
	ArrayFinal   []int
	ArrayAtual   []int
	Alfabeto     string
	TamPassword  int
	Payload      string
	Md5BytesHope [16]byte
	NumberJob    float64
}

//InputModelGeneric -> Representação da entrada do usuário no sistema
type InputModelGeneric struct {
	Alfabeto           string
	MinCaractere       int
	MaxCaractere       int
	TotalJobsParalelo  int
	Md5BytesHopeString string
	Salt               string
}

//ProgressoGeneric progresso de um job
type ProgressoGeneric struct {
	NumberJob     float64
	UltimaPalavra string
	KeysPerSecond float64
	Progresso     float64
	TotalChaves   float64
	TotalTestada  float64
	PasswordFound string
	Status        string
	DtStart       time.Time
}
