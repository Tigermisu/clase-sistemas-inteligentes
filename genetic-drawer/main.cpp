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

long getAverageAptitude(vector<ArtPiece> pop) {
	long avgAptitude = 0;

	evaluateArtisticAppeal(pop);

	for (ArtPiece p : pop) {
		avgAptitude += p.aptitude;
	}

	return avgAptitude / pop.size();
	
}

int main(int argc, char** argv) {
	int populationCount = 50,
		maxGenerations = 500;
	if (argc != 3) {
		cout << "Usage: genetic-drawer [populationCount] [generations]" << endl;
		cout << "Using default parameters: 50 art pieces with 500 generations\n" << endl;
	} else {
		populationCount = atoi(argv[1]);
		maxGenerations = atoi(argv[2]);
	}

	srand(time(NULL)); // Seed the RNG

	vector<ArtPiece> populationRed = generateBeautyFromNoise(populationCount),
		populationGreen = generateBeautyFromNoise(populationCount),
		populationBlue = generateBeautyFromNoise(populationCount);

	Mat initialSampleRed = populationRed[0].toImage(),
		initialSampleGreen = populationGreen[0].toImage(),
		initialSampleBlue = populationBlue[0].toImage(),
		initialSampleComposite = weaveMasterpiece(initialSampleRed, initialSampleGreen, initialSampleBlue, false);

	cout << "Initial average aptitude for red layer: " << getAverageAptitude(populationRed) << endl;
	cout << "Initial average aptitude for green layer: " << getAverageAptitude(populationGreen) << endl;
	cout << "Initial average aptitude for blue layer: " << getAverageAptitude(populationBlue) << endl;

	for (int i = 1; i <= maxGenerations; i++) {
		evaluateArtisticAppeal(populationRed);
		populationRed = naturallySelectArt(populationRed);

		evaluateArtisticAppeal(populationGreen);
		populationGreen = naturallySelectArt(populationGreen);

		evaluateArtisticAppeal(populationBlue);
		populationBlue = naturallySelectArt(populationBlue);
		
		drawProgressBar(i, maxGenerations);
	}

	cout << "\nFinal average aptitude for red layer: " << getAverageAptitude(populationRed) << endl;
	cout << "Final average aptitude for green layer: " << getAverageAptitude(populationGreen) << endl;
	cout << "Final average aptitude for blue layer: " << getAverageAptitude(populationBlue) << endl;

	ArtPiece redRelic = extractOpusMagnum(populationRed),
		greenRelic = extractOpusMagnum(populationGreen),
		blueRelic = extractOpusMagnum(populationBlue);

	Mat masterpiece = weaveMasterpiece(redRelic.toImage(), greenRelic.toImage(), blueRelic.toImage(), true);

	imshow("Ideal", ArtPiece::getIdealImage());
	imshow("Initial Sample", initialSampleComposite);
	imshow("Opus Magnum", masterpiece);
	waitKey(0);
}