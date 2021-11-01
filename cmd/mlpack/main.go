package main

import (
	"fmt"

	"mlpack.org/v1/mlpack"
)

func main() {
	// Download dataset.
	mlpack.DownloadFile("https://www.mlpack.org/datasets/covertype-small.data.csv.gz",
		"data.csv.gz")
	mlpack.DownloadFile("https://www.mlpack.org/datasets/covertype-small.labels.csv.gz",
		"labels.csv.gz")
	// Extract/Unzip the dataset.
	mlpack.UnZip("data.csv.gz", "data.csv")
	dataset, _ := mlpack.Load("data.csv")
	mlpack.UnZip("labels.csv.gz", "labels.csv")
	labels, _ := mlpack.Load("labels.csv")
	// Split the dataset using mlpack.
	params := mlpack.PreprocessSplitOptions()
	params.InputLabels = labels
	params.TestRatio = 0.3
	params.Verbose = true
	test, test_labels, train, train_labels :=
		mlpack.PreprocessSplit(dataset, params)
	// Train a random forest.
	rf_params := mlpack.RandomForestOptions()
	rf_params.NumTrees = 10
	rf_params.MinimumLeafSize = 3
	rf_params.PrintTrainingAccuracy = true
	rf_params.Training = train
	rf_params.Labels = train_labels
	rf_params.Verbose = true
	rf_model, _, _ := mlpack.RandomForest(rf_params)
	// Predict the labels of the test points.
	rf_params_2 := mlpack.RandomForestOptions()
	rf_params_2.Test = test
	rf_params_2.InputModel = &rf_model
	rf_params_2.Verbose = true
	_, predictions, _ := mlpack.RandomForest(rf_params_2)
	// Now print the accuracy.
	rows, _ := predictions.Dims()
	var sum int = 0
	for i := 0; i < rows; i++ {
		if predictions.At(i, 0) == test_labels.At(i, 0) {
			sum = sum + 1
		}
	}
	fmt.Print(sum, " correct out of ", rows, " (",
		(float64(sum)/float64(rows))*100, "%).\n")
}
