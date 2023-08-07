use keylight_control::keylight::{Keylight, KeylightError, KeylightRestAdapter, LightCommand};
use tauri::State;

use crate::model::AppState;

#[tauri::command]
pub fn discover_keylights(state: State<AppState>) -> Result<Vec<Keylight>, KeylightError> {
    let mut keylight_control = state
        .keylight_control
        .lock()
        .map_err(|_err| KeylightError::CommandError(String::from("test")))?;
    keylight_control.load_keylights()?;
    load_lights(&mut keylight_control.lights, &state.adapter)
}

#[tauri::command]
pub fn refresh_lights(state: State<AppState>) -> Result<Vec<Keylight>, KeylightError> {
    let mut keylight_control = state
        .keylight_control
        .lock()
        .map_err(|_err| KeylightError::CommandError(String::from("test")))?;
    load_lights(&mut keylight_control.lights, &state.adapter)
}

fn load_lights(
    lights: &mut Vec<Keylight>,
    adapter: &KeylightRestAdapter,
) -> Result<Vec<Keylight>, KeylightError> {
    for light in lights.iter_mut() {
        light.lights(adapter)?;
    }
    Ok(lights.clone())
}

#[tauri::command]
pub fn set_light(state: State<AppState>, command: LightCommand) -> Result<(), KeylightError> {
    let mut keylight_control = state
        .keylight_control
        .lock()
        .map_err(|_err| KeylightError::CommandError(String::from("test")))?;
    let light = keylight_control.find_keylight_mut(&command.id);
    match light {
        Some(light) => light.set_light(command, &state.adapter),
        None => Err(KeylightError::KeylightDoesNotExist(command.id)),
    }
}