#include <opencv2/core.hpp>
#include <opencv2/imgcodecs.hpp>
#include <opencv2/highgui.hpp>
#include <opencv2/imgproc.hpp>
#include <iostream>
#include <vector>

#include "GeneticArt.h"
#include "ArtPiece.h"



using namespace cv;
using namespace std;

void drawProgressBar(int progress, int limit) {
	string progressBar = "   [";
	for (char i = 0; i < 20; i++) {
		if (i < progress * 20 / limit) {
			progressBar += "#";
		} else {
			progressBar += " ";
		}
	}
	progressBar += "]";
	cout << progressBar << " (" << progress << "/" << limit << ")\r";
}

int main(int argc, char** argv) {
	int populationCount = 50,
		maxGenerations = 500;
	if (argc != 3) {
		cout << "Using default parameters: 50 art pieces with 500 generations" << endl;
	} else {
		populationCount = atoi(argv[1]);
		maxGenerations = atoi(argv[2]);
	}

	vector<ArtPiece> population = generateBeautyFromNoise(populationCount);

	Mat initialSample = population[0].toImage();

	evaluateArtisticAppeal(population);

	int avgAptitude = 0;

	for (ArtPiece p : population) {
		avgAptitude += p.aptitude;
	}

	cout << "Initial average aptitude: " << (avgAptitude / population.size()) << endl;

	for (int i = 0; i <= maxGenerations; i++) {
		evaluateArtisticAppeal(population);

		population = naturallySelectArt(population);

		evaluateArtisticAppeal(population);			

		drawProgressBar(i, maxGenerations);
	}


	evaluateArtisticAppeal(population);


	avgAptitude = 0;

	for (ArtPiece p : population) {
		avgAptitude += p.aptitude;
	}

	cout << "\nFinal average aptitude: " << (avgAptitude / population.size()) << endl;

	ArtPiece relic = extractOpusMagnum(population);
	imshow("Ideal", ArtPiece::getIdealImage());
	imshow("Initial Sample", initialSample);
	imshow("Opus Magnum", relic.toImage());
	waitKey(0);
}