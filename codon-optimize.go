package main

import (
	"fmt"
	"os"

	"github.com/Open-Science-Global/poly/io/fasta"
	"github.com/Open-Science-Global/poly/transform/codon"
	"github.com/spf13/cobra"
)

func main() {
	Execute()
}

var (
	input      string
	output     string
	codonTable string
)

var rootCmd = &cobra.Command{
	Use:   "codon-optimize",
	Short: "A github action to codon optimize a list of sequences.",
	Long:  "A github action to codon optimize a list of sequences.",
	Run: func(cmd *cobra.Command, args []string) {
		Script(input, output, codonTable)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "A fasta file with all the sequences that will be read")
	rootCmd.PersistentFlags().StringVarP(&output, "ouput", "o", "", "Path and file name where an output fasta file will be written")
	rootCmd.PersistentFlags().StringVarP(&codonTable, "codonTable", "t", "", "A codon table in a JSON file format")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("ouput")
	rootCmd.MarkFlagRequired("codonTable")
}

func Script(input string, outputPath string, codonTablePath string) {
	sequences := fasta.Read(input)
	codonTable := codon.ReadCodonJSON(codonTablePath)

	var optimizedFasta []fasta.Fasta
	for _, sequence := range sequences {
		optimized := CodonOptimization(sequence.Sequence, codonTable)
		optimizedFasta = append(optimizedFasta, fasta.Fasta{sequence.Name, optimized})
	}
	fmt.Println(optimizedFasta)
	fasta.Write(optimizedFasta, outputPath)
}

// run this to set the Codon Optimization function
func CodonOptimization(enzymeSequence string, codonTable codon.Table) string {

	// Optimize sequence using the protein sequence and codon table
	optimizedSequence, _ := codon.Optimize(enzymeSequence, codonTable)

	// Lets check if the codon optimization actually works by making some checks:
	// First one is if both codon sequences are different
	if optimizedSequence == enzymeSequence {
		fmt.Println("ERROR: Source sequence and optimized sequence are identical. Check your inputs and run again.")
		os.Exit(0)
	}

	// Check if both translated sequences are equal
	protein, _ := codon.Translate(optimizedSequence, codon.GetCodonTable(11))
	if protein != enzymeSequence {
		fmt.Println("ERROR: Protein sequences are not identical. Codon optimization does not change amino acid values.")
		os.Exit(0)
	}
	return optimizedSequence
}
