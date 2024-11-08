# CouchTube

CouchTube is a self-hostable YouTube frontend designed to simulate a TV channel experience. It dynamically loads YouTube videos from a predefined list of channels and schedules playback based on the current time. Users can also submit their custom video lists through a JSON file URL.

I hope CouchTube will be a community-driven project, where people can create and share their own channel JSON lists. Feel free to submit pull requests with new channel lists to enhance the default set in this repo.

The project is in its early days of development. There will probably be many issues and bugs. Please use [issues](https://github.com/ozencb/couchtube/issues) to report them.

CouchTube is inspired by [ytch.xyz](https://ytch.xyz/).


---

## Getting Started

### Using Docker

To run CouchTube using Docker, you can either use a `docker-compose.yml` file:

```yaml
version: "3.8"

services:
  couchtube:
    image: ghcr.io/ozencb/couchtube:latest
    container_name: couchtube_app
    ports:
      - "8081:8081"
    restart: unless-stopped
```

Or start it directly with a `docker run` command:

```sh
docker run -d \
  --name couchtube_app \
  -p 8081:8081 \
  --restart unless-stopped \
  ghcr.io/ozencb/couchtube:latest
```

### Building From Source

Ensure you have Golang 1.22 or higher installed.

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/ozencb/couchtube.git
   cd couchtube
   ```

2. **Install Go Dependencies**:
   ```sh
   go mod tidy
   ```

3. **Run the Application**:
   ```sh
   go run main.go
   ```
   The server will start on `http://localhost:8081`.

4. **Access the Application**:
   Open a browser and go to `http://localhost:8081` to access CouchTube.

On the first run, CouchTube will create a `couchtube.db` SQLite database file, initialize necessary tables, and populate them with any default channels found in `default-channels.json`.

---

## Usage

### Custom JSON Format for Channel and Video Lists

You can create custom JSON files to specify channels and video lists.

#### JSON Structure

Create your JSON file using the following format:

```json
{
  "channels": [
    {
      "name": "Channel Name",
      "videos": [
        {
          "url": "https://www.youtube.com/watch?v=VIDEO_ID",
          "segmentStart": 10,
          "segmentEnd": 300
        },
        {
          "url": "https://www.youtube.com/watch?v=ANOTHER_VIDEO_ID",
          "segmentStart": 0,
          "segmentEnd": 200
        }
      ]
    },
    {
      "name": "Another Channel Name",
      "videos": [
        {
          "url": "https://www.youtube.com/watch?v=DIFFERENT_VIDEO_ID",
          "segmentStart": 0,
          "segmentEnd": 150
        }
      ]
    }
  ]
}
```

#### Field Descriptions

- **channels**: An array of channel objects. Each channel contains:
  - **name**: The channel name (string).
  - **videos**: An array of video objects containing:
    - **url**: The URL of the YouTube video (required).
    - **segmentStart**: The start time (in seconds) within the video where playback begins.
    - **segmentEnd**: The end time (in seconds) within the video where playback ends.

#### Example JSON File

Save your custom JSON file using the above structure or make it accessible through a URL.

### Uploading Custom JSON

Within the CouchTube application, click the settings icon (gear icon) to submit a URL pointing to your custom JSON file. This URL should contain the JSON with channels and videos you want CouchTube to use.

---

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

---

## License

This project is licensed under the GNU General Public License.

---

## Additional Notes

1. **Database**: Ensure you don't have an existing `couchtube.db` file to avoid database conflicts.
2. **Error Handling**: The app includes basic error handling but might need enhancements for production use.
3. **Video Availability**: Videos marked private, restricted, or disabled for embedding may not play. CouchTube attempts to handle such errors by skipping to the next available video.

---

Enjoy using CouchTube!
```