# Quick-name
The Program will clip all videos in the FolderPath. With each video clip, it will use Google speech-to-text API to transcribe the audio. 
Once transcribed, GEMINI will analyze the transcription to find out the video's title.


For the google voice to text to work set the "GOOGLE_APPLICATION_CREDENTIALS="key.json" in env file and add the the json file to the project

For GIMINI to work set the GEMINI_API_KEY in the env file

The prompt can be changed to fit your needs


