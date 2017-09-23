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

void evolveCommunity(vector<ArtPiece> &r, vector<ArtPiece> &g, vector<ArtPiece> &b, int tgtGenerations, string populationName) {
	cout << "Initial average aptitude for red layer (" << populationName << "): " << getAverageAptitude(r) << endl;
	cout << "Initial average aptitude for green layer (" << populationName << "): " << getAverageAptitude(g) << endl;
	cout << "Initial average aptitude for blue layer (" << populationName << "): " << getAverageAptitude(b) << endl;



	
	for (int i = 1; i <= tgtGenerations; i++) {
		evaluateArtisticAppeal(r);
		r = naturallySelectArt(r);

		evaluateArtisticAppeal(g);
		g = naturallySelectArt(g);

		evaluateArtisticAppeal(b);
		b = naturallySelectArt(b);

		drawProgressBar(i, tgtGenerations);
	}

	cout << "\nFinal average aptitude for red layer (" << populationName << "): " << getAverageAptitude(r) << endl;
	cout << "Final average aptitude for green layer (" << populationName << "): " << getAverageAptitude(g) << endl;
	cout << "Final average aptitude for blue layer (" << populationName << "): " << getAverageAptitude(b) << endl << endl;
}

void appraiseMasterpiece(ArtPiece &r, ArtPiece &g, ArtPiece &b) {
	long aptitude = (r.aptitude + g.aptitude + b.aptitude) / 3;

	cout << "Masterpiece. Aptitude: " << aptitude << endl;
}

int main(int argc, char** argv) {
	int populationCount = 20,
		maxGenerations = 300,
		colorIntensity = 255,
		postProcessingIntensity = 19;

	cout << "Welcome to Chris' Genetic Art Generator.\nThis simple program creates Pollock-style art out of the blue!." << endl;
	cout << "Params: genetic-drawer [populationCount] [generations] [colorIntensity 0-255] [postProcessing 0-39]" << endl << endl;
	
	if (argc == 5) {
		populationCount = atoi(argv[1]);
		maxGenerations = atoi(argv[2]);
		colorIntensity = atoi(argv[3]);
		postProcessingIntensity = atoi(argv[4]);

		colorIntensity > 255 ? colorIntensity = 255 : colorIntensity;
	} else if (argc == 1) {
		cout << "Using default parameters: 20 art pieces with 300 generations,\nfull color intensity (255) and average post-processing (19)" << endl << endl;
	} else {
		cout << "Please specify all parameters or leave blank to use defaults." << endl;
		return 1;
	}

	srand(time(NULL)); // Seed the RNG

	// Generate vector of vectors (format: RGB,RGB,RGB,RGB)
	vector<vector<ArtPiece>> population;
	for (int i = 0; i < 12; i++) {
		population.push_back(generateBeautyFromNoise(populationCount));
	}

	Mat initialSampleRed = population[0][0].toImage(colorIntensity),
		initialSampleGreen = population[0][1].toImage(colorIntensity),
		initialSampleBlue = population[0][2].toImage(colorIntensity),
		initialSampleComposite = weaveMasterpiece(initialSampleRed, initialSampleGreen, initialSampleBlue, false);

	for (int i = 0; i < 12; i+=3) {
		string name;
		switch (i) {
		case 0: name = "Alpha Renaissance"; break;
		case 3: name = "Beta Naturalism"; break;
		case 6: name = "Gamma Romantics"; break;
		case 9: name = "Delta Revolutionaries"; break;
		}
		evolveCommunity(population[i], population[i + 1], population[i + 2], maxGenerations, name);
	}

	vector<vector<ArtPiece>> easternSociety = introduceDistantWorlds(population[0], population[1], population[2],
																	 population[3], population[4], population[5]);
	cout << "Introducing eastern society." << endl << endl;

	evolveCommunity(easternSociety[0], easternSociety[1], easternSociety[2], maxGenerations, "Eastern Society");

	vector<vector<ArtPiece>> westernSociety = introduceDistantWorlds(population[6], population[7], population[8],
																	population[9], population[10], population[11]);

	cout << "Introducing western society." << endl << endl;

	evolveCommunity(westernSociety[0], westernSociety[1], westernSociety[2], maxGenerations, "Western Society");

	vector<vector<ArtPiece>> globalizedWorld = introduceDistantWorlds(easternSociety[0], easternSociety[1], easternSociety[2],
																	westernSociety[0], westernSociety[1], westernSociety[2]);

	cout << "Introducing the globalized world." << endl << endl;

	evolveCommunity(globalizedWorld[0], globalizedWorld[1], globalizedWorld[2], maxGenerations, "Globalized World");

	cout << "Generating masterpiece." << endl << endl;

	ArtPiece redRelic = extractOpusMagnum(globalizedWorld[0]),
		greenRelic = extractOpusMagnum(globalizedWorld[1]),
		blueRelic = extractOpusMagnum(globalizedWorld[2]);

	appraiseMasterpiece(redRelic, greenRelic, blueRelic);

	Mat masterpiece = weaveMasterpiece(redRelic.toImage(colorIntensity), 
		greenRelic.toImage(colorIntensity), blueRelic.toImage(colorIntensity), postProcessingIntensity);

	imshow("Aptitude Function", ArtPiece::getIdealImage());
	imshow("Initial Sample", initialSampleComposite);
	imshow("Opus Magnum", masterpiece);
	waitKey(0);

	return 0;
}

