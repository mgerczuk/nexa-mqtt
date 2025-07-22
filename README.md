# nexa-mqtt
![License](https://img.shields.io/github/license/mgerczuk/nexa-mqtt) ![GitHub last commit](https://img.shields.io/github/last-commit/mgerczuk/nexa-mqtt) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mgerczuk/nexa-mqtt)

`nexa-mqtt` is a standalone application designed to retrieve data and metrics from your Growatt NEXA 2000 home battery used in balcony power plants. It publishes this information to an MQTT broker, making it easily accessible for Home Assistant or other applications. It is a fork of https://github.com/mtrossbach/noah-mqtt.

The application features Home Assistant auto-discovery, allowing your NEXA devices to be automatically recognized and integrated with Home Assistant via the MQTT integration.

# ![HomeAssistant screenshot](/assets/ha-screenshot.png)

🌟 If you find my project helpful, please consider giving me a star on GitHub! Your support motivates me to improve and delve deeper into enhancing the project. Thank you!

---

# Configuration

`nexa-mqtt` supports three API modes:

*   **`app`**: (previous default) This mode utilizes the Shine App APIs.  These APIs offer faster data updates and support setting parameters. However, they are the least stable, as they are prone to change with new app updates.  They are also subject to strict rate limits, which may result in IP bans.
*   **`web`**: This mode uses the Growatt Website APIs. These APIs provide a more stable way to fetch data.  Setting parameters is not supported in this mode.
*   **`web+app`**: (current default) This mode combines the best of both worlds. It uses the Growatt Website APIs for data fetching (for stability) and the App APIs for setting parameters.


You can configure `nexa-mqtt` using the following environment variables:

| Environment Variable               | Description                                                                             | Default                        |
|:-----------------------------------|:----------------------------------------------------------------------------------------|:-------------------------------| 
| `LOG_LEVEL`                        | Sets the logging level of the application                                               | INFO                           |
| `POLLING_INTERVAL`                 | Time in seconds between fetching new status data                                        | 30                             |
| `BATTERY_DETAILS_POLLING_INTERVAL` | Time in seconds between fetching battery details (per battery SoC & temperature).       | 180                            |
| `PARAMETER_POLLING_INTERVAL`       | Time in seconds between fetching parameter data (system-output-power, charging limits). | 180                            |
| `GROWATT_API_MODE`                 | Growatt API mode, either `app`, `web`, `web+app`                                        | web+app                        |
| `GROWATT_USERNAME`                 | Your Growatt account username (required)                                                | -                              |
| `GROWATT_PASSWORD`                 | Your Growatt account password (required)                                                | -                              |
| `GROWATT_SERVER_URL_WEB`           | Growatt server url for web apis                                                         | https://openapi.growatt.com    |
| `GROWATT_SERVER_URL_APP`           | Growatt server url for app apis                                                         | https://server-api.growatt.com |
| `MQTT_HOST`                        | Address of your MQTT broker (required)                                                  | -                              |
| `MQTT_PORT`                        | Port number of your MQTT broker                                                         | 1883                           |
| `MQTT_CLIENT_ID`                   | Identifier for the MQTT client                                                          | nexa-mqtt                      |
| `MQTT_USERNAME`                    | Username for connecting to your MQTT broker                                             | -                              |
| `MQTT_PASSWORD`                    | Password for connecting to your MQTT broker                                             | -                              |
| `MQTT_TOPIC_PREFIX`                | Prefix for MQTT topics used by nexa-mqtt                                                | nexa2mqtt                      |
| `HOMEASSISTANT_TOPIC_PREFIX`       | Prefix for topics used by Home Assistant                                                | homeassistant                  |
| `HOMEASSISTANT_SWITCH_AS_SELECT`   | Publish 'switch' entities as 'select'. Set to 'True' for OpenHAB, see below             | false                          |

Adjust these settings to fit your environment and requirements.

---

# Data provided by nexa-mqtt

## Published Topics

The following MQTT topics are used by `nexa-mqtt` to publish data:

### 1. General Device Data
- **Topic:** `nexa2mqtt/{DEVICE_SERIAL}`
- **Description:** This topic contains general data about the device.
- **Example:** `nexa2mqtt/0ABC00AA15AA00AA`
- **Example Payload:**
```json
{
  "output_w": -398, // current output power in watts
                    // is negative when solar or battery power is delivered to the grid
                    // is positive when battery is charged from grid
  "solar_w": 102, // current solar generation power in watts
  "soc": 40, // current state of charge of the whole appliance
  "charge_w": 0, // current charging power in watts
  "discharge_w": 314, // current discharge power in watts
  "battery_num": 2, // number of batteries
  "generation_total_kwh": 319.8, // total energy generation
  "generation_today_kwh": 3.1, // engery generation today
  "work_mode": "load_first", // current work mode: load_first or battery_first
  "status": "on_grid" // connectivity status: offline, smart_self_use, fault, on_grid or off_grid
}
```

### 2. Battery Information
- **Topic:** `nexa2mqtt/{DEVICE_SERIAL}/BAT{BAT_NR}`
- **Description:** This topic contains information about the device's batteries. Replace `{BAT_NR}` with the battery number (e.g., BAT0, BAT1, BAT2, etc.).
- **Example:** `nexa2mqtt/0ABC00AA15AA00AA/BAT0`
- **Example Payload:**
```json
{
   "serial": "0ABC00AA15AA00AA", // battery serial number
   "soc": 42, // current state of charge of this battery
   "temp": 26 // current temperatur of this battery
}
```

### 3. Device Configuration
- **Topic:** `nexa2mqtt/{DEVICE_SERIAL}/parameters`
- **Description:** This topic contains the current configuration parameters of the device.
- **Example:** `nexa2mqtt/0ABC00AA15AA00AA/parameters`
- **Example Payload:**
```json
{
   "charging_limit": 100, // battery charging limit in percent, between 70 and 100
   "discharge_limit": 10, // battery discharge limit in percent, between 0 and 30
   "default_output_w": 150, // desired system AC output power in watts, between 0 and 1000
   "default_mode": "load_first", // or battery_first
   "allow_grid_charging": "OFF", // ON when battery may be charged from grid
   "grid_connection_control": "OFF", // ON for off-grid mode
   "ac_couple_power_control": "OFF" // ON for 1000W max. AC output. 
                                    // Note: this may be forbidden when connected to public grid!
}
```

The value pairs `charging_limit`, `discharge_limit` and `default_output_w`, `default_mode` are set together. If one of them is missing in the payload the cached previous value is used. A debounce timer of 500 ms is used to combine payloads with individual properties to a combined payload. That means that any setting of a value is executed after a delay of 500 ms.

## Setting Device Parameters

You can update the device's parameter settings by posting a message to the following topic:

- **Topic:** `nexa2mqtt/{DEVICE_SERIAL}/parameters/set`
- **Description:** Send configuration settings to this topic to update the device's parameters.
- **Example:** `nexa2mqtt/1234567890/parameters/set`
- **Example Payload:**
```json
{
   "charging_limit": 100, // battery charging limit in percent, between 70 and 100
   "discharge_limit": 9, // battery discharge limit in percent, between 0 and 30
   "default_output_w": 800, // system output power in watts, between 0 and 800 
   "default_mode": "load_first" // or battery_first
}
```


---

# Run the application standalone

## Option 1: Running `nexa-mqtt` with Docker

To run the latest version of `nexa-mqtt` using Docker, follow these steps:

1. **Install Docker**: Ensure Docker is installed on your system. You can download Docker Desktop from [Docker’s official website](https://www.docker.com/products/docker-desktop).

2. **Open a Terminal**:
   - **Windows**: Use Command Prompt or PowerShell.
   - **Linux/macOS**: Use the Terminal.

3. **Execute the Docker Command**: Run the following command, replacing the placeholders with your actual values:

   ```
   docker run --name nexa-mqtt -e GROWATT_USERNAME=myusername -e GROWATT_PASSWORD=mypassword -e MQTT_HOST=localhost -e MQTT_PORT=1883 ghcr.io/mgerczuk/nexa-mqtt:latest
   ```
   
- Replace myusername with your Growatt username.
- Replace mypassword with your Growatt password.
- Replace localhost with the hostname or IP address of your MQTT broker.
- Replace 1883 with the port number your MQTT broker uses (default is 1883).

The application will connect to your MQTT broker and retrieve all metrics and data for your NEXA devices.

## Option 2: Downloading and running a Debian package

1. **Download the deb package file**: Go to the [Releases](https://github.com/mgerczuk/nexa-mqtt/releases) page of the repository and download the .deb file for your operating system and system architecture.

2. **Install the package**

   ```sh
   sudo apt install -f <deb-file>
   ```

When there is an update simply download the new deb package file and install with the same install command.

nexa-mqtt is started and will be started automatically after a reboot. Check with `journalctl -t nexa-mqtt` if there are any problems, e.g. user name or password errors.

You can modify the environment variables by executing

   ```sh
   sudo systemctl edit nexa-mqtt
   sudo systemctl daemon-reload
   sudo systemctl restart nexa-mqtt
   ```
To uninstall the package execute `sudo apt remove nexa-mqtt`.

## Option 3: Downloading and running a prebuilt binary

If you prefer not to compile the binary yourself, you can download a prebuilt version:

1. **Download the Binary**: Go to the [Releases](https://github.com/mgerczuk/nexa-mqtt/releases) page of the repository and download the prebuilt binary for your operating system and system architecture.

2. **Extract the Binary**: If the binary is compressed (e.g., in a zip or tar file), extract it to a directory of your choice.

3. **Run the Application**: Open a terminal in the directory containing the binary and run it using the appropriate command for your OS, setting the necessary environment variables:

   - **Windows** (Command Prompt):

     ```sh
     set GROWATT_USERNAME=myusername
     set GROWATT_PASSWORD=mypassword
     set MQTT_HOST=localhost
     set MQTT_PORT=1883
     nexa-mqtt.exe
     ```

   - **Windows** (PowerShell):

     ```sh
     $env:GROWATT_USERNAME=„myusername“
     $env:GROWATT_PASSWORD=„mypassword“
     $env:MQTT_HOST=„localhost“
     $env:MQTT_PORT=„1883“
     .\nexa-mqtt.exe
     ```

   - **Linux/macOS**:

     ```sh
     GROWATT_USERNAME=myusername GROWATT_PASSWORD=mypassword MQTT_HOST=localhost MQTT_PORT=1883 ./nexa-mqtt
     ```

Again, replace `myusername`, `mypassword`, `localhost`, and `1883` with your actual Growatt account details and MQTT broker information.

## Option 4: Compiling the binary yourself

To compile the binary yourself, ensure you have Go installed on your machine:

1. **Install Go**: Download and install the latest version of Go from [the official Go website](https://golang.org/dl/).

2. **Clone the Repository**: Open a terminal and run the following command to clone the repository:
        
        git clone https://github.com/mgerczuk/nexa-mqtt.git
        cd nexa-mqtt

3. **Build the application**:

        go build -o nexa-mqtt cmd/nexa-mqtt/main.go

Afterwards follow the instructions for running the application from option 2.

---

# Integration into HomeAssistant

## Run standalone (Home Assistant Container, Home Assistant Core)
`nexa-mqtt` interacts with Home Assistant by publishing data from your Growatt NEXA 2000 home battery to an MQTT broker. This setup allows Home Assistant to subscribe to and integrate this data seamlessly into its ecosystem.

![Home Assistant Integration](./assets/nexa-mqtt-ha-dark.drawio.png#gh-dark-mode-only)
![Home Assistant Integration](./assets/nexa-mqtt-ha.drawio.png#gh-light-mode-only)

If you’re already using MQTT with other integrations like zigbee2mqtt or AhoyDTU, you already have the MQTT integration configured and active. In this case, you can skip step 1 and 2 as your existing setup should work with `nexa-mqtt`.

The following integration process for `nexa-mqtt` with Home Assistant works for all installation methods, regardless of how Home Assistant is installed—whether it’s through Home Assistant OS, Home Assistant Supervised, or Home Assistant Container. 

1. **Set Up an MQTT Broker**:  
   Ensure you have an MQTT broker running, such as [Mosquitto](https://mosquitto.org/), and that it’s accessible from both nexa-mqtt and Home Assistant.

2. **Check MQTT Integration in Home Assistant**:  
   - Navigate to **Settings** > **Devices & Services** in Home Assistant.
   - Click **Add Integration** and select „MQTT“.
   - Enter your MQTT broker details (hostname, port, username, password).
   - Test the connection to ensure it’s working correctly.

3. **Run nexa-mqtt**:  
   Start `nexa-mqtt` using the appropriate configuration for your MQTT broker.

4. **Verify Device Discovery**:  
   Check **Devices** and **Entities** under **Settings** > **Devices & Services** in Home Assistant to confirm that your Noah devices are automatically discovered.

By following these steps, `nexa-mqtt` will communicate with Home Assistant via your MQTT broker, also supporting automatic device discovery. If you already have MQTT set up, it should integrate seamlessly with your existing configuration.

## Run as Home Assistant add-on (Home Assistant OS, Home Assistant Supervised)

If you are using Home Assistant OS or Home Assistant Supervised you can run `nexa-mqtt` as a Home Assistant add-on, which provides seamless integration with your Home Assistant setup.
This option leverages the add-on system to manage and run `nexa-mqtt` directly on your Home Assistant instance.

#### Steps to Use the Home Assistant Add-on
0. **Prerequisite:**
   - Have the Mosquitto Add-on installed and running -or- have a separate MQTT running
   - Home Assistant MQTT integration enabled

1. **Add the Repository:**
   - Open your Home Assistant web interface.
   - Navigate to **Settings** > **Add-ons** > **Add-on Store**.
   - Click on the three-dot menu in the top right corner and select **Repositories**.
   - Add the following URL: `https://github.com/mgerczuk/hassio-addons`.

[![Open your Home Assistant instance and show the add add-on repository dialog with a specific repository URL pre-filled.](https://my.home-assistant.io/badges/supervisor_add_addon_repository.svg)](https://my.home-assistant.io/redirect/supervisor_add_addon_repository/?repository_url=https%3A%2F%2Fgithub.com%2Fmgerczuk%2Fhassio-addons)

2. **Install the Add-on:**
   - Search for the `nexa-mqtt` add-on within the Add-on Store.
   - Click on the add-on and select **Install**.

3. **Configure the Add-on:**
   - After installation, configure the add-on settings by providing your **Growatt username** and **Growatt password** and setup the other options as needed.
   - If you do not use the Mosquitto Add-on, please also define your MQTT settings
4. **Start the Add-on:**
   - Click **Start** to launch the `nexa-mqtt` add-on.

The Home Assistant add-on provides an easy and integrated way to run `nexa-mqtt`, allowing you to manage it directly from the Home Assistant interface.

For more detailed information and updates, visit the [repository](https://github.com/mgerczuk/hassio-addons).

---

# Integration into OpenHAB

`nexa-mqtt` interacts with OpenHAB by publishing data from your Growatt NEXA 2000 home battery to an MQTT broker. This setup allows OpenHAB to subscribe to and integrate this data seamlessly into its ecosystem.

`nexa-mqtt` uses Home Assistant auto-discover. OpenHAB has some problems with this, especially you cannot properly set a '**Switch**' channel data imported from HA auto-discovery. Be sure to set the environment variable `HOMEASSISTANT_SWITCH_AS_SELECT` to `True` when you use OpenHAB. The 'switch' entities will then be imported as **String** channels.

![Home Assistant Integration](./assets/nexa-mqtt-openhab-dark.drawio.png#gh-dark-mode-only)
![Home Assistant Integration](./assets/nexa-mqtt-openhab.drawio.png#gh-light-mode-only)

1. **Set Up an MQTT Broker**:  
   Ensure you have an MQTT broker running, such as [Mosquitto](https://mosquitto.org/), and that it’s accessible from both nexa-mqtt and OpenHAB.

2. **Check MQTT Integration in OpenHAB**: 
   If not already done install the MQTT Binding from the Add-on Store

3. **Run nexa-mqtt**:  
   Start `nexa-mqtt` using the appropriate configuration for your MQTT broker and the environment variable `HOMEASSISTANT_SWITCH_AS_SELECT` to `True`.

4. **Device Discovery**:  
   OpenHAB should now show a new thing in the inbox on the **Settings**/**Things** page. Click on the inbox, select the entry an chose **Add as Thing**. Enter a name for the new thing. The thing should now be **online** and you can start linking items to the channels.

   In case the thing does not go online first check the MQTT broker and the MQTT binding. If that is Ok, try to disable/enable the thing and/or restart `nexa-mqtt`.

