# National Park Service Webcam Gallery

This is a simple Go/Gin web app to display a gallery of the various National
Park Service webcams as discovered through the NPS API. This was a recreational
programming project done on a Friday evening to explore Gin a bit.

I've hosted the application publicly using Railway:

<https://npswebcams-production.up.railway.app/>

[![Deploy on Railway](https://img.shields.io/badge/Deploy%20on-Railway-%230A0A0A.svg)](https://railway.app/)

## Running locally

An API key can be obtained from the [NPS Developer Resources "Get Started"
page](https://www.nps.gov/subjects/developer/get-started.htm).

## To run

```bash
export NPS_API_KEY=<your key>
go run .
```

Then visit <http://localhost:8080/> in your browser.

## Screenshot

![Screenshot of web gallery](readme_assets/preview.png)

## Diagram

![Sequence diagram](readme_assets/mermaid_diagram.png)

## License

MIT
