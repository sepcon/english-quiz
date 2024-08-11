package main

import (
	"encoding/json"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/sepcon/quizprob/pkg/model/quiz_service"
)

var questionsJson = `
[
  {
    "id": "q1",
    "ask": "What is the capital of France?",
    "options": {
      "A": "Berlin",
      "B": "Madrid",
      "C": "Paris",
      "D": "Rome"
    }
  },
  {
    "id": "q2",
    "ask": "Which planet is known as the Red Planet?",
    "options": {
      "A": "Earth",
      "B": "Mars",
      "C": "Jupiter",
      "D": "Saturn"
    }
  },
  {
    "id": "q3",
    "ask": "What is the largest ocean on Earth?",
    "options": {
      "A": "Atlantic Ocean",
      "B": "Indian Ocean",
      "C": "Arctic Ocean",
      "D": "Pacific Ocean"
    }
  },
  {
    "id": "q4",
    "ask": "Who wrote 'Romeo and Juliet'?",
    "options": {
      "A": "William Shakespeare",
      "B": "Charles Dickens",
      "C": "Mark Twain",
      "D": "Jane Austen"
    }
  },
  {
    "id": "q5",
    "ask": "What is the chemical symbol for water?",
    "options": {
      "A": "H2O",
      "B": "O2",
      "C": "CO2",
      "D": "NaCl"
    }
  },
  {
    "id": "q6",
    "ask": "Which country is known as the Land of the Rising Sun?",
    "options": {
      "A": "China",
      "B": "Japan",
      "C": "Thailand",
      "D": "India"
    }
  },
  {
    "id": "q7",
    "ask": "What is the hardest natural substance on Earth?",
    "options": {
      "A": "Gold",
      "B": "Iron",
      "C": "Diamond",
      "D": "Platinum"
    }
  },
  {
    "id": "q8",
    "ask": "Who painted the Mona Lisa?",
    "options": {
      "A": "Vincent van Gogh",
      "B": "Pablo Picasso",
      "C": "Leonardo da Vinci",
      "D": "Claude Monet"
    }
  },
  {
    "id": "q9",
    "ask": "What is the smallest prime number?",
    "options": {
      "A": "0",
      "B": "1",
      "C": "2",
      "D": "3"
    }
  },
  {
    "id": "q10",
    "ask": "Which element has the chemical symbol 'O'?",
    "options": {
      "A": "Oxygen",
      "B": "Gold",
      "C": "Osmium",
      "D": "Oganesson"
    }
  },
  {
    "id": "q11",
    "ask": "What is the largest planet in our solar system?",
    "options": {
      "A": "Earth",
      "B": "Mars",
      "C": "Jupiter",
      "D": "Saturn"
    }
  },
  {
    "id": "q12",
    "ask": "Who is known as the Father of Computers?",
    "options": {
      "A": "Albert Einstein",
      "B": "Isaac Newton",
      "C": "Charles Babbage",
      "D": "Nikola Tesla"
    }
  },
  {
    "id": "q13",
    "ask": "What is the main ingredient in guacamole?",
    "options": {
      "A": "Tomato",
      "B": "Avocado",
      "C": "Onion",
      "D": "Pepper"
    }
  },
  {
    "id": "q14",
    "ask": "Which country is the largest by land area?",
    "options": {
      "A": "Canada",
      "B": "China",
      "C": "Russia",
      "D": "United States"
    }
  },
  {
    "id": "q15",
    "ask": "What is the speed of light?",
    "options": {
      "A": "300,000 km/s",
      "B": "150,000 km/s",
      "C": "450,000 km/s",
      "D": "600,000 km/s"
    }
  },
  {
    "id": "q16",
    "ask": "Who discovered penicillin?",
    "options": {
      "A": "Marie Curie",
      "B": "Alexander Fleming",
      "C": "Louis Pasteur",
      "D": "Gregor Mendel"
    }
  },
  {
    "id": "q17",
    "ask": "What is the capital of Australia?",
    "options": {
      "A": "Sydney",
      "B": "Melbourne",
      "C": "Canberra",
      "D": "Brisbane"
    }
  },
  {
    "id": "q18",
    "ask": "Which gas is most abundant in the Earth's atmosphere?",
    "options": {
      "A": "Oxygen",
      "B": "Carbon Dioxide",
      "C": "Nitrogen",
      "D": "Hydrogen"
    }
  },
  {
    "id": "q19",
    "ask": "What is the largest mammal in the world?",
    "options": {
      "A": "Elephant",
      "B": "Blue Whale",
      "C": "Giraffe",
      "D": "Hippopotamus"
    }
  },
  {
    "id": "q20",
    "ask": "Who wrote 'To Kill a Mockingbird'?",
    "options": {
      "A": "Harper Lee",
      "B": "F. Scott Fitzgerald",
      "C": "Ernest Hemingway",
      "D": "J.D. Salinger"
    }
  },
  {
    "id": "q21",
    "ask": "What is the smallest country in the world?",
    "options": {
      "A": "Monaco",
      "B": "San Marino",
      "C": "Vatican City",
      "D": "Liechtenstein"
    }
  },
  {
    "id": "q22",
    "ask": "Which planet is closest to the sun?",
    "options": {
      "A": "Venus",
      "B": "Earth",
      "C": "Mercury",
      "D": "Mars"
    }
  },
  {
    "id": "q23",
    "ask": "What is the main language spoken in Brazil?",
    "options": {
      "A": "Spanish",
      "B": "Portuguese",
      "C": "French",
      "D": "English"
    }
  },
  {
    "id": "q24",
    "ask": "Who developed the theory of relativity?",
    "options": {
      "A": "Isaac Newton",
      "B": "Galileo Galilei",
      "C": "Albert Einstein",
      "D": "Niels Bohr"
    }
  },
  {
    "id": "q25",
    "ask": "What is the tallest mountain in the world?",
    "options": {
      "A": "K2",
      "B": "Kangchenjunga",
      "C": "Mount Everest",
      "D": "Lhotse"
    }
  },
  {
    "id": "q26",
    "ask": "Which organ is responsible for pumping blood throughout the body?",
    "options": {
      "A": "Lungs",
      "B": "Liver",
      "C": "Kidneys",
      "D": "Heart"
    }
  },
  {
    "id": "q27",
    "ask": "What is the largest desert in the world?",
    "options": {
      "A": "Sahara Desert",
      "B": "Arabian Desert",
      "C": "Gobi Desert",
      "D": "Antarctic Desert"
    }
  },
  {
    "id": "q28",
    "ask": "Who was the first person to walk on the moon?",
    "options": {
      "A": "Yuri Gagarin",
      "B": "Buzz Aldrin",
      "C": "Neil Armstrong",
      "D": "Michael Collins"
    }
  },
  {
    "id": "q29",
    "ask": "What is the main ingredient in traditional Japanese miso soup?",
    "options": {
      "A": "Soybeans",
      "B": "Rice",
      "C": "Fish",
      "D": "Seaweed"
    }
  },
  {
    "id": "q30",
    "ask": "Which country hosted the 2016 Summer Olympics?",
    "options": {
      "A": "China",
      "B": "Brazil",
      "C": "United Kingdom",
      "D": "Japan"
    }
  }
]`

var correctAnswersJson = `{
  "q1": {
    "choice": "C",
    "score": 1
  },
  "q2": {
    "choice": "B",
    "score": 1
  },
  "q3": {
    "choice": "D",
    "score": 1
  },
  "q4": {
    "choice": "A",
    "score": 1
  },
  "q5": {
    "choice": "A",
    "score": 1
  },
  "q6": {
    "choice": "B",
    "score": 1
  },
  "q7": {
    "choice": "C",
    "score": 1
  },
  "q8": {
    "choice": "C",
    "score": 1
  },
  "q9": {
    "choice": "C",
    "score": 1
  },
  "q10": {
    "choice": "A",
    "score": 1
  },
  "q11": {
    "choice": "C",
    "score": 1
  },
  "q12": {
    "choice": "C",
    "score": 1
  },
  "q13": {
    "choice": "B",
    "score": 1
  },
  "q14": {
    "choice": "C",
    "score": 1
  },
  "q15": {
    "choice": "A",
    "score": 1
  },
  "q16": {
    "choice": "B",
    "score": 1
  },
  "q17": {
    "choice": "C",
    "score": 1
  },
  "q18": {
    "choice": "C",
    "score": 1
  },
  "q19": {
    "choice": "B",
    "score": 1
  },
  "q20": {
    "choice": "A",
    "score": 1
  },
  "q21": {
    "choice": "C",
    "score": 1
  },
  "q22": {
    "choice": "C",
    "score": 1
  },
  "q23": {
    "choice": "B",
    "score": 1
  },
  "q24": {
    "choice": "C",
    "score": 1
  },
  "q25": {
    "choice": "C",
    "score": 1
  },
  "q26": {
    "choice": "D",
    "score": 1
  },
  "q27": {
    "choice": "D",
    "score": 1
  },
  "q28": {
    "choice": "C",
    "score": 1
  },
  "q29": {
    "choice": "A",
    "score": 1
  },
  "q30": {
    "choice": "B",
    "score": 1
  }
}
`

type Questions = []quiz_service.Question
type CorrectAnswer struct {
	Choice quiz_service.ChoiceType `json:"choice,omitempty"`
	Score  quiz.ScoreType          `json:"score,omitempty"`
}
type CorrectAnswers = map[quiz_service.QuestionIDType]CorrectAnswer

func MakeQuestionBank() (questions Questions, answers CorrectAnswers) {
	json.Unmarshal([]byte(questionsJson), &questions)
	json.Unmarshal([]byte(correctAnswersJson), &answers)
	return
}
