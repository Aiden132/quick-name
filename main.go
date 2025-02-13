package main

import (
	"context"
	"fmt"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/mowshon/moviego"
	"google.golang.org/api/option"
)

var FolderPath string = "videos/"



func TranscribeAudio(path string) string {
	Transcribed := ""
	ctx := context.Background()

	// Instantiates a client
	client, err := speech.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	data, err := os.ReadFile(path) // Read the file content
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	config := speechpb.RecognitionConfig{
		Model:                 "latest_long",
		Encoding:              speechpb.RecognitionConfig_MP3,
		SampleRateHertz:       48000,
		AudioChannelCount:     2,
		EnableWordTimeOffsets: true,
		EnableWordConfidence:  true,
		LanguageCode:          "en-US",
	}

	audio := speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Content{Content: data}, // Send file content
	}

	request := speechpb.LongRunningRecognizeRequest{
		Config: &config,
		Audio:  &audio,
	}

	op, err := client.LongRunningRecognize(ctx, &request)
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		log.Fatalf("failed to wait for long-running operation: %v", err)
	}
	// Prints the results
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			Transcribed += alt.Transcript + "\n"
			// fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
	return Transcribed
}

func nameShow() {
	godotenv.Load()

	ctx := context.Background()

	prompt := `
		This is an episode of SHOWNAME Season X use this transcription
		and tell me what episode it could be. put it in the format SHOWNAMESXXEXX and then tell me why
	`

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY") // SET UP THE API KEY IN A ENV FILE

	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(TranscribeAudio("temp/test.mp3")+prompt))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}

func ClipVideo(path string) {
	folder, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Can open folder: %v %v\v", path, err)
		return
	}
	for _, file := range folder {
		if file.IsDir() {
			ClipVideo(path + file.Name() + "/")
			continue
		}
		videoString := FolderPath + file.Name()
		first, _ := moviego.Load(videoString)

		err := first.SubClip(0, 120).Output("temp/test.mp3").Run()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(videoString, "Done")
		nameShow()
		os.Remove("temp/test.mp3")
	}
}

func main() {

	ClipVideo(FolderPath)

}
