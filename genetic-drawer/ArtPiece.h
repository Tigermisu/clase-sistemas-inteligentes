#pragma once
#include <opencv2\core.hpp>

using namespace cv;

class ArtPiece
{
public:
	ArtPiece();
	static const int geneLength = 2500;
	long aptitude;
	bool genes[geneLength];

	Mat toImage(unsigned char);
	static Mat getIdealImage(void);

};

Mat ArtPiece::getIdealImage() {
	int dimensions = sqrt(geneLength);
	Mat img(dimensions, dimensions, CV_8UC1);

	int maxDistance = 3 * dimensions / 10;

	for (int i = 0; i < dimensions; i++) {
		unsigned char *imgRow = img.ptr<unsigned char>(i);
		int iDist = pow(abs(i - dimensions / 2), 2);
		for (int j = 0; j < dimensions; j++) {
			int distance = sqrt(pow(abs(j - dimensions / 2), 2) + iDist);
			if (distance < maxDistance) {
				imgRow[j] = 255;
			} else {
				imgRow[j] = 0;
			}
		}
	}

	resize(img, img, Size(500, 500), 0, 0);

	return img;	
}

Mat ArtPiece::toImage(unsigned char colorIntensity) {
	int dimensions = sqrt(geneLength);
	Mat img(dimensions, dimensions, CV_8UC1);

	for (int i = 0; i < dimensions; i++) {
		unsigned char *imgRow = img.ptr<unsigned char>(i);
		for (int j = 0; j < dimensions; j++) {
			imgRow[j] = genes[i * dimensions + j] ? colorIntensity : 0;
		}
	}

	resize(img, img, Size(500, 500), 0, 0);

	return img;
}

ArtPiece::ArtPiece() {
	aptitude = 0;
}
