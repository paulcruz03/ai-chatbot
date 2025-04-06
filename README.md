# Chatbot AI project

Sample Project to upskill AI Prompt development. For this project, I used Go lang for fast development compared to NodeJS. For Ai, I used Google's `gemini-2.0-flash-exp` model. 

Sharing this for those who want to upskill also. Setup:
1. Install Go lang. Checkout their [installation guide](https://go.dev/dl/)
2. Clone this project
3. Install the go packages
- You could use this for package update
```
go mod tidy
```
- But if you want to just install the packages
```
go mod download
```
4. Setup `.env` file for the Google's Gemini Key
```
GEMINI_API_KEY=${key}
```
5. Run the app and open `ws_tester.html` on the project folder
```
go run .
```

- to build the app, run this cmd:
```
docker build -t ${image name} .
```

- to run the build image, run this cmd:
```
docker run -p 8080:8080 ${image name}
```
then visit http://localhost:8080/

### To do
- Add External Integration (Separate Project)