#pragma once
#include <opencv2/core.hpp>
#include <iostream>
#include <stdlib.h>
#include <algorithm> 
#include <vector>
#include <time.h>

#include "ArtPiece.h"

using namespace cv;
using namespace std;

vector<ArtPiece> generateBeautyFromNoise(int populationSize) {
	vector<ArtPiece> v = {};

	for (int i = 0; i < populationSize; i++) {
		ArtPiece p{};
		for (int j = 0; j < p.geneLength; j++) {
			p.genes[j] = (rand() % 2) == 1 ? true : false;
		}

		v.push_back(p);
	}

	return v;
}

void evaluateArtisticAppeal(vector<ArtPiece> &population) {
	int dimensions = sqrt(ArtPiece::geneLength),
		maxDistance = 3 * dimensions / 10;

	for (ArtPiece &p : population) {
		long aptitude = 1;
		for (int i = 0; i < p.geneLength; i++) {
			int row = i / dimensions,
				col = i % dimensions,
				distance = sqrt(pow(abs(row - dimensions / 2), 2) + pow(abs(col - dimensions / 2), 2));

			if (distance < maxDistance) {
				aptitude += (p.genes[i] ? 3 : -2);
			} else {
				aptitude += (p.genes[i] ? -3 : 3);
			}
		}
		p.aptitude = aptitude;
	}
}

vector<ArtPiece> naturallySelectArt(vector<ArtPiece> population) {
	int populationSize = population.size(),
		geneLength = population[0].geneLength;


	long aptitudeSum = 0;

	vector<ArtPiece> parents = {};
	vector<ArtPiece> children = {};

	// Calculate aptitude sum
	for (ArtPiece p : population) {
		aptitudeSum += p.aptitude;
	}

	// Choose parents
	for (int i = 0; i < populationSize; i++) {
		long rouletteNumber = rand() % aptitudeSum,
			runningSum = 0;

		for (ArtPiece p: population) {
			runningSum += p.aptitude;
			if (rouletteNumber < runningSum) {
				parents.push_back(p);
				break;
			}
		}
	}

	// Make couples 
	random_shuffle(parents.begin(), parents.end());

	// Make offspring
	for (int i = 0; i < populationSize; i+=2) {
		int crossPoint = rand() % (geneLength - 1);
		ArtPiece firstChild{},
			secondChild{};
		for (int j = 0; j < geneLength; j+= 1) {
			if (j <= crossPoint) {
				firstChild.genes[j] = parents[i].genes[j];
				secondChild.genes[j] = parents[i + 1].genes[j];
			} else {
				firstChild.genes[j] = parents[i + 1].genes[j];
				secondChild.genes[j] = parents[i].genes[j];
			}
		}

		children.push_back(firstChild);
		children.push_back(secondChild);
	}

	// Mutation
	// Make offspring
	for (ArtPiece &p: children) {
		if (rand() % 100 < 5) {
			int mutatedGenes = rand() % (geneLength / 10);
			for (int i = 0; i < mutatedGenes; i++) {
				int rdnGene = rand() % geneLength;
				p.genes[rdnGene] = !p.genes[rdnGene];
			}
		}
	}

	return children;
}

ArtPiece extractOpusMagnum(vector<ArtPiece> population) {
	evaluateArtisticAppeal(population);
	long maxAptitude = INT32_MIN;
	ArtPiece mostFit;

	for (ArtPiece p : population) {
		if (p.aptitude > maxAptitude) {
			maxAptitude = p.aptitude;
			mostFit = p;
		}
	}

	return mostFit;
}

Mat weaveMasterpiece(Mat red, Mat green, Mat blue, int postProcessingIntensity) {
	Mat masterpiece;
	std::vector<Mat> channels;

	channels.push_back(blue);
	channels.push_back(green);
	channels.push_back(red);


	merge(channels, masterpiece);

	if (postProcessingIntensity > 0) {
		if (postProcessingIntensity > 39) postProcessingIntensity = 39;
		else if (postProcessingIntensity % 2 == 0) postProcessingIntensity--;
		medianBlur(masterpiece, masterpiece, postProcessingIntensity);
	}

	return masterpiece;
}

vector<ArtPiece> bloodyCombatToDeath(vector<ArtPiece> champions, vector<ArtPiece> challengers) {
	vector<ArtPiece> contestants, combatants, victors;
	long aptitudeSum = 0;
	int combatantsSize = champions.size();

	// Prepare combatants
	for (int i = 0; i < champions.size(); i++) {
		contestants.push_back(champions[0]);
		contestants.push_back(challengers[0]);
	}

	evaluateArtisticAppeal(contestants);

	for (ArtPiece c : contestants) {
		aptitudeSum += c.aptitude;
	}

	// Choose the 'strongest' n combatants out of 2n
	for (int i = 0; i < combatantsSize; i++) {
		long rouletteNumber = rand() % aptitudeSum,
			runningSum = 0;

		for (ArtPiece c : contestants) {
			runningSum += c.aptitude;
			if (rouletteNumber < runningSum) {
				combatants.push_back(c);
				break;
			}
		}
	}

	// Tournament to death!
	for (int i = 0; i < combatantsSize; i++) {
		if (combatants[i].aptitude > combatants[(i + 1) % combatantsSize].aptitude) {
			victors.push_back(combatants[i]);
		} else {
			victors.push_back(combatants[(i + 1) % combatantsSize]);
		}
	}

	return victors;
}

// These are references just to avoid the overhead of making copies of these huge Art Pieces.
vector<vector<ArtPiece>> introduceDistantWorlds(vector<ArtPiece> &r, vector<ArtPiece> &g, vector<ArtPiece> &b,
												vector<ArtPiece> &rAlt, vector<ArtPiece> &gAlt, vector<ArtPiece> &bAlt) {
	vector<vector<ArtPiece>> braveNewWorld;

	braveNewWorld.push_back(bloodyCombatToDeath(b, bAlt));
	braveNewWorld.push_back(bloodyCombatToDeath(g, gAlt));
	braveNewWorld.push_back(bloodyCombatToDeath(r, rAlt));

	return braveNewWorld;
}