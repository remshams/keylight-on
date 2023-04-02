use crate::keylight::{DiscoveredKeylight, Keylight};

pub trait KeylightFinder {
    fn discover(self) -> Vec<DiscoveredKeylight>;
}

pub struct KeylightControl {
    lights: Vec<Keylight>,
}

impl KeylightControl {
    pub fn discover_lights() -> KeylightControl {
        KeylightControl { lights: vec![] }
    }
}
