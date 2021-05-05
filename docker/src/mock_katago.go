package main

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type Runner struct {
	stop bool
}

func (runner *Runner) ShowInfo() {
	os.Stderr.WriteString("This is mock katago.\n")
	os.Stderr.WriteString("Ready to accept gpt commands\n")
}

func (runner *Runner) runGetRules() {
	os.Stdout.WriteString("= ")
	os.Stdout.WriteString("{\"friendlyPassOk\":false,\"hasButton\":false,\"ko\":\"POSITIONAL\",\"scoring\":\"AREA\",\"suicide\":true,\"tax\":\"NONE\",\"whiteHandicapBonus\":\"0\"}\n")

}
func (runner *Runner) runAnalyze() {
	os.Stdout.WriteString("=\n")
	runner.stop = false
	for {
		if runner.stop {
			break
		}
		os.Stdout.WriteString("info move D4 visits 12 utility -0.209635 winrate 0.398815 scoreMean -0.820066 scoreStdev 17.2358 scoreLead -0.820066 scoreSelfplay -1.27913 prior 0.163632 lcb 0.384575 utilityLcb -0.249509 order 1 pv D4 Q16 D16 Q4 info move Q16 visits 11 utility -0.199381 winrate 0.40231 scoreMean -0.825123 scoreStdev 17.465 scoreLead -0.825123 scoreSelfplay -1.25644 prior 0.160565 lcb 0.393602 utilityLcb -0.223766 order 0 pv Q16 D16 C17 info move Q4 visits 10 utility -0.199752 winrate 0.401927 scoreMean -0.822762 scoreStdev 17.4283 scoreLead -0.822762 scoreSelfplay -1.25834 prior 0.15481 lcb 0.389543 utilityLcb -0.234427 order 2 pv Q4 Q16 D16 info move D16 visits 8 utility -0.208521 winrate 0.39795 scoreMean -0.881401 scoreStdev 17.3401 scoreLead -0.881401 scoreSelfplay -1.34429 prior 0.154493 lcb 0.383158 utilityLcb -0.24994 order 3 pv D16 Q4 D4 Q16 info move D17 visits 7 utility -0.214748 winrate 0.396242 scoreMean -0.914599 scoreStdev 17.4673 scoreLead -0.914599 scoreSelfplay -1.40377 prior 0.0326206 lcb 0.366583 utilityLcb -0.297793 order 4 pv D17 D4 Q16 Q4 info move D3 visits 7 utility -0.211478 winrate 0.396986 scoreMean -0.879298 scoreStdev 17.4636 scoreLead -0.879298 scoreSelfplay -1.34756 prior 0.0318391 lcb 0.368854 utilityLcb -0.290248 order 5 pv D3 D16 Q4 Q16 R17 info move C4 visits 7 utility -0.210888 winrate 0.395878 scoreMean -0.895391 scoreStdev 17.4358 scoreLead -0.895391 scoreSelfplay -1.39123 prior 0.0313944 lcb 0.368486 utilityLcb -0.287586 order 6 pv C4 Q4 D16 Q16 R17 info move Q3 visits 7 utility -0.207753 winrate 0.3975 scoreMean -0.874237 scoreStdev 17.4234 scoreLead -0.874237 scoreSelfplay -1.35895 prior 0.0311447 lcb 0.373249 utilityLcb -0.275655 order 7 pv Q3 Q16 D4 D16 C17 info move R16 visits 7 utility -0.206759 winrate 0.398851 scoreMean -0.859451 scoreStdev 17.4381 scoreLead -0.859451 scoreSelfplay -1.32023 prior 0.0295987 lcb 0.379029 utilityLcb -0.262263 order 8 pv R16 D16 Q4 D4 C3 info move C16 visits 7 utility -0.208346 winrate 0.398106 scoreMean -0.903949 scoreStdev 17.4617 scoreLead -0.903949 scoreSelfplay -1.35086 prior 0.0267378 lcb 0.348045 utilityLcb -0.35355 order 9 pv C16 Q16 D4 Q4 info move Q17 visits 6 utility -0.217281 winrate 0.393785 scoreMean -0.926662 scoreStdev 17.4145 scoreLead -0.926662 scoreSelfplay -1.42341 prior 0.0348316 lcb 0.359717 utilityLcb -0.312672 order 10 pv Q17 Q4 D16 D4 C3 info move R4 visits 6 utility -0.208199 winrate 0.39871 scoreMean -0.857564 scoreStdev 17.4445 scoreLead -0.857564 scoreSelfplay -1.33299 prior 0.0345577 lcb 0.346534 utilityLcb -0.35429 order 11 pv R4 D4 Q16 D16 info move D5 visits 5 utility -0.249156 winrate 0.38421 scoreMean -1.09883 scoreStdev 17.3595 scoreLead -1.09883 scoreSelfplay -1.63614 prior 0.00857516 lcb 0.314687 utilityLcb -0.443821 order 12 pv D5 Q4 D16 info move E16 visits 5 utility -0.23331 winrate 0.387685 scoreMean -1.10023 scoreStdev 17.2901 scoreLead -1.10023 scoreSelfplay -1.61906 prior 0.00803781 lcb 0.242894 utilityLcb -0.638723 order 13 pv E16 D4 C16 info move P4 visits 4 utility -0.249219 winrate 0.383006 scoreMean -1.1422 scoreStdev 17.3152 scoreLead -1.1422 scoreSelfplay -1.70275 prior 0.00808218 lcb 0.195671 utilityLcb -0.773757 order 14 pv P4 Q16 D4info\n")
		time.Sleep(time.Second)
	}
}

func (runner *Runner) runGenmove() {
	os.Stderr.WriteString("CHAT:Visits 5733 Winrate 39.62% ScoreLead -0.9 ScoreStdev 17.4 PV D16 Q4 Q16 D4 C3 C4 D3 E4 F2 R17 Q17 R16 R14\n")
	os.Stdout.WriteString("= D16\n")
}
func (runner *Runner) runOthers() {
	os.Stdout.WriteString("=\n")
}
func (runner *Runner) RunCmd(cmd string) {
	runner.Stop()
	if strings.Contains(cmd, "analyze") {
		runner.runAnalyze()
	} else if strings.Contains(cmd, "rule") {
		runner.runGetRules()
	} else if strings.Contains(cmd, "genmove") {
		runner.runGenmove()
	} else if strings.Contains(cmd, "quit") {
		os.Exit(0)
	} else {
		runner.runOthers()
	}
}
func (runner *Runner) Stop() {
	runner.stop = true
}
func main() {
	runner := Runner{}
	runner.ShowInfo()
	go func() {
		for {
			var reader = bufio.NewReader(os.Stdin)
			message, err := reader.ReadString('\n')
			if err != nil {
				// log.Printf("ERROR failed to read from stdin. %+v\n", err)
				os.Exit(0)
			}
			go runner.RunCmd(message)
		}

	}()
	for {
		time.Sleep(time.Second)
	}
}
