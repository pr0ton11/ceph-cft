package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

// Allowed section prefixes in the configuration file
var ALLOWED_SECTION_PREFIX []string = []string{"global", "mon", "mgr", "osd", "mds", "client"}

// This application is a simple utility to manage the Ceph configuration file by defining configuration values as environment variables.
// The application will load the current configuration file, update the values based on the environment variables, and write the updated configuration file.
// Format of the environment variables: CEPH_{SECTION}_{KEY}={VALUE}
// Example: CEPH_GLOBAL_CLUSTER_NETWORK=10.0.0.0/24
// If the section is not defined, the application will use the default section (GLOBAL).

func main() {
	// Retrieve the configuration path
	configPath, ok := os.LookupEnv("CFT_CONFIG_PATH")
	// Otherwise, use the default path to the configuration file
	if !ok {
		configPath = "/etc/ceph/ceph.conf"
	}

	// Load the current configuration
	current, err := ini.Load(configPath)
	if err != nil {
		log.Printf("Error loading the current configuration file %s: %v", configPath, err)
		current = ini.Empty()
	}

	// Iterate over the environment variables
	for _, env := range os.Environ() {
		// Workaround for the issue that dots are not allowed in environment variables
		// Replace the double underscore with a dot
		env = strings.Replace(env, "__", ".", -1)
		// Convert the environment variable to lower case for simplicity
		env = strings.ToLower(env)
		// Split the environment variable into key and value
		assignment := strings.Split(env, "=")
		// Assign the key
		key := assignment[0]
		// Assign the value
		value := assignment[1]

		// Check if the environment variable is a ceph configuration variable
		if strings.HasPrefix(key, "ceph_") {
			log.Printf("Found key: %s with value: %s", key, value)
			// Split the key further into section and configuration key
			keyParts := strings.Split(key, "_")
			// Set flag to detect if the section is defined
			sectionDetected := false
			// Check if the section is defined
			detectedSection := "global"
			// Only allow prefixes if they are defined in the allowed section prefixes
			for _, section := range ALLOWED_SECTION_PREFIX {
				if strings.HasPrefix(keyParts[1], section) {
					sectionDetected = true
					detectedSection = keyParts[1]
				}
			}
			// Output a warning if the section is not defined
			if !sectionDetected {
				log.Printf("No section deteced in key %s, using default %s", key, detectedSection)
			} else {
				log.Printf("Detected section: %s for key %s", detectedSection, key)
			}
			// Cleanup the key by removing both the prefix and the section if defined
			if sectionDetected {
				key = strings.Join(keyParts[2:], "_")
			} else {
				key = strings.Join(keyParts[1:], "_")
			}

			// If the value contains a space, wrap it in quotes
			if strings.Contains(value, " ") {
				if !strings.HasPrefix(value, "\"") && !strings.HasSuffix(value, "\"") {
					value = "\"" + value + "\""
				}
			}

			// Update the configuration
			log.Printf("Setting key: %s with value: %s in section: %s", key, value, detectedSection)
			current.Section(detectedSection).Key(key).SetValue(value)
		}

		// Write the updated configuration
		err = current.SaveTo(configPath)
		if err != nil {
			log.Printf("Error saving the updated configuration file %s: %v", configPath, err)
		}
	}

}
