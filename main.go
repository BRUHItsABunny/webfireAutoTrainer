package main

import (
	"context"
	"fmt"
	gokhttp "github.com/BRUHItsABunny/gOkHttp"
	"github.com/BRUHItsABunny/webfireAutoTrainer/browser"
	"github.com/BRUHItsABunny/webfireAutoTrainer/client"
	"github.com/davecgh/go-spew/spew"
	"time"
)

func main() {
	chromeBrowser := browser.NewBrowser(context.Background())
	defer chromeBrowser.Close()

	err := chromeBrowser.LoginScreen()
	if err != nil {
		fmt.Println("Failed to visit the WebFireTraining site")
		panic(fmt.Errorf("chromeBrowser.LoginScreen: %w", err))
	}
	fmt.Println("Waiting until you have logged in...")

	err = chromeBrowser.AwaitLoginHomeScreen()
	if err != nil {
		fmt.Println("Failed to log in")
		panic(fmt.Errorf("chromeBrowser.AwaitLoginHomeScreen: %w", err))
	}
	fmt.Println("Logged in, scraping needed values for each course and class...")

	courses, err := chromeBrowser.GetCourses()
	if err != nil {
		fmt.Println("Failed to scrape courses and classes")
		panic(fmt.Errorf("chromeBrowser.GetCourses: %w", err))
	}
	fmt.Println("Scraped the needed data, copying cookie data...")

	laravelSession, xsrfToken, err := chromeBrowser.GetCookieData()
	if err != nil {
		fmt.Println("Failed to scrape courses and classes")
		panic(fmt.Errorf("chromeBrowser.GetCourses: %w", err))
	}
	fmt.Println("Preparation done, start auto training")

	hClient, err := gokhttp.NewHTTPClient()
	if err != nil {
		panic(err)
	}
	wftClient := client.NewWFTClient(hClient, laravelSession, xsrfToken)
	// Complete each class and the exam with 100% score
	for _, course := range courses {
		exam := course.Classes["EXAM"]
		for className, courseClass := range course.Classes {
			fmt.Println(spew.Sdump(courseClass))
			if className == exam.Name || courseClass.Completed {
				continue
			}
			res, err := wftClient.StartClass(context.Background(), courseClass)
			if err != nil {
				panic(err)
			}
			fmt.Println(fmt.Sprintf("Started %s - %s: %s", course.Name, courseClass.Name, string(res)))
			time.Sleep(time.Second)
			res, err = wftClient.InitClass(context.Background(), courseClass)
			if err != nil {
				panic(err)
			}
			fmt.Println(fmt.Sprintf("Initialized %s - %s: %s", course.Name, courseClass.Name, string(res)))
			time.Sleep(time.Second)
			res, err = wftClient.FinishClass(context.Background(), courseClass)
			if err != nil {
				panic(err)
			}
			fmt.Println(fmt.Sprintf("Finished %s - %s: %s", course.Name, courseClass.Name, string(res)))
			time.Sleep(time.Second)
			res, err = wftClient.ExitClass(context.Background(), courseClass)
			if err != nil {
				panic(err)
			}
			fmt.Println(fmt.Sprintf("Exited %s - %s: %s", course.Name, courseClass.Name, string(res)))
			time.Sleep(time.Second)
		}
		fmt.Println(spew.Sdump(exam))
		if exam.Completed {
			continue
		}
		res, err := wftClient.StartClass(context.Background(), exam)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Started %s - %s: %s", course.Name, exam.Name, string(res)))
		time.Sleep(time.Second)
		res, err = wftClient.InitClass(context.Background(), exam)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Initialized %s - %s: %s", course.Name, exam.Name, string(res)))
		time.Sleep(time.Second)
		res, err = wftClient.FinishExam(context.Background(), exam, 100.00)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Finished %s - %s: %s", course.Name, exam.Name, string(res)))
		time.Sleep(time.Second)
		res, err = wftClient.ExitClass(context.Background(), exam)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Exited %s - %s: %s", course.Name, exam.Name, string(res)))
		time.Sleep(time.Second)
	}
	fmt.Println("done")
}
