mod adapter;
mod control;
mod finder;
mod model;

pub use adapter::KeylightRestAdapter;
pub use control::{KeylightAdapter, KeylightControl, KeylightFinder};
pub use finder::ZeroConfKeylightFinder;
pub use model::{Keylight, KeylightMetadata, Light};

#[cfg(test)]
mod keylight_mocks {

    use super::{model::KeylightError, *};

    pub fn create_metadata_fixture() -> KeylightMetadata {
        KeylightMetadata {
            name: String::from("first"),
            ip: String::from("102.168.1.1"),
            port: 1234,
        }
    }

    pub fn create_metadata_list_fixture() -> Vec<KeylightMetadata> {
        vec![
            KeylightMetadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
            KeylightMetadata {
                name: String::from("second"),
                ip: String::from("102.168.1.2"),
                port: 4567,
            },
            KeylightMetadata {
                name: String::from("first"),
                ip: String::from("102.168.1.1"),
                port: 1234,
            },
        ]
    }

    pub fn create_lights_fixture() -> Vec<Light> {
        vec![
            Light {
                on: true,
                brightness: 2,
                temperature: 2,
            },
            Light {
                on: true,
                brightness: 4,
                temperature: 4,
            },
        ]
    }

    pub struct MockKeylightFinder {
        pub metadata: Vec<KeylightMetadata>,
    }

    impl KeylightFinder for MockKeylightFinder {
        type Output = Vec<KeylightMetadata>;

        fn discover(&self) -> Self::Output {
            self.metadata.clone()
        }
    }

    impl MockKeylightFinder {
        pub fn new(metadata: Vec<KeylightMetadata>) -> MockKeylightFinder {
            MockKeylightFinder { metadata }
        }
    }

    pub struct MockKeylightAdapter {
        pub lights: Result<Vec<Light>, KeylightError>,
        pub set_lights_result: Result<(), KeylightError>,
    }

    impl KeylightAdapter for MockKeylightAdapter {
        fn lights(&self, _ip: &str) -> Result<Vec<Light>, KeylightError> {
            self.lights.clone()
        }
        fn set_lights(&self, _ip: &str, _lights: &[Light]) -> Result<(), KeylightError> {
            self.set_lights_result.clone()
        }
    }

    impl MockKeylightAdapter {
        pub fn new(
            lights: Result<Vec<Light>, KeylightError>,
            set_lights_result: Option<Result<(), KeylightError>>,
        ) -> MockKeylightAdapter {
            MockKeylightAdapter {
                lights,
                set_lights_result: set_lights_result.unwrap_or(Ok(())),
            }
        }
    }

    pub fn create_keylight_fixture(
        keylight_adapter: &MockKeylightAdapter,
        lights: Option<Vec<Light>>,
    ) -> Keylight<MockKeylightAdapter> {
        Keylight::new(
            &keylight_adapter,
            create_metadata_fixture(),
            lights.or(None),
        )
    }
}
