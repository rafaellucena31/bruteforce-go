package main

//f32cd21c2cf422762bef2fcd150edb35 = 012Rafael
//abcdefghijklmnopqrstuvxzywJKR
import (
	"encoding/hex"
	"fmt"
	"math"
	"sync"

	tm "github.com/buger/goterm"
	hmz "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	bfworker "github.com/rafaellucena31/bf/core"
	models "github.com/rafaellucena31/bf/models"
)

func clearPartialScreen(x int, y int) {
	tm.MoveCursor(x, y)
	for i := y; i < tm.Height(); i++ {

		tm.Printf("%s\n", lineBlank(tm.Width()))
	}

	tm.MoveCursor(x, y)

}

func lineBlank(width int) string {
	x := ""
	for i := 0; i < width; i++ {
		x = x + " "
	}
	return x
}
func podeEscreverProgresso(JobsAtivos []models.ProgressoGeneric) bool {
	for i := 0; i < len(JobsAtivos); i++ {
		jobtmp := JobsAtivos[i]
		if jobtmp.Status == "Password Found" {
			return false
		}
	}
	return true
}
func main() {

	tm.Clear()

	inputModel := models.InputModelGeneric{}
	passwordFound := false

	for {
		tm.MoveCursor(1, 1)

		logo := `
██████╗ ██████╗ ██╗   ██╗████████╗███████╗    ███████╗ ██████╗ ██████╗  ██████╗███████╗    ███████╗██╗   ██╗███████╗████████╗███████╗███╗   ███╗
██╔══██╗██╔══██╗██║   ██║╚══██╔══╝██╔════╝    ██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔════╝    ██╔════╝╚██╗ ██╔╝██╔════╝╚══██╔══╝██╔════╝████╗ ████║
██████╔╝██████╔╝██║   ██║   ██║   █████╗      █████╗  ██║   ██║██████╔╝██║     █████╗      ███████╗ ╚████╔╝ ███████╗   ██║   █████╗  ██╔████╔██║
██╔══██╗██╔══██╗██║   ██║   ██║   ██╔══╝      ██╔══╝  ██║   ██║██╔══██╗██║     ██╔══╝      ╚════██║  ╚██╔╝  ╚════██║   ██║   ██╔══╝  ██║╚██╔╝██║
██████╔╝██║  ██║╚██████╔╝   ██║   ███████╗    ██║     ╚██████╔╝██║  ██║╚██████╗███████╗    ███████║   ██║   ███████║   ██║   ███████╗██║ ╚═╝ ██║
╚═════╝ ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚══════╝    ╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝╚══════╝    ╚══════╝   ╚═╝   ╚══════╝   ╚═╝   ╚══════╝╚═╝     ╚═╝
`

		green := color.New(color.FgHiGreen).SprintFunc()
		greeInput := color.New(color.FgGreen).SprintFunc()
		tm.Printf("%s\n\n", green(logo))

		yellow := color.New(color.FgYellow).SprintFunc()
		cyan := color.New(color.FgHiCyan).SprintFunc()
		tm.Printf("%s: %s\n", yellow("Author"), cyan("Rafael Lucena"))
		tm.Printf("%s: %s\n", yellow("Email"), cyan("rafael.lucena@gmail.com"))
		tm.Printf("%s: %s\n", yellow("Versão"), cyan("0.1"))
		tm.Printf("\n\n\n")

		tm.Flush()

		tm.Printf("%s", greeInput("Digite/cole o token MD5 (Conteudo Encriptado/HASH): "))
		tm.Flush()
		fmt.Scanln(&inputModel.Md5BytesHopeString)

		tm.Printf("%s", greeInput("Digite/cole o salt, caso exista: "))
		tm.Flush()
		fmt.Scanln(&inputModel.Salt)

		tm.Printf("%s", greeInput("Digite/cole o alfabeto: "))
		tm.Flush()
		fmt.Scanln(&inputModel.Alfabeto)

		tm.Printf("%s", greeInput("Digite/cole o min: "))
		tm.Flush()
		fmt.Scanln(&inputModel.MinCaractere)

		tm.Printf("%s", greeInput("Digite/cole o max: "))
		tm.Flush()
		fmt.Scanln(&inputModel.MaxCaractere)

		tm.Printf("%s", greeInput("Digite/cole o total de jobs paralelos: "))
		tm.Flush()
		fmt.Scanln(&inputModel.TotalJobsParalelo)

		clearPartialScreen(1, 16)

		BruteForceParameters := tm.NewTable(0, 10, 5, ' ', 0)
		tm.Printf("\n%s\n\n", greeInput("Parâmetros para descobrir o password"))
		tm.Flush()
		fmt.Fprintf(BruteForceParameters, "Alfabeto\tMin - Max\tTheads\n")
		fmt.Fprintf(BruteForceParameters, "%s\t%d - %d\t%d\n", inputModel.Alfabeto, inputModel.MinCaractere, inputModel.MaxCaractere, inputModel.TotalJobsParalelo)
		tm.Println(BruteForceParameters)
		tm.Flush()

		yellowBold := color.New(color.FgYellow).Add(color.Bold).SprintFunc()
		greenBold := color.New(color.FgGreen).Add(color.Bold).SprintFunc()
		tm.Printf("\n%s %s\n\n", yellowBold("Processo Iniciado..."), greenBold(" Good Luck! "))
		tm.Flush()
		clearPartialScreen(1, 28)

		var (
			semaphoreSize = inputModel.TotalJobsParalelo
		)
		sem := make(chan struct{}, semaphoreSize)
		var wg sync.WaitGroup

		TamanhoCargaTrabalho := 100000000
		contador := inputModel.MinCaractere
		//arrayInicial := make([]int, contador)
		arrayAtual := make([]int, contador)
		arrayFinal := make([]int, inputModel.MaxCaractere)
		ultimoIndiceCilindro := len(inputModel.Alfabeto) - 1
		for i := 0; i < len(arrayFinal); i++ {
			arrayFinal[i] = ultimoIndiceCilindro
		}
		payload := inputModel.Salt + "##hash##"
		chQuit, chProgresso := make(chan struct{}), make(chan models.ProgressoGeneric)

		jobCounter := float64(0)

		go func() {
			for !bfworker.IntsEquals(arrayAtual, arrayFinal) {

				arraySegmento, alcancouMaximo := bfworker.RotacionarArraySafe(arrayAtual, TamanhoCargaTrabalho-1, ultimoIndiceCilindro)
				//fmt.Printf("%v - %v \n", arrayAtual, arrayFinal)
				arrayInicialCopy := make([]int, len(arrayAtual))
				copy(arrayInicialCopy, arrayAtual)

				arrayAtualCopy := make([]int, len(arrayAtual))
				copy(arrayAtualCopy, arrayAtual)

				arraySegmentoCopy := make([]int, len(arraySegmento))
				copy(arraySegmentoCopy, arraySegmento)
				var hash [16]byte
				slice, _ := hex.DecodeString(inputModel.Md5BytesHopeString)
				copy(hash[:], slice[:16])

				job := models.JobGeneric{
					Alfabeto:     inputModel.Alfabeto,
					ArrayInicial: arrayInicialCopy,
					ArrayAtual:   arrayAtualCopy,
					ArrayFinal:   arraySegmentoCopy,
					Payload:      payload,
					NumberJob:    math.Floor(math.Mod(jobCounter, float64(inputModel.TotalJobsParalelo))),
					Md5BytesHope: hash,
				}
				sem <- struct{}{}
				wg.Add(1)
				jobCounter = jobCounter + 1
				go func() {
					//fmt.Printf("\n--------aqui2\n")
					bfworker.StartBruteForce(job, chQuit, chProgresso)

					//fmt.Printf("\n--------aqui2\n")
					// release semaphore
					<-sem
					wg.Done()
				}()
				if alcancouMaximo {
					if !bfworker.IntsEquals(arrayAtual, arrayFinal) && contador < inputModel.MaxCaractere {
						contador = contador + 1
						arrayAtual = make([]int, contador)
					} else if contador == inputModel.MaxCaractere {
						arrayAtual = arrayFinal
					}
				} else {
					arraySegmentoCopy := make([]int, len(arraySegmento))
					copy(arraySegmentoCopy, arraySegmento)
					arrayAtual, _ = bfworker.RotacionarArraySafe(arraySegmentoCopy, 1, ultimoIndiceCilindro)
				}

			}
			wg.Wait()
			if !passwordFound {
				tm.Printf("\n\nPassword Not Found\n\n")
				tm.Flush()
				close(chQuit)
			}

		}()

		JobsAtivos := make([]models.ProgressoGeneric, inputModel.TotalJobsParalelo)
		tm.Printf("\n\nMonitorando\n\n")
		tm.Flush()

	monitor:
		for {

			select {
			case <-chQuit:
				break monitor
			case prog := <-chProgresso:

				if JobsAtivos[int(prog.NumberJob)].Status != "Password Found" {
					passwordFound = true
					JobsAtivos[int(prog.NumberJob)] = prog
				}
				if JobsAtivos[int(prog.NumberJob)].Status == "Password Found" {
					passwordFound = true
					clearPartialScreen(1, 34)
					tokenInfoTable := tm.NewTable(0, 15, 5, ' ', 0)
					fmt.Fprintf(tokenInfoTable, "JOB\tK/S\tTotal Chaves\tTotal Testadas\tProgress\tUltimo Teste\tPassword Found\tStatus\n")
					fmt.Fprintf(tokenInfoTable, "%f\t%s\t%s\t%s\t%f %%\t%s\t%s\t%s\n",
						prog.NumberJob,
						hmz.FormatFloat("#,###", prog.KeysPerSecond),
						hmz.FormatFloat("#,###", prog.TotalChaves),
						hmz.FormatFloat("#,###", prog.TotalTestada),
						prog.Progresso*100,
						prog.UltimaPalavra,
						prog.PasswordFound,
						prog.Status)
					tm.Println(tokenInfoTable)
					tm.Flush()
					//close(chProgresso)
				}
				if podeEscreverProgresso(JobsAtivos) {
					clearPartialScreen(1, 34)
					tokenInfoTable := tm.NewTable(0, 15, 5, ' ', 0)
					fmt.Fprintf(tokenInfoTable, "JOB\tK/S\tTotal Chaves\tTotal Testadas\tProgress\tUltimo Teste\tPassword Found\tStatus\n")
					for i := 0; i < len(JobsAtivos); i++ {
						jobtmp := JobsAtivos[i]
						fmt.Fprintf(tokenInfoTable, "%f\t%s\t%s\t%s\t%f %%\t%s\t%s\t%s\n",
							jobtmp.NumberJob,
							hmz.FormatFloat("#,###", jobtmp.KeysPerSecond),
							hmz.FormatFloat("#,###", jobtmp.TotalChaves),
							hmz.FormatFloat("#,###", jobtmp.TotalTestada),
							jobtmp.Progresso*100,
							jobtmp.UltimaPalavra,
							jobtmp.PasswordFound,
							jobtmp.Status)

					}
					tm.Println(tokenInfoTable)
					tm.Flush()
				}

			}
		}
		//close(chProgresso)
		break
	}
}
