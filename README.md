# WebFireAutoTrainer
A tool that automatically completes all the courses on [WebFireTraining](https://webfiretraining.com/) with 100% score in less than 10 minutes depending on how fast your connection and device are.

### Introduction

At my weekend job we are suddenly made to complete all kinds of training, a lot of it double and very basic.

I heard someone spent more than 8 hours completing these, and I figured that bunnies can do this faster.

So I fired up Charles proxy and tried out one lesson/class of one course and noticed the site never actually validates your answers on the server side - aka it was `client authorative`.

In order of execution, the only requests that seemed to matter where:
* [GETPARAMS](https://github.com/BRUHItsABunny/webfireAutoTrainer/blob/main/_media/WFAT_class_getparams.png)
* [PUTPARAMS:I](https://github.com/BRUHItsABunny/webfireAutoTrainer/blob/main/_media/WFAT_class_init.png)
* [PUTPARAMS:P](https://github.com/BRUHItsABunny/webfireAutoTrainer/blob/main/_media/WFAT_class_progress.png)
* [EXITUA](https://github.com/BRUHItsABunny/webfireAutoTrainer/blob/main/_media/WFAT_class_exit.png)

Shortly after I started replaying these requests in Charles with the next lesson/class and saw how effective it was.

Do you want to pass all the quizzes in minutes rather than hours? Forget about answers, just use feeling!

https://user-images.githubusercontent.com/53124399/217393683-5e70b473-d13f-41ce-9421-4137149d1572.mp4
Source: [Uncle Roger on Giphy](https://giphy.com/clips/mrnigelng-cooking-show-uncle-roger-5HDCKts1Re8NxcwTSq)


### How does it work?
This program opens a Chrome window and interacts with it through the [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/), it then visits the login page and waits for you to login.

After it detects you have logged in it will start scraping all the values needed to replay the requests mentioned above.

Finally, it loops through all the courses and classes to replay those requests if you haven't completed them yet.
