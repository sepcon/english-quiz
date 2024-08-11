# Real-Time Vocabulary Quiz Coding Challenge

## Overview
Welcome to the Real-Time Quiz coding challenge! Your task is to create a technical solution for a real-time quiz feature for an English learning application. This feature will allow users to answer questions in real-time, compete with others, and see their scores updated live on a leaderboard.



## Prerequisites

*   **Go:** This project requires Go to be installed. You can download and install it from the official website: https://golang.org/dl/

## Getting Started

1.  **Clone the repository:**

    ```bash
    git clone [https://github.com/sepcon/english-quiz.git](https://github.com/sepcon/english-quiz.git)
    cd english-quiz
    ```

2.  **Install dependencies (if any):**

    *   If your project uses any external Go modules, navigate to the project root and run:

    ```bash
    go mod download
    ```

3.  **Build and run the applications:** 
    ```bash
    go run ./cmd/events  & #EventService
    go run ./cmd/scoring & #ScoringService
    go run ./cmd/quiz    & #QuizService
    ```

    *   This will start the required services and detach them from terminal
4. **Test the app with CLI client**
    * Open each terminal for a client, to join the quiz with the `userid=someone`, lets run the command
    ```bash
   go run ./cmd/client -userid=someone
    ```
   * Please execute the above command with different `userid` for seeing the realtime update results 