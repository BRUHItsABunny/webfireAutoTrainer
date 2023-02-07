package browser

import (
	"context"
	"errors"
	"fmt"
	"github.com/BRUHItsABunny/webfireAutoTrainer/api"
	"github.com/BRUHItsABunny/webfireAutoTrainer/constants"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"net/url"
	"strings"
)

type WFTBrowser struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc
	tasksCtx    context.Context
	tasksCancel context.CancelFunc
}

func NewBrowser(ctx context.Context) *WFTBrowser {
	result := &WFTBrowser{}
	result.allocCtx, result.allocCancel = chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	result.tasksCtx, result.tasksCancel = chromedp.NewContext(result.allocCtx, chromedp.WithLogf(log.Printf))
	return result
}

func (b *WFTBrowser) Close() {
	b.allocCancel()
	b.tasksCancel()
}

func (b *WFTBrowser) LoginScreen() error {
	err := chromedp.Run(b.tasksCtx, chromedp.Navigate(constants.EndpointLogin))
	if err != nil {
		return fmt.Errorf("chromedp.Run: %w", err)
	}
	return nil
}

func (b *WFTBrowser) AwaitLoginHomeScreen() error {
	err := chromedp.Run(b.tasksCtx,
		// Function waits here for you to login - program doesn't stop until you login or close the browser
		chromedp.WaitVisible(`btn-footer pull-right`, chromedp.BySearch),
	)
	if err != nil {
		return fmt.Errorf("chromedp.Run: %w", err)
	}
	return nil
}

func (b *WFTBrowser) AwaitCoursesVisible() error {
	err := chromedp.Run(b.tasksCtx, chromedp.WaitVisible(`courseList`, chromedp.BySearch))
	if err != nil {
		return fmt.Errorf("chromedp.Run: %w", err)
	}
	return nil
}

func (b *WFTBrowser) getCourseList() (*cdp.Node, error) {
	var result []*cdp.Node
	err := chromedp.Run(b.tasksCtx, chromedp.Nodes(`courseList`, &result, chromedp.BySearch))
	if err != nil {
		return nil, fmt.Errorf("chromedp.Run: %w", err)
	}
	if len(result) < 1 {
		return nil, errors.New("list of courses not found")
	}
	return result[0], nil
}

func (b *WFTBrowser) getCourseNodes() ([]*cdp.Node, error) {
	var courseNodes []*cdp.Node
	courseList, err := b.getCourseList()
	if err != nil {
		return nil, fmt.Errorf("b.getCourseList: %w", err)
	}

	err = chromedp.Run(b.tasksCtx, chromedp.Nodes(`li a`, &courseNodes, chromedp.ByQueryAll, chromedp.FromNode(courseList)))
	if err != nil {
		return nil, fmt.Errorf("chromedp.Run: %w", err)
	}
	return courseNodes, nil
}

func (b *WFTBrowser) visitAndAwaitCourse(courseURL string) error {
	err := chromedp.Run(b.tasksCtx,
		chromedp.Navigate(courseURL),
		chromedp.WaitVisible(`row outlines`, chromedp.BySearch),
	)
	if err != nil {
		return fmt.Errorf("chromedp.Run: %w", err)
	}
	return nil
}

func (b *WFTBrowser) getCourseName() (string, error) {
	var result string
	err := chromedp.Run(b.tasksCtx, chromedp.Text(`heading2 color1`, &result, chromedp.BySearch))
	if err != nil {
		return result, fmt.Errorf("chromedp.Run: %w", err)
	}
	return result, nil
}

func (b *WFTBrowser) getClassesList() (*cdp.Node, error) {
	var result []*cdp.Node
	err := chromedp.Run(b.tasksCtx, chromedp.Nodes(`row outlines`, &result, chromedp.BySearch))
	if err != nil {
		return nil, fmt.Errorf("chromedp.Run: %w", err)
	}
	if len(result) < 1 {
		return nil, errors.New("list of classes not found")
	}
	return result[0], nil
}

func (b *WFTBrowser) getCourseClasses(classesList *cdp.Node) ([]*cdp.Node, error) {
	var result []*cdp.Node
	err := chromedp.Run(b.tasksCtx, chromedp.Nodes(`ul li`, &result, chromedp.ByQueryAll, chromedp.FromNode(classesList)))
	if err != nil {
		return nil, fmt.Errorf("chromedp.Run: %w", err)
	}
	return result, nil
}

func (b *WFTBrowser) parseCourseClass(courseClass *cdp.Node) (*api.WFTClass, error) {
	var (
		classCheckmark []*cdp.Node
		classHref      []*cdp.Node
		className      string
	)
	err := chromedp.Run(b.tasksCtx,
		chromedp.Nodes(`img`, &classCheckmark, chromedp.ByQueryAll, chromedp.FromNode(courseClass), chromedp.AtLeast(0)),
		chromedp.Nodes(`a`, &classHref, chromedp.ByQueryAll, chromedp.FromNode(courseClass)),
		chromedp.Text(`a`, &className, chromedp.ByQueryAll, chromedp.FromNode(courseClass)),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp.Run: %w", err)
	}

	result := &api.WFTClass{
		Name:      className,
		Completed: len(classCheckmark) > 0,
	}

	fakeURL := "https://webfiretraining.com" + strings.Split(classHref[0].AttributeValue("onclick"), "'")[1]
	parsedFakeURL, err := url.Parse(fakeURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}
	fakeUrlSplit := strings.Split(parsedFakeURL.Path, "/")
	result.Referer = fmt.Sprintf("https://webfiretraining.com/webfireadmin/code/content/%s/%s/lms/AICCComm.html", fakeUrlSplit[4], fakeUrlSplit[5])
	result.SessionID = parsedFakeURL.Query().Get("AICC_sid")
	return result, nil
}

func (b *WFTBrowser) parseCourseClasses() (map[string]*api.WFTClass, error) {
	result := map[string]*api.WFTClass{}
	classesList, err := b.getClassesList()
	if err != nil {
		return nil, fmt.Errorf("b.getClassesList: %w", err)
	}

	unparsedCourseClasses, err := b.getCourseClasses(classesList)
	if err != nil {
		return nil, fmt.Errorf("b.getCourseClasses: %w", err)
	}

	for classIdx, unparsedClass := range unparsedCourseClasses {
		parsedClass, err := b.parseCourseClass(unparsedClass)
		if err != nil {
			return nil, fmt.Errorf("b.parseCourseClass[%d]: %w", classIdx, err)
		}
		result[parsedClass.Name] = parsedClass
	}
	return result, nil
}

func (b *WFTBrowser) parseCourse(course *cdp.Node) (*api.WFTCourse, error) {
	result := &api.WFTCourse{URL: course.AttributeValue("href")}
	err := b.visitAndAwaitCourse(result.URL)
	if err != nil {
		return nil, fmt.Errorf("b.visitAndAwaitCourse: %w", err)
	}

	result.Name, err = b.getCourseName()
	if err != nil {
		return nil, fmt.Errorf("b.getCourseName: %w", err)
	}

	result.Classes, err = b.parseCourseClasses()
	if err != nil {
		return nil, fmt.Errorf("b.parseCourseClasses: %w", err)
	}
	return result, nil
}

func (b *WFTBrowser) GetCourses() (map[string]*api.WFTCourse, error) {
	result := map[string]*api.WFTCourse{}

	unparsedCourses, err := b.getCourseNodes()
	if err != nil {
		return result, fmt.Errorf("b.getCourseNodes: %w", err)
	}

	for courseIdx, unparsedCourse := range unparsedCourses {
		parsedCourse, err := b.parseCourse(unparsedCourse)
		if err != nil {
			return result, fmt.Errorf("b.parseCourse[%d]: %w", courseIdx, err)
		}
		result[parsedCourse.Name] = parsedCourse
	}

	return result, nil
}

func (b *WFTBrowser) GetCookieData() (string, string, error) {
	var (
		laravelSession, xsrfToken string
	)
	err := chromedp.Run(b.tasksCtx,
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, cookie := range cookies {
				if cookie.Name == "laravel_session" {
					laravelSession = cookie.Value
				}
				if cookie.Name == "XSRF-TOKEN" {
					xsrfToken = cookie.Value
				}
			}
			return nil
		}),
	)
	if err != nil {
		return laravelSession, xsrfToken, fmt.Errorf("chromedp.Run: %w", err)
	}
	return laravelSession, xsrfToken, nil
}
