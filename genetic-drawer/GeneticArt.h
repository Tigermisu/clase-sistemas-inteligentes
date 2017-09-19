#pragma once
#include <SFML/Graphics.hpp>
#include <opencv2/core.hpp>
#include <iostream>
#include <stdlib.h>
#include <algorithm> 
#include <vector>
#include <time.h>

using namespace cv;
using namespace std;

Mat createRandomImage(int height, int width) {
	/*
	Mat image(height, width, CV_8UC3);

	for (int i = 0; i < height; i++) {
		char* row = image.ptr<char>(i);
		for (int j = 0; j < width * 3; j++)	{
			row[j] = rand() % 256;
		}
	}
	*/

	Mat image(height, width, CV_8UC3);
	randu(image, 0, 255);

	return image;
}

vector<Mat> generateBeautyFromNoise(int populationSize, int imageHeight, int imageWidth) {
	vector<Mat> v = {};
	
	srand(time(NULL)); // Seed the RNG

	for (int i = 0; i < populationSize; i++) {
		v.push_back(createRandomImage(imageHeight, imageWidth));
	}

	return v;
}

vector<int> evaluateArtisticAppeal(vector<Mat> population) {
	vector<int> v = {};

	for (Mat image: population)	{
		int aptitude = 0,
			height = image.size().height,
			width = image.size().width;
		for (int i = 0; i < height; i++) {
			unsigned char* row = image.ptr<unsigned char>(i);
			/*
			for (int j = 0; j < width * 3; j += 6) {
				aptitude += 10000 / (1 + abs(row[j] - row[j + 3]) + abs(row[j + 1] - row[j + 4]) + abs(row[j + 2] - row[j + 5]));
			}
			*/
			for (int j = 0; j < width * 3; j += 3) {
				aptitude += (row[j] * 2 - (row[j + 2] + row[j + 1])) / 10; // Prefer red pixels
			}
		}
		
		if (aptitude < 0) aptitude = 1;

		v.push_back(aptitude);
	}

	return v;
}

void weaveChildWithEuphoricLove(Mat dad, Mat mom, Mat firstborn, Mat theSecondOne, int row, int col) {
	unsigned char* dadRow = dad.ptr<unsigned char>(row);
	unsigned char* momRow = mom.ptr<unsigned char>(row);
	unsigned char* firstRow = firstborn.ptr<unsigned char>(row);
	unsigned char* secondRow = theSecondOne.ptr<unsigned char>(row);

	for (int j = col * 3; j < col * 3 + 3; j++) {
		firstRow[j] = dadRow[j];
		secondRow[j] = momRow[j];
	}
}

vector<Mat> naturallySelectArt(vector<Mat> population, vector<int> aptitudes) {
	int aptitudeSum = 0,
		populationSize = population.size(),
		imageWidth = population[0].cols,
		imageHeight = population[0].cols,
		imageSize = imageWidth * imageHeight;

	vector<Mat> parents = {};
	vector<Mat> children = {};

	
	for (int aptitude : aptitudes) {
		aptitudeSum += aptitude;
	}

	// Choose parents
	for (int i = 0; i < populationSize; i++) {
		int rouletteNumber = rand() % aptitudeSum,
			runningSum = 0;

		for (int j = 0; j < populationSize; j++) {
			runningSum += aptitudes[j];
			if (rouletteNumber < runningSum) {
				parents.push_back(population[j]);
			}
		}
	}

	// Make couples 
	random_shuffle(parents.begin(), parents.end());

	// Make offspring
	for (int i = 0; i < populationSize; i+=2) {
		int crossPoint = rand() % imageSize;
		Mat firstChild(imageHeight, imageWidth, CV_8UC3),
			secondChild(imageHeight, imageWidth, CV_8UC3);
		for (int j = 0; j < imageSize; j++) {
			int row = j / imageWidth,
				column = j % imageWidth;
			if (j > crossPoint) {
				weaveChildWithEuphoricLove(parents[i], parents[i+1], firstChild, secondChild, row, column);
			} else {
				weaveChildWithEuphoricLove(parents[i], parents[i+1], secondChild, firstChild, row, column);
			}
		}

		children.push_back(firstChild);
		children.push_back(secondChild);
	}

	// There is no mutation in art (yet)

	return children;
}

Mat extractOpusMagnum(vector<Mat> population) {
	vector<int> aptitudes = evaluateArtisticAppeal(population);
	int maxAptitude = INT32_MIN;
	Mat mostFit;

	for (int i = 0; i < population.size(); i++) {
		if (aptitudes[i] > maxAptitude) {
			maxAptitude = aptitudes[i];
			mostFit = population[i];
		}
	}

	return mostFit;
}