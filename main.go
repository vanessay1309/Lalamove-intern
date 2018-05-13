//This project serves as a solution to the challenge for Lalamove 2018 Tech Internship programme
package main

import (
	"context"
	"fmt"
	"os"
	"bufio"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

//This funnction returns boolean value in whether the two semvers have same major and minor version
func ifEqual(a *semver.Version, b *semver.Version) bool{
	//Convert semver into int for comparison
	slicedA := (*a).Slice()
	slicedB := (*b).Slice()

	if ((slicedA[0]==slicedB[0]) && (slicedA[1]==slicedB[1])){
		return true
	}else{
		return false
	}
}

//This function returns the index of slice if the major and minor function has been recorded in slice, -1 otherwise
func ifExists(currentSlice []*semver.Version, a *semver.Version) int{
	//If slice is empty return not exist
	if (len(currentSlice)==0){
		return -1
	}else{
		//else loop to check whether the version is recorded
		for i :=0; i<len(currentSlice); i++ {
			if (ifEqual(a,currentSlice[i])){
				return i
			}
		}
		return -1
	}
	return -1
}

//This is a Bubble sort for slice in descending version
func sortSlice(currentSlice []*semver.Version) []*semver.Version {
	for i :=0; i<len(currentSlice)-1; i++ {
		for j :=0; j<len(currentSlice)-i-1; j++ {
		 if ((*currentSlice[j]).LessThan(*currentSlice[j+1])){
				currentSlice[j], currentSlice[j+1] = currentSlice[j+1], currentSlice[j]
			}
		}
	}
	return currentSlice
	}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version

	for i :=0; i<len(releases); i++ {
		//Only include versions that are above minVersion
		if (minVersion.LessThan(*releases[i])){
			//Update Slice
			if ((ifExists(versionSlice, releases[i]))<0){
				//if not exist in slice, append to slice
				versionSlice = append(versionSlice, releases[i])
			}else{
				//if version exists in slice and but is not the latest also append
				index := ifExists(versionSlice, releases[i])
				if ((*versionSlice[index]).LessThan(*releases[i])){
				versionSlice = append(versionSlice, releases[i])
				}
			}
		}
	}

	//Sort slice
	versionSlice = sortSlice(versionSlice);
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please include filename for input")
		return
	}
	file, errF := os.Open(os.Args[1])
	if errF != nil {
		fmt.Println("Please confirm your filename is correct")
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	for scanner.Scan(){
			//For each line, obtain the respository and minVersion
			res := (strings.Split(scanner.Text(), ","))[0]
			res1 := (strings.Split(res, "/"))[0]
			res2 := (strings.Split(res, "/"))[1]
			minVersion :=  semver.New((strings.Split(scanner.Text(), ","))[1])


		// Github
		client := github.NewClient(nil)
		ctx := context.Background()
		opt := &github.ListOptions{PerPage: 10}
		releases, _, err := client.Repositories.ListReleases(ctx, res1, res2, opt)
		if err != nil {
			// make use of Error handling, so that the program can continue (i.e it might be a 404 error for this respository but not for the next one)
			//https://golang.org/doc/effective_go.html#errors
			//panic(err)
			fmt.Println(err)
		}

		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)
		fmt.Printf("latest versions of %s: %s \n",res, versionSlice)
	}
}
